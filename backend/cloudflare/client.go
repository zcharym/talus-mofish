package cloudflare

import (
	"context"
	"fmt"
	"strings"
	"time"

	cf "github.com/cloudflare/cloudflare-go/v7"
	"github.com/cloudflare/cloudflare-go/v7/accounts"
	"github.com/cloudflare/cloudflare-go/v7/d1"
	"github.com/cloudflare/cloudflare-go/v7/kv"
	"github.com/cloudflare/cloudflare-go/v7/option"
	"github.com/cloudflare/cloudflare-go/v7/workers"
)

const listPageSize = 100

// Client talks to the Cloudflare API via cloudflare-go/v7.
type Client struct {
	AccountID string
	APIToken  string
	api       *cf.Client
	now       func() time.Time
}

// NewClient builds a Cloudflare API client for the given account.
// Extra request options (base URL, HTTP client, retries) are for tests.
func NewClient(accountID, apiToken string, opts ...option.RequestOption) *Client {
	accountID = strings.TrimSpace(accountID)
	apiToken = strings.TrimSpace(apiToken)
	c := &Client{
		AccountID: accountID,
		APIToken:  apiToken,
		now:       time.Now,
	}
	if apiToken == "" {
		return c
	}
	sdkOpts := append([]option.RequestOption{option.WithAPIToken(apiToken)}, opts...)
	c.api = cf.NewClient(sdkOpts...)
	return c
}

// Configured reports whether account ID and API token are both set.
func (c *Client) Configured() bool {
	if c == nil {
		return false
	}
	return strings.TrimSpace(c.AccountID) != "" && strings.TrimSpace(c.APIToken) != ""
}

func (c *Client) accountID() string {
	if c == nil {
		return ""
	}
	return strings.TrimSpace(c.AccountID)
}

func (c *Client) clock() time.Time {
	if c != nil && c.now != nil {
		return c.now()
	}
	return time.Now()
}

func (c *Client) sdk() (*cf.Client, error) {
	if c == nil || c.api == nil {
		return nil, fmt.Errorf("cloudflare client is not configured")
	}
	return c.api, nil
}

func formatTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.UTC().Format(time.RFC3339)
}

// Ping verifies the token can read the configured account.
func (c *Client) Ping(ctx context.Context) (PingResult, error) {
	if c == nil {
		return PingResult{}, fmt.Errorf("cloudflare client is nil")
	}
	if !c.Configured() {
		return PingResult{}, fmt.Errorf("cloudflare account ID and API token are required")
	}
	account, err := c.GetAccount(ctx)
	if err != nil {
		return PingResult{}, err
	}
	return PingResult{
		OK:          true,
		AccountID:   account.ID,
		AccountName: account.Name,
	}, nil
}

// GetAccount loads the configured account's id and name.
func (c *Client) GetAccount(ctx context.Context) (Account, error) {
	api, err := c.sdk()
	if err != nil {
		return Account{}, err
	}
	res, err := api.Accounts.Get(ctx, accounts.AccountGetParams{
		AccountID: cf.F(c.accountID()),
	})
	if err != nil {
		return Account{}, err
	}
	id := strings.TrimSpace(res.ID)
	if id == "" {
		id = c.accountID()
	}
	return Account{ID: id, Name: strings.TrimSpace(res.Name)}, nil
}

// ListWorkers returns Workers scripts on the account.
func (c *Client) ListWorkers(ctx context.Context) ([]Worker, error) {
	api, err := c.sdk()
	if err != nil {
		return nil, err
	}
	page, err := api.Workers.Scripts.List(ctx, workers.ScriptListParams{
		AccountID: cf.F(c.accountID()),
	})
	if err != nil {
		return nil, err
	}
	out := make([]Worker, 0, len(page.Result))
	for _, item := range page.Result {
		name := strings.TrimSpace(item.ID)
		if name == "" {
			continue
		}
		out = append(out, Worker{
			Name:       name,
			CreatedOn:  formatTime(item.CreatedOn),
			ModifiedOn: formatTime(item.ModifiedOn),
			Cron:       []string{},
		})
	}
	return out, nil
}

// ListD1 returns D1 databases on the account.
func (c *Client) ListD1(ctx context.Context) ([]D1Database, error) {
	api, err := c.sdk()
	if err != nil {
		return nil, err
	}
	page, err := api.D1.Database.List(ctx, d1.DatabaseListParams{
		AccountID: cf.F(c.accountID()),
		PerPage:   cf.F(float64(listPageSize)),
	})
	if err != nil {
		return nil, err
	}
	items := make([]D1Database, 0, len(page.Result))
	for _, item := range page.Result {
		id := strings.TrimSpace(item.UUID)
		name := strings.TrimSpace(item.Name)
		if id == "" && name == "" {
			continue
		}
		items = append(items, D1Database{ID: id, Name: name})
	}
	return items, nil
}

