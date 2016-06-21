package main

import (
	"./service"
	"errors"
	"fmt"
	"github.com/ChimeraCoder/anaconda"
	"github.com/spf13/viper"
	"log"
	"net/url"
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
	initConfig()

	anaconda.SetConsumerKey(viper.GetString("twitter.consumer_key"))
	anaconda.SetConsumerSecret(viper.GetString("twitter.consumer_secret"))
	botName := viper.GetString("twitter.bot_name")

	api := anaconda.NewTwitterApi(
		viper.GetString("twitter.access_token"),
		viper.GetString("twitter.access_token_secret"),
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
				if text := pattern(status, botName); text != "" {
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

func initConfig() {
	viper.SetConfigFile("config.json")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalln(err.Error())
	}
}

func pattern(status anaconda.Tweet, botName string) string {
	if strings.HasPrefix(status.Text, fmt.Sprintf("@%s :", botName)) {
		return command(status)
	} else {
		return ""
	}
}

func command(status anaconda.Tweet) string {
	arr := strings.Split(status.Text, " ")
	if len(arr) > 1 {
		command, query := arr[1], arr[2:len(arr)]
		switch command {
		case ":m":
			return service.ITunes(query)
		case ":y":
			return service.YouTube(query)
		case ":x":
			return service.Xvideos(query)
		case ":sc":
			return service.Soundcloud(query)
		case ":up":
			if status.User.ScreenName == viper.GetString("twitter.owner_name") {
				if len(status.Entities.Media) > 0 {
					return service.DropboxUpload(status.Entities.Media[0].Media_url)
				} else {
					return service.DropboxUpload(status.Entities.Urls[0].Expanded_url)
				}
			} else {
				return onError(errors.New("ダメだよ。"))
			}
		default:
			return ""
		}
	} else {
		switch status.Text {
		case ":go":
			return fmt.Sprintf(
				"ʕ ◔ϖ◔ʔ [%s]",
				MyTime{time.Now()}.fmt(),
			)
		case ":test":
			return fmt.Sprintf(
				"process is running. [%s]",
				MyTime{time.Now()}.fmt(),
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
		MyTime{time.Now()}.fmt(),
	)
}

func reply(user string, status string) string {
	return fmt.Sprintf("@%s %s", user, status)
}
