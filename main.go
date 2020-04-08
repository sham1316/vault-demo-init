package main

import (
	"fmt"
	"github.com/hashicorp/vault/api"
	"log"
	"net/http"
	"os"
	"time"
)

var httpClient = &http.Client{
	Timeout: 10 * time.Second,
}

var tokenCreateRequest = &api.TokenCreateRequest{
	Policies:       []string{"reader"},
	TTL:            "10m",
	ExplicitMaxTTL: "10m",
	NumUses:        1,
}

func main() {
	vaultAddr := os.Getenv("VAULT_ADDR")
	vaultTokenPath := os.Getenv("VAULT_TOKEN_PATH")
	vaultMasterToken := os.Getenv("VAULT_MASTER_TOKEN")
	fmt.Println("VAULT_ADDR:", vaultAddr)
	fmt.Println("VAULT_TOKEN_PATH:", vaultTokenPath)
	fmt.Println("VAULT_MASTER_TOKEN:", vaultMasterToken)

	var err error
	var client *api.Client
	if client, err = api.NewClient(&api.Config{Address: vaultAddr, HttpClient: httpClient}); err != nil {
		panic(err)
	}
	client.SetToken(vaultMasterToken)

	var secret *api.Secret
	if secret, err = client.Auth().Token().Create(tokenCreateRequest); err != nil {
		panic(err)
	}
	var file *os.File
	if file, err = os.OpenFile(vaultTokenPath, os.O_CREATE, 666); err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	file.Write([]byte(secret.Auth.ClientToken))
}
