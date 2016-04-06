## Minio Go media player.
HTML5 media player using [Minio-Go library](https://github.com/minio/minio-go).

 - [Prerequisites](#prerequisites)
 - [Installing media-player](#installing-media-player)
 - [Running media-player](#running-media-player)
 - [Additional links](#additional-links)

## Prerequisites

 - Amazon s3 account or a running instance of Minio Server.
   [Click here for setting up Minio server](https://github.com/minio/minio#install-).
 - Keep your media files in the S3 or Minio bucket.

## Installing media-player

If you do not have a working Golang environment, please follow [Install Golang](./INSTALLGO.md).

```sh
$ go get github.com/minio/minio-go-media-player/media-player
```

## Running media-player

### Environment variables.
Set Access key and Secret key environment variables.

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

### Create bucket and copy media assets.

Following example uses [mc(Minio Client)](https://github.com/minio/mc) to create a bucket.
```sh
$ mc mb <aliasname>/<bucket-name>
$ mc cp Music/*.mp3 <aliasname>/<bucket-name>
```

### Run media player.

Now we are all set to run the `media-player` example.

```sh
$ media-player -b <bucket-name>
2016/04/02 17:24:54 Starting media player, please visit your browser at http://localhost:8080
```

- `-b` sets the bucket name to use and its mandatory.

### Optional arguments.

- Endpoint defaults to 's3.amazonaws.com', to set a custom endpoint use `-e`.

```sh
$ media-player -b <bucket-name> -e play.minio.io:9000
2016/04/02 17:24:54 Starting media player, please visit your browser at http://localhost:8080
```

- By default we always make secure SSL connections, enable insecure with `-i` option.

```sh
$ media-player -b <bucket-name> -e localhost:9000 -i
```

## Additional Links
- [Minio Go Library for Amazon S3 compatible cloud storage](www.github.com/minio/minio-go)
- [Minio Go API Reference](https://github.com/minio/minio-go/blob/master/API.md)
- [More API examples](https://github.com/minio/minio-go#example)
