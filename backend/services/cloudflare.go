package services

import (
	"context"
	"fmt"
	"time"

	"github.com/songwei.ma/talus-mofish/backend/cloudflare"
	"github.com/songwei.ma/talus-mofish/backend/storage"
	"github.com/songwei.ma/talus-mofish/backend/types"
)

const dashboardTimeout = 25 * time.Second

// CloudflareService exposes a read-only Cloudflare operations snapshot to the Agent window.
type CloudflareService struct {
	config *storage.ConfigStore
	client *cloudflare.Client
}

// NewCloudflareService creates the Cloudflare Wails service.
func NewCloudflareService(cfg *storage.ConfigStore) *CloudflareService {
	return &CloudflareService{config: cfg}
}

func (s *CloudflareService) apiClient() *cloudflare.Client {
	if s.client != nil {
		return s.client
	}
	cfg := s.config.Get().Cloudflare
	return cloudflare.NewClient(cfg.AccountID, cfg.APIToken)
}

// IsConfigured reports whether an account ID and API token are saved.
func (s *CloudflareService) IsConfigured() bool {
	return types.CloudflareConfigured(s.config.Get().Cloudflare)
}

// Ping verifies the saved token can read the configured account.
func (s *CloudflareService) Ping() (cloudflare.PingResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dashboardTimeout)
	defer cancel()
	result, err := s.apiClient().Ping(ctx)
	if err != nil {
		return cloudflare.PingResult{}, err
	}
	return result, nil
}

// GetDashboard returns a single read-only snapshot of Workers, analytics, D1, and KV.
func (s *CloudflareService) GetDashboard() (cloudflare.Dashboard, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dashboardTimeout)
	defer cancel()
	dash, err := s.apiClient().Dashboard(ctx)
	if err != nil {
		return cloudflare.Dashboard{}, fmt.Errorf("cloudflare dashboard: %w", err)
	}
	if dash.Workers == nil {
		dash.Workers = []cloudflare.Worker{}
	}
	if dash.Known == nil {
		dash.Known = []cloudflare.KnownService{}
	}
	if dash.D1 == nil {
		dash.D1 = []cloudflare.D1Database{}
	}
	if dash.KV == nil {
		dash.KV = []cloudflare.KVNamespace{}
	}
	return dash, nil
}
