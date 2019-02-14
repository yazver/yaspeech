package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/yazver/yaspeech/auth"

	"github.com/yazver/yaspeech"
)

func main() {
	oauth := os.Getenv("OAUTH")
	folderID := os.Getenv("FOLDER_ID")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	//yaspeech.InitIAMTokenContext(ctx, oauth)
	token := auth.NewIAMTokenContext(ctx, &auth.YandexAccount{OAuth: oauth})
	tts := yaspeech.NewTextToSpeech(folderID, token)

	audio, err := tts.Synthesize("Съешь ещё этих мягких французских булок, да выпей же чаю.")
	if err != nil {
		fmt.Println("Unable to convert:", err)
		return
	}
	if err := ioutil.WriteFile("audio.ogg", audio, 0644); err != nil {
		return
	}
}
