package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
)

// curl 'https://play.golang.org/compile' -H 'origin: https://play.golang.org' -H 'accept-encoding: gzip, deflate' -H 'x-requested-with: XMLHttpRequest' -H 'accept-language: en-US,en;q=0.8,sv;q=0.6,nb;q=0.4' -H 'user-agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/46.0.2490.86 Safari/537.36' -H 'content-type: application/x-www-form-urlencoded; charset=UTF-8' -H 'accept: application/json, text/javascript, */*; q=0.01' -H 'referer: https://play.golang.org/' -H 'authority: play.golang.org' -H 'cookie: _ga=GA1.2.2084321319.1441558441; __utma=110886291.2084321319.1441558441.1448479986.1448487895.37; __utmc=110886291; __utmz=110886291.1448487895.37.32.utmcsr=google|utmccn=(organic)|utmcmd=organic|utmctr=(not%20provided); __utmt=1; __utmt_b=1; __utma=43835492.1239961236.1447220376.1448369578.1448570945.4; __utmb=43835492.3.9.1448570954414; __utmc=43835492; __utmz=43835492.1447612074.2.2.utmcsr=google|utmccn=(organic)|utmcmd=organic|utmctr=(not%20provided)' --data 'version=2&body=package+main%0A%0Aimport+%22fmt%22%0A%0Afunc+main()+%7B%0A%09fmt.Println(%22Hello%2C+playground%22)%0A%7D%0A' --compressed

// {"Errors":"","Events":[{"Message":"Hello, playground\n","Kind":"stdout","Delay":0}]}

type Event struct {
	Message string `json:"Message"`
	Kind    string `json:"Kind"`
	Delay   int    `json:"Delay"`
}
type PlayResponse struct {
	Errors string  `json:"Errors"`
	Events []Event `json:"Events"`
}

const compileUrl = "https://play.golang.org/compile"

func runGoPkg(code string) (*PlayResponse, error) {
	resp, err := http.PostForm(
		compileUrl,
		url.Values{
			"version": {"2"},
			"body":    {code},
		})

	if resp.StatusCode == 200 {
		var response PlayResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		if err != nil {
			return nil, err
		} else {
			return &response, nil
		}
	}

	return nil, errors.New("Failed to run Go code, status: " + resp.Status)
}
