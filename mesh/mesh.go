package mesh

import (
	"github.com/fananchong/v-micro/transport"
)

// mesh 插件是否有必要存在，有待商榷。
// 该插件与 selector 插件有一定的重合。
// 可以考虑自己的 selector 插件包含该部分功能。
// 先做 selector 插件

// Mesh mesh 用于治理微服务间连接关系，宏观上看可以把 Registry + Mesh + [未来的流量控制器] 称之为`服务网格`
type Mesh interface {
	GetPeer(name string) *transport.Socket
}
