package news

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

const defaultNewsURL = "https://content.guardianapis.com/search"

var errMissingAPIKey = errors.New("api key is required")

type NewsResponse struct {
	Response struct {
		Status      string     `json:"status"`
		Total       int        `json:"total"`
		StartIndex  int        `json:"startIndex"`
		CurrentPage int        `json:"currentPage"`
		Pages       int        `json:"pages"`
		Results     []NewsItem `json:"results"`
	} `json:"response"`
}

type NewsItem struct {
	Id                 string `json:"id"`
	WebTitle           string `json:"webTitle"`
	WebUrl             string `json:"webUrl"`
	ApiUrl             string `json:"apiUrl"`
	SectionId          string `json:"sectionId"`
	SectionName        string `json:"sectionName"`
	WebPublicationDate string `json:"webPublicationDate"`
}

type Client struct {
	apiKey string
	apiURL string
	client *http.Client
}

type ClientOption func(*Client)

func WithAPIKey(apiKey string) ClientOption {
	return func(c *Client) {
		c.apiKey = apiKey
	}
}

func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) {
		c.apiURL = baseURL
	}
}

func WithHTTPClient(client *http.Client) ClientOption {
	return func(c *Client) {
		c.client = client
	}
}

func NewClient(opts ...ClientOption) (*Client, error) {
	c := &Client{
		apiURL: defaultNewsURL,
		client: http.DefaultClient,
	}

	for _, opt := range opts {
		opt(c)
	}
	if c.apiKey == "" {
		return nil, errMissingAPIKey
	}
	return c, nil
}

func (c *Client) Search(term string) (*NewsResponse, error) {
	params := url.Values{}
	if term != "" {
		params.Add("q", term)
	}

	params.Add("api-key", c.apiKey)

	urlStr := fmt.Sprintf("%s?%s", c.apiURL, params.Encode())
	log.Println(urlStr)
	resp, err := c.client.Get(urlStr)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch news: %w", err)
	}
	defer resp.Body.Close()

	var news NewsResponse
	if err := json.NewDecoder(resp.Body).Decode(&news); err != nil {
		return nil, fmt.Errorf("fail	ed to decode response: %w", err)
	}
	return &news, nil
}

// FetchPage fetches a page of news
func (c *Client) FetchPage(query string, page int) ([]NewsItem, error) {
	params := url.Values{}
	if query != "" {
		params.Add("q", query)
	}
	params.Add("page", strconv.Itoa(page))
	params.Add("api-key", c.apiKey)

	urlStr := fmt.Sprintf("%s?%s", c.apiURL, params.Encode())
	log.Println(urlStr)
	resp, err := c.client.Get(urlStr)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch news: %w", err)
	}
	defer resp.Body.Close()

	var news NewsResponse
	if err := json.NewDecoder(resp.Body).Decode(&news); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	return news.Response.Results, nil
}
