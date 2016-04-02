package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/minio/minio-go"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

var (
	addr = flag.String("http", ":8080", "http listen address")

	bucket          = flag.String("b", "", "bucket name for operations on the object storage")
	endpoint        = flag.String("e", "s3.amazonaws.com", "object storage endpoint, defaults to amazon s3")
	enable_insecure = flag.Bool("i", false, "false for https connection, true for http only")
)

var (
	access_key = ""
	secret_key = ""
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
	presignedURL, err := storageClient.PresignedGetObject(bucket, objectName, time.Duration(ttl)*time.Second, reqParams)
	if err != nil {
		return "", err
	}
	return presignedURL, nil
}

// returns a new client for s3/Minio operations.
func newStorageClient(endpoint, accessKey, secretKey string, enableInsecure bool) (*minio.Client, error) {
	storageClient, err := minio.New(endpoint, accessKey, secretKey, enableInsecure)
	if err != nil {
		return nil, err
	}
	return storageClient, nil
}

// asserts for empty string and logs a warning.
func assertEmpty(value string, msg string) {
	if value == "" {
		log.Print(msg)
	}
}

func main() {
	flag.Parse()

	access_key = os.Getenv("AWS_ACCESS_KEY")
	secret_key = os.Getenv("AWS_SECRET_KEY")

	assertEmpty(access_key, "Env variable 'AWS_ACCESS_KEY' not set")
	assertEmpty(secret_key, "Env variable 'AWS_SECRET_KEY' not set")

	// Handler to serve the index page containing the player.
	http.Handle("/", http.FileServer(http.Dir("./web")))

	// End point for list object operation.
	// Called when player in the front end is initialized.
	http.HandleFunc("/list", ListObjects)

	// Given point which recieves the object name and returns presigned URL in the response.
	http.HandleFunc("/getpresign", GetPresignedURL)
	http.ListenAndServe(*addr, nil)
}

// Handler for ListsObjects from the Object Storage server and bucket configured above.
func ListObjects(w http.ResponseWriter, r *http.Request) {
	storageClient, err := newStorageClient(*endpoint, access_key, secret_key, *enable_insecure)
	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), 500)
		return
	}
	objects, err := listObjectsFromBucket(storageClient, *bucket)
	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), 500)
		return
	}
	// generating presigned url for the first object in the list.
	// presigned URL will be generated on the fly for the other objects when they are played.
	if len(objects) > 0 {
		presignedURL, err := getPreSignedUrl(storageClient, *bucket, objects[0].Key, 1000, nil)
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
	storageClient, err := newStorageClient(*endpoint, access_key, secret_key, *enable_insecure)
	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), 500)
		return
	}
	presignedURL, err := getPreSignedUrl(storageClient, *bucket, objectName, 1000, nil)
	// Gernerate presigned get object url.
	if err != nil {
		log.Print(err.Error())
		http.Error(w, err.Error(), 500)
		return
	}
	log.Println(presignedURL)
	w.Write([]byte(presignedURL))
}
