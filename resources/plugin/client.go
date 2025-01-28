package plugin

import (
	"context"
	"encoding/json"
	"fmt"
	news "github.com/wprzechrzta/cq-source-guardian/internal"

	"github.com/cloudquery/plugin-sdk/v4/message"
	"github.com/cloudquery/plugin-sdk/v4/plugin"
	"github.com/cloudquery/plugin-sdk/v4/scheduler"
	"github.com/cloudquery/plugin-sdk/v4/schema"
	"github.com/cloudquery/plugin-sdk/v4/transformers"
	"github.com/rs/zerolog"
	"github.com/wprzechrzta/cq-source-guardian/client"
	"github.com/wprzechrzta/cq-source-guardian/resources/services"
)

type Client struct {
	logger    zerolog.Logger
	config    client.Spec
	tables    schema.Tables
	scheduler *scheduler.Scheduler
	services  *news.Client

	plugin.UnimplementedDestination
}

func Configure(_ context.Context, logger zerolog.Logger, spec []byte, opts plugin.NewClientOptions) (plugin.Client, error) {
	if opts.NoConnection {
		return &Client{
			logger: logger.With().Str("module", "news").Logger(),

			tables: getTables(),
		}, nil
	}

	config := &client.Spec{}
	logger.Info().Msg("loading plugin configuration")
	if err := json.Unmarshal(spec, config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal spec: %w", err)
	}
	config.SetDefaults()
	logger.Info().Msgf("---config: %+v", config)
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("failed to validate spec: %w", err)
	}

	newsClient, err := news.NewClient(news.WithAPIKey(config.Key))
	if err != nil {
		return nil, err
	}

	return &Client{
		logger:   logger.With().Str("module", "news").Logger(),
		config:   *config,
		tables:   getTables(),
		services: newsClient,
		scheduler: scheduler.NewScheduler(
			scheduler.WithLogger(logger),
			scheduler.WithConcurrency(3)),
	}, nil
}

func (c *Client) Sync(ctx context.Context, options plugin.SyncOptions, res chan<- message.SyncMessage) error {
	tt, err := c.tables.FilterDfs(options.Tables, options.SkipTables, options.SkipDependentTables)
	if err != nil {
		return err
	}

	schedulerClient := client.New(c.logger, c.config, c.services)

	eff := c.scheduler.Sync(ctx, schedulerClient, tt, res, scheduler.WithSyncDeterministicCQID(options.DeterministicCQID))
	if eff != nil {
		return fmt.Errorf("failed to sync: %w", eff)
	}
	return nil
}

func (c *Client) Tables(_ context.Context, options plugin.TableOptions) (schema.Tables, error) {
	tt, err := c.tables.FilterDfs(options.Tables, options.SkipTables, options.SkipDependentTables)
	if err != nil {
		return nil, err
	}

	return tt, nil
}

func (*Client) Close(_ context.Context) error {
	return nil
}

func getTables() schema.Tables {
	tables := schema.Tables{
		services.NewsTable(),
	}
	if err := transformers.TransformTables(tables); err != nil {
		panic(err)
	}
	for _, t := range tables {
		schema.AddCqIDs(t)
	}
	return tables
}
