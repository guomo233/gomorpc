package gomorpc

import "math/rand"

type LoadBalancer interface {
	Get([]string) string
}

type RandomLoadBalancer struct {}

type RoundRabinLoadBalancer struct {
	next int
}

func (r *RandomLoadBalancer) Get(hosts []string) string {
	idx := rand.Intn(len(hosts))
	return hosts[idx]
}

var _ LoadBalancer = (*RandomLoadBalancer)(nil)

func (r *RoundRabinLoadBalancer) Get(hosts []string) string {
	if r.next < len(hosts) {
		r.next += 1
		return hosts[r.next - 1]
	}

	r.next = 1
	return hosts[0]
}

var _ LoadBalancer = (*RoundRabinLoadBalancer)(nil)
