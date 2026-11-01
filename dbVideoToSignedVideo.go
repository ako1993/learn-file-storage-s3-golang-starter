package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/bootdotdev/learn-file-storage-s3-golang-starter/internal/database"
)

func (cfg *apiConfig) dbVideoToSignedVideo(video database.Video) (database.Video, error) {
	vid_url := video.VideoURL
	if vid_url == nil {
		return video, fmt.Errorf("video URL pointer is nil")
	}
	fmt.Printf("VIDEO URL:%v", vid_url)
	parts := strings.Split(*vid_url, ",")
	fmt.Printf("PARTS: %v", parts)
	bucket := parts[0]
	fmt.Printf("BUCKET: %v", bucket)
	key := parts[1]
	duration := time.Minute * 5
	presigned, err := generatePresignedURL(cfg.s3Client, bucket, key, duration)
	if err != nil {
		fmt.Printf("ERROR GENERATING PRESIGNED URL:%v", err)
		return video, err
	}
	video.VideoURL = &presigned
	return video, nil
}
