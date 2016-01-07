package qfetch

import (
	"bufio"
	"fmt"
	"github.com/qiniu/api.v6/auth/digest"
	"github.com/qiniu/api.v6/conf"
	"github.com/qiniu/api.v6/rs"
	"github.com/syndtr/goleveldb/leveldb"
	"net/url"
	"os"
	"strings"
	"sync"
)

func Fetch(job, filePath, bucket, accessKey, secretKey string, worker int, zone string) {
	//open file
	fh, openErr := os.Open(filePath)
	if openErr != nil {
		fmt.Println("Open resource file error,", openErr)
		return
	}
	defer fh.Close()

	//open leveldb
	proFile := fmt.Sprintf(".%s.job", job)
	ldb, lerr := leveldb.OpenFile(proFile, nil)

	if lerr != nil {
		fmt.Println("Open fetch progress file error,", lerr)
		return
	}
	defer ldb.Close()
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
	wg := sync.WaitGroup{}
	var seq int = 0
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
			fmt.Println("Invalid resource line,", line)
			continue
		}

		resUrl := items[0]
		resKey := ""

		if len(items) == 1 {
			resUri, pErr := url.Parse(resUrl)
			if pErr != nil {
				fmt.Println("Invalid resource url", resUrl)
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

		//otherwise fetch it
		seq += 1
		if seq%worker == 0 {
			wg.Wait()
		}

		wg.Add(1)
		go func() {
			defer wg.Done()

			fErr := client.Fetch(nil, bucket, resKey, resUrl)
			if fErr == nil {
				ldb.Put([]byte(resUrl), []byte(resKey), nil)
			} else {
				fmt.Println("Fetch", resUrl, " error due to", fErr)
			}
		}()
	}
	wg.Wait()
}
