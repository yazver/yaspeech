package main

import (
	"fmt"
	"os"

	"github.com/yazver/yaspeech"
)

func main() {
	oauth := os.Getenv("OAUTH")
	folderID := os.Getenv("FOLDER_ID")

	// ctx, cancel := context.WithCancel(context.Background())
	// defer cancel()
	// yaspeech.InitIAMTokenContext(ctx, oauth)
	yaspeech.InitIAMToken(oauth)
	stt := yaspeech.NewSpeechToText(folderID)

	file, err := os.Open("audio.ogg")
	if err != nil {
		fmt.Println("Unable to open file:", err)
		return
	}
	text, err := stt.Recognize(file)
	if err != nil {
		fmt.Println("Unable to convert:", err)
		return
	}
	fmt.Println(text)
}
