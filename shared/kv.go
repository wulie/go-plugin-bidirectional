package shared

import (
	"context"
	"github.com/hashicorp/go-plugin"
	"github.com/wulie/go-plugin-bidirectional/proto"
	"google.golang.org/grpc"
)

type Counter interface {
	Put(key string, value int64, a Add) error
	Get(key string) (int64, error)
}

type Add interface {
	Sum(a, b int64) (int64, error)
}

type CounterPlugin struct {
	plugin.NetRPCUnsupportedPlugin
	Impl Counter
}

func (c *CounterPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	proto.RegisterCounterServer(s, &CounterGRPCServer{
		broker: broker,
		Impl:   c.Impl,
	})
	return nil
}

func (c *CounterPlugin) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, client *grpc.ClientConn) (interface{}, error) {
	return &CounterGRPCClient{
		client: proto.NewCounterClient(client),
		broker: broker,
	}, nil
}

var Handshake = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "BASIC_PLUGIN",
	MagicCookieValue: "hello",
}

// PluginMap is the map of plugins we can dispense.
var PluginMap = map[string]plugin.Plugin{
	"counter": &CounterPlugin{},
}
