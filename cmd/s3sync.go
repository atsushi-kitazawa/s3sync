package main

import (
	//"bytes"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"s3sync/cmd/util"
	"s3sync/configs"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	_ "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

var (
    upload flag.FlagSet
    uploadForce bool
    uploadMode string

    download flag.FlagSet
    downloadForce bool
    downloadMode string

    ls flag.FlagSet
)

func init() {
    upload := flag.NewFlagSet("upload", flag.ExitOnError)
    upload.BoolVar(&uploadForce, "f", false, "upload force")
    upload.StringVar(&uploadMode, "m", "", "upload mode")

    download := flag.NewFlagSet("download", flag.ExitOnError)
    download.BoolVar(&downloadForce, "f", false, "download force")
    download.StringVar(&downloadMode, "m", "", "download mode")

    ls := flag.NewFlagSet("ls", flag.ExitOnError)
    _ = ls
}

func main() {
    _main()
}

func _main() {
    if len(os.Args) < 2 {
	log.Fatalf("expected 'upload' or 'download' subcommands.")
    }

    // load setting.yaml
    s := configs.Load("./configs/setting.yaml")

    switch os.Args[1] {
    case "upload":
	upload.Parse(os.Args[2:])
	fmt.Println("upload subcommand.")
	fmt.Println("   force:", uploadForce)
	fmt.Println("   mode:", uploadMode)
	s3sync_upload(s)
    case "download":
	download.Parse(os.Args[2:])
	fmt.Println("download subcommand")
	fmt.Println("   force:", &downloadForce)
	fmt.Println("   mode:", &downloadMode)
	s3sync_download(s)
    case "ls":
	ls.Parse(os.Args[2:])
	fmt.Println("ls subcommand")
	list := s3sync_ls(s)
	for _, f := range list {
	    fmt.Println(f)
	}
    default:
	fmt.Println("expected 'upload' or 'download', 'ls' subcommands.")
	os.Exit(1)
    }
}

// maybe use BatchUploadObject
func s3sync_upload(s *configs.Setting) {
    client := client(s)
    uploader := manager.NewUploader(client, func(u *manager.Uploader) {
	u.BufferProvider = manager.NewBufferedReadSeekerWriteToPool(25 * 1024 * 1024)
    })

    // upload
    uploadlist := util.ListDir(s.LocalDir)
    //fmt.Println("list>", uploadlist)
    for _, file := range uploadlist {
	body, err := os.Open(s.LocalDir + file)
	if err != nil {
	    panic(err)
	}

	_, err = uploader.Upload(context.TODO(), &s3.PutObjectInput {
	    Bucket: aws.String(s.BucketName),
	    Key: aws.String(s.S3Dir + file),
	    Body: body,
	})
	if err != nil {
	    panic(err)
	}
    }

    // delete
    s3List := s3sync_ls(s)
    //fmt.Println("s3list>", s3List)
    delList := make([]types.ObjectIdentifier, 0)
    for _, val := range s3List {
	if !contains(uploadlist, strings.TrimPrefix(val, s.S3Dir)) {
	    object := types.ObjectIdentifier{
		Key: aws.String(val),
	    }
	    delList = append(delList, object)
	}
    }
    //fmt.Println(delList)

    if len(delList) == 0 {
	return
    }

    input := &s3.DeleteObjectsInput {
        Bucket: aws.String(s.BucketName),
        Delete: &types.Delete {
            delList,
            true,
        },
    }
    _, err := client.DeleteObjects(context.TODO(), input)
    if err != nil {
        panic(err)
    }
}
// maybe use BatchDownloadObject
func s3sync_download(s *configs.Setting) {
    client := client(s)

    downloader := manager.NewDownloader(client, func(d *manager.Downloader) {
	d.PartSize = 64 * 1024 * 1024
    })

    // get download list from ls. This List contains directory.
    downloadList := s3sync_ls(s)

    // remove directory in list
    downloadListRemoveDir := make([]string, 0)
    for _, d := range downloadList {
	if !strings.HasSuffix(d, "/") {
	    downloadListRemoveDir = append(downloadListRemoveDir, d)
	}
    }

    // create file
    for _, s3File := range downloadListRemoveDir {
	localPath := s.LocalDir + strings.TrimPrefix(s3File, s.S3Dir)
	os.MkdirAll(filepath.Dir(localPath), 0777)
	localFile, err := os.Create(localPath)
	if err != nil {
	    panic(err)
	}
	defer localFile.Close()

	n ,err := downloader.Download(context.TODO(), localFile, &s3.GetObjectInput {
	    Bucket: aws.String(s.BucketName),
	    Key: aws.String(s3File),
	})
	if err != nil {
	    panic(err)
	}
	fmt.Printf("download %d bytes\n", n)
    }

    localList := util.ListDir(s.LocalDir)
    //fmt.Println("localList>", localList)
    //fmt.Println("downloadList>", downloadList)
    //delList := make([]string, 0)
    for _, val := range localList {
	if !contains(downloadList, s.S3Dir + val) {
	    //fmt.Println("val>", val)
	    if err := os.Remove(s.LocalDir + val); err != nil {
		fmt.Printf("failed remove %s\n", val)
	    }
	}
    }
}

func s3sync_ls(s *configs.Setting) []string {
    client := client(s)

    output, err := client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input {
	Bucket: aws.String(s.BucketName),
	Prefix: aws.String(s.S3Dir),
    })
    if err != nil {
	panic(err)
    }

    ret := make([]string, 0)
    for _, object := range output.Contents {
	//fmt.Printf("key=%s, size=%d \n", aws.ToString(object.Key), object.Size)
	ret = append(ret, aws.ToString(object.Key))
    }

    return ret
}

func client(s *configs.Setting) *s3.Client {
    client := s3.New(s3.Options{
	Region: s.Credential.Region,
	Credentials: credentials.NewStaticCredentialsProvider(s.Credential.Apikey, s.Credential.Secretkey, ""),
    })
    return client
}

func contains(s []string, item string) bool {
    //fmt.Println("item>", item)
    for _, val := range s {
	if val == item {
	    return true
	}
    }
    return false
}
