package main

type SearchResult struct {
	Page         int32   `json:"page"`
	PerPage      int32   `json:"per_page"`
	TotalResults int32   `json:"total_results"`
	NextPage     string  `json:"next_page"`
	Photos       []Photo `json:photos`
}

type Photo struct {
	Id             int32       `json:"id"`
	Width          int32       `json:"width"`
	Height         int32       `json:"height"`
	Url            string      `json:"url"`
	Photographer   string      `json:"photographer"`
	PhotographerId int32       `json:"photographer_id"`
	AvgColor       string      `json:"avg_color"`
	Src            PhotoSource `json:"src"`
	Liked          bool        `json:"liked"`
}

type PhotoSource struct {
	Original  string `json:"original"`
	Large2x   string `json:"large2x"`
	Large     string `json:"large"`
	Medium    string `json:"medium"`
	Small     string `json:"small"`
	Portrait  string `json:"portrait"`
	Landscape string `json:"lanscape"`
	Tiny      string `json:"tiny"`
}
