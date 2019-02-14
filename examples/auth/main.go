package main

import (
	"fmt"
	"log"
	"os"

	"github.com/yazver/yaspeech/auth"
)

func main() {
	oauth := os.Getenv("OAUTH")
	folderID := os.Getenv("FOLDER_ID")
	serviceAccountName := os.Getenv("SERVICE_ACCOUNT_NAME")

	token := auth.NewIAMToken(&auth.YandexAccount{OAuth: oauth})

	// Getting a list of service accounts for specified folder.
	accounts, err := auth.ServiceAccounts(folderID, token)
	if err != nil {
		log.Fatalln("Unable get accounts:", err)
	}
	fmt.Printf("Service accounts list: %v\n", accounts)
	accountID, ok := accounts[serviceAccountName]
	if !ok {
		log.Fatalf("The account \"%s\" don't exist.\n", serviceAccountName)
	}

	// Getting the key data of the service account. Can be executed once.
	keyID, privateKey, err := auth.CreateServiceAccountKey(accountID, token)
	if err != nil {
		log.Fatalln("Unable get the key data:", err)
	}
	fmt.Printf("Key ID: %s\nPrivate key:\n%s", keyID, privateKey)

	// Getting IAM token.
	token = auth.NewIAMToken(&auth.ServiceAccount{ID: accountID, KeyID: keyID, PrivateKey: privateKey})
	iam, err := token.Get()
	if err != nil {
		log.Fatalln("Unable get accounts:", err)
	}
	fmt.Printf("IAM token: %s\n", iam)
}
