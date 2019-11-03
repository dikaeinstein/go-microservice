package loadbalancer

import (
	"net/url"
	"testing"
)

func TestLoadBalancerUpdateEndpoints(t *testing.T) {
	endpoints := []url.URL{
		url.URL{Host: "www.google.com"},
		url.URL{Host: "www.google.co.uk"},
	}
	lb := New(&RandomStrategy{}, endpoints)

	newEndpoints := []url.URL{
		url.URL{Host: "https://github.com"},
		url.URL{Host: "https://gitlab.com"},
	}
	lb.UpdateEndpoints(newEndpoints)

	endpoint := lb.GetEndpoint()

	if endpoint != newEndpoints[0] && endpoint != newEndpoints[1] {
		t.Errorf("updated endpoint should be either %v or %v",
			newEndpoints[0], newEndpoints[1])
	}
}

func TestLoadBalancerPicksAnEndpoint(t *testing.T) {
	endpoints := []url.URL{
		url.URL{Host: "www.google.com"},
		url.URL{Host: "www.google.co.uk"},
	}
	lb := New(&RandomStrategy{}, endpoints)

	endpoint := lb.GetEndpoint()

	if endpoint != endpoints[0] && endpoint != endpoints[1] {
		t.Errorf("endpoint should be either %v or %v",
			endpoints[0], endpoints[1])
	}
}
