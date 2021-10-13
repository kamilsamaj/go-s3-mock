package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"io/ioutil"
	"log"
)

type S3GetObjectAPI interface {
	GetObject(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error)
}

func GetObjectFromS3(ctx context.Context, api S3GetObjectAPI, bucket, key string) ([]byte, error) {
	object, err := api.GetObject(ctx, &s3.GetObjectInput{
		Bucket: &bucket,
		Key:    &key,
	})
	if err != nil {
		return nil, err
	}
	defer object.Body.Close()

	return ioutil.ReadAll(object.Body)
}

func main() {
	cfg, err := config.LoadDefaultConfig(context.Background())

	if err != nil {
		log.Fatalln(err)
	}

	s3Client := s3.NewFromConfig(cfg)
	bytes, err := GetObjectFromS3(context.Background(), s3Client, "accolade-api-swaggers", "accolade_jarvis_v1.yml")
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(string(bytes))
}

