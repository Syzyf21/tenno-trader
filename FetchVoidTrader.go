package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const voidTraderEndpoint = "https://api.warframestat.us/pc/voidTrader"

func fetchBaroInventory() (*VoidTraderResponse, error) {
	client := &http.Client{Timeout: 10 * time.Second}

	req, err := http.NewRequest(http.MethodGet, voidTraderEndpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("Error while building request: %w", err)
	}
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Error while fetching baroo inventory: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Error while reading response: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Warframestat API returned status %d: %s", resp.StatusCode, string(body))
	}

	var out VoidTraderResponse
	if err := json.Unmarshal(body, &out); err != nil {
		return nil, fmt.Errorf("Error while decoding response: %w", err)
	}
	return &out, nil
}

func (v *VoidTraderResponse) activationDate() (time.Time, error) {
	t, err := time.Parse(time.RFC3339, v.Activation)
	if err != nil {
		return time.Time{}, fmt.Errorf("Error while parsing activation date %q: %w", v.Activation, err)
	}
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC), nil
}
