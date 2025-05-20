package metrics

type Clients struct {
	Prometheus    PrometheusClient
	MetricsServer MetricsServerClient
	Kubecost      KubecostClient
}

type Option func(*Clients)

func NewClients(opts ...Option) *Clients {
	c := &Clients{}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

func WithPrometheus(url string) Option {
	return func(c *Clients) {
		c.Prometheus = *NewPrometheusClient(url)
	}
}

func WithMetricsServer(url string) Option {
	return func(c *Clients) {
		c.MetricsServer = NewMetricsServerClient(url)
	}
}

func WithKubecost(url, apiKey string) Option {
	return func(c *Clients) {
		c.Kubecost = NewKubecostClient(url, apiKey)
	}
}
