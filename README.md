#qfetch

###简介
qfetch是一个数据迁移工具，利用七牛提供的[fetch](http://developer.qiniu.com/docs/v6/api/reference/rs/fetch.html)功能来抓取指定文件列表中的文件。在文件列表中，你只需要提供资源的外链地址和要保存在七牛空间中的文件名就可以了。

###下载

**建议下载最新版本**

|版本     |支持平台|链接|
|--------|---------|----|
|qfetch v1.1|Linux, Windows, Mac OSX|[下载](http://7rfgu2.com1.z0.glb.clouddn.com/qfetch-v1.1.zip)|

###使用
该工具是一个命令行工具，需要指定相关的参数来运行。

```
Usage of qfetch:
  -ak="HCALkwxJcWd_8UlXCb6QWdA-pEZj1FXXSK0G1lMw": qiniu access key
  -sk="B0dP7eMztCMnmDiZfdrKXt69_q54fogZs2b1qAMx": qiniu secret key
  -bucket="demo": qiniu bucket
  -file="/home/jemy/Documents/resource.list": resource file to fetch
  -job="fetch-test": job name to record the progress
  -worker=0: max goroutine in a worker group
```

上面所有的参数都是必须指定的。

|命令|描述|
|--------|---------|
|ak|七牛账号的AccessKey，可以从七牛的后台获取|
|sk|七牛账号的SecretKey，可以从七牛的后台获取|
|bucket|文件抓取后存储的空间，为空间的名字|
|file|抓取任务列表文件的本地路径，内容为多个由需要抓取的资源外链和对应的保存在七牛空间中的文件名组成的行|
|job|任务的名称，指定这个参数主要用来将抓取成功的文件放在本地数据库中，便于后面核对|
|worker|抓取的并发数量，可以适当地指定较大的并发请求数量来提高批量抓取的效率，可根据实际带宽和文件平均大小来计算得出|


**模式一:**

上面的`file`参数指定的索引文件的行格式如下：

```
文件链接1\t保存名称1
文件链接2\t保存名称2
文件链接3\t保存名称3
...
```

其中`\t`表示Tab分隔符号。

例如：

```
http://img.abc.com/0/000/484/0000484193.fid	2009-10-14/2922168_b.jpg
http://img.abc.com/0/000/553/0000553777.fid	2009-07-01/2270194_b.jpg
http://img.abc.com/0/000/563/0000563511.fid	2009-03-01/1650739_s.jpg
http://img.abc.com/0/000/563/0000563514.fid	2009-05-01/1953696_m.jpg
http://img.abc.com/0/000/563/0000563515.fid	2009-02-01/1516376_s.jpg
```

上面的方式最终抓取保存在空间中的文件名字是：

```
2009-10-14/2922168_b.jpg
2009-07-01/2270194_b.jpg
2009-03-01/1650739_s.jpg
2009-05-01/1953696_m.jpg
2009-02-01/1516376_s.jpg
```

**模式二:**

上面的`file`参数指定的索引文件的行格式如下：

```
文件链接1
文件链接2
文件链接3
...
```

上面的方式也是支持的，这种方式的情况下，文件保存的名字将从指定的文件链接里面自动解析。

例如：

```
http://img.abc.com/0/000/484/0000484193.fid
http://img.abc.com/0/000/553/0000553777.fid
http://img.abc.com/0/000/563/0000563511.fid
http://img.abc.com/0/000/563/0000563514.fid
http://img.abc.com/0/000/563/0000563515.fid
```

其抓取后保存在空间中的文件名字是：

```
0/000/484/0000484193.fid
0/000/553/0000553777.fid
0/000/563/0000563511.fid
0/000/563/0000563514.fid
0/000/563/0000563515.fid
```


###日志
抓取成功的文件在本地都会写入以`job`参数指定的值为名称的本地leveldb数据库中。该leveldb名称以`.`开头，所以在Linux或者Mac系统下面是个隐藏文件。在整个文件索引都抓取完成后，可以使用[leveldb](https://github.com/jemygraw/leveldb)工具来导出所有的成功的文件列表，和原来的列表比较，就可以得出失败的抓取列表。上面的方法也可以被用来验证抓取的完整性。

###示例
抓取指令为：

```
qfetch -ak='x98pdzDw8dtwM-XnjCwlatqwjAeed3lwyjcNYqjv' -sk='OCCTbp-zhD8x_spN0tFx4WnMABHxggvveg9l9m07' 
-bucket='image' -file='diff.txt' -worker=300 -job='diff'  | tee diff.log
```

上面的指令抓取文件索引`diff.txt`里面的文件，存储到空间`piccenter`里面，并发请求数量`300`，任务的名称叫做`diff`，成功列表日志文件名称是`.diff.job`。另外由于该命令打印的报警日志输出到终端，所以可以使用`tee`命令将内容复制一份到日志文件中。

导出成功列表：

```
leveldb -export='.diff.job' >> list.txt
```

注意，任务的leveldb的名字时`.diff.job`。
