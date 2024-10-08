package googleauth

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type TokenManager struct {
	config    *oauth2.Config
	tokenFile string
	token     *oauth2.Token
	ctx       context.Context
}

// NewTokenManager create a new token manager
func (t *TokenManager) Init() *TokenManager {
	token, err := t.loadToken("token.json")
	if err != nil {
		fmt.Printf("No token or corrupted file: %v", err)
	}

	return &TokenManager{
		ctx:    context.Background(),
		token:  token,
		config: &oauth2.Config{},
	}
}

func (t *TokenManager) AskToken() {

	ctx := context.Background()

	authURL := t.config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Visit this link and get an authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Error when reading authorization code: %v", err)
	}

	newToken, err := t.config.Exchange(ctx, authCode)
	if err != nil {
		log.Fatalf("Erreur lors de l'échange du token OAuth: %v", err)
	}

	t.token = newToken
	t.saveToken("token.json")
}

// LoadToken load token from file
func (t *TokenManager) loadToken(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	token := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(token)
	return token, err
}

// SetTokenFile set token file path
func (t *TokenManager) SetTokenFile(path string) {
	t.tokenFile = path
}

// SaveToken save token to file
func (t *TokenManager) saveToken(path string) {
	fmt.Printf("Enregistrement du token dans %s\n", path)
	f, err := os.Create(path)
	if err != nil {
		log.Fatalf("Impossible de créer le fichier de token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(t.token)
}

func (t *TokenManager) SetConfigFromSecret(secret string, scope ...string) {

	b, err := os.ReadFile(secret)
	if err != nil {
		log.Fatalf("Error during reading file %s: %v", secret, err)
	}

	// config, err := google.ConfigFromJSON(b, scope)
	config, err := google.ConfigFromJSON(b, scope...)

	if err != nil {
		log.Fatalf("Erreur during OAuth2 configuration: %v", err)
	}

	t.config = config
}

func (t *TokenManager) SetTokenFromFile(tokenPath string) {
	f, err := os.Open(tokenPath)
	if err != nil {
		log.Fatalf("erreur lors du chargement du token : %v", err)
	}
	defer f.Close()
	token := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(token)
	if err != nil {
		log.Fatalf("erreur lors du décodage du token : %v", err)
	}
	t.token = token
	fmt.Println("Token successfully loaded")
}

// LoadToken load token from file typically names client_secret_xxxxx.json, downloaded from Google Cloud.
func loadConfigFromJSON(tokenFile, scope string) (*oauth2.Config, error) {
	b, err := os.ReadFile(tokenFile)
	if err != nil {
		log.Fatalf("Erreur lors de la lecture du fichier client_secret.json: %v", err)
	}
	config, err := google.ConfigFromJSON(b, scope)
	if err != nil {
		log.Fatalf("Erreur lors de la configuration OAuth2: %v", err)
	}
	return config, nil
}

func (t *TokenManager) GetConfig() *oauth2.Config {
	return t.config
}

func (t *TokenManager) GetToken() *oauth2.Token {
	return t.token
}

func (t *TokenManager) IsTokenValid() bool {
	fmt.Println("Still up 1")
	fmt.Printf("Token: %v\n", t.token)
	return t.token.Valid()
}

func (t *TokenManager) RefreshToken() error {
	fmt.Println("Still up 3")
	if !t.IsTokenValid() {
		ctx := context.Background()
		newToken, err := t.config.TokenSource(ctx, t.token).Token()
		if err != nil {
			return fmt.Errorf("failed to refresh token: %v", err)
		}

		t.token = newToken
		t.saveToken("token.json")

		log.Println("Token refresh and saved successfully.")
	} else {
		log.Println("Token is valid, no need to refresh.")
	}

	return nil
}
