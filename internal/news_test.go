package news_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)


func TestSearch(t *testing.T) {
	apiKey := os.Getenv("GUARDIAN_API_KEY")
	client, err := news.NewClient(news.WithAPIKey(apiKey))
	assert.NoError(t, err)

	news, err := client.Search("test")
	assert.NoError(t, err)
	assert.NotNil(t, news)
}
