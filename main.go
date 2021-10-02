package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
	"time"

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

// search Pexels for any topic
func (c *Context) SearchPhotos(query string, perPage, page int) (*SearchResult, error) {
	url := fmt.Sprintf(PhotoApi+"/search?query=%s&per_page=%d&page=%d", query, perPage, page)
	resp, err := c.requestDoWithAuth("GET", url)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

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

//  Receive real-time photos curated by the Pexels team
func (c *Context) CuratedPhotos(perPage, page int) (*CuratedResult, error) {
	url := fmt.Sprintf(PhotoApi+"/curated?per_page=%d&page=%d", perPage, page)
	resp, err := c.requestDoWithAuth("GET", url)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	var res CuratedResult
	err = json.Unmarshal(data, &res)

	if err != nil {
		return nil, err
	}

	return &res, nil
}

// Retrieve a specific Photo from its id.
func (c *Context) GetPhoto(id int32) (*Photo, error) {
	url := fmt.Sprintf(PhotoApi+"/photos/%d", id)
	resp, err := c.requestDoWithAuth("GET", url)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	var res Photo
	err = json.Unmarshal(data, &res)

	if err != nil {
		return nil, err
	}

	return &res, nil
}

// Get random photo
func (c *Context) GetRandomPhoto() (*Photo, error) {
	rand.Seed(time.Now().Unix())
	randNum := rand.Intn(1001)
	res, err := c.CuratedPhotos(1, randNum)
	if err == nil && len(res.Photos) == 1 {
		return &res.Photos[0], nil
	}

	return nil, err
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

	fmt.Println(result)

}
