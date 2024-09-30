package fetcher

import (
	"context"
	"fmt"
	"log"
	"strings"

	"golang.org/x/oauth2"
	"google.golang.org/api/option"
	"google.golang.org/api/youtubeanalytics/v2"
)

type Fetcher struct {
	ytaService *youtubeanalytics.Service
}

// func Init() *Fetcher {
// 	return &Fetcher{
// 	}
// }

func (f *Fetcher) InitYTAnalytics(config *oauth2.Config, token *oauth2.Token) error {

	ctx := context.Background()

	service, err := youtubeanalytics.NewService(ctx, option.WithTokenSource(config.TokenSource(ctx, token)))

	if err != nil {
		log.Fatalf("Unable to create YouTube Analytics service: %v", err)
	}

	f.ytaService = service

	return nil
}

func (f *Fetcher) GetVideoStats(videoID string, metrics []string) {

	startDate := "2000-08-01"
	endDate := "2023-09-01"

	metricsString := strings.Join(metrics, ",")

	call := f.ytaService.Reports.Query().
		Ids("channel==MINE").
		StartDate(startDate).
		EndDate(endDate).
		Filters("video==" + videoID).
		Metrics(metricsString)
		// Pour avoir le genre, il faut ajouter un autre filtre .Dimensions("gender")

	response, err := call.Do()
	if err != nil {
		log.Fatalf("Erreur lors de la récupération des statistiques: %v", err)
	}

	fmt.Printf("Statistiques pour %s:\n", videoID)

	for _, row := range response.Rows {
		for i, v := range row {
			fmt.Printf("- %s: %v\n", metrics[i], v)
		}
	}
}
