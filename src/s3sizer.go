/*
Author: Samarth Kanungo
Project: s3sizer
Description: This project aims to calculate size of directories of a S3 Bucket
Language: GO (1.16)
*/

package main

import (
	"flag"
	"fmt"
	client "sams3/pkg/S3/client"
	s3ops "sams3/pkg/S3/s3operations"
)

var (
	bkt    *string
	flevel *int
)

func init() {
	bkt = flag.String("b", "test", "current environment")
	flevel = flag.Int("l", 1, "port number")
}

func formatResult(bucketsize map[string]interface{}, dirsize map[string]interface{}, filesize map[string]interface{}) {
	for key, value := range bucketsize {
		fmt.Println(key, ": ", value)
	}
	fmt.Println("\n--------Directory-Sizes------------")
	for key, value := range dirsize {
		fmt.Println(key, ": ", value)
	}
	fmt.Println("\n-----------file-Sizes--------------")
	for key, value := range filesize {
		fmt.Println(key, ": ", value)
	}
}

func main() {

	flag.Parse()
	fmt.Println("\n------------------------------------")
	fmt.Println("Bucket: ", *bkt)
	fmt.Println("Folder Level: ", *flevel)

	// Bucket and folderlevel are user inputs
	bucket := *bkt
	folderLevel := *flevel

	// Create S3 client to perform operations on S3 buckets.
	s3client := client.S3client()
	// Get details of all objects in a S3 bucket
	details := s3ops.FetchS3details(bucket, s3client)
	// Get Directories at corrosponding levels.
	dirs, files := s3ops.Directories(details, folderLevel)
	// Get details of Directory and its corrosponding size.
	dirsize := s3ops.Dirsize(details, folderLevel, dirs)
	// Get total bucket size
	bucketSize := s3ops.Totalbucketsize(details)
	// Get file size
	filesize := s3ops.Filesize(details, folderLevel, files)
	// Print formatted result
	formatResult(bucketSize, dirsize, filesize)
}
