## Minio Go Mp3 player. 
 HTML5 based Media using [Minio-Go library](https://github.com/minio/minio-go).  
 
 - [Prerequisites](#prerequisites)
 - [Downloading the sample code](#downloading-the-sample-code)
 - [Run the Sample Code](#run-the-sample-code)
 - [Additional links](#additional-links)
 
## Prerequisites
 - Install Minio-go library
 - Keep your media files in the S3 or [Minio](www.minio.io) bucket. 

Note : If you do not have a working Golang environment, please follow [Install Golang](./INSTALLGO.md).

```sh
$ go get github.com/minio/minio-go
```
## Downloading the sample code
- To view or to download the code, go to:
  [https://github.com/hackintoshrao/minio-go-media-player](https://github.com/hackintoshrao/minio-go-media-player)

## Run the sample code
1. Set Access key and Secret key

    On Linux, OS X or Unix:

    ~~~~
    export AWS_ACCESS_KEY='your-access-key'
    export AWS_SECRET_KEY='your-secret-key'
    ~~~~

    On Windows:

    ~~~~
    set AWS_ACCESS_KEY=your-access-key
    set AWS_SECRET_ACCESS=your-secret-key
    ~~~~


2.  Execute the following commands

    ~~~~
    go run player-minio.go -b <bucket-name>
    ~~~~
    
`-b`  sets the bucket name , and its mandatory. 

`-e` sets the endpoint,  defaults to  s3.amazonaws.com. 
     set the endpoint to localhost:9000 testing with [Minio](www.github.com/minio/minio) locally.
     set it to the ip address of the host if Minio server is run remotely.
   
`-i` sets the enable_insecure flag. set to `false` by default. Set it to `true` only for insecure connection. 
