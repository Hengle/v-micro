package selector

import (
	"math/rand"
	"time"

	"github.com/fananchong/v-micro/registry"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// StatefulRandom Random is a random strategy algorithm for node selection
func StatefulRandom(services []*registry.Service) Next {
	var nodes []*registry.Node

	for _, service := range services {
		nodes = append(nodes, service.Nodes...)
	}
	return func() (*registry.Node, error) {
		if len(nodes) == 0 {
			return nil, ErrNoneAvailable
		}

		i := rand.Int() % len(nodes)
		node := nodes[i]
		nodes = append(nodes[:i], nodes[i+1:]...)
		return node, nil
	}
}

// Random is a random strategy algorithm for node selection
func Random(services []*registry.Service) Next {
	var nodes []*registry.Node

	for _, service := range services {
		nodes = append(nodes, service.Nodes...)
	}

	return func() (*registry.Node, error) {
		if len(nodes) == 0 {
			return nil, ErrNoneAvailable
		}

		i := rand.Int() % len(nodes)
		return nodes[i], nil
	}
}

// StatefulRoundRobin RoundRobin is a roundrobin strategy algorithm for node selection
func StatefulRoundRobin(services []*registry.Service) Next {
	var nodes []*registry.Node

	for _, service := range services {
		nodes = append(nodes, service.Nodes...)
	}

	var i = rand.Int() % len(nodes)
	origin := i
	flag := false
	return func() (*registry.Node, error) {
		if len(nodes) == 0 || flag {
			return nil, ErrNoneAvailable
		}

		node := nodes[i]

		i = (i + 1) % len(nodes)
		if i == origin {
			flag = true
		}

		return node, nil
	}
}

// RoundRobin is a roundrobin strategy algorithm for node selection
func RoundRobin(services []*registry.Service) Next {
	var nodes []*registry.Node

	for _, service := range services {
		nodes = append(nodes, service.Nodes...)
	}

	var i = rand.Int()

	return func() (*registry.Node, error) {
		if len(nodes) == 0 {
			return nil, ErrNoneAvailable
		}

		node := nodes[i%len(nodes)]
		i++

		return node, nil
	}
}
