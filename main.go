package main

import (
	"fmt"
	"net/http"

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
	HTTPClient     http.Client
	RemainingTimes int32
}

func NewContext(token string) *Context {
	c := http.Client()
	return &Context{
		Token:      token,
		HTTPClient: c,
	}
}

func main() {
	err := envconfig.Process("pexels", &opt)
	if err != nil {
		fmt.Printf("Failed to parse command line arguments: %s", err.Error())
	}

	ctx := Context{
		Token: opt.Token,
	}

}
