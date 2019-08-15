// Package cache provides a registry cache
package cache

import (
	"math"
	"math/rand"
	"sync"
	"time"

	"github.com/fananchong/v-micro/common/log"
	"github.com/fananchong/v-micro/registry"
)

// Cache is the registry cache interface
type Cache interface {
	// embed the registry interface
	registry.Registry
	// stop the cache watcher
	Stop()
}

type cache struct {
	registry.Registry
	opts Options

	// registry cache
	sync.RWMutex
	cache  map[string][]*registry.Service
	inited sync.Map

	exit chan bool
}

var (
	// DefaultTTL default ttl
	DefaultTTL = 30 * time.Second
)

func backoff(attempts int) time.Duration {
	if attempts == 0 {
		return time.Duration(0)
	}
	return time.Duration(math.Pow(10, float64(attempts))) * time.Millisecond
}

func (c *cache) quit() bool {
	select {
	case <-c.exit:
		return true
	default:
		return false
	}
}

func (c *cache) del(service string) {
	delete(c.cache, service)
}

func (c *cache) get(service string) (services []*registry.Service, err error) {
	// watch service if not watched
	if _, ok := c.inited.LoadOrStore(service, 1); !ok {
		log.Infof("watch service: %s", service)
		// ask the registry
		if services, err = c.Registry.GetService(service); err != nil {
			log.Error(err)
			return
		}
		for _, s := range services {
			c.update(&registry.Result{Action: "create", Service: s}, false)
		}

		go c.run(service)
		go c.updateCache(service)
	}

	// read lock
	c.RLock()
	defer c.RUnlock()

	// check the cache first
	services = c.cache[service]

	// make a copy
	return registryCopy(services), nil
}

func (c *cache) updateCache(service string) {
	for {
		// exit early if already dead
		if c.quit() {
			return
		}

		select {
		case <-time.After(DefaultTTL):
			// ask the registry
			newservices, err := c.Registry.GetService(service)
			if err != nil {
				log.Error(err)
				continue
			}

			func() {
				c.Lock()
				defer c.Unlock()

				// delete invaild node
				c.deleteOldNode(service, newservices)

				// update vaild node
				for _, s := range newservices {
					c.update(&registry.Result{Action: "update", Service: s}, false)
				}
			}()
		}
	}
}

func (c *cache) deleteOldNode(service string, newservices []*registry.Service) {
	// find invaild node, delete it
	oldservices := registryCopy(c.cache[service])
	for _, oldservice := range oldservices {
		// find same version
		var newservice *registry.Service
		for _, service := range newservices {
			if service.Version == oldservice.Version {
				newservice = service
				break
			}
		}
		if newservice != nil {
			// delete vaild node
			for i := len(oldservice.Nodes) - 1; i >= 0; i-- {
				old := oldservice.Nodes[i]
				for _, new := range newservice.Nodes {
					if old.ID == new.ID {
						oldservice.Nodes = append(oldservice.Nodes[:i], oldservice.Nodes[i+1:]...)
						break
					}
				}
			}
		}
		c.update(&registry.Result{Action: "delete", Service: oldservice}, false)
	}
}

func (c *cache) set(service string, services []*registry.Service) {
	c.cache[service] = services
}

