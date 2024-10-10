package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/gil-x/goyoutubestats/fetcher"
	"github.com/gil-x/goyoutubestats/googleauth"
)

// Pour definir la liste des scopes
var scope = [...]string{
	"https://www.googleapis.com/auth/youtube.readonly",
	"https://www.googleapis.com/auth/yt-analytics.readonly",
	"https://www.googleapis.com/auth/yt-analytics-monetary.readonly",
}

// Pour les fichiers a utiliser
var (
	token  string
	client string
)

var (
	channelID string
	videoID   string
)

// Need to use an authorization code in order to get a token
// Need to compile in order to use flags
func main() {

	flag.StringVar(&token, "tok", "./token.json", "give the json token path")
	flag.StringVar(&client, "cli", "./client_secret.json", "give the json client path")
	flag.StringVar(&channelID, "chan", "UCIr96U-QJwY2plydsdPbj_A", "give a channel id")
	flag.StringVar(&videoID, "vid", "PQdJCKUpXS8", "give a video id")
	flag.Parse()

	fmt.Printf("Request with\n - Token: %s\n - Client: %s\n - Scopes: %v\n", token, client, scope)

	var tokenManager googleauth.TokenManager
	tokenManager.Init(token)

	tokenManager.SetConfigFromSecret(client, scope[:]...)

	if _, err := os.Stat(token); os.IsNotExist(err) {
		tokenManager.AskToken(token)
	} else if err != nil {
		log.Fatalln(err)
	}

	tokenManager.SetTokenFromFile(token)

	if !tokenManager.IsTokenValid() {
		fmt.Println("Token outdated")
		if err := tokenManager.RefreshToken(token); err != nil {
			log.Fatalln(err)
		}
	}

	var fetcher fetcher.Fetcher
	err := fetcher.InitYTAnalytics(tokenManager.GetConfig(), tokenManager.GetToken())
	if err != nil {
		log.Fatalln(err)
	}

	fetcher.GetVideoStats(channelID, videoID, []string{
		"views",
		// "likes",
		// "dislikes",
		// "comments",
		// "shares",
		// "estimatedMinutesWatched",
		// "averageViewDuration",
		// "estimatedRevenue",
	})

}
