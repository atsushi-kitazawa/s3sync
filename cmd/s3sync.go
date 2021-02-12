package main

import (
    "bytes"
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
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

    if len(os.Args) < 2 {
	fmt.Println("expected 'upload' or 'download' subcommands.")
	os.Exit(1)
    }

    switch os.Args[1] {
    case "upload":
	upload.Parse(os.Args[2:])
	fmt.Println("upload subcommand.")
	fmt.Println("   force:", *forceUpload)
	fmt.Println("   mode:", *modeUpload)
	s3sync_upload()
    case "download":
	download.Parse(os.Args[2:])
	fmt.Println("download subcommand")
	fmt.Println("   force:", *forceDownload)
	fmt.Println("   mode:", *modeDownload)
	s3sync_download()
    default:
	fmt.Println("expected 'upload' or 'download' subcommands.")
	os.Exit(1)
    }

    ls()
}

func s3sync_upload() {
    client := client()
    uploader := manager.NewUploader(client, func(u *manager.Uploader) {
	u.BufferProvider = manager.NewBufferedReadSeekerWriteToPool(25 * 1024 * 1024)
    })

    _, err := uploader.Upload(context.TODO(), &s3.PutObjectInput {
	Bucket: aws.String("s3sync-cli"),
	Key: aws.String("key1"),
	Body: bytes.NewReader([]byte("hoge uga ahe")),
    })
    if err != nil {
	panic(err)
    }
}

func s3sync_download() {
    f, err := os.Create("download.txt")
    if err != nil {
	panic(err)
    }
    defer f.Close()

    client := client()

    downloader := manager.NewDownloader(client, func(d *manager.Downloader) {
	d.PartSize = 64 * 1024 * 1024
    })
    n ,err := downloader.Download(context.TODO(), f, &s3.GetObjectInput {
	Bucket: aws.String("s3sync-cli"),
	Key: aws.String("key1"),
    })
    if err != nil {
	panic(err)
    }
    fmt.Println("download %d bytes\n", n)
}

func ls() {
    client := client()

    output, err := client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input {
	Bucket: aws.String("s3sync-cli"),
    })
    if err != nil {
	panic(err)
    }

    for _, object := range output.Contents {
	fmt.Printf("key=%s, size=%d \n", aws.ToString(object.Key), object.Size)
    }
}

func client() *s3.Client {
    cfg, err := config.LoadDefaultConfig(context.TODO())
    if err != nil {
	panic(err)
    }

    client := s3.NewFromConfig(cfg)
    return client
}

func loadConf() {
    // load
    // bucket name
    // local directory
    // s3 bucket name
}
