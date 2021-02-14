package main

import (
	//"bytes"
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"s3sync/cmd/util"
	"s3sync/configs"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func main() {
    _main()
}

func _main() {
    upload := flag.NewFlagSet("upload", flag.ExitOnError)
    forceUpload := upload.Bool("f", false, "force upload")
    modeUpload := upload.String("m", "", "mode")

    download := flag.NewFlagSet("download", flag.ExitOnError)
    forceDownload := download.Bool("f", false, "force download")
    modeDownload := download.String("m", "", "mode")

    ls := flag.NewFlagSet("ls", flag.ExitOnError)

    if len(os.Args) < 2 {
	fmt.Println("expected 'upload' or 'download' subcommands.")
	os.Exit(1)
    }

    // load setting.yaml
    s := configs.Load("./configs/setting.yaml")

    //ls(s)
    //os.Exit(0)

    switch os.Args[1] {
    case "upload":
	upload.Parse(os.Args[2:])
	fmt.Println("upload subcommand.")
	fmt.Println("   force:", *forceUpload)
	fmt.Println("   mode:", *modeUpload)
	s3sync_upload(s)
    case "download":
	download.Parse(os.Args[2:])
	fmt.Println("download subcommand")
	fmt.Println("   force:", *forceDownload)
	fmt.Println("   mode:", *modeDownload)
	s3sync_download(s)
    case "ls":
	ls.Parse(os.Args[2:])
	fmt.Println("ls subcommand")
	s3sync_ls(s)
    default:
	fmt.Println("expected 'upload' or 'download', 'ls' subcommands.")
	os.Exit(1)
    }
}

// maybe use BatchUploadObject
func s3sync_upload(s *configs.Setting) {
    client := client()
    uploader := manager.NewUploader(client, func(u *manager.Uploader) {
	u.BufferProvider = manager.NewBufferedReadSeekerWriteToPool(25 * 1024 * 1024)
    })

    list := util.ListDir(s.LocalDir)
    for _, file := range list {
	body, err := os.Open(s.LocalDir + file)
	if err != nil {
	    panic(err)
	}
	//fmt.Println(file)

	_, err = uploader.Upload(context.TODO(), &s3.PutObjectInput {
	    Bucket: aws.String(s.BucketName),
	    Key: aws.String(s.S3Dir + file),
	    Body: body,
	})
	if err != nil {
	    panic(err)
	}
    }
}

// maybe use BatchDownloadObject
func s3sync_download(s *configs.Setting) {
    client := client()

    downloader := manager.NewDownloader(client, func(d *manager.Downloader) {
	d.PartSize = 64 * 1024 * 1024
    })

    // get download list from ls
    downloadList := s3sync_ls(s)

    // remove directory
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
}

func s3sync_ls(s *configs.Setting) []string {
    client := client()

    output, err := client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input {
	Bucket: aws.String(s.BucketName),
	Prefix: aws.String(s.S3Dir),
    })
    if err != nil {
	panic(err)
    }

    ret := make([]string, 0)
    for _, object := range output.Contents {
	fmt.Printf("key=%s, size=%d \n", aws.ToString(object.Key), object.Size)
	ret = append(ret, aws.ToString(object.Key))
    }

    return ret
}

func client() *s3.Client {
    cfg, err := config.LoadDefaultConfig(context.TODO())
    if err != nil {
	panic(err)
    }

    client := s3.NewFromConfig(cfg)
    return client
}
