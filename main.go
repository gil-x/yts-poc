package main

import (
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
		fmt.Println("Still up 2")
		var fetcher fetcher.Fetcher
		fetcher.InitYTAnalytics(tokenManager.GetConfig(), tokenManager.GetToken())
		fetcher.GetVideoStats("PQdJCKUpXS8", []string{
			"views",
			"likes",
			"dislikes",
			"comments",
			"shares",
			"estimatedMinutesWatched",
			"averageViewDuration",
			// "estimatedRevenue", Marche pas, peut être parce que j'ai rien de monétisé.
		})
	}
}

// CASE 3: I got a client_secret.json from google but no token file
func case3() {
	var tokenManager googleauth.TokenManager
	tokenManager.Init()
	tokenManager.SetConfigFromSecret("client_secret.json", "https://www.googleapis.com/auth/youtube")
	tokenManager.AskToken()
}

// Need to use an authorization code in order to get a token
func main() {
	// case1()
	case2()
	// case3()
}
