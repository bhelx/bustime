package lib

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/bhelx/bustime"
)

type BusErr struct {
	Rt  string `json:"rt"`
	Msg string `json:"msg"`
}

type BustimeData struct {
	Vehicles []Vehicle `json:"vehicle"`
	Errors   []BusErr  `json:"error"`
}

type BustimeResponse struct {
	Data BustimeData `json:"bustime-response"`
}

type Scraper struct {
	client  *http.Client
	storage Storage
	config  *bustime.Config
}

func NewScraper(config *bustime.Config) *Scraper {
	var client *http.Client
	tr := &http.Transport{
		MaxIdleConnsPerHost: 1024,
		TLSHandshakeTimeout: 0 * time.Second,
	}
	client = &http.Client{Transport: tr}
	storage, err := NewStorage()
	if err != nil {
		fmt.Println(err)
	}
	return &Scraper{
		client,
		storage,
		config,
	}
}

func (s *Scraper) Start() {
	for {
		result := s.fetch()
		for _, vehicle := range result.Vehicles {
			err := s.storage.Store(&vehicle)
			if err != nil {
				fmt.Println(err)
			}
		}
		fmt.Printf("Found %d vehicles\n", len(result.Vehicles))
		time.Sleep(5 * time.Second)
	}
}

func (v *Scraper) fetch() *BustimeData {
	key := v.config.Key
	baseURL := v.config.Url
	url := fmt.Sprintf("%s?key=%s&tmres=m&rtpidatafeed=bustime&format=json", baseURL, key)
	resp, err := v.client.Get(url)
	if resp.Body != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		log.Fatal(err)
	}
	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	result := &BustimeResponse{}
	json.Unmarshal(body, result)

	return &result.Data
}

func (v *Scraper) Close() {
	v.storage.Close()
	v.client.CloseIdleConnections()
}
