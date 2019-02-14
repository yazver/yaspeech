package main

import (
	"fmt"
	"os"

	"github.com/yazver/yaspeech/auth"

	"github.com/yazver/yaspeech"
)

func main() {
	oauth := os.Getenv("OAUTH")
	folderID := os.Getenv("FOLDER_ID")

	// ctx, cancel := context.WithCancel(context.Background())
	// defer cancel()
	// token := auth.NewIAMTokenContext(ctx, &auth.YandexAccount{OAuth: oauth})
	token := auth.NewIAMToken(&auth.YandexAccount{OAuth: oauth})
	stt := yaspeech.NewSpeechToText(folderID, token)

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
