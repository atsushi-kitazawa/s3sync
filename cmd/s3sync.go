package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
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

    switch os.Args[1] {
    case "upload":
	upload.Parse(os.Args[2:])
	fmt.Println("upload subcommand.")
	fmt.Println("   force:", *forceUpload)
	fmt.Println("   mode:", *modeUpload)
    case "download":
	download.Parse(os.Args[2:])
	fmt.Println("download subcommand")
	fmt.Println("   force:", *forceDownload)
	fmt.Println("   mode:", *modeDownload)
    default:
	fmt.Println("expected 'upload' or 'download' subcommands.")
	os.Exit(1)
    }

    ls()
}

func ls() {
    cfg, err := config.LoadDefaultConfig(context.TODO())
    if err != nil {
	panic(err)
    }

    client := s3.NewFromConfig(cfg)

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
