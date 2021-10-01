package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/kelseyhightower/envconfig"
)

const (
	PhotoApi = "https://api.pexels.com/v1"
	VideoApi = "https://api.pexels.com/videos"
)

// Configaration options from environment variables
var opt struct {
	Token string `default:"563492ad6f91700001000001580423ce8af84b89bd62c4749afa4d0e"`
}

type Context struct {
	Token          string
	HTTPCl         http.Client
	RemainingTimes int32
}

func (c *Context) SearchPhotos(query string, perPage, page int) (*SearchResult, error) {
	url := fmt.Sprintf(PhotoApi+"/search?query=%s&per_page=%d&page=%d", query, perPage, page)
	resp, err := c.requestDoWithAuth("GET", url)
	defer resp.Body.Close()

	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	var res SearchResult
	err = json.Unmarshal(data, &res)

	if err != nil {
		return nil, err
	}

	return &res, nil
}

func (c *Context) requestDoWithAuth(method, url string) (*http.Response, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", c.Token)
	res, err := c.HTTPCl.Do(req)

	if err != nil {
		return nil, err
	}

	times, err := strconv.Atoi(res.Header.Get("X-Ratelimit-Remaining"))
	if err != nil {
		return nil, err
	}

	c.RemainingTimes = int32(times)

	return res, nil

}

func NewContext(token string) *Context {
	c := http.Client{}
	return &Context{
		Token:  token,
		HTTPCl: c,
	}
}

func main() {
	err := envconfig.Process("pexels", &opt)
	if err != nil {
		fmt.Printf("Failed to parse command line arguments: %s", err.Error())
	}

	ctx := NewContext(opt.Token)

	result, err := ctx.SearchPhotos("waves", 15, 2)

	if err != nil {
		fmt.Errorf("Search errpr:%v", err)
	}

	if result.Page == 0 {
		fmt.Errorf("Search result wrong")
	}

	fmt.Printf(result)

}
