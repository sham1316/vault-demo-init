package main

import (
	"fmt"
	"github.com/hashicorp/vault/api"
	"io/ioutil"
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

	err = ioutil.WriteFile(vaultTokenPath, []byte(secret.Auth.ClientToken), 0644)
	if err != nil {
		panic(err)
	}

	b, err := ioutil.ReadFile(vaultTokenPath)
	if err != nil {
		panic(err)
	}
	fmt.Println("VAULT_TOKEN:", string(b))
}
