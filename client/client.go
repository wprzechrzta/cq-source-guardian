package client

import (
	"github.com/cloudquery/plugin-sdk/v4/schema"
	news "github.com/wprzechrzta/cq-source-guardian/internal"

	"github.com/rs/zerolog"
)

type Client struct {
	Logger zerolog.Logger
	Spec   Spec
	News   *news.Client
}

func (c *Client) ID() string {
	return "news"
}

func New(logger zerolog.Logger, spec Spec, services *news.Client) schema.ClientMeta {
	return &Client{
		Logger: logger,
		Spec:   spec,
		News:   services,
	}
}
