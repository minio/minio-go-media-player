## Minio Go media player.
HTML5 media player using [Minio-Go library](https://github.com/minio/minio-go).

 - [Prerequisites](#prerequisites)
 - [Downloading the sample code](#downloading-the-sample-code)
 - [Run the Sample Code](#run-the-sample-code)
 - [Additional links](#additional-links)

## Prerequisites

 - Amazon s3 account or a running instance of Minio Server.
   [Click here for setting up Minio server](https://github.com/minio/minio#install-).
 - Keep your media files in the S3 or Minio bucket.

## Downloading and installing code.

<blockquote>
If you do not have a working Golang environment, please follow [Install Golang](./INSTALLGO.md).
</blockquote>

```sh
$ go get github.com/minio/minio-go-media-player/media-player
```

## Running `media-player`.
1. Set Access key and Secret key environment variables.

- On `bash`
```
export AWS_ACCESS_KEY='your-access-key'
export AWS_SECRET_KEY='your-secret-key'
```

- On `tcsh`
```
setenv AWS_ACCESS_KEY 'your-access-key'
setenv AWS_SECRET_KEY 'your-secret-key'
```

- On windows command prompt.

```
set AWS_ACCESS_KEY=your-access-key
set AWS_SECRET_KEY=your-secret-key
```

2.  Execute the following commands

```sh
$ media-player -b <bucket-name>
2016/04/02 17:24:54 Starting media player, please visit your browser at http://localhost:8080
```

- `-b` sets the bucket name to use and its mandatory.

- `-e` sets the endpoint, defaults to s3.amazonaws.com.

Set a custom endpoint for example 'play.minio.io:9000' to use Minio public server.

```sh
$ media-player -b testbucket -e play.minio.io:9000
2016/04/02 17:24:54 Starting media player, please visit your browser at http://localhost:8080
```

- `-i` By default we always make secure SSL connections, enable this option if your endpoint is on an insecure connection.

## Additional Links
- [Minio Go Library for Amazon S3 compatible cloud storage](www.github.com/minio/minio-go)
- [Minio Go API Reference](https://github.com/minio/minio-go/blob/master/API.md)
- [More API examples](https://github.com/minio/minio-go#example)
