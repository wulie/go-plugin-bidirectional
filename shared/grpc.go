package shared

import (
	"context"
	"github.com/hashicorp/go-plugin"
	"github.com/wulie/go-plugin-bidirectional/proto"
	"google.golang.org/grpc"
)

type CounterGRPCClient struct {
	client proto.CounterClient
	broker *plugin.GRPCBroker
}

func (c *CounterGRPCClient) Put(key string, value int64, a Add) error {
	addServer := &AddGRPCServer{Impl: a}
	var s *grpc.Server
	serverFunc := func(opts []grpc.ServerOption) *grpc.Server {
		s = grpc.NewServer(opts...)
		proto.RegisterAddHelperServer(s, addServer)
		return s
	}
	id := c.broker.NextId()
	go c.broker.AcceptAndServe(id, serverFunc)
	_, err := c.client.Put(context.Background(), &proto.PutRequest{
		AddServer: id,
		Key:       key,
		Value:     value,
	})
	return err
}

func (c *CounterGRPCClient) Get(key string) (int64, error) {
	resp, err := c.client.Get(context.Background(), &proto.GetRequest{Key: key})
	if err != nil {
		return 0, err
	}
	return resp.Value, nil
}

type CounterGRPCServer struct {
	broker *plugin.GRPCBroker
	Impl   Counter
}

func (c *CounterGRPCServer) Get(ctx context.Context, request *proto.GetRequest) (*proto.GetResponse, error) {
	get, err := c.Impl.Get(request.Key)
	if err != nil {
		return nil, err
	}
	return &proto.GetResponse{Value: get}, nil
}

func (c *CounterGRPCServer) Put(ctx context.Context, request *proto.PutRequest) (*proto.Empty, error) {
	conn, err := c.broker.Dial(request.AddServer)
	if err != nil {
		return nil, err
	}
	client := proto.NewAddHelperClient(conn)

	err = c.Impl.Put(request.Key, request.Value, &AddGRPCClient{client: client})
	return &proto.Empty{}, err
}

type AddGRPCClient struct {
	client proto.AddHelperClient
}

func (a *AddGRPCClient) Sum(x, y int64) (int64, error) {
	sum, err := a.client.Sum(context.Background(), &proto.SumRequest{
		A: x,
		B: y,
	})
	if err != nil {
		return 0, err
	}
	return sum.R, err
}

type AddGRPCServer struct {
	Impl Add
}

func (a *AddGRPCServer) Sum(c context.Context, request *proto.SumRequest) (*proto.SumResponse, error) {
	sum, err := a.Impl.Sum(request.A, request.B)
	if err != nil {
		return nil, err
	}
	return &proto.SumResponse{R: sum}, nil
}
