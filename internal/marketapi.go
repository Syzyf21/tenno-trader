package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	marketBaseURL      = "https://api.warframe.market"
	MarketRequestPause = 350 * time.Millisecond
)

type itemsResponse struct {
	Items []MarketItem `json:"data"`
}

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

func FetchMarketItemIndex() (map[string]string, error) {
	body, err := marketGet("/v2/items")
	if err != nil {
		return nil, err
	}

	var parsed itemsResponse
	if err := json.Unmarshal(body, &parsed); err != nil {
		return nil, fmt.Errorf("Error while parsing item list: %w", err)
	}

	index := make(map[string]string, len(parsed.Items))
	for _, it := range parsed.Items {
		index[NormalizeItemName(it.I18n.En.Name)] = it.Slug
	}
	return index, nil
}

func NormalizeItemName(s string) string {
	var b strings.Builder
	for _, r := range strings.ToLower(s) {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			b.WriteRune(r)
		}
	}
	return b.String()
}

func FetchItemStatistics(urlName string) ([]StatEntry, error) {
	body, err := marketGet("/v1/items/" + urlName + "/statistics")
	if err != nil {
		return nil, err
	}

	var parsed statisticsResponse
	if err := json.Unmarshal(body, &parsed); err != nil {
		return nil, fmt.Errorf("Error while parsing statistics for %s: %w", urlName, err)
	}
	return parsed.Payload.StatisticsClosed.Days90, nil
}

func AverageInWindow(entries []StatEntry, window AnalysisWindow, isMax bool) (avgPrice, avgVolume float64, count int) {
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
		if e.ModRank != 0 && !isMax {
			continue
		}
		if e.ModRank == 0 && isMax {
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
