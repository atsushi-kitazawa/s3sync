package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func main() {
    _main()
}

func _main() {
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
