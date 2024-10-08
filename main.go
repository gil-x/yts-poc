package main

import (
	"flag"
	"fmt"

	"github.com/gil-x/goyoutubestats/fetcher"
	"github.com/gil-x/goyoutubestats/googleauth"
)

// CASE 1: I got an valid token from google
func case1() {
	var tokenManager googleauth.TokenManager
	tokenManager.Init()
	tokenManager.SetTokenFromFile("token.json")
	if tokenManager.IsTokenValid() {
		fmt.Println("Token valid, no need to ask")
	} else {
		fmt.Println("Token not valid, try Case 2")
	}
}

// CASE 2: I got an outdated token from google
func case2() {
	var tokenManager googleauth.TokenManager
	tokenManager.Init()
	tokenManager.SetConfigFromSecret("client_secret.json", "https://www.googleapis.com/auth/youtube")
	tokenManager.SetTokenFromFile("token.json")
	if tokenManager.IsTokenValid() {
		fmt.Println("Token valid, no need to ask")
	} else {
		fmt.Println("Token outdated")
		tokenManager.RefreshToken()
	}
	var fetcher fetcher.Fetcher
	fetcher.InitYTAnalytics(tokenManager.GetConfig(), tokenManager.GetToken())
}

// CASE 3: I got a client_secret.json from google but no token file
func case3() {
	var tokenManager googleauth.TokenManager
	tokenManager.Init()
	tokenManager.SetConfigFromSecret("client_secret.json", "https://www.googleapis.com/auth/youtube.readonly", "https://www.googleapis.com/auth/yt-analytics.readonly", "https://www.googleapis.com/auth/yt-analytics-monetary.readonly", "https://www.googleapis.com/auth/youtubepartner.readonly")
	tokenManager.AskToken()
}

// Need to use an authorization code in order to get a token
// Need to compile in order to use flags
func main() {

	useCase := flag.Int("case", 1, "Use case (1, 2, 3)")
	flag.Parse()
	fmt.Printf("Use case: %d\n", *useCase)

	switch *useCase {
	case 1:
		case1()
	case 2:
		case2()
	case 3:
		case3()
	}
}
