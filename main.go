package main

import (
	"./service"
	"fmt"
	"github.com/ChimeraCoder/anaconda"
	"github.com/joho/godotenv"
	"log"
	"net/url"
	"os"
	"strings"
	"time"
)

type MyTime struct {
	time time.Time
}

func (t MyTime) fmt() string {
	return t.time.In(time.FixedZone("Asia/Tokyo", 9*60*60)).Format(time.RFC3339)
}

func main() {
	err := godotenv.Load("go.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	anaconda.SetConsumerKey(os.Getenv("TWITTER_CONSUMER_KEY"))
	anaconda.SetConsumerSecret(os.Getenv("TWITTER_CONSUMER_SECRET"))

	api := anaconda.NewTwitterApi(
		os.Getenv("TWITTER_OAUTH_TOKEN"),
		os.Getenv("TWITTER_OAUTH_TOKEN_SECRET"),
	)
	defer api.Close()

	stream := api.UserStream(url.Values{})

	defer func() {
		err := recover()
		if err != nil {
			api.PostTweet(fmt.Sprint(err), url.Values{})
		}
	}()

	for {
		select {
		case item := <-stream.C:
			switch status := item.(type) {
			case anaconda.Tweet:
				fmt.Printf("%s: %s\n", status.User.ScreenName, status.Text)
				if text := pattern(status); text != "" {
					v := url.Values{"in_reply_to_status_id": []string{status.IdStr}}
					if _, err := api.PostTweet(reply(status.User.ScreenName, text), v); err != nil {
						api.PostTweet(reply(status.User.ScreenName, onError(err)), v)
					}
				}
			default:
			}
		}
	}
}

func pattern(status anaconda.Tweet) string {
	if strings.HasPrefix(status.Text, ":") {
		return command(status)
	} else {
		return ""
	}
}

func command(status anaconda.Tweet) string {
	arr := strings.Split(status.Text, " ")
	if len(arr) > 1 {
		command, query := arr[0], arr[1:len(arr)]
		switch command {
		case ":m":
			return service.ITunes(query)
		case ":y":
			return service.YouTube(query)
		case ":x":
			return service.Xvideos(query)
		default:
			return ""
		}
	} else {
		switch status.Text {
		case ":go":
			return fmt.Sprintf(
				"ʕ ◔ϖ◔ʔ [%s]",
				&MyTime{time.Now()}.fmt(),
			)
		case ":test":
			return fmt.Sprintf(
				"process is running. [%s]",
				&MyTime{time.Now()}.fmt(),
			)
		default:
			return ""
		}
	}
}

func onError(err error) string {
	return fmt.Sprintf(
		"%s [%s]",
		err.Error(),
		&MyTime{time.Now()}.fmt(),
	)
}

func reply(user string, status string) string {
	return fmt.Sprintf("@%s %s", user, status)
}
