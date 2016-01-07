package service

import (
	"./google"
	"./itunes"
	"fmt"
	"time"
)

func ITunes(query []string) string {
	if res, err := itunes.Search(query); err != nil {
		return err.Error()
	} else {
		if len(res) > 0 {
			return fmt.Sprintf(
				"%s / %s - %s",
				res[0].TrackName,
				res[0].ArtistName,
				res[0].PreviewUrl,
			)
		} else {
			return fmt.Sprintf(
				"検索結果が0件でした。 [%s]",
				time.Now().String(),
			)
		}
	}
}

func YouTube(query []string) string {
	if res, err := google.YouTube(query); err != nil {
		return err.Error()
	} else {
		if len(res) > 0 {
			return fmt.Sprintf(
				"%s - %s%s",
				res[0].Snippet.Title,
				"https://www.youtube.com/watch?v=",
				res[0].Id.VideoId,
			)
		} else {
			return fmt.Sprintf(
				"検索結果が0件でした。 [%s]",
				time.Now().String(),
			)
		}
	}
}

func onError(err error) string {
	return fmt.Sprintf(
		"%s [%s]",
		err.Error(),
		time.Now().String(),
	)
}
