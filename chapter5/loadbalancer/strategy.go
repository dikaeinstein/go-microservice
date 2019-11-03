package loadbalancer

import (
	"math/rand"
	"net/url"
	"time"
)

// Strategy is an interface to be implemented by loadbalancing
// strategies like round robin or random.
type Strategy interface {
	NextEndpoint() url.URL
	SetEndpoints([]url.URL)
}

// RandomStrategy implements Strategy for random endpoint selection
type RandomStrategy struct {
	endpoints []url.URL
}

// SetEndpoints sets the available endpoints for use by the strategy
func (rs *RandomStrategy) SetEndpoints(endpoints []url.URL) {
	rs.endpoints = endpoints
}

// NextEndpoint returns an endpoint using the random strategy
func (rs *RandomStrategy) NextEndpoint() url.URL {
	src := rand.NewSource(time.Now().UnixNano())
	r := rand.New(src)

	return rs.endpoints[r.Intn(len(rs.endpoints))]
}

// RoundRobinStrategy implements Strategy for round robin endpoint selection
type RoundRobinStrategy struct {
	endpoints []url.URL
	next      int
}

// SetEndpoints sets the available endpoints for use by the strategy
func (rrs *RoundRobinStrategy) SetEndpoints(endpoints []url.URL) {
	rrs.endpoints = endpoints
}

// NextEndpoint returns an endpoint using the round robin strategy
func (rrs *RoundRobinStrategy) NextEndpoint() url.URL {
	endpoint := rrs.endpoints[rrs.next]
	rrs.next = (rrs.next + 1) % len(rrs.endpoints)

	return endpoint
}
