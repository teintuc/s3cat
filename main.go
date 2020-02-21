package main

import (
	"bytes"
	"flag"
	"io"
	"log"
	"net/url"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

var s3key string

func init() {
	flag.StringVar(&s3key, "key", "", "s3 key") // s3://bucket/file
	flag.Parse()
}

func main() {
	bucket, err := url.Parse(s3key)
	if err != nil {
		log.Fatalf("Can't parse s3key, %v", err)
	}

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	downloader := s3manager.NewDownloader(sess)

	buff := &aws.WriteAtBuffer{}
	numBytes, err := downloader.Download(buff,
		&s3.GetObjectInput{
			Bucket: aws.String(bucket.Host),
			Key:    aws.String(bucket.Path),
		})
	if err != nil {
		log.Fatalf("Unable to download item %v", err)
	}
	if numBytes <= 0 {
		log.Fatalf("Received 0 bytes")
	}

	// Print directly the s3 file
	bytes, err := io.Copy(os.Stdout, bytes.NewReader(buff.Bytes()))
	if err != nil {
		log.Fatalf("Unable to copy item %v", err)
	}
	if bytes <= 0 {
		log.Fatalf("Copied 0 bytes")
	}

}
