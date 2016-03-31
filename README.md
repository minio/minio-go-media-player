## Minio Go Mp3 player. 
 HTML5 based Mp3 player with the Media being served from Amazon S3 or [Minio Server](https://github.com/minio/minio) bucket. 
 Web server written in Golang serves the player and interacts with S3/Minio using [Minio-Go-SDK](https://github.com/minio/minio-go).  
 
 Here is how to use the player. 
 1. Upload your MP3 to a single bucket in S3 or Minio Server. 
 2. Open [player-minio.go](https://github.com/hackintoshrao/minio-go-media-player/blob/master/player-minio.go).
 3. Edit the [configuration](https://github.com/hackintoshrao/minio-go-media-player/blob/master/player-minio.go#L19) as directed in the code.(Basically setting 
 up your service address,access and secret keys).
 4. Run the server. 
    `go run player-minio.go`
 5. Visit `localhost:8080` on your browser.     
