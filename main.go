package main

import (
	"fmt"
	"github.com/ChimeraCoder/anaconda"
	"github.com/joho/godotenv"
	"log"
	"net/url"
	"os"
	"time"
	"strings"
)

func main() {
	err := godotenv.Load("twitter_oauth.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	anaconda.SetConsumerKey(os.Getenv("TWITTER_CONSUMER_KEY"))
	anaconda.SetConsumerSecret(os.Getenv("TWITTER_CONSUMER_KEY"))

	api := anaconda.NewTwitterApi(os.Getenv("TWITTER_OAUTH_TOKEN"), os.Getenv("TWITTER_OAUTH_TOKEN_SECRET"))
	defer api.Close()

	v := url.Values{}
	stream := api.UserStream(v)

	for {
		select {
		case item := <-stream.C:
			switch status := item.(type) {
			case anaconda.Tweet:
				fmt.Printf("%s: %s\n", status.User.ScreenName, status.Text)
				api.PostTweet(pattern(status), v)
			default:
			}
		}
	}
}

func pattern(status anaconda.Tweet) string {
	if strings.HasPrefix(status.Text, ":") {
		return command(status)
	} else {
		switch status.Text {
		case "Go": return fmt.Sprintf("@%s ʕ ◔ϖ◔ʔ", status.User.ScreenName)
		default: return ""
		}
	}
}

func command(status anaconda.Tweet) string {
	switch status.Text {
	case ":test": return fmt.Sprintf("@%s process is running. %s", status.User.ScreenName, time.Now().String())
	default: return ""
	}
}
