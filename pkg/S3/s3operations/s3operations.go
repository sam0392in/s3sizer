/*
Author: 	 Samarth Kanungo
Description: This package aims to provide extended capabilities to get directory size in S3.
*/

package s3operations

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	byteconverter "code.cloudfoundry.org/bytefmt"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

/*
This function will read all the objects from S3 buckets with their corrosponding sizes and return a object list
*/
func FetchS3details(bucket string, client *s3.Client) []*s3.ListObjectsV2Output {

	var data []*s3.ListObjectsV2Output
	params := &s3.ListObjectsV2Input{
		Bucket: &bucket,
	}
	totalObjects := 0
	paginator := s3.NewListObjectsV2Paginator(client, params)
	for paginator.HasMorePages() {
		output, err := paginator.NextPage(context.TODO())
		if err != nil {
			fmt.Println(err)
		}
		totalObjects += len(output.Contents)
		data = append(data, output)
	}
	fmt.Println("total objects:", totalObjects)
	return data
}

/*
This function accepts list of absolute path of objects and return relative paths of objects depending on directory level given by user
Maximum 3 levels are allowed
*/
func SelectDirLevel(paths []string, dirLevel int) (string, int) {
	var lvl string
	var totalPaths int
	if len(paths) > 2 {
		switch dirLevel {
		case 1:
			lvl = paths[0]
			totalPaths = 1
		case 2:
			lvl = paths[0] + "/" + paths[1]
			totalPaths = 2
		case 3:
			lvl = paths[0] + "/" + paths[1] + "/" + paths[2]
			totalPaths = 3
		}
	} else if len(paths) > 1 {
		switch dirLevel {
		case 1:
			lvl = paths[0]
			totalPaths = 1
		case 2:
			lvl = paths[0] + "/" + paths[1]
			totalPaths = 2
		}
	} else {
		switch dirLevel {
		case 1:
			lvl = paths[0]
			totalPaths = 1
		}
	}
	return lvl, totalPaths
}

/*
This function accept list of all objects and return directories at certain levels
*/
func Directories(details []*s3.ListObjectsV2Output, dirLevel int) ([]string, []string) {
	var directories []string
	var files []string
	var dircache string
	paths := make([]string, 4)
	for _, data := range details {
		for _, obj := range data.Contents {
			paths = strings.Split(string(*obj.Key), "/")
			lvl, totalPaths := SelectDirLevel(paths, dirLevel)
			if dircache == "" || dircache != lvl {
				if len(paths) > totalPaths {
					directories = append(directories, lvl)
				} else {
					files = append(files, lvl)
				}
			}
			dircache = lvl
		}
	}
	return directories, files
}

/*
This function returns total size of bucket (+-20KB acceptable error)
*/
func Totalbucketsize(details []*s3.ListObjectsV2Output) map[string]interface{} {
	var bucketSize int64
	totalbucketsize := map[string]interface{}{}
	for _, data := range details {
		for _, obj := range data.Contents {
			bucketSize += obj.Size
		}
	}
	totalbucketsize["Total Bucket Size"] = byteconverter.ByteSize(uint64(bucketSize))
	return totalbucketsize
}

/*
This function returns size of directories and its size
*/
func Dirsize(details []*s3.ListObjectsV2Output, dirLevel int, folders []string) map[string]interface{} {
	result := map[string]interface{}{}
	var totalSize int64
	var dir string
	for _, dir = range folders {
		totalSize = 0
		for _, data := range details {
			for _, obj := range data.Contents {
				paths := strings.Split(string(*obj.Key), "/")
				lvl, _ := SelectDirLevel(paths, dirLevel)
				re := regexp.MustCompile(dir)
				matchPath := re.MatchString(lvl)

				if matchPath {
					totalSize = totalSize + obj.Size
				}
			}
			result[dir] = byteconverter.ByteSize(uint64(totalSize))
		}
	}
	return result
}

/*
This functions will return files with corrosponding sizes in a particular level
*/
func Filesize(details []*s3.ListObjectsV2Output, dirLevel int, files []string) map[string]interface{} {
	result := map[string]interface{}{}
	var totalSize int64
	var f string
	for _, f = range files {
		totalSize = 0
		for _, d := range details {
			for _, obj := range d.Contents {
				paths := strings.Split(string(*obj.Key), "/")
				lvl, _ := SelectDirLevel(paths, dirLevel)
				re := regexp.MustCompile(f)
				matchPath := re.MatchString(lvl)
				if matchPath {
					totalSize = totalSize + obj.Size
				}
			}
			result[f] = byteconverter.ByteSize(uint64(totalSize))
		}
	}
	return result
}
