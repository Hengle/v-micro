package mdns

import (
	"testing"
	"time"

	"github.com/fananchong/v-micro/registry"
)

func TestMDNS(t *testing.T) {
	testData := []*registry.Service{
		{
			Name:    "test1",
			Version: "1.0.1",
			Nodes: []*registry.Node{
				{
					ID:      "test1-1",
					Address: "10.0.0.1:10001",
					Metadata: map[string]string{
						"foo": "bar",
					},
				},
			},
		},
		{
			Name:    "test2",
			Version: "1.0.2",
			Nodes: []*registry.Node{
				{
					ID:      "test2-1",
					Address: "10.0.0.2:10002",
					Metadata: map[string]string{
						"foo2": "bar2",
					},
				},
			},
		},
		{
			Name:    "test3",
			Version: "1.0.3",
			Nodes: []*registry.Node{
				{
					ID:      "test3-1",
					Address: "10.0.0.3:10003",
					Metadata: map[string]string{
						"foo3": "bar3",
					},
				},
			},
		},
	}

	// new registry
	opt := []registry.Option{Port(5354)}
	r := NewRegistry(opt...)

	for _, service := range testData {
		// register service
		if err := r.Register(service); err != nil {
			t.Fatal(err)
		}

		// get registered service
		s, err := r.GetService(service.Name)
		if err != nil {
			t.Fatal(err)
		}

		if len(s) != 1 {
			t.Fatalf("Expected one result for %s got %d", service.Name, len(s))

		}

		if s[0].Name != service.Name {
			t.Fatalf("Expected name %s got %s", service.Name, s[0].Name)
		}

		if s[0].Version != service.Version {
			t.Fatalf("Expected version %s got %s", service.Version, s[0].Version)
		}

		if len(s[0].Nodes) != 1 {
			t.Fatalf("Expected 1 node, got %d", len(s[0].Nodes))
		}

		node := s[0].Nodes[0]

		if node.ID != service.Nodes[0].ID {
			t.Fatalf("Expected node id %s got %s", service.Nodes[0].ID, node.ID)
		}

		if node.Address != service.Nodes[0].Address {
			t.Fatalf("Expected node address %s got %s", service.Nodes[0].Address, node.Address)
		}
	}

	services, err := r.ListServices()
	if err != nil {
		t.Fatal(err)
	}

	for _, service := range testData {
		var seen bool
		for _, s := range services {
			if s.Name == service.Name {
				seen = true
				break
			}
		}
		if !seen {
			t.Fatalf("Expected service %s got nothing", service.Name)
		}

		// deregister
		if err := r.Deregister(service); err != nil {
			t.Fatal(err)
		}

		time.Sleep(time.Millisecond * 5)

		// check its gone
		s, _ := r.GetService(service.Name)
		if len(s) > 0 {
			t.Fatalf("Expected nothing got %+v", s[0])
		}
	}

}
