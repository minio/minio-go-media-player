# Go Music Player App [![Gitter](https://badges.gitter.im/Join%20Chat.svg)](https://gitter.im/minio/minio?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)

![minio_GO1](https://github.com/minio/minio-go-media-player/blob/master/docs/screenshots/minio-go1.jpg?raw=true)

本示例将会手把手（限女生）指导你如何用Golang构建一个简单的音乐播放器。在这个app中，我们会向你展示如何从Minio Server上获取你的音频文件。你可以通过[这里](https://github.com/minio/minio-go-media-player)获取完整的代码，代码是以Apache 2.0 License发布的。

## 1. 前提条件

* 从[这里](https://docs.minio.io/docs/minio-client-quickstart-guide)下载并安装mc。
* 从[这里](https://docs.minio.io/docs/minio )下载并安装Minio Server。
* 一个Golang的开发环境。如果没有的话，请参考[如何安装Golang](https://docs.minio.io/docs/how-to-install-golang)。

## 2. 依赖

* 你的`playlist`存储桶中需要有音频（mp3）文件。


## 3. 安装`media-player`

参考下面示例，使用“go get”下载示例代码，“go get”会安装所需要的依赖。

```sh
go get -u github.com/minio/minio-go-media-player/media-player
```

现在`media-player`已经按耐不住，想唱首歌给你听了。

## 4. 运行media-player

#### 环境变量`bash`

设置环境变量`Access key`和`Secret key`。

```sh
export AWS_ACCESS_KEY=Q3AM3UQ867SPQQA43P2F
export AWS_SECRET_KEY=zuf+tfteSlswRu7BJ86wekitnifILbZam1KYY3TG
```

#### 环境变量`tcsh`

```sh
setenv AWS_ACCESS_KEY Q3AM3UQ867SPQQA43P2F
setenv AWS_SECRET_KEY zuf+tfteSlswRu7BJ86wekitnifILbZam1KYY3TG
```

#### 环境变量`Windows command prompt`

```sh
set AWS_ACCESS_KEY=Q3AM3UQ867SPQQA43P2F
set AWS_SECRET_KEY=zuf+tfteSlswRu7BJ86wekitnifILbZam1KYY3TG
```


### 设置存储桶

1. 我们已经创建了一个公开的Minio Server(https://play.minio.io:9000) 供大家进行开发和测试。调用`mc mb`命令，在`play.minio.io:9000`上创建一个名叫`media-assets`的存储桶。 

```sh
mc mb play/media-assets
```

2. 将你心爱的音乐上传到这个存储桶中。我们可以用mc来将本地的mp3文件上传到play上的存储桶中。

```sh
mc cp ~/Music/*.mp3 play/media-assets
```
**注意** : 我们已经在play.minio.io上创建了`media-assets`这个存储桶，并将本示例用到的资源上传到这个存储桶了。

### 运行音乐播放器

现在已经万事俱备。使用`-b`命令行参数指定音乐文件所在的存储桶。

```sh
cd $GOPATH/bin
./media-player -b media-assets
2016/04/02 17:24:54 Starting media player, please visit your browser at http://localhost:8080
```
现在如果你访问http://localhost:8080 ，你应该可以看到这个示例程序。

### 可选参数

1. Endpoint默认指向的是'https://play.minio.io:9000' ，想换一个的话可以使用`-e`。

```sh
cd $GOPATH/bin
./media-player -b <bucket-name> -e https://s3.amazonaws.com
2016/04/02 17:24:54 Starting media player, please visit your browser at http://localhost:8080
```

2.  如果想指定本地运行的Minio Server。

```sh
media-player -b <bucket-name> -e http://localhost:9000
```

## 5. 创建播放列表

这个播放器第一件要做的事情就是通过调用[ListObjects](https://docs.minio.io/docs/golang-client-api-reference#ListObjects)方法，获取指定的存储桶中的音频资源，并创建一个播放列表。这些对象将会以JSON格式返回，做为音乐播放器的播放列表的数据。

下面的流程图和示例代码提供了如何实现的概述。

![minio_GO2](https://github.com/minio/minio-go-media-player/blob/master/docs/screenshots/minio-go2.jpg?raw=true)


```go
 for objectInfo := range api.minioClient.ListObjects(*bucketName, "", isRecursive, doneCh) {
  if objectInfo.Err != nil {
		http.Error(w, objectInfo.Err.Error(), http.StatusInternalServerError)
			return
		}
		objectName := objectInfo.Key // object name.
		playListEntry := mediaPlayList{
			Key: objectName,
		}
		playListEntries = append(playListEntries, playListEntry)
	}
	playListEntriesJSON, err := json.Marshal(playListEntries)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Successfully wrote play list in json.
	w.Write(playListEntriesJSON)
```

## 6. Streaming Media

当你点击浏览器上的播放器时，播放器调用[PresignedGetObject](https://docs.minio.io/docs/golang-client-api-reference#PresignedGetObject) ，生成了一个安全的URL。

播放器使用PresignedGetObject生成的安全URL从服务器直接获取流。 以下流程图和示例代码提供了如何实现的概述。

![minio_GO3](https://github.com/minio/minio-go-media-player/blob/master/docs/screenshots/minio-go3.jpg?raw=true)


```go
 // GetPresignedURLHandler - generates presigned access URL for an object.
func (api mediaHandlers) GetPresignedURLHandler(w http.ResponseWriter, r *http.Request) {
	// The object for which the presigned URL has to be generated is sent as a query
	// parameter from the client.
	objectName := r.URL.Query().Get("objName")
	if objectName == "" {
		http.Error(w, "No object name set, invalid request.", http.StatusBadRequest)
		return
	}
	presignedURL, err := api.storageClient.PresignedGetObject(*bucketName, objectName, 1000*time.Second, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write([]byte(presignedURL))
}
```

## 7. 了解更多

- [Using `minio-go` client SDK with Minio Server](https://docs.minio.io/docs/golang-client-quickstart-guide)
- [Minio Golang Client SDK API Reference](https://docs.minio.io/docs/golang-client-api-reference)
