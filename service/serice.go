package service

import (
	"./itunes"
	"fmt"
	"time"
)

func ITunes(query []string) string {
	if res, err := itunes.Search(query); err != nil {
		return fmt.Sprintf(
			"%s [%s]",
			err.Error(),
			time.Now().String(),
		)
	} else {
		if len(res.Results) > 0 {
			return fmt.Sprintf(
				"%s - %s | %s",
				res.Results[0].TrackName,
				res.Results[0].ArtistName,
				res.Results[0].PreviewUrl,
			)
		} else {
			return fmt.Sprintf(
				"検索結果が0件でした。 [%s]",
				time.Now().String(),
			)
		}
	}
}
