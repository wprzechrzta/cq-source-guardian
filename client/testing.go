package client

import (
	"context"
	"github.com/cloudquery/plugin-sdk/v4/plugin"
	"github.com/cloudquery/plugin-sdk/v4/scheduler"
	"github.com/cloudquery/plugin-sdk/v4/schema"
	"github.com/cloudquery/plugin-sdk/v4/transformers"
	"github.com/rs/zerolog"
	news "github.com/wprzechrzta/cq-source-guardian/internal"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

func TestHelper(t *testing.T, table *schema.Table, ts *httptest.Server) {
	table.IgnoreInTests = false
	t.Helper()

	l := zerolog.New(zerolog.NewTestWriter(t)).Output(
		zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.StampMicro},
	).Level(zerolog.DebugLevel).With().Timestamp().Logger()

	sched := scheduler.NewScheduler(scheduler.WithLogger(l))
	spec := &Spec{
		Key: "test-key",
	}
	spec.SetDefaults()
	if err := spec.Validate(); err != nil {
		t.Fatalf("failed to validate spec: %v", err)
	}

	newsClient, err := news.NewClient(
		news.WithAPIKey("dummy-key"),
		news.WithBaseURL(ts.URL),
		news.WithHTTPClient(ts.Client()))
	if err != nil {
		t.Fatal(err)
	}

	c := New(l, *spec, newsClient)
	tables := schema.Tables{table}
	if err := transformers.TransformTables(tables); err != nil {
		t.Fatal(err)
	}
	messages, err := sched.SyncAll(context.Background(), c, tables)
	if err != nil {
		t.Fatalf("failed to sync: %v", err)
	}
	plugin.ValidateNoEmptyColumns(t, tables, messages)

}
