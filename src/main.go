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

	flag.Usage = func() {
		fmt.Println(`Usage of qfetch:
  -ak="": qiniu access key
  -sk="": qiniu secret key
  -bucket="": qiniu bucket
  -job="": job name to record the progress
  -file="": resource list file to fetch
  -worker=0: max goroutine in a worker group
  -zone="z0": qiniu zone, z0 or z1`)
	}

	flag.StringVar(&job, "job", "", "job name to record the progress")
	flag.IntVar(&worker, "worker", 0, "max goroutine in a worker group")
	flag.StringVar(&file, "file", "", "resource file to fetch")
	flag.StringVar(&bucket, "bucket", "", "qiniu bucket")
	flag.StringVar(&accessKey, "ak", "", "qiniu access key")
	flag.StringVar(&secretKey, "sk", "", "qiniu secret key")
	flag.StringVar(&zone, "zone", "z0", "qiniu zone, z0 or z1")

	flag.Parse()

	if accessKey == "" {
		fmt.Println("AccessKey is not set")
		return
	}

	if secretKey == "" {
		fmt.Println("SecretKey is not set")
		return
	}

	if bucket == "" {
		fmt.Println("Bucket is not set")
		return
	}

	if job == "" {
		fmt.Println("Invalid job name")
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

	if worker <= 0 {
		fmt.Println("Invalid worker")
		return
	}

	if zone == "" {
		fmt.Println("Zone is not set")
		return
	}

	qfetch.Fetch(job, file, bucket, accessKey, secretKey, worker, zone)
}
