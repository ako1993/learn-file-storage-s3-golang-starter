package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
)

type FFProbeData struct {
	Streams []struct {
		CodecName string `json:"codec_name"`
		CodecType string `json:"codec_type"`
		Width     int    `json:"width,omitempty"`
		Height    int    `json:"height,omitempty"`
		Duration  string `json:"duration,omitempty"`
	} `json:"streams"`
	Format struct {
		Filename       string `json:"filename"`
		Duration       string `json:"duration"`
		Size           string `json:"size"`
		BitRate        string `json:"bit_rate"`
		FormatLongName string `json:"format_long_name"`
	} `json:"format"`
}

func getVideoAspectRatio(filePath string) (string, error) {
	cmd := exec.Command("ffprobe", "-v", "error", "-print_format", "json", "-show_streams", filePath)
	var b bytes.Buffer
	cmd.Stdout = &b
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Error running command to get aspect ratio:%v", err)
		return "", err
	}
	var vid_info FFProbeData
	err = json.Unmarshal(b.Bytes(), &vid_info)
	if err != nil {
		fmt.Printf("Could not unmarshal bytes to Json:%v", err)
		return "", err
	}
	var videoWidth, videoHeight int
	for _, stream := range vid_info.Streams {
		if stream.CodecType == "video" {
			videoWidth = stream.Width
			videoHeight = stream.Height
			break // Found the video stream, no need to continue
		}
	}
	w, h, err := nearestStandardRatio(videoWidth, videoHeight)
	fmt.Printf("VIDEO WIDTH:%v\n", videoWidth)
	fmt.Printf("VIDEO HEIGHT:%v\n", videoHeight)
	if err != nil {
		fmt.Printf("Error getting aspect ratio:%v", err)
		return "", err
	}
	aspect_ratio := fmt.Sprintf("%v:%v", w, h)
	return aspect_ratio, nil
}
