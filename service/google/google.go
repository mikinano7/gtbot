package google

import (
	"google.golang.org/api/youtube/v3"
	"google.golang.org/api/googleapi/transport"
	"net/http"
	"strings"
	"github.com/spf13/viper"
)

func YouTube(query []string) ([]*youtube.SearchResult, error) {
	developerKey := viper.GetString("google.developer_key")

	client := &http.Client{
		Transport: &transport.APIKey{Key: developerKey},
	}
	svc, err := youtube.New(client)
	if err != nil {
		return nil, err
	}

	call := svc.Search.List("snippet").Q(strings.Join(query, " ")).MaxResults(3).Order("rating")
	res, err := call.Do()
	if err != nil {
		return nil, err
	}

	var arr []*youtube.SearchResult
	for _, item := range res.Items {
		switch item.Id.Kind {
		case "youtube#video": arr = append(arr, item)
		}
	}

	return arr, err
}
