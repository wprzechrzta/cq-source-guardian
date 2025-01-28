package services

import (
	"context"
	"github.com/cloudquery/plugin-sdk/v4/schema"
	"github.com/cloudquery/plugin-sdk/v4/transformers"
	"github.com/wprzechrzta/cq-source-guardian/client"
	news "github.com/wprzechrzta/cq-source-guardian/internal"
)

const maxPages = 3

func NewsTable() *schema.Table {
	return &schema.Table{
		Name:     "guardian_news",
		Resolver: fetchNews,
		Transform: transformers.TransformWithStruct(&news.NewsItem{},
			transformers.WithPrimaryKeys("Id"),
		),
	}
}

func fetchNews(_ context.Context, meta schema.ClientMeta, parent *schema.Resource, res chan<- any) error {
	cl := meta.(*client.Client)

	term := "deep seek r1"
	resp, err := cl.News.Search(term)
	if err != nil {
		return err
	}

	total := 0
	for i := resp.Response.StartIndex; i < resp.Response.StartIndex+maxPages; i++ {
		itemPage, err := cl.News.FetchPage(term, i)
		if err != nil {
			cl.Logger.Error().Err(err).Msgf("failed to fetch news page %d", i)
		}
		total += len(itemPage)
		for _, item := range itemPage {
			res <- item
		}
	}
	cl.Logger.Printf("Fetched %d new items", total)
	return nil
}
