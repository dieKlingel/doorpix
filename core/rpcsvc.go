package core

//go:generate

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"sync"

	"github.com/dieklingel/doorpix/core/internal/doorpix"
	"github.com/dieklingel/doorpix/core/internal/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type RPCService struct {
	proto.UnimplementedCoreServer

	Config  doorpix.Config
	Emitter doorpix.Emit

	server *grpc.Server
}

func (service *RPCService) Init() error {
	slog.Debug("init rpc service")

	service.server = grpc.NewServer()
	reflection.Register(service.server)

	proto.RegisterCoreServer(service.server, service)
	slog.Debug("successfully initialized rpc service")
	return nil
}

func (service *RPCService) Exec(ctx context.Context, wg *sync.WaitGroup) error {
	slog.Debug("exec rpc service")

	go func() {
		addr := fmt.Sprintf("%s:%d", service.Config.RPC.Host, service.Config.RPC.Port)
		socket, err := net.Listen("tcp", addr)
		if err != nil {
			slog.Error("failed to listen", "error", err)
			return
		}
		if err := service.server.Serve(socket); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
			slog.Error("rpc server error", "error", err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		<-ctx.Done()
		slog.Debug("shutting down rpc service")
		service.server.GracefulStop()
		slog.Debug("successfully shut down rpc service")
	}()

	return nil
}

func (service *RPCService) Emit(ctx context.Context, in *proto.EmitRequest) (*proto.EmitResponse, error) {
	if in.Type == "" {
		return nil, errors.New("type is required")
	}

	eventtype := doorpix.EventType(in.Type)
	eventdata := make(map[string]any)
	for k, v := range in.Data {
		eventdata[k] = v
	}

	service.Emitter(eventtype, eventdata)

	return &proto.EmitResponse{}, nil
}
