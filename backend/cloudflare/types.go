package cloudflare

// Account is a Cloudflare account identity from GET /accounts/{id}.
type Account struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Worker is one Workers script plus optional 24h analytics and cron triggers.
type Worker struct {
	Name        string   `json:"name"`
	CreatedOn   string   `json:"createdOn"`
	ModifiedOn  string   `json:"modifiedOn"`
	Requests24h float64  `json:"requests24h"`
	Errors24h   float64  `json:"errors24h"`
	CPUTimeP99  float64  `json:"cpuTimeP99"`
	Cron        []string `json:"cron"`
}

// D1Database is a D1 instance from GET /accounts/{id}/d1/database.
type D1Database struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// KVNamespace is a Workers KV namespace.
type KVNamespace struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

// KnownService highlights a Worker this repo deploys when it exists on the account.
type KnownService struct {
	ID           string  `json:"id"`
	Title        string  `json:"title"`
	Present      bool    `json:"present"`
	BindingsNote string  `json:"bindingsNote"`
	Worker       *Worker `json:"worker,omitempty"`
}

// KPIs are account-level 24h summary figures.
type KPIs struct {
	WorkerCount int     `json:"workerCount"`
	Requests24h float64 `json:"requests24h"`
	Errors24h   float64 `json:"errors24h"`
	ErrorRate   float64 `json:"errorRate"`
}

// DashboardErrors holds scoped failures so one missing permission does not hide the rest.
type DashboardErrors struct {
	Account   string `json:"account"`
	Workers   string `json:"workers"`
	Analytics string `json:"analytics"`
	D1        string `json:"d1"`
	KV        string `json:"kv"`
}

// Dashboard is the read-only Agent snapshot returned by CloudflareService.GetDashboard.
type Dashboard struct {
	Configured bool            `json:"configured"`
	FetchedAt  string          `json:"fetchedAt"`
	Account    Account         `json:"account"`
	KPIs       KPIs            `json:"kpis"`
	Workers    []Worker        `json:"workers"`
	Known      []KnownService  `json:"known"`
	D1         []D1Database    `json:"d1"`
	KV         []KVNamespace   `json:"kv"`
	Errors     DashboardErrors `json:"errors"`
}

// PingResult is the Test connection payload.
type PingResult struct {
	OK          bool   `json:"ok"`
	AccountID   string `json:"accountId"`
	AccountName string `json:"accountName"`
}

type knownWorker struct {
	ID    string
	Title string
	Note  string
}

var knownWorkers = []knownWorker{
	{ID: "talus-auth", Title: "Talus Auth", Note: "D1 database AUTH_DB"},
	{ID: "echo-watch", Title: "Echo Watch", Note: "KV namespace"},
}
