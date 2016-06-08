package qfetch

import (
	"bufio"
	"fmt"
	"github.com/qiniu/api.v6/auth/digest"
	"github.com/qiniu/api.v6/conf"
	"github.com/qiniu/api.v6/rs"
	"github.com/qiniu/rpc"
	"github.com/syndtr/goleveldb/leveldb"
	"log"
	"net/url"
	"os"
	"strings"
	"sync"
)

var once sync.Once
var fetchTasks chan func()

func doFetch(tasks chan func()) {
	for {
		task := <-tasks
		task()
	}
}

func Fetch(job, filePath, bucket, accessKey, secretKey string, worker int, zone, logFile string) {
	//open file
	fh, openErr := os.Open(filePath)
	if openErr != nil {
		fmt.Println("Open resource file error,", openErr)
		return
	}
	defer fh.Close()

	logFh, openErr := os.Create(logFile)
	if openErr != nil {
		log.SetOutput(os.Stdout)
	} else {
		log.SetOutput(logFh)
		defer logFh.Close()
	}

	//open leveldb
	proFile := fmt.Sprintf(".%s.job", job)
	notFoundFile := fmt.Sprintf(".%s.404.job", job)
	ldb, lerr := leveldb.OpenFile(proFile, nil)

	if lerr != nil {
		fmt.Println("Open fetch progress file error,", lerr)
		return
	}
	defer ldb.Close()

	ldbNotFound, lerr := leveldb.OpenFile(notFoundFile, nil)
	if lerr != nil {
		fmt.Println("Open fetch not found file error,", lerr)
	}
	defer ldbNotFound.Close()

	//fetch prepare
	switch zone {
	case "bc":
		conf.IO_HOST = "http://iovip-z1.qbox.me"
	case "aws":
		conf.IO_HOST = "http://iovip.gdipper.com"
	default:
		conf.IO_HOST = "http://iovip.qbox.me"
	}

	mac := digest.Mac{
		accessKey, []byte(secretKey),
	}
	client := rs.New(&mac)

	once.Do(func() {
		fetchTasks = make(chan func(), worker)
		for i := 0; i < worker; i++ {
			go doFetch(fetchTasks)
		}
	})

	fetchWaitGroup := sync.WaitGroup{}

	//scan each line
	bReader := bufio.NewScanner(fh)
	bReader.Split(bufio.ScanLines)
	for bReader.Scan() {
		line := strings.TrimSpace(bReader.Text())
		if line == "" {
			continue
		}

		items := strings.Split(line, "\t")
		if !(len(items) == 1 || len(items) == 2) {
			log.Printf("Invalid resource line `%s`", line)
			continue
		}

		resUrl := items[0]
		resKey := ""

		if len(items) == 1 {
			resUri, pErr := url.Parse(resUrl)
			if pErr != nil {
				log.Printf("Invalid resource url `%s`", resUrl)
				continue
			}
			resKey = resUri.Path
			if strings.HasPrefix(resKey, "/") {
				resKey = resKey[1:]
			}
		} else if len(items) == 2 {
			resKey = items[1]
		}

		//check from leveldb whether it is done
		val, exists := ldb.Get([]byte(resUrl), nil)
		if exists == nil && string(val) == resKey {
			continue
		}

		nfVal, nfExists := ldbNotFound.Get([]byte(resUrl), nil)
		if nfExists == nil && string(nfVal) == resKey {
			continue
		}

		//otherwise fetch it
		fetchWaitGroup.Add(1)
		fetchTasks <- func() {
			defer fetchWaitGroup.Done()

			_, fErr := client.Fetch(nil, bucket, resKey, resUrl)
			if fErr == nil {
				ldb.Put([]byte(resUrl), []byte(resKey), nil)
			} else {
				if v, ok := fErr.(*rpc.ErrorInfo); ok {
					if v.Code == 404 {
						ldbNotFound.Put([]byte(resUrl), []byte(resKey), nil)
					}
					log.Printf("Fetch %s error due to `%s`", resUrl, v.Err)
				} else {
					log.Printf("Fetch %s error due to `%s`", resUrl, fErr)
				}
			}
		}
	}

	fetchWaitGroup.Wait()
}
