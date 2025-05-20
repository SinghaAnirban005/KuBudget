package models

type DashboardRequest struct {
	Type         string `json:"type"`
	ExportToFile bool   `json:"exportToFile"`
}

type DashboardPanel struct {
	Title   string        `json:"title"`
	Type    string        `json:"type"`
	Targets []PanelTarget `json:"targets"`
}

type PanelTarget struct {
	Expr   string `json:"expr"`
	Legend string `json:"legend"`
}

type NodeUtilizationDashboard struct {
	Title  string           `json:"title"`
	Panels []DashboardPanel `json:"panels"`
}

type CostBreakdownDashboard struct {
	Title  string           `json:"title"`
	Panels []DashboardPanel `json:"panels"`
}
