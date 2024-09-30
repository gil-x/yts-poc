package client

import (
	"log"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type Client struct {
	API    string
	ID     string
	Secret string
	config *oauth2.Config
}

func (c *Client) New(secret, scope string) error {
	b, err := os.ReadFile(secret)
	if err != nil {
		log.Fatalf("Error during reading file %s: %v", secret, err)
	}

	config, err := google.ConfigFromJSON(b, scope)
	if err != nil {
		log.Fatalf("Erreur during OAuth2 configuration: %v", err)
	}

	c.config = config

	return nil
}
