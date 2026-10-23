package main

import (
	"fmt"
	"os/exec"
)

func processVideoForFastStart(filePath string) (string, error) {
	output_file_path := fmt.Sprintf("%v.processing", filePath)
	cmd := exec.Command("ffmpeg", "-i", filePath, "-c", "copy", "-movflags", "faststart", "-f", "mp4", output_file_path)
	err := cmd.Run()
	if err != nil {
		fmt.Printf("ERROR EXECUTING FAST START COMMAND:%v", err)
		return "", err
	}
	return output_file_path, nil
}
