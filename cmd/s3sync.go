package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"os"
	"s3sync/configs"

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
    default:
	fmt.Println("expected 'upload' or 'download' subcommands.")
	os.Exit(1)
    }
}

func s3sync_upload(s *configs.Setting) {
    client := client()
    uploader := manager.NewUploader(client, func(u *manager.Uploader) {
	u.BufferProvider = manager.NewBufferedReadSeekerWriteToPool(25 * 1024 * 1024)
    })

    _, err := uploader.Upload(context.TODO(), &s3.PutObjectInput {
	Bucket: aws.String(s.BucketName),
	Key: aws.String("test/bbb.txt"),
	Body: bytes.NewReader([]byte("hoge uga ahe")),
    })
    if err != nil {
	panic(err)
    }
}

func s3sync_download(s *configs.Setting) {
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
	Bucket: aws.String(s.BucketName),
	Key: aws.String("test/aaa.txt"),
    })
    if err != nil {
	panic(err)
    }
    fmt.Printf("download %d bytes\n", n)
}

func ls(s *configs.Setting) {
    client := client()

    output, err := client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input {
	Bucket: aws.String(s.BucketName),
	Prefix: aws.String(s.S3Dir),
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
