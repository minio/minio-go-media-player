# Go Music Player App

 ![Screenshot](./assets/1.png)

 This document will guide you through the code to build a simple media player in Golang. In this app, we show you how to retrieve your media from Minio server. The media files inside a bucket are listed as the playlist, secure URLs are generated on demand whenever we play a song. Full code is available here: https://github.com/minio/minio-go-media-player, released under Apache 2.0 License.

## 1. Prerequisites
* Install mc  from [here](https://docs.minio.io/docs/minio-client-quick-start-guide).
* Install Minio Server from [here](https://docs.minio.io/docs/minio ).
* A working Golang environment. If you do not have a working Golang environment, please follow - [How to install Golang?](/docs/how-to-install-golang)

## 2. Dependencies
* Media files (mp3) for your playlist bucket.

## 3. Installing `media-player`

Let's go ahead and use 'go get' to fetch the example as shown below, 'go get' will install all necessary dependencies as needed.

```sh
$ go get -u github.com/minio/minio-go-media-player/media-player
```
Now `media-player` is ready to be used.
## 4. Running media-player

#### Environment Variables `bash`
Set Access key and Secret key environment variables.
```sh
$ export AWS_ACCESS_KEY=Q3AM3UQ867SPQQA43P2F
$ export AWS_SECRET_KEY=zuf+tfteSlswRu7BJ86wekitnifILbZam1KYY3TG
```

#### Environment Variables `tcsh`
```sh
$ setenv AWS_ACCESS_KEY Q3AM3UQ867SPQQA43P2F
$ setenv AWS_SECRET_KEY zuf+tfteSlswRu7BJ86wekitnifILbZam1KYY3TG
```

#### Environment Variables `Windows command prompt`
```sh
> set AWS_ACCESS_KEY=Q3AM3UQ867SPQQA43P2F
> set AWS_SECRET_KEY=zuf+tfteSlswRu7BJ86wekitnifILbZam1KYY3TG
```



### Set Up Bucket

1. We've created a public minio server at https://play.minio.io:9000 for developers to use as a sandbox.  Create a bucket called `media-assets` on `play.minio.io:9000`. We are going to use `mc mb` command to accomplish this.
 ```sh
$ mc mb play/media-assets
```
2. Upload your media assets into this bucket. You can again use mc to do this. Let's move few mp3's from the local disk into the bucket we created on play.
```sh
$ mc cp ~/Music/*.mp3 play/media-assets
```
**Note** : We have already created a `media-assets` bucket on play.minio.io and copied the assets used in this example.

### Run Media Player
Now we are all set to run the `media-player` example. Use `-b` command line option to specify the bucket we have already created in the previous step.
```sh
$ cd $GOPATH/bin
$ ./media-player -b media-assets
2016/04/02 17:24:54 Starting media player, please visit your browser at http://localhost:8080
```
Now if you visit http://localhost:8080  you should be able to see the example application.

### Optional Arguments

1. Endpoint defaults to 'https://play.minio.io:9000', to set a custom endpoint use `-e
```sh
$ cd $GOPATH/bin
$ ./media-player -b <bucket-name> -e https://s3.amazonaws.com
2016/04/02 17:24:54 Starting media player, please visit your browser at http://localhost:8080
```
2.  For using an endpoint as Minio server running locally.
```sh
 $ media-player -b <bucket-name> -e http://localhost:9000
```


## 5. Building Playlist
The first thing the player does is build a playlist, by using [ListObjects](https://docs.minio.io/v1.0/docs/golang-api-reference#ListObjects) method, to lists all the media assets in the media bucket specified. These objects will be rendered as a playlist for the media player as shown in the player image above. Each object is collected and sent to the browser in JSON format.

The following flow diagram and sample code provides an overview on how this is achieved.

 ![Screenshot](./assets/2.png)

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

When an user clicks to play media on the browser, secure URLs are generated on demand by the media player. In order to do this the player uses [PresignedGetObject](https://docs.minio.io/docs/golang-client-api-reference#PresignedGetObject).

The secure URL generated by PresignedGetObject is used by the player to stream the media from the server directly. The following flow diagram and sample code provide an overview on how this is achieved.

 ![Screenshot](./assets/3.png)

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
## 7. Explore Further

- [Using `minio-go` client SDK with Minio Server](/docs/golang-client-quickstart-guide)
- [Minio Golang Client SDK API Reference](/docs/golang-client-api-reference)