// ListKV returns Workers KV namespaces on the account.
func (c *Client) ListKV(ctx context.Context) ([]KVNamespace, error) {
	api, err := c.sdk()
	if err != nil {
		return nil, err
	}
	page, err := api.KV.Namespaces.List(ctx, kv.NamespaceListParams{
		AccountID: cf.F(c.accountID()),
		PerPage:   cf.F(float64(listPageSize)),
	})
	if err != nil {
		return nil, err
	}
	items := make([]KVNamespace, 0, len(page.Result))
	for _, item := range page.Result {
		id := strings.TrimSpace(item.ID)
		title := strings.TrimSpace(item.Title)
		if id == "" && title == "" {
			continue
		}
		items = append(items, KVNamespace{ID: id, Title: title})
	}
	return items, nil
}

// ListSchedules returns cron expressions attached to a Worker script.
func (c *Client) ListSchedules(ctx context.Context, scriptName string) ([]string, error) {
	scriptName = strings.TrimSpace(scriptName)
	if scriptName == "" {
		return nil, fmt.Errorf("script name is required")
	}
	api, err := c.sdk()
	if err != nil {
		return nil, err
	}
	res, err := api.Workers.Scripts.Schedules.Get(ctx, scriptName, workers.ScriptScheduleGetParams{
		AccountID: cf.F(c.accountID()),
	})
	if err != nil {
		return nil, err
	}
	out := make([]string, 0, len(res.Schedules))
	for _, item := range res.Schedules {
		cron := strings.TrimSpace(item.Cron)
		if cron == "" {
			continue
		}
		out = append(out, cron)
	}
	return out, nil
}

type analyticsRow struct {
	ScriptName  string
	Requests24h float64
	Errors24h   float64
	CPUTimeP99  float64
}

type graphqlResponse struct {
	Data struct {
		Viewer struct {
			Accounts []struct {
				Rows []graphqlGroup `json:"workersInvocationsAdaptive"`
			} `json:"accounts"`
		} `json:"viewer"`
	} `json:"data"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors"`
}

type graphqlGroup struct {
	Dimensions struct {
		ScriptName string `json:"scriptName"`
	} `json:"dimensions"`
	Sum struct {
		Requests float64 `json:"requests"`
		Errors   float64 `json:"errors"`
	} `json:"sum"`
	Quantiles struct {
		CPUTimeP99 float64 `json:"cpuTimeP99"`
	} `json:"quantiles"`
}

const workersAnalyticsQuery = `
query WorkerAnalytics($accountTag: string!, $datetimeStart: Time!, $datetimeEnd: Time!) {
  viewer {
    accounts(filter: {accountTag: $accountTag}) {
      workersInvocationsAdaptive(
        limit: 10000
        filter: { datetime_geq: $datetimeStart, datetime_leq: $datetimeEnd }
        orderBy: [sum_requests_DESC]
      ) {
        dimensions { scriptName }
        sum { requests errors }
        quantiles { cpuTimeP50 cpuTimeP99 }
      }
    }
  }
}
`

// WorkerAnalytics returns 24h invocation totals grouped by script name.
// Workers GraphQL has no *Groups node — use workersInvocationsAdaptive.
func (c *Client) WorkerAnalytics(ctx context.Context) (map[string]analyticsRow, error) {
	api, err := c.sdk()
	if err != nil {
		return nil, err
	}
	end := c.clock().UTC()
	start := end.Add(-24 * time.Hour)

	var parsed graphqlResponse
	err = api.Post(ctx, "graphql", map[string]any{
		"query": workersAnalyticsQuery,
		"variables": map[string]any{
			"accountTag":    c.accountID(),
			"datetimeStart": start.Format(time.RFC3339),
			"datetimeEnd":   end.Format(time.RFC3339),
		},
	}, &parsed)
	if err != nil {
		return nil, err
	}
	if len(parsed.Errors) > 0 && strings.TrimSpace(parsed.Errors[0].Message) != "" {
		return nil, fmt.Errorf("cloudflare graphql: %s", parsed.Errors[0].Message)
	}

	rows := make(map[string]analyticsRow)
	if len(parsed.Data.Viewer.Accounts) == 0 {
		return rows, nil
	}
	for _, group := range parsed.Data.Viewer.Accounts[0].Rows {
		name := strings.TrimSpace(group.Dimensions.ScriptName)
		if name == "" {
			continue
		}
		existing := rows[name]
		existing.ScriptName = name
		existing.Requests24h += group.Sum.Requests
		existing.Errors24h += group.Sum.Errors
		if group.Quantiles.CPUTimeP99 > existing.CPUTimeP99 {
			existing.CPUTimeP99 = group.Quantiles.CPUTimeP99
		}
		rows[name] = existing
	}
	return rows, nil
}
