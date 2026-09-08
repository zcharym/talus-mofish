package cloudflare

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/cloudflare/cloudflare-go/v7/option"
)

type apiEnvelope[T any] struct {
	Success bool         `json:"success"`
	Errors  []apiMessage `json:"errors"`
	Result  T            `json:"result"`
}

type apiMessage struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type accountResult struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type scriptResult struct {
	ID         string `json:"id"`
	CreatedOn  string `json:"created_on"`
	ModifiedOn string `json:"modified_on"`
}

type d1Result struct {
	UUID string `json:"uuid"`
	Name string `json:"name"`
}

type kvResult struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

type cronItem struct {
	Cron string `json:"cron"`
}

type scheduleResult struct {
	Schedules []cronItem `json:"schedules"`
}

func TestClientPingSendsBearerAndParsesAccount(t *testing.T) {
	t.Parallel()

	var gotAuth, gotPath string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		gotPath = r.URL.Path
		writeOK(w, accountResult{ID: "acc-1", Name: "Taluship"})
	}))
	t.Cleanup(server.Close)

	client := testClient(server, "acc-1", "test-token")
	got, err := client.Ping(context.Background())
	if err != nil {
		t.Fatalf("Ping: %v", err)
	}
	if gotAuth != "Bearer test-token" {
		t.Fatalf("Authorization = %q", gotAuth)
	}
	if gotPath != "/accounts/acc-1" {
		t.Fatalf("path = %q", gotPath)
	}
	if !got.OK || got.AccountID != "acc-1" || got.AccountName != "Taluship" {
		t.Fatalf("ping = %+v", got)
	}
}

func TestClientPingRequiresCredentials(t *testing.T) {
	t.Parallel()

	client := &Client{}
	_, err := client.Ping(context.Background())
	if err == nil || !strings.Contains(err.Error(), "required") {
		t.Fatalf("error = %v", err)
	}
}

func TestClientMapsJSONError(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusUnauthorized, apiEnvelope[json.RawMessage]{
			Success: false,
			Errors:  []apiMessage{{Code: 10000, Message: "Authentication error"}},
		})
	}))
	t.Cleanup(server.Close)

	client := testClient(server, "acc-1", "bad")
	_, err := client.Ping(context.Background())
	if err == nil || !strings.Contains(err.Error(), "10000") || !strings.Contains(err.Error(), "Authentication error") {
		t.Fatalf("error = %v", err)
	}
}

func TestDashboardNotConfigured(t *testing.T) {
	t.Parallel()

	client := &Client{}
	dash, err := client.Dashboard(context.Background())
	if err != nil {
		t.Fatalf("Dashboard: %v", err)
	}
	if dash.Configured {
		t.Fatal("expected unconfigured dashboard")
	}
	if len(dash.Known) != 2 {
		t.Fatalf("known = %d", len(dash.Known))
	}
}

