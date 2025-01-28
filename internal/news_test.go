package news

import (
	"os"
	"testing"
)

func TestSearch(t *testing.T) {
	apiKey := os.Getenv("GUARDIAN_API_KEY")
	t.Logf("client: %+v", apiKey)
	client, err := NewClient(WithAPIKey(apiKey))
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	resp, err := client.Search("deep seek r1")
	if err != nil {
		t.Fatalf("failed to search: %v", err)
	}
	t.Logf("fetched articles: %d, pages: %d", len(resp.Response.Results), resp.Response.Pages)
}

func TestFetchPage(t *testing.T) {
	apiKey := os.Getenv("GUARDIAN_API_KEY")
	t.Logf("client: %+v", apiKey)
	client, err := NewClient(WithAPIKey(apiKey))
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	query := "deep seek r1"
	resp, err := client.Search(query)
	if err != nil {
		t.Fatalf("failed to search: %v", err)
	}
	maxCount := 3
	startIndex := resp.Response.StartIndex
	if resp.Response.Pages < maxCount {
		maxCount = resp.Response.Pages
	}
	results := []NewsItem{}
	for i := startIndex; i < startIndex+maxCount; i++ {
		page, err := client.FetchPage(query, i)
		if err != nil {
			t.Fatalf("failed to fetch page: %v", err)
		}
		results = append(results, page...)
	}

	t.Logf("fetched articles: %d", len(results))
}
