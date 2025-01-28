package services

import (
	"encoding/json"
	"github.com/cloudquery/plugin-sdk/v4/faker"
	"github.com/wprzechrzta/cq-source-guardian/client"
	news "github.com/wprzechrzta/cq-source-guardian/internal"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewsTable(t *testing.T) {
	var newsData news.NewsResponse
	if err := faker.FakeObject(&newsData); err != nil {
		t.Fatal(err)
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		d, _ := json.Marshal(newsData)
		_, _ = w.Write(d)
	}))
	defer ts.Close()

	client.TestHelper(t, NewsTable(), ts)
}
