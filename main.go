package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
)

const (
	photoApi = "https://api.pexels.com/v1"
	videoApi = "https://api.pexels.com/videos"
)

type Client struct {
	Token         string
	hc            http.Client
	RemainingTime int32
}

type PhotSource struct {
	Original  string `json:"original"`
	Large     string `json:"large"`
	Large2x   string `json:"large2x"`
	Medium    string `json:"medium"`
	Small     string `json:"small"`
	Portrait  string `json:"portrait"`
	Square    string `json:"square"`
	Landscape string `json:"landscape"`
	Tiny      string `json:"tiny"`
}
type Photo struct {
	ID              int32      `json:"id"`
	Width           int32      `json:"width"`
	Height          int32      `json:"height"`
	URL             string     `json:"url"`
	Photographer    string     `json:"photographer"`
	PhotographerURL string     `json:"photographer_url"`
	Src             PhotSource `json:"src"`
}

type SearchResult struct {
	Page         int32   `json:"page"`
	PerPage      int32   `json:"per_page"`
	TotalResults int32   `json:"total_results"`
	NextPage     string  `json:"next_page"`
	Photos       []Photo `json:"photos"`
}

type CuratedResult struct {
	Page     int32   `json:"page"`
	PerPage  int32   `json:"per_page"`
	NextPage string  `json:"next_page"`
	Photos   []Photo `json:"photos"`
}

type Video struct {
	ID            int32          `json:"id"`
	Width         int32          `json:"width"`
	Height        int32          `json:"height"`
	URL           string         `json:"url"`
	Image         string         `json:"image"`
	FullRes       interface{}    `json:"full_res"`
	Duration      float64        `json:"duration"`
	VideoFiles    []VideoFile    `json:"video_files"`
	VideoPictures []VideoPicture `json:"video_pictures"`
}
type VideoSearchResult struct {
	Page         int32   `json:"page"`
	PerPage      int32   `json:"per_page"`
	TotalResults int32   `json:"total_results"`
	NextPage     string  `json:"next_page"`
	Videos       []Video `json:"videos"`
}

type PopularVideos struct {
	Page         int32   `json:"page"`
	PerPage      int32   `json:"per_page"`
	TotalResults int32   `json:"total_results"`
	Url          string  `json:"url"`
	Videos       []Video `json:"videos"`
}

type VideoFile struct {
	ID       int32  `json:"id"`
	Quality  string `json:"quality"`
	FileType string `json:"file_type"`
	Width    int32  `json:"width"`
	Height   int32  `json:"height"`
	Link     string `json:"link"`
}

type VideoPicture struct {
	ID      int32   `json:"id"`
	Picture string  `json:"picture"`
	Nr      int32   `json:"nr"`
	Url     string  `json:"url"`
	Videos  []Video `json:"videos"`
}

func NewClient(token string) *Client {
	c := http.Client{}
	return &Client{
		Token: token,
		hc:    c,
	}
}

// function that search photos by query
func (c *Client) SearchPhotos(query string, perPage int32, page int32) (*SearchResult, error) {
	api := fmt.Sprintf(photoApi+"/search?query=%s&per_page=%d&page=%d", query, perPage, page)
	res, err := c.requestDoWithAut("GET", api)
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, fmt.Errorf("no response received from API")
	}
	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var result SearchResult
	err = json.Unmarshal(data, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// function that get random photos

func (c *Client) GetRandomPhoto() (*Photo, error) {
	rand.Seed(time.Now().UnixNano())
	randNum := rand.Int31n(1001)
	res, err := c.CuratedPhoto(1, randNum)

	if err == nil && len(res.Photos) == 1 {
		return &res.Photos[0], nil
	}
	return nil, err
}

// function that gives you photos with pagination
func (c *Client) CuratedPhoto(perPage, page int32) (*CuratedResult, error) {
	api := fmt.Sprintf(photoApi+"/curated?per_page=%d&page=%d", perPage, page)

	res, err := c.requestDoWithAut("GET", api)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var result CuratedResult
	err = json.Unmarshal(data, &result)
	return &result, err
}

// get photo by id
func (c *Client) GetPhoto(id int32) (*Photo, error) {
	api := fmt.Sprintf(photoApi+"/%d", id)
	res, err := c.requestDoWithAut("GET", api)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	var result Photo
	err = json.Unmarshal(data, &result)
	return &result, err
}

// SearchVideos function that search videos by query

func (c *Client) SearchVideo(query, parPage, page int) (*VideoSearchResult, error) {
	api := fmt.Sprintf(videoApi+"/search?query=%s&per_page=%d&page=%d", query, parPage, page)
	res, err := c.requestDoWithAut("GET", api)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var result VideoSearchResult
	err = json.Unmarshal(data, &result)
	return &result, err
}

// PopularVideos function that get popular videos
func (c *Client) PopularVideo(perPage, page int32) (*PopularVideos, error) {
	api := fmt.Sprintf(videoApi+"/popular?per_page=%d&page=%d", perPage, page)
	res, err := c.requestDoWithAut("GET", api)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var result PopularVideos
	err = json.Unmarshal(data, &result)
	return &result, err
}

// Get random videos function
func (c *Client) GetRandomVideos() (*Video, error) {
	rand.Seed(time.Now().Unix())
	randNum := rand.Int31n(1001)
	res, err := c.PopularVideo(1, randNum)

	if err == nil && len(res.Videos) == 1 {
		return &res.Videos[0], nil
	}
	return nil, err
}

func (c *Client) GetRemainingRequest() int32 {
	return c.RemainingTime
}

// function that do the request
func (c *Client) requestDoWithAut(method, url string) (*http.Response, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", c.Token)
	res, err := c.hc.Do(req)
	if err != nil {
		return nil, err
	}

	// Check if the response is successful (2xx status code)
	if res.StatusCode >= 200 && res.StatusCode < 300 {
		// Attempt to get the X-RateLimit-Remaining header if the response is successful
		remainingTime := res.Header.Get("X-RateLimit-Remaining")
		if remainingTime != "" {
			times, err := strconv.Atoi(remainingTime)
			if err != nil {
				return nil, err
			}
			c.RemainingTime = int32(times)
		} else {
			fmt.Println("X-RateLimit-Remaining header is missing")
		}
	} else {
		// Handle non-2xx responses
		fmt.Printf("API request failed with status code: %d\n", res.StatusCode)
		return res, fmt.Errorf("API request failed with status code: %d", res.StatusCode)
	}

	return res, nil
}

func main() {
	os.Setenv("PEXELS_API_KEY", "EX7x2cXTq1iCfy9BtfkmLPuw69vU5XIv817OF0Xna2J3pAYuEWn2ISu3")
	Token := os.Getenv("PEXELS_API_KEY")

	var c = NewClient(Token)

	result, err := c.CuratedPhoto(15, 1)

	if err != nil {
		fmt.Printf("Search photos error: %v", err)
	}
	if result.Page == 0 {
		fmt.Println("No photos found")
		return
	}
	fmt.Printf("%v", result)
}
