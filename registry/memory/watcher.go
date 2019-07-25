package memory

import (
	"errors"

	"github.com/fananchong/v-micro/registry"
)

type watcherImpl struct {
	id   string
	wo   registry.WatchOptions
	res  chan *registry.Result
	exit chan bool
}

func (m *watcherImpl) Next() (*registry.Result, error) {
	for {
		select {
		case r := <-m.res:
			if len(m.wo.Service) > 0 && m.wo.Service != r.Service.Name {
				continue
			}
			return r, nil
		case <-m.exit:
			return nil, errors.New("Watcher stopped")
		}
	}
}

func (m *watcherImpl) Stop() {
	select {
	case <-m.exit:
		return
	default:
		close(m.exit)
	}
}
