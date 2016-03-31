package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/minio/minio-go"
	"log"
	"net/http"
	"net/url"
	"time"
)

var (
	addr = flag.String("http", ":8080", "http listen address")
)

// STORAGE options.
// STORAGE = "s3.amazonaws.com" if s3 is used.
// STORAGE = "localhost:<port>" for testing on localhost Minio server.
// STORAGE = "<IP Address>:<port>i" for testing on remote Minio server.

// ACCESSKEYID = "<Your Access KEY ID of S3/Minio service>"
// SECRETACCESSKEY = "<Your Secret Key of S3/Minio service>"
// BUCKET_NAME = "<Your-Bucket-Name">
// ENABLE_INSECURE = true , for http only connection.
// ENABLE_INSECURE = false, for https connection.
const (
	// points to S3 by default.
	STORAGE         = "s3.amazonaws.com"
	ACCESSKEYID     = ""
	SECRETACCESSKEY = ""
	BUCKET_NAME     = ""
	// set to false by default.
	ENABLE_INSECURE = false
)

// The playlist for the music player on the browser.
// Used for sending the response from listObjects for the player.
type playlist struct {
	Key string
	URL string
}

// ListsObjects from fiven storageClient and bucket.
func listObjectsFromBucket(storageClient *minio.Client, bucket string) ([]playlist, error) {
	// Create a done channel to control 'ListObjects' go routine.
	doneCh := make(chan struct{})
	// Indicate to our routine to exit cleanly upon return.
	defer close(doneCh)
	var objectInfos []playlist
	// List all objects from a bucket-name with a matching prefix.
	for objectInfo := range storageClient.ListObjects(bucket, "", true, doneCh) {
		if objectInfo.Err != nil {
			fmt.Println(objectInfo.Err)
			return objectInfos, objectInfo.Err
		}
		res := playlist{Key: objectInfo.Key}
		objectInfos = append(objectInfos, res)

	}
	return objectInfos, nil
}

// return given URL query parameter from the http request.
func getUrlQueryParam(r *http.Request, param string) string {
	return r.URL.Query().Get(param)
}

// Get a presigned access URL for the given object for the specified ttl period.
func getPreSignedUrl(storageClient *minio.Client, bucket, objectName string, ttl int, reqParams url.Values) (string, error) {
	presignedURL, err := storageClient.PresignedGetObject(BUCKET_NAME, objectName, time.Duration(ttl)*time.Second, reqParams)
	if err != nil {
		return "", err
	}
	return presignedURL, nil
}

// returns a new client for s3/Minio operations.
func newStorageClient(storage, accessKey, secretKey string, enableInsecure bool) (*minio.Client, error) {
	storageClient, err := minio.New(storage, accessKey, secretKey, enableInsecure)
	if err != nil {
		return nil, err
	}
	return storageClient, nil
}

func main() {
	flag.Parse()
	// Handler to serve the index page containing the player.
	http.HandleFunc("/", Index)
	// For accessing the static content.
	http.HandleFunc("/jplayer/", func(w http.ResponseWriter, r *http.Request) {
		log.Print("serve static : " + r.URL.Path)
		http.ServeFile(w, r, r.URL.Path[1:])
	})
	// End point for list object operation.
	// Called when player in the front end is initialized.
	http.HandleFunc("/list", ListObjects)
	// Given point which recieves the object name and returns presigned URL in the response.
	http.HandleFunc("/getpresign", GetPresignedURL)
	http.ListenAndServe(*addr, nil)
}

// Handler serving the index page.
func Index(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./index.html")
	log.Print("index page served")
}

// Handler for ListsObjects from the Object Storage server and bucket configured above.
func ListObjects(w http.ResponseWriter, r *http.Request) {
	storageClient, err := newStorageClient(STORAGE, ACCESSKEYID, SECRETACCESSKEY, ENABLE_INSECURE)
	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), 500)
		return
	}
	objects, err := listObjectsFromBucket(storageClient, BUCKET_NAME)
	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), 500)
		return
	}
	// generating presigned url for the first object in the list.
	// presigned URL will be generated on the fly for the other objects when they are played.
	if len(objects) > 0 {
		presignedURL, err := getPreSignedUrl(storageClient, BUCKET_NAME, objects[0].Key, 1000, nil)
		// Gernerate presigned get object url.
		if err != nil {
			log.Print(err.Error())
			http.Error(w, err.Error(), 500)
			return
		}
		objects[0].URL = presignedURL
	}
	json.NewEncoder(w).Encode(objects)
}

// Handler for obtaining Presigned access URL for given object.
func GetPresignedURL(w http.ResponseWriter, r *http.Request) {
	// The object for which the presigned URL has to be generated is sent as a query
	// parameter from the client.
	objectName := getUrlQueryParam(r, "objname")
	// Set request parameters
	storageClient, err := newStorageClient(STORAGE, ACCESSKEYID, SECRETACCESSKEY, ENABLE_INSECURE)
	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), 500)
		return
	}
	presignedURL, err := getPreSignedUrl(storageClient, BUCKET_NAME, objectName, 1000, nil)
	// Gernerate presigned get object url.
	if err != nil {
		log.Print(err.Error())
		http.Error(w, err.Error(), 500)
		return
	}
	log.Println(presignedURL)
	w.Write([]byte(presignedURL))
}
