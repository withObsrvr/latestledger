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
	Network     string     `json:"network"`
	GeneratedAt time.Time  `json:"generated_at"`
	Header      Header     `json:"header"`
	Hero        Hero       `json:"hero"`
	Meta        Meta       `json:"meta"`
	Provenance  Provenance `json:"provenance"`
}

type Header struct {
	LatestLedgerSequence int64     `json:"latest_ledger_sequence"`
	LatestLedgerClosedAt time.Time `json:"latest_ledger_closed_at"`
}

type Hero struct {
	Health      Health       `json:"health"`
	LatestLedger LatestLedger `json:"latest_ledger"`
	Cadence     Cadence      `json:"cadence"`
	Contracts   Contracts    `json:"contracts"`
	Soroban     Soroban      `json:"soroban"`
	Trends      Trends       `json:"trends"`
	TTL         TTL          `json:"ttl"`
	ActivityMix ActivityMix  `json:"activity_mix"`
}

type Health struct {
	Status       string `json:"status"`
	LoadBand     string `json:"load_band"`
	ActivityBand string `json:"activity_band"`
}

type LatestLedger struct {
	Sequence         int64     `json:"sequence"`
	ClosedAt         time.Time `json:"closed_at"`
	TransactionCount int64     `json:"transaction_count"`
	OperationCount   int64     `json:"operation_count"`
}

type Cadence struct {
	AvgCloseSeconds       float64 `json:"avg_close_seconds"`
	TxPerLedgerRecentAvg  float64 `json:"tx_per_ledger_recent_avg"`
	OpsPerLedgerRecentAvg float64 `json:"ops_per_ledger_recent_avg"`
}

type Contracts struct {
	Active24h int64 `json:"active_24h"`
}

type Soroban struct {
	InstructionPct float64 `json:"instruction_pct"`
	ReadWritePct   float64 `json:"read_write_pct"`
}

type Trends struct {
	TxVs24hAvgPct   float64 `json:"tx_vs_24h_avg_pct"`
	AnomalyDetected bool    `json:"anomaly_detected"`
}

type TTL struct {
	ExpiringContractCount int64 `json:"expiring_contract_count"`
	WorstRemainingHours   int64 `json:"worst_remaining_hours"`
	WorstRemainingLedgers int64 `json:"worst_remaining_ledgers"`
}

type ActivityMix struct {
	SwapTx24h         int64 `json:"swap_tx_24h"`
	ContractCallTx24h int64 `json:"contract_call_tx_24h"`
}

type Meta struct {
	LatestLedgerAgeSeconds int64 `json:"latest_ledger_age_seconds"`
}

type Provenance struct {
	Route      string   `json:"route"`
	DataSource string   `json:"data_source"`
	Partial    bool     `json:"partial"`
	Warnings   []string `json:"warnings"`
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

	url := fmt.Sprintf("%s/%s/api/v1/home/summary?_=%d", strings.TrimRight(c.BaseURL, "/"), network, time.Now().UnixNano())
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return stats, err
	}
	req.Header.Set("Authorization", "Api-Key "+c.APIKey)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Pragma", "no-cache")

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
