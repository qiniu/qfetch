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
	var zone string
	var logFile string
	var checkExists bool

	flag.Usage = func() {
		fmt.Println(`Usage of qfetch:
  -ak="": qiniu access key
  -sk="": qiniu secret key
  -bucket="": qiniu bucket
  -job="": job name to record the progress
  -file="": resource list file to fetch
  -worker=0: max goroutine in a worker group
  -check-exists: check whether file exists in bucket
  -log="": fetch failed log file
  -zone="nb": qiniu zone, nb or bc or aws`)
	}

	flag.StringVar(&job, "job", "", "job name to record the progress")
	flag.IntVar(&worker, "worker", 0, "max goroutine in a worker group")
	flag.StringVar(&file, "file", "", "resource file to fetch")
	flag.StringVar(&bucket, "bucket", "", "qiniu bucket")
	flag.StringVar(&accessKey, "ak", "", "qiniu access key")
	flag.StringVar(&secretKey, "sk", "", "qiniu secret key")
	flag.StringVar(&zone, "zone", "nb", "qiniu zone, nb or bc or aws")
	flag.StringVar(&logFile, "log", "", "fetch failed log file")
	flag.BoolVar(&checkExists, "check-exists", false, "check whether file exists in bucket")

	flag.Parse()

	if accessKey == "" {
		fmt.Println("Error: accessKey is not set")
		return
	}

	if secretKey == "" {
		fmt.Println("Error: secretKey is not set")
		return
	}

	if bucket == "" {
		fmt.Println("Error: bucket is not set")
		return
	}

	if job == "" {
		fmt.Println("Error: job name is not set")
		return
	}

	if file == "" {
		fmt.Println("Error: resource file to fetch not set")
		return
	}
	_, ferr := os.Stat(file)
	if ferr != nil {
		fmt.Println(fmt.Sprintf("Error: file '%s' not exist", file))
		return
	}

	if worker <= 0 {
		fmt.Println("Error: worker count must larger than zero")
		return
	}

	if !(zone == "nb" || zone == "bc" || zone == "aws") {
		fmt.Println("Error: zone is incorrect")
		return
	}

	qfetch.Fetch(job, checkExists, file, bucket, accessKey, secretKey, worker, zone, logFile)
}
