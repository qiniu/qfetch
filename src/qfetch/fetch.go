package qfetch

import (
	"bufio"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
	"sync"

	"github.com/qiniu/api.v6/auth/digest"
	"github.com/qiniu/api.v6/rs"
	"github.com/qiniu/rpc"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

var once sync.Once
var fetchTasks chan func()

func doFetch(tasks chan func()) {
	for {
		task := <-tasks
		task()
	}
}

func Fetch(mac *digest.Mac, job string, checkExists bool, fileListPath, bucket string, logFile string, worker int) {
	//open file list to fetch
	fh, openErr := os.Open(fileListPath)
	if openErr != nil {
		fmt.Println("Open resource file error,", openErr)
		return
	}
	defer fh.Close()

	//try open log file
	if logFile != "" {
		logFh, openErr := os.Create(logFile)
		if openErr != nil {
			log.SetOutput(os.Stdout)
		} else {
			log.SetOutput(logFh)
			defer logFh.Close()
		}
	} else {
		log.SetOutput(os.Stdout)
		defer os.Stdout.Sync()
	}

	ldbWOpt := opt.WriteOptions{
		Sync: true,
	}

	//open leveldb success and not found
	successLdbPath := fmt.Sprintf(".%s.job", job)
	notFoundLdbPath := fmt.Sprintf(".%s.404.job", job)

	successLdb, lerr := leveldb.OpenFile(successLdbPath, nil)
	if lerr != nil {
		fmt.Println("Open fetch progress file error,", lerr)
		return
	}
	defer successLdb.Close()

	notFoundLdb, lerr := leveldb.OpenFile(notFoundLdbPath, nil)
	if lerr != nil {
		fmt.Println("Open fetch not found file error,", lerr)
		return
	}
	defer notFoundLdb.Close()

	client := rs.New(mac)
	//init work group
	once.Do(func() {
		fetchTasks = make(chan func(), worker)
		for i := 0; i < worker; i++ {
			go doFetch(fetchTasks)
		}
	})

	fetchWaitGroup := sync.WaitGroup{}

	//scan each line and add task
	bReader := bufio.NewScanner(fh)
	bReader.Split(bufio.ScanLines)
	for bReader.Scan() {
		line := strings.TrimSpace(bReader.Text())
		if line == "" {
			continue
		}

		items := strings.Split(line, "\t")
		if !(len(items) == 1 || len(items) == 2) {
			log.Printf("Invalid resource line `%s`\n", line)
			continue
		}

		resUrl := items[0]
		resKey := ""

		if len(items) == 1 {
			resUri, pErr := url.Parse(resUrl)
			if pErr != nil {
				log.Printf("Invalid resource url `%s`\n", resUrl)
				continue
			}
			resKey = resUri.Path
			if strings.HasPrefix(resKey, "/") {
				resKey = resKey[1:]
			}
		} else if len(items) == 2 {
			resKey = items[1]
		}

		//check from leveldb success whether it is done
		val, exists := successLdb.Get([]byte(resUrl), nil)
		if exists == nil && string(val) == resKey {
			log.Printf("Skip url fetched `%s` => `%s`\n", resUrl, resKey)
			continue
		}

		//check from leveldb not found whether it meet 404
		nfVal, nfExists := notFoundLdb.Get([]byte(resUrl), nil)
		if nfExists == nil && string(nfVal) == resKey {
			log.Printf("Skip url 404 `%s` => `%s`\n", resUrl, resKey)
			continue
		}

		//check whether file already exists in bucket
		if checkExists {
			if entry, err := client.Stat(nil, bucket, resKey); err == nil && entry.Hash != "" {
				successLdb.Put([]byte(resUrl), []byte(resKey), &ldbWOpt)
				log.Printf("Skip url exists `%s` => `%s`\n", resUrl, resKey)
				continue
			}
		}

		//otherwise fetch it
		fetchWaitGroup.Add(1)
		fetchTasks <- func() {
			defer fetchWaitGroup.Done()

			_, fErr := client.Fetch(nil, bucket, resKey, resUrl)
			if fErr == nil {
				successLdb.Put([]byte(resUrl), []byte(resKey), nil)
			} else {
				if v, ok := fErr.(*rpc.ErrorInfo); ok {
					if v.Code == 404 {
						notFoundLdb.Put([]byte(resUrl), []byte(resKey), &ldbWOpt)
					}
					log.Printf("Fetch `%s` error due to `%s`\n", resUrl, v.Err)
				} else {
					log.Printf("Fetch `%s` error due to `%s`\n", resUrl, fErr)
				}
			}
		}
	}

	//wait for all the fetch done
	fetchWaitGroup.Wait()
}
