package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	marketBaseURL = "https://api.warframe.market"
	// warframe.market asks integrators to stay at or under ~3 requests/sec.
	marketRequestPause = 350 * time.Millisecond
)

// itemsResponse mirrors GET /v2/items
type itemsResponse struct {
	Payload struct {
		Items []MarketItem `json:"data"`
	}
}

// statisticsResponse mirrors GET /v1/items/{url_name}/statistics
// Note: this is the legacy v1 endpoint. It is marked deprecated by
// warframe.market but remains functional and is the only endpoint that
// currently exposes historical daily price/volume statistics.
type statisticsResponse struct {
	Payload struct {
		StatisticsClosed struct {
			Days90  []StatEntry `json:"90days"`
			Hours48 []StatEntry `json:"48hours"`
		} `json:"statistics_closed"`
	} `json:"payload"`
}

var marketHTTPClient = &http.Client{Timeout: 20 * time.Second}

func marketGet(path string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, marketBaseURL+path, nil)
	if err != nil {
		return nil, fmt.Errorf("Error while building request for %s: %w", path, err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Platform", "pc")
	req.Header.Set("Language", "en")
	req.Header.Set("User-Agent", "Tenno-Trader/1.0 (+desktop app)")

	resp, err := marketHTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Error while making request for %s: %w", path, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Error while reading response for %s: %w", path, err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Error:%s returned status %d: %s", path, resp.StatusCode, string(body))
	}
	return body, nil
}

func fetchMarketItemIndex() (map[string]string, error) {
	body, err := marketGet("/v2/items")
	if err != nil {
		return nil, err
	}
	fmt.Printf("Body %v", body)

	var parsed itemsResponse
	if err := json.Unmarshal(body, &parsed); err != nil {
		return nil, fmt.Errorf("decode item list: %w", err)
	}

	index := make(map[string]string, len(parsed.Payload.Items))
	for _, it := range parsed.Payload.Items {
		index[normalizeItemName(it.I18n.En.Name)] = it.Slug
	}
	return index, nil
}

// normalizeItemName lowercases a name and strips everything but letters and
// digits, so "Akbolto Prime Blueprint" and "akbolto_prime_blueprint" style
// variants both hash to the same key.
func normalizeItemName(s string) string {
	var b strings.Builder
	for _, r := range strings.ToLower(s) {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			b.WriteRune(r)
		}
	}
	return b.String()
}

// fetchItemStatistics retrieves historical daily statistics for a single
// item, identified by its warframe.market url_name (e.g. "akbolto_prime_set").
func fetchItemStatistics(urlName string) ([]StatEntry, error) {
	body, err := marketGet("/v1/items/" + urlName + "/statistics")
	if err != nil {
		return nil, err
	}

	var parsed statisticsResponse
	if err := json.Unmarshal(body, &parsed); err != nil {
		return nil, fmt.Errorf("decode statistics for %s: %w", urlName, err)
	}
	return parsed.Payload.StatisticsClosed.Days90, nil
}

// averageInWindow computes the average avg_price and average volume across
// all daily entries whose date falls within [window.Start, window.End]
// (inclusive, compared by calendar day in UTC).
func averageInWindow(entries []StatEntry, window AnalysisWindow) (avgPrice, avgVolume float64, count int) {
	var sumPrice, sumVolume float64
	for _, e := range entries {
		t, err := time.Parse(time.RFC3339, e.Datetime)
		if err != nil {
			continue
		}
		day := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
		if day.Before(window.Start) || day.After(window.End) {
			continue
		}
		sumPrice += e.AvgPrice
		sumVolume += float64(e.Volume)
		count++
	}
	if count == 0 {
		return 0, 0, 0
	}
	return sumPrice / float64(count), sumVolume / float64(count), count
}
