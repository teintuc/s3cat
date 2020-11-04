package main

import (
	"bytes"
	"flag"
	"fmt"
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
var out string

func init() {
	flag.StringVar(&s3key, "key", "", "s3 key") // s3://bucket/file
	flag.StringVar(&out, "out", "", "(optional) output file. Stdout by default")
	flag.Parse()
}

func writeToFile(buff *aws.WriteAtBuffer, outPath string) {
	outrsc, err := os.Create(out)
	if err != nil {
		log.Fatalf("Dest file: %v", err)
	}
	defer outrsc.Close()

	n2, err := outrsc.Write(buff.Bytes())
	if err != nil {
		log.Fatalf("%v", err)
	}
	fmt.Printf("wrote %d bytes\n", n2)
}

func writeToStdout(buff *aws.WriteAtBuffer) {
	// Print directly the s3 file
	bytes, err := io.Copy(os.Stdout, bytes.NewReader(buff.Bytes()))
	if err != nil {
		log.Fatalf("Unable to copy item %v", err)
	}
	if bytes <= 0 {
		log.Fatalf("Copied 0 bytes")
	}
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

	// If we have an output file we write to it.
	// Otherwise we write to stdout.
	if len(out) > 0 {
		writeToFile(buff, out)
	} else {
		writeToStdout(buff)
	}
}
