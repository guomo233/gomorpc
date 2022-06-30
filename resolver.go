package gomorpc

import "fmt"

// domain should be like: "domain:hostname"
type Resolver interface {
	Resolve(string) ([]string, error)
}

type StaticResolver struct {
	Hosts map[string][]string
}

type resolverRegister []Resolver

func (r *StaticResolver) Register(hosts map[string][]string) {
	r.Hosts = hosts
}

func (r *StaticResolver) Resolve(domain string) ([]string, error) {
	hosts, ok := r.Hosts[domain]
	if !ok {
		return nil, fmt.Errorf("domain %s not found", domain)
	}

	return hosts, nil
}

func (r resolverRegister) Resolve(domain string) ([]string, error) {
	for _, resolver := range r {
		hosts, err := resolver.Resolve(domain)
		if len(hosts) > 0 && err == nil {
			return hosts, nil
		}
	}

	return nil, fmt.Errorf("domain %s not found", domain)
}

var _ Resolver = (*StaticResolver)(nil)

var _resolverRegister resolverRegister

func RegisterResolver(r Resolver) {
	_resolverRegister = append(_resolverRegister, r)
}
