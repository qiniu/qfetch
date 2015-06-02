package main

import (
	"flag"
	"fmt"
	"os"
	"qfetch"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	var job string
	var worker int
	var file string
	var bucket string
	var accessKey string
	var secretKey string

	flag.StringVar(&job, "job", "fetch-test", "job name to record the progress")
	flag.IntVar(&worker, "worker", 0, "max goroutine in a worker group")
	flag.StringVar(&file, "file", "/home/jemy/Documents/resource.list", "resource file to fetch")
	flag.StringVar(&bucket, "bucket", "demo", "qiniu bucket")
	flag.StringVar(&accessKey, "ak", "HCALkwxJcWd_8UlXCb6QWdA-pEZj1FXXSK0G1lMw", "qiniu access key")
	flag.StringVar(&secretKey, "sk", "B0dP7eMztCMnmDiZfdrKXt69_q54fogZs2b1qAMx", "qiniu secret key")

	flag.Parse()

	if job == "" {
		fmt.Println("Invalid job name")
		return
	}
	if worker <= 0 {
		fmt.Println("Invalid worker")
		return
	}
	if file == "" {
		fmt.Println("Invalid resource file name")
		return
	}
	_, ferr := os.Stat(file)
	if ferr != nil {
		fmt.Println(fmt.Sprintf("File `%s' not exist", file))
		return
	}

	if bucket == "" {
		fmt.Println("Bucket is not set")
		return
	}

	if accessKey == "" {
		fmt.Println("AccessKey is not set")
		return
	}

	if secretKey == "" {
		fmt.Println("SecretKey is not set")
		return
	}
	qfetch.Fetch(job, file, bucket, accessKey, secretKey, worker)
}
