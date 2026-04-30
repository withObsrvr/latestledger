package app

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type Network string

const (
	NetworkMainnet Network = "mainnet"
	NetworkTestnet Network = "testnet"
)

func ParseNetwork(value string) Network {
	switch strings.ToLower(value) {
	case "testnet":
		return NetworkTestnet
	default:
		return NetworkMainnet
	}
}

type StatsResponse struct {
	GeneratedAt   time.Time `json:"generated_at"`
	DataFreshness string    `json:"data_freshness"`
	Ledger        Ledger    `json:"ledger"`

	Transactions24h *Transactions24h `json:"transactions_24h,omitempty"`
	Operations24h   *Operations24h   `json:"operations_24h,omitempty"`
}

type Ledger struct {
	LatestSequence       int64     `json:"latest_sequence"`
	LatestHash           string    `json:"latest_hash"`
	ClosedAt             time.Time `json:"closed_at"`
	ProtocolVersion      int       `json:"protocol_version"`
	AvgCloseTimeSeconds  float64   `json:"avg_close_time_seconds"`
}

type Transactions24h struct {
	Total            int64 `json:"total"`
	Successful       int64 `json:"successful"`
	Failed           int64 `json:"failed"`
	SorobanCount     int64 `json:"soroban_count"`
	TotalFeesCharged int64 `json:"total_fees_charged"`
}

type Operations24h struct {
	Total          int64 `json:"total"`
	SorobanOpCount int64 `json:"soroban_op_count"`
}

type Client struct {
	HTTP    *http.Client
	BaseURL string
	APIKey  string
}

func NewClient(apiKey string) *Client {
	return &Client{
		HTTP: &http.Client{Timeout: 10 * time.Second},
		BaseURL: "https://gateway.withobsrvr.com/lake/v1",
		APIKey: apiKey,
	}
}

func (c *Client) NetworkStats(ctx context.Context, network Network) (StatsResponse, error) {
	var stats StatsResponse
	if c.APIKey == "" {
		return stats, fmt.Errorf("OBSRVR_API_KEY is not configured")
	}

	url := fmt.Sprintf("%s/%s/api/v1/silver/stats/network", strings.TrimRight(c.BaseURL, "/"), network)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return stats, err
	}
	req.Header.Set("Authorization", "Api-Key "+c.APIKey)
	req.Header.Set("Accept", "application/json")

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return stats, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return stats, fmt.Errorf("gateway returned %s", resp.Status)
	}
	if err := json.NewDecoder(resp.Body).Decode(&stats); err != nil {
		return stats, err
	}
	return stats, nil
}
