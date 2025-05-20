package models

type CostData struct {
	Total     float64         `json:"total"`
	Breakdown CostBreakdown   `json:"breakdown"`
	Trends    []CostTrendItem `json:"trends"`
}

type CostBreakdown struct {
	ByNamespace map[string]float64 `json:"byNamespace"`
	ByNode      map[string]float64 `json:byNode`
}

type CostTrendItem struct {
	Date string  `json:"date"`
	Cost float64 `json:"cost"`
}