func TestDashboardAggregatesWorkersAnalyticsAndStorage(t *testing.T) {
	t.Parallel()

	fixed := time.Date(2026, 9, 7, 9, 0, 0, 0, time.UTC)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/accounts/acc-1":
			writeOK(w, accountResult{ID: "acc-1", Name: "Taluship"})
		case r.Method == http.MethodGet && r.URL.Path == "/accounts/acc-1/workers/scripts":
			writeOK(w, []scriptResult{
				{ID: "talus-auth", CreatedOn: "2026-01-01T00:00:00Z", ModifiedOn: "2026-09-01T00:00:00Z"},
				{ID: "echo-watch", CreatedOn: "2026-02-01T00:00:00Z", ModifiedOn: "2026-09-02T00:00:00Z"},
			})
		case r.Method == http.MethodGet && strings.HasSuffix(r.URL.Path, "/schedules"):
			if strings.Contains(r.URL.Path, "echo-watch") {
				writeOK(w, scheduleResult{Schedules: []cronItem{{Cron: "*/5 * * * *"}}})
				return
			}
			writeOK(w, scheduleResult{Schedules: []cronItem{}})
		case r.Method == http.MethodGet && r.URL.Path == "/accounts/acc-1/d1/database":
			writeOK(w, []d1Result{{UUID: "d1-1", Name: "talus-auth"}})
		case r.Method == http.MethodGet && r.URL.Path == "/accounts/acc-1/storage/kv/namespaces":
			writeOK(w, []kvResult{{ID: "kv-1", Title: "echo-watch-kv"}})
		case r.Method == http.MethodPost && r.URL.Path == "/graphql":
			body, _ := io.ReadAll(r.Body)
			if !strings.Contains(string(body), "workersInvocationsAdaptive") {
				t.Errorf("graphql body missing workersInvocationsAdaptive query")
			}
			if strings.Contains(string(body), "workersInvocationsAdaptiveGroups") {
				t.Errorf("graphql body used removed Groups node")
			}
			writeJSON(w, http.StatusOK, map[string]any{
				"data": map[string]any{
					"viewer": map[string]any{
						"accounts": []map[string]any{{
							"workersInvocationsAdaptive": []map[string]any{
								{
									"dimensions": map[string]any{"scriptName": "talus-auth"},
									"sum":        map[string]any{"requests": 1200.0, "errors": 12.0},
									"quantiles":  map[string]any{"cpuTimeP99": 1500.0},
								},
								{
									"dimensions": map[string]any{"scriptName": "echo-watch"},
									"sum":        map[string]any{"requests": 300.0, "errors": 0.0},
									"quantiles":  map[string]any{"cpuTimeP99": 800.0},
								},
							},
						}},
					},
				},
			})
		default:
			t.Errorf("unexpected %s %s", r.Method, r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	t.Cleanup(server.Close)

	client := testClient(server, "acc-1", "token")
	client.now = func() time.Time { return fixed }

	dash, err := client.Dashboard(context.Background())
	if err != nil {
		t.Fatalf("Dashboard: %v", err)
	}
	if !dash.Configured || dash.Account.Name != "Taluship" {
		t.Fatalf("account = %+v", dash.Account)
	}
	if dash.FetchedAt != "2026-09-07T09:00:00Z" {
		t.Fatalf("fetchedAt = %q", dash.FetchedAt)
	}
	if dash.KPIs.WorkerCount != 2 || dash.KPIs.Requests24h != 1500 || dash.KPIs.Errors24h != 12 {
		t.Fatalf("kpis = %+v", dash.KPIs)
	}
	if dash.KPIs.ErrorRate != 12.0/1500.0 {
		t.Fatalf("error rate = %v", dash.KPIs.ErrorRate)
	}
	if len(dash.Workers) != 2 || dash.Workers[0].Name != "echo-watch" {
		t.Fatalf("workers = %+v", dash.Workers)
	}
	watch := dash.Workers[0]
	if watch.Requests24h != 300 || len(watch.Cron) != 1 || watch.Cron[0] != "*/5 * * * *" {
		t.Fatalf("echo-watch = %+v", watch)
	}
	if len(dash.D1) != 1 || dash.D1[0].Name != "talus-auth" {
		t.Fatalf("d1 = %+v", dash.D1)
	}
	if len(dash.KV) != 1 || dash.KV[0].Title != "echo-watch-kv" {
		t.Fatalf("kv = %+v", dash.KV)
	}
	if len(dash.Known) != 2 || !dash.Known[0].Present || dash.Known[0].ID != "talus-auth" {
		t.Fatalf("known = %+v", dash.Known)
	}
	if dash.Errors != (DashboardErrors{}) {
		t.Fatalf("errors = %+v", dash.Errors)
	}
}

func TestDashboardPartialPermissionKeepsWorkers(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/accounts/acc-1":
			writeOK(w, accountResult{ID: "acc-1", Name: "Taluship"})
		case r.URL.Path == "/accounts/acc-1/workers/scripts":
			writeOK(w, []scriptResult{{ID: "talus-auth"}})
		case strings.HasSuffix(r.URL.Path, "/schedules"):
			writeOK(w, scheduleResult{Schedules: []cronItem{}})
		case r.URL.Path == "/accounts/acc-1/d1/database":
			writeJSON(w, http.StatusForbidden, apiEnvelope[json.RawMessage]{
				Success: false,
				Errors:  []apiMessage{{Code: 9109, Message: "Unauthorized to access requested resource"}},
			})
		case r.URL.Path == "/accounts/acc-1/storage/kv/namespaces":
			writeOK(w, []kvResult{})
		case r.URL.Path == "/graphql":
			writeJSON(w, http.StatusForbidden, apiEnvelope[json.RawMessage]{
				Success: false,
				Errors:  []apiMessage{{Message: "Analytics not permitted"}},
			})
		default:
			t.Errorf("unexpected %s %s", r.Method, r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	t.Cleanup(server.Close)

	client := testClient(server, "acc-1", "token")
	dash, err := client.Dashboard(context.Background())
	if err != nil {
		t.Fatalf("Dashboard: %v", err)
	}
	if dash.Errors.Account != "" || dash.Errors.Workers != "" {
		t.Fatalf("unexpected core errors: %+v", dash.Errors)
	}
	if dash.Errors.D1 == "" || !strings.Contains(dash.Errors.D1, "Unauthorized") {
		t.Fatalf("d1 error = %q", dash.Errors.D1)
	}
	if dash.Errors.Analytics == "" {
		t.Fatalf("analytics error empty: %+v", dash.Errors)
	}
	if len(dash.Workers) != 1 || dash.Workers[0].Name != "talus-auth" {
		t.Fatalf("workers = %+v", dash.Workers)
	}
	if len(dash.D1) != 0 {
		t.Fatalf("d1 = %+v", dash.D1)
	}
}

func TestConfigured(t *testing.T) {
	t.Parallel()

	if NewClient("", "token").Configured() {
		t.Fatal("empty account should not be configured")
	}
	if NewClient("acc", "").Configured() {
		t.Fatal("empty token should not be configured")
	}
	if !NewClient(" acc ", " token ").Configured() {
		t.Fatal("trimmed credentials should be configured")
	}
}

func testClient(server *httptest.Server, accountID, token string) *Client {
	return NewClient(accountID, token,
		option.WithBaseURL(server.URL+"/"),
		option.WithHTTPClient(server.Client()),
		option.WithMaxRetries(0),
	)
}

func writeOK(w http.ResponseWriter, result any) {
	writeJSON(w, http.StatusOK, apiEnvelope[any]{
		Success: true,
		Result:  result,
	})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
