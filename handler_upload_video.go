package main

import (
	"crypto/rand"
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/bootdotdev/learn-file-storage-s3-golang-starter/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerUploadVideo(w http.ResponseWriter, r *http.Request) {
	const maxMemory = 10 << 30
	r.Body = http.MaxBytesReader(w, r.Body, maxMemory)
	videoIDString := r.PathValue("videoID")
	videoID, err := uuid.Parse(videoIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid ID", err)
		return
	}
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find JWT", err)
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT", err)
		return
	}
	video, err := cfg.db.GetVideo(videoID)
	if userID != video.UserID {
		respondWithError(w, http.StatusUnauthorized, "User ID does not match", err)
		return
	}
	vid_data, _, err := r.FormFile("video")
	if err != nil {
		fmt.Printf("Could not get video form file:%v", err)
		return
	}
	defer vid_data.Close()
	media_type, _, err := mime.ParseMediaType("video/mp4")
	if err != nil {
		fmt.Printf("Could not parse media type:%v", err)
		return
	}
	if media_type != "video/mp4" {
		fmt.Printf("Invalid type. Not video/mp4")
		return
	}
	f, err := os.CreateTemp("", "tubely-upload.mp4")
	if err != nil {
		fmt.Printf("Error creating temporary file:%v", err)
		return
	}
	defer os.Remove(f.Name())
	defer f.Close()
	_, err = io.Copy(f, vid_data)
	if err != nil {
		fmt.Printf("Could not copy video to temp file:%v", err)
		return
	}
	aspect_ratio, err := getVideoAspectRatio(f.Name())
	if err != nil {
		fmt.Printf("Could not get video aspect ratio after saving to disk:%v", err)
		return
	}
	var prefix string
	if aspect_ratio == "16:9" {
		prefix = "landscape"
	} else if aspect_ratio == "9:16" {
		prefix = "portrait"
	} else {
		prefix = "other"
	}
	fmt.Printf("HERE IS THE ASPECT RATIO: %v\n", aspect_ratio)
	_, err = f.Seek(0, io.SeekStart)
	if err != nil {
		fmt.Printf("Could not reset pointer to the beginning of the video file:%v", err)
		return
	}
	output_video_file, err := processVideoForFastStart(f.Name())
	if err != nil {
		fmt.Printf("COULD NOT PROCESS VIDEO FOR FAST START:%v", err)
		return
	}
	body, err := os.Open(output_video_file)
	if err != nil {
		fmt.Printf("COULD NOT OPEN OUTPUT VIDEO FILE")
		return
	}
	key := make([]byte, 32)
	rand.Read(key)
	key_ := fmt.Sprintf("%v/%v.mp4", prefix, key)
	put_object_params := s3.PutObjectInput{
		Bucket:      &cfg.s3Bucket,
		Key:         &key_,
		Body:        body,
		ContentType: &media_type,
	}
	_, err = cfg.s3Client.PutObject(r.Context(), &put_object_params)
	if err != nil {
		fmt.Printf("Could not execute put object s3 function:%v", err)
		return
	}
	updated_vid_url := fmt.Sprintf("https://%v.s3.%v.amazonaws.com/%v", cfg.s3Bucket, cfg.s3Region, key_)
	video.VideoURL = &updated_vid_url
	cfg.db.UpdateVideo(video)
}
