package cloudflare

import (
	"context"
	"sort"
	"sync"
	"time"
)

// Dashboard gathers a read-only operations snapshot. Missing permissions
// are recorded on Dashboard.Errors instead of failing the whole call.
func (c *Client) Dashboard(ctx context.Context) (Dashboard, error) {
	if c == nil {
		return Dashboard{}, nil
	}
	if !c.Configured() {
		return Dashboard{
			Configured: false,
			Workers:    []Worker{},
			Known:      knownSnapshot(nil),
			D1:         []D1Database{},
			KV:         []KVNamespace{},
		}, nil
	}

	dash := Dashboard{
		Configured: true,
		FetchedAt:  c.clock().UTC().Format(time.RFC3339),
		Workers:    []Worker{},
		D1:         []D1Database{},
		KV:         []KVNamespace{},
	}

	var (
		account    Account
		accountErr error
		workers    []Worker
		workersErr error
		analytics  map[string]analyticsRow
		analErr    error
		d1         []D1Database
		d1Err      error
		kv         []KVNamespace
		kvErr      error
	)

	var wg sync.WaitGroup
	wg.Add(5)
	go func() {
		defer wg.Done()
		account, accountErr = c.GetAccount(ctx)
	}()
	go func() {
		defer wg.Done()
		workers, workersErr = c.ListWorkers(ctx)
	}()
	go func() {
		defer wg.Done()
		analytics, analErr = c.WorkerAnalytics(ctx)
	}()
	go func() {
		defer wg.Done()
		d1, d1Err = c.ListD1(ctx)
	}()
	go func() {
		defer wg.Done()
		kv, kvErr = c.ListKV(ctx)
	}()
	wg.Wait()

	if accountErr != nil {
		dash.Errors.Account = accountErr.Error()
		dash.Account = Account{ID: c.accountID()}
	} else {
		dash.Account = account
	}

	if workersErr != nil {
		dash.Errors.Workers = workersErr.Error()
		workers = []Worker{}
	}
	if analErr != nil {
		dash.Errors.Analytics = analErr.Error()
		analytics = map[string]analyticsRow{}
	}
	if d1Err != nil {
		dash.Errors.D1 = d1Err.Error()
		d1 = []D1Database{}
	} else if d1 == nil {
		d1 = []D1Database{}
	}
	if kvErr != nil {
		dash.Errors.KV = kvErr.Error()
		kv = []KVNamespace{}
	} else if kv == nil {
		kv = []KVNamespace{}
	}

	workers = mergeAnalytics(workers, analytics)
	c.attachSchedules(ctx, workers)

	sort.Slice(workers, func(i, j int) bool {
		return workers[i].Name < workers[j].Name
	})

	var requests, errors float64
	for _, worker := range workers {
		requests += worker.Requests24h
		errors += worker.Errors24h
	}
	dash.Workers = workers
	dash.D1 = d1
	dash.KV = kv
	dash.KPIs = KPIs{
		WorkerCount: len(workers),
		Requests24h: requests,
		Errors24h:   errors,
		ErrorRate:   errorRate(requests, errors),
	}
	dash.Known = knownSnapshot(workers)
	return dash, nil
}

func mergeAnalytics(workers []Worker, analytics map[string]analyticsRow) []Worker {
	if analytics == nil {
		analytics = map[string]analyticsRow{}
	}
	seen := make(map[string]int, len(workers))
	for i := range workers {
		if workers[i].Cron == nil {
			workers[i].Cron = []string{}
		}
		seen[workers[i].Name] = i
		row, ok := analytics[workers[i].Name]
		if !ok {
			continue
		}
		workers[i].Requests24h = row.Requests24h
		workers[i].Errors24h = row.Errors24h
		workers[i].CPUTimeP99 = row.CPUTimeP99
	}
	for name, row := range analytics {
		if _, ok := seen[name]; ok {
			continue
		}
		workers = append(workers, Worker{
			Name:        name,
			Requests24h: row.Requests24h,
			Errors24h:   row.Errors24h,
			CPUTimeP99:  row.CPUTimeP99,
			Cron:        []string{},
		})
	}
	return workers
}

func (c *Client) attachSchedules(ctx context.Context, workers []Worker) {
	if len(workers) == 0 {
		return
	}
	var wg sync.WaitGroup
	for i := range workers {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			crons, err := c.ListSchedules(ctx, workers[index].Name)
			if err != nil || crons == nil {
				workers[index].Cron = []string{}
				return
			}
			workers[index].Cron = crons
		}(i)
	}
	wg.Wait()
}

func knownSnapshot(workers []Worker) []KnownService {
	byName := make(map[string]Worker, len(workers))
	for _, worker := range workers {
		byName[worker.Name] = worker
	}
	out := make([]KnownService, 0, len(knownWorkers))
	for _, known := range knownWorkers {
		item := KnownService{
			ID:           known.ID,
			Title:        known.Title,
			BindingsNote: known.Note,
		}
		if worker, ok := byName[known.ID]; ok {
			copy := worker
			item.Present = true
			item.Worker = &copy
		}
		out = append(out, item)
	}
	return out
}

func errorRate(requests, errors float64) float64 {
	if requests <= 0 {
		return 0
	}
	return errors / requests
}
