package itunes

import (
	"fmt"
	"encoding/json"
	"net/http"
	"io/ioutil"
	"io"
	"strings"
)

const (
	searchUrl = "https://itunes.apple.com/search?version=2"
)

type Request struct {
	Term string
	Country string
	Media string
	Entity string
	Limit int
	Lang string
	Explicit bool
}

type Response struct {
	Results []Result `json:"results"`
}

type Result struct {
	WrapperType string `json:"wrapperType"`
	Kind string `json:"kind"`
	TrackName string `json:"trackName"`
	ArtistName string `json:"artistName"`
	CollectionName string `json:"collectionName"`
	ArtworkUrl100 string `json:"artworkUrl100"`
	PreviewUrl string `json:"previewUrl"`
	TrackTimeMillis int64 `json:"trackTimeMillis"`
}

func Search(query []string) (*Response, error) {
	term := strings.Join(query, "+")
	req := newRequest(term, 50)
	jpUrl := fmt.Sprintf("%s&country=%s&lang=%s", searchUrl, req.Country, req.Lang) //jp
	musicUrl := fmt.Sprintf("%s&media=%s&entity=%s&explicit=%t", jpUrl, req.Media, req.Entity, req.Explicit) //music
	requestUrl := fmt.Sprintf("%s&term=%s&limit=%d", musicUrl, term, req.Limit) //req

	client := http.Client{}
	response, err := client.Get(requestUrl)
	defer response.Body.Close()
	if err != nil {
		return nil, err
	} else {
		return parse(response.Body)
	}
}

func parse(body io.Reader) (*Response, error) {
	res, err := ioutil.ReadAll(body)

	response := new(Response)
	if err = json.Unmarshal(res, response); err != nil {
		return nil, err
	} else {
		return response, nil
	}
}

func newRequest(term string, limit int) Request {
	return Request{
		Term:term,
		Country:"JP",
		Media:"music",
		Entity:"musicTrack",
		Limit:limit,
		Lang:"ja_jp",
		Explicit:true,
	}
}
