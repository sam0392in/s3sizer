# S3 Sizer
![logo](./pkg/S3/sams3.png "icon")

- This project aims to calculate size of directories in S3 bucket.

## Usage

- -b: bucket name
- -l: folder level
```
 ./s3sizer -b sam-codetransfer-cloudops -l 1
```
 ### Output
```
$ go run s3sizer.go -b sam-codetransfer-cloudops -l 2
------------------------------------
Bucket:  sam-codetransfer-cloudops
Folder Level:  2
total objects: 7
Total Bucket Size :  33.2K

--------Directory-Sizes------------
sams3/pkg :  19.7K
sams3/src :  1.6K

--------file-Sizes------------
sams3/README.md :  996B
sams3/go.mod :  230B
sams3/go.sum :  10.8K

```


### Get the details of particular directory

```
$ go run s3sizer.go -b sam-codetransfer-cloudops -l 2 | grep go.
sams3/go.mod :  230B
sams3/go.sum :  10.8K

```
