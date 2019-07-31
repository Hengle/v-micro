package mesh

import (
	"github.com/fananchong/v-micro/transport"
)

// Mesh mesh 用于治理微服务间连接关系，宏观上看可以把 Registry + Mesh + [未来的流量控制器] 称之为`服务网格`
type Mesh interface {
	GetPeer(name string) *transport.Socket
}
