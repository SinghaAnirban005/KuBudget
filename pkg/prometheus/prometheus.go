package pkg

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/SinghaAnirban005/KuBudget/internal"
)

type Client struct {
	baseURL string
	client  *http.Client
}

type QueryResult struct {
	Status string `json:"status"`
	Data   struct {
		ResultType string                `json:"resultType"`
		Result     []internal.SamplePair `json:"result"`
	} `json:"data"`
}

type RangeQueryResult struct {
	Status string `json:"status"`
	Data   struct {
		ResultType string        `json:"resultType"`
		Result     []MatrixValue `json:"result"`
	} `json:"data"`
}

type MatrixValue struct {
	Metric internal.MetricPoint  `json:"metric"`
	Values []internal.SamplePair `json:"values"`
}

func NewClient(baseURL string) *Client {
	return &Client{
		baseURL: baseURL,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *Client) Query(ctx context.Context, query string) (*QueryResult, error) {
	u, err := url.Parse(c.baseURL + "/api/v1/query")
	if err != nil {
		return nil, err
	}

	params := url.Values{}
	params.Add("query", query)
	u.RawQuery = params.Encode()

	req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("prometheus query failed with status %d", resp.StatusCode)
	}

	var result QueryResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *Client) QueryRange(ctx context.Context, query string, start, end time.Time, step time.Duration) (*RangeQueryResult, error) {
	u, err := url.Parse(c.baseURL + "/api/v1/query_range")
	if err != nil {
		return nil, err
	}

	params := url.Values{}
	params.Add("query", query)
	params.Add("start", strconv.FormatInt(start.Unix(), 10))
	params.Add("end", strconv.FormatInt(end.Unix(), 10))
	params.Add("step", strconv.Itoa(int(step.Seconds()))+"s")

	u.RawQuery = params.Encode()

	req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(req)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("prometheus range query failed with status %d", resp.StatusCode)
	}

	var result RangeQueryResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, err
}

func (c *Client) GetCPUUsage(ctx context.Context, namespace string, pod string) (float64, error) {
	query := fmt.Sprintf(`rate(container_cpu_usage_seconds_total{namespace="%s",pod="%s"}[5m])`, namespace, pod)

	result, err := c.Query(ctx, query)
	if err != nil {
		return 0, err
	}

	if len(result.Data.Result) == 0 {
		return 0, nil
	}
	return float64(result.Data.Result[0].Value), nil
}

func (c *Client) GetMemoryUsage(ctx context.Context, namespace, pod string) (float64, error) {
	query := fmt.Sprintf(`container_memory_usage_bytes{namespace="%s",pod="%s"}`, namespace, pod)
	result, err := c.Query(ctx, query)
	if err != nil {
		return 0, nil
	}

	if len(result.Data.Result) == 0 {
		return 0, nil
	}

	return float64(result.Data.Result[0].Value), nil
}

func (c *Client) GetNetworkIO(ctx context.Context, namespace, pod string) (float64, float64, error) {
	rxQuery := fmt.Sprintf(`rate(container_network_receive_bytes_total{namespace="%s",pod="%s"}[5m])`, namespace, pod)
	rxResult, err := c.Query(ctx, rxQuery)
	if err != nil {
		return 0, 0, err
	}

	txQuery := fmt.Sprintf(`rate(container_network_transmit_bytes_total{namespace="%s",pod="%s"}[5m])`, namespace, pod)
	txResult, err := c.Query(ctx, txQuery)
	if err != nil {
		return 0, 0, err
	}

	var rxBytes, txBytes float64
	if len(rxResult.Data.Result) > 0 {
		rxBytes = float64(rxResult.Data.Result[0].Value)
	}
	if len(txResult.Data.Result) > 0 {
		txBytes = float64(txResult.Data.Result[0].Value)
	}

	return rxBytes, txBytes, nil
}
