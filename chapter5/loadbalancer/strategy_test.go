package loadbalancer

import (
	"net/url"
	"testing"
)

func TestRoundRobinStrategy(t *testing.T) {
	endpoints := []url.URL{
		url.URL{Host: "www.google.com"},
		url.URL{Host: "www.google.co.uk"},
		url.URL{Host: "https://github.com"},
		url.URL{Host: "https://gitlab.com"},
	}

	rrs := &RoundRobinStrategy{}
	rrs.SetEndpoints(endpoints)

	for _, endpoint := range endpoints {
		if endpoint != rrs.NextEndpoint() {
			t.Errorf("endpoint should be %v", endpoint)
		}
	}
}
