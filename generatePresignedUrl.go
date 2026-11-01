package main

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func generatePresignedURL(s3Client *s3.Client, bucket, key string, expireTime time.Duration) (string, error) {
	client := s3.NewPresignClient(s3Client)
	params := s3.GetObjectInput{
		Bucket: &bucket,
		Key:    &key,
	}
	req, err := client.PresignGetObject(context.Background(), &params, s3.WithPresignExpires(expireTime))
	if err != nil {
		fmt.Printf("ERROR CREATING PRESIGNED REQUEST:%v", err)
		return "", err
	}
	Url := req.URL
	return Url, nil
}
