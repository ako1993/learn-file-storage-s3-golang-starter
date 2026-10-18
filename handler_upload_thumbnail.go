package main

import (
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/bootdotdev/learn-file-storage-s3-golang-starter/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerUploadThumbnail(w http.ResponseWriter, r *http.Request) {
	const maxMemory = 10 << 20
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

	fmt.Println("uploading thumbnail for video", videoID, "by user", userID)

	r.ParseMultipartForm(maxMemory)
	file, header, err := r.FormFile("thumbnail")
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Unable to parse form file", err)
		return
	}
	defer file.Close()
	media_type := header.Header.Get("Content-Type")
	media_type_, _, err := mime.ParseMediaType(media_type)
	fmt.Println(media_type_)
	if media_type_ != "image/jpeg" || media_type_ != "image/png" {
		respondWithError(w, 401, "Wrong type!", err)
		return
	}
	data, err := io.ReadAll(file)
	if err != nil {
		fmt.Printf("Error reading file data:%v", err)
		return
	}
	video, err := cfg.db.GetVideo(videoID)
	if userID != video.UserID {
		respondWithError(w, http.StatusUnauthorized, "User ID does not match", err)
		return
	}
	thumbnail_extension := strings.TrimPrefix(media_type, "image/")
	id_and_extension := fmt.Sprintf("%v.%v", video.ID.String(), thumbnail_extension)
	full_url := filepath.Join(cfg.assetsRoot, id_and_extension)
	new_file, err := os.Create(full_url)
	if err != nil {
		fmt.Printf("Error creating file:%v", err)
		return
	}
	data_string := string(data)
	reader := strings.NewReader(data_string)
	_, err = io.Copy(new_file, reader)
	if err != nil {
		fmt.Printf("Could not copy data to file:%v", err)
		return
	}
	first_half := fmt.Sprintf("http://localhost:%v", cfg.port)
	thumbnail_url := fmt.Sprintf("%v/%v", first_half, full_url)
	video.ThumbnailURL = &thumbnail_url
	err = cfg.db.UpdateVideo(video)
	if err != nil {
		fmt.Printf("Could not update thumbnail URL in Database:%v", err)
		return
	}
	respondWithJSON(w, http.StatusOK, video)
}
