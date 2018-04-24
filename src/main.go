package main

import (
	"flag"
	"fmt"
	"os"
	"qfetch"
	"runtime"

	"github.com/qiniu/api.v6/conf"
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

	//support qos, should be specified both
	var rsHost string
	var ioHost string

	flag.Usage = func() {
		fmt.Println(`Usage of qfetch:
  -ak="": qiniu access key
  -sk="": qiniu secret key
  -bucket="": qiniu bucket
  -job="": job name to record the progress
  -file="": resource list file to fetch
  -worker=1: max goroutine in a worker group
  -check-exists: check whether file exists in bucket
  -log="": fetch runtime log file
  -zone="": qiniu zone, support [nb, bc, hn, as0, na0]
  -rs-host="": rs host to support specified qos system
  -io-host="": io host to support specified qos sytem
 version 1.8`)
	}

	flag.StringVar(&job, "job", "", "job name to record the progress")
	flag.IntVar(&worker, "worker", 1, "max goroutine in a worker group")
	flag.StringVar(&file, "file", "", "resource file to fetch")
	flag.StringVar(&bucket, "bucket", "", "qiniu bucket")
	flag.StringVar(&accessKey, "ak", "", "qiniu access key")
	flag.StringVar(&secretKey, "sk", "", "qiniu secret key")
	flag.StringVar(&zone, "zone", "", "qiniu zone, support [nb, bc, hn, as0, na0]")
	flag.StringVar(&logFile, "log", "", "fetch runtime log file")
	flag.BoolVar(&checkExists, "check-exists", false, "check whether file exists in bucket")
	flag.StringVar(&rsHost, "rs-host", "", "rs host to support specified qos system")
	flag.StringVar(&ioHost, "io-host", "", "io host to support specified qos system")

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

	if (rsHost != "" && ioHost == "") || (rsHost == "" && ioHost != "") {
		fmt.Println("Error: rs host and io host should be specified together")
		return
	}

	if rsHost != "" && ioHost != "" && zone != "" {
		fmt.Println("Error: if you specified rs host and io host, zone should be empty")
		return
	}

	if zone != "" && !(zone == "nb" || zone == "bc" || zone == "hn" || zone == "as0" || zone == "na0") {
		fmt.Println("Error: zone is incorrect")
		return
	}

	conf.IO_HOST = ioHost
	conf.RS_HOST = rsHost

	qfetch.Fetch(job, checkExists, file, bucket, accessKey, secretKey, worker, zone, logFile)
}