func (c *cache) update(res *registry.Result, lock bool) {
	if res == nil || res.Service == nil {
		return
	}

	if lock {
		c.Lock()
		defer c.Unlock()
	}

LABEL:
	services, ok := c.cache[res.Service.Name]
	if !ok {
		c.cache[res.Service.Name] = make([]*registry.Service, 0)
		goto LABEL
	}

	// existing service found
	var service *registry.Service
	var index int
	for i, s := range services {
		if s.Version == res.Service.Version {
			service = s
			index = i
		}
	}

	switch res.Action {
	case "create", "update":
		if service == nil {
			c.set(res.Service.Name, append(services, res.Service))
			for _, cur := range res.Service.Nodes {
				log.Infof("service register, node: %v", cur)
			}
			return
		}

		// log
		for _, node := range res.Service.Nodes {
			var seen bool
			for _, cur := range service.Nodes {
				if cur.ID == node.ID {
					seen = true
					break
				}
			}
			if !seen {
				log.Infof("service register, node: %v", node)
			}
		}

		// append old nodes to new service
		for _, cur := range service.Nodes {
			var seen bool
			for _, node := range res.Service.Nodes {
				if cur.ID == node.ID {
					seen = true
					break
				}
			}
			if !seen {
				res.Service.Nodes = append(res.Service.Nodes, cur)
			}
		}

		services[index] = res.Service
		c.set(res.Service.Name, services)
	case "delete":
		if service == nil {
			return
		}

		var nodes []*registry.Node

		// filter cur nodes to remove the dead one
		for _, cur := range service.Nodes {
			var seen bool
			for _, del := range res.Service.Nodes {
				if del.ID == cur.ID {
					seen = true
					break
				}
			}
			if !seen {
				nodes = append(nodes, cur)
			} else {
				log.Infof("service deregister, node: %s", cur.ID)
			}
		}

		// still got nodes, save and return
		if len(nodes) > 0 {
			service.Nodes = nodes
			services[index] = service
			c.set(service.Name, services)
			return
		}

		// zero nodes left

		// only have one thing to delete
		// nuke the thing
		if len(services) == 1 {
			c.del(service.Name)
			return
		}

		// still have more than 1 service
		// check the version and keep what we know
		var srvs []*registry.Service
		for _, s := range services {
			if s.Version != service.Version {
				srvs = append(srvs, s)
			}
		}

		// save
		c.set(service.Name, srvs)
	}
}

// run starts the cache watcher loop
// it creates a new watcher if there's a problem
func (c *cache) run(service string) {

	var a, b int

	for {
		// exit early if already dead
		if c.quit() {
			return
		}

		// jitter before starting
		j := rand.Int63n(100)
		time.Sleep(time.Duration(j) * time.Millisecond)

		// create new watcher
		w, err := c.Registry.Watch(
			registry.WatchService(service),
		)

		if err != nil {
			log.Error(err)
			if c.quit() {
				return
			}

			d := backoff(a)

			if a > 3 {
				log.Info("rcache: ", err, " backing off ", d)
				a = 0
			}

			time.Sleep(d)
			a++

			continue
		}

		// reset a
		a = 0

		// watch for events
		if err := c.watch(w); err != nil {
			log.Error(err)
			if c.quit() {
				return
			}

			d := backoff(b)

			if b > 3 {
				log.Info("rcache: ", err, " backing off ", d)
				b = 0
			}

			time.Sleep(d)
			b++

			continue
		}

		// reset b
		b = 0
	}
}

// watch loops the next event and calls update
// it returns if there's an error
func (c *cache) watch(w registry.Watcher) error {
	// used to stop the watch
	stop := make(chan bool)

	// manage this loop
	go func() {
		defer w.Stop()

		select {
		// wait for exit
		case <-c.exit:
			return
		// we've been stopped
		case <-stop:
			return
		}
	}()

	for {
		res, err := w.Next()
		if err != nil {
			log.Error(err)
			close(stop)
			return err
		}
		c.update(res, true)
	}
}

func (c *cache) GetService(service string) ([]*registry.Service, error) {
	// get the service
	services, err := c.get(service)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	// if there's nothing return err
	if len(services) == 0 {
		return nil, registry.ErrNotFound
	}

	// return services
	return services, nil
}

func (c *cache) Stop() {
	select {
	case <-c.exit:
		return
	default:
		close(c.exit)
	}
}

func (c *cache) String() string {
	return "rcache"
}

// New returns a new cache
func New(r registry.Registry, opts ...Option) Cache {
	rand.Seed(time.Now().UnixNano())
	options := Options{
		TTL: DefaultTTL,
	}

	for _, o := range opts {
		o(&options)
	}

	return &cache{
		Registry: r,
		opts:     options,
		cache:    make(map[string][]*registry.Service),
		exit:     make(chan bool),
	}
}

// Copy makes a copy of services
func registryCopy(current []*registry.Service) []*registry.Service {
	services := make([]*registry.Service, len(current))
	for i, service := range current {
		// copy service
		s := new(registry.Service)
		*s = *service

		// copy nodes
		nodes := make([]*registry.Node, len(service.Nodes))
		for j, node := range service.Nodes {
			n := new(registry.Node)
			*n = *node
			nodes[j] = n
		}
		s.Nodes = nodes

		// append service
		services[i] = s
	}

	return services
}
