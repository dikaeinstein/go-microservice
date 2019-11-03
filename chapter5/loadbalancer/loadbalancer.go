package loadbalancer

import (
	"net/url"
)

// LoadBalancer returns endpoints for downstream calls
type LoadBalancer struct {
	strategy Strategy
}

// New creates a new loadbalancer and setting the given strategy
func New(strategy Strategy, endPoints []url.URL) *LoadBalancer {
	strategy.SetEndpoints(endPoints)
	return &LoadBalancer{strategy}
}

// GetEndpoint gets an endpoint based on the given strategy
func (l *LoadBalancer) GetEndpoint() url.URL {
	return l.strategy.NextEndpoint()
}

// UpdateEndpoints updates the endpoints available to the strategy
func (l *LoadBalancer) UpdateEndpoints(urls []url.URL) {
	l.strategy.SetEndpoints(urls)
}
