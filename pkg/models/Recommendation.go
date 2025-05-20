package models

type Recommendation struct {
	Type         string  `json:"type"`
	ResourceType string  `json:"resourceType"`
	Namespace    string  `json:"namespace,omitempty"`
	Message      string  `json:"message"`
	Savings      float64 `json:"savings"`
}
