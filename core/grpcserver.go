package core

//go:generate

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"strings"

	"github.com/dieklingel/doorpix/core/internal/eventemitter"
	"github.com/dieklingel/doorpix/core/internal/proto"
	"github.com/dieklingel/doorpix/core/internal/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type GRPCServerProps struct {
	Address string
	Port    int
}

type GRPCServer struct {
	props GRPCServerProps

	eventemitter *eventemitter.EventEmitter
	server       *grpc.Server
	ctx          service.Context
	methods      *rpcMethods
}

func NewGRPCServer(eventemitter *eventemitter.EventEmitter, props GRPCServerProps) *GRPCServer {
	server := &GRPCServer{
		props:        props,
		eventemitter: eventemitter,
		ctx:          service.NewContext(context.Background()),
		methods:      &rpcMethods{},
	}
	server.methods.server = server
	return server
}

func (server *GRPCServer) Start() {
	server.server = grpc.NewServer()
	reflection.Register(server.server)

	proto.RegisterCoreServer(server.server, server.methods)

	server.exec()
}

func (server *GRPCServer) Stop() {
	server.server.Stop()
	server.ctx.Wait()
}

func (server *GRPCServer) exec() {

	server.ctx.Lock()
	go func() {
		defer server.ctx.Unlock()

		addr := fmt.Sprintf("%s:%d", server.props.Address, server.props.Port)
		socket, err := net.Listen("tcp", addr)
		if err != nil {
			slog.Error("failed to listen", "error", err)
			return
		}
		if err := server.server.Serve(socket); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
			slog.Error("rpc server error", "error", err)
		}
	}()

}

type rpcMethods struct {
	proto.UnimplementedCoreServer
	server *GRPCServer
}

// Emit is a gRPC method that triggers an event.
func (rpc *rpcMethods) Emit(ctx context.Context, req *proto.EmitRequest) (*emptypb.Empty, error) {
	eventPath := req.EventPath
	eventData := make(map[string]any)
	for key, value := range req.EventData {
		eventData[key] = value
	}

	if eventPath == "" {
		return nil, fmt.Errorf("event path is required")
	}

	eventPath = fmt.Sprintf("events/rpc/%s", eventPath)
	err := rpc.server.eventemitter.Emit(eventPath, eventData)
	return &emptypb.Empty{}, err
}

// GetKioskDefinition is a gRPC method that retrieves the kiosk definition.
func (rpc *rpcMethods) GetKioskDefinition(context.Context, *emptypb.Empty) (*proto.KioskDefinitionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetKioskDefinition not implemented")
}

// OnStateUpdate is a gRPC method that streams state updates.
func (rpc *rpcMethods) Listen(req *proto.ListenRequest, stream grpc.ServerStreamingServer[proto.ListenResponse]) error {
	allowedPathPrefixes := []string{
		"state",
	}

	eventPathIsAllowed := false
	for _, allowedPrefix := range allowedPathPrefixes {
		if strings.HasPrefix(req.EventPath, allowedPrefix) {
			eventPathIsAllowed = true
			break
		}
	}
	if !eventPathIsAllowed {
		return status.Errorf(codes.PermissionDenied, "the path %s is not allowed", req.EventPath)
	}

	listener := rpc.server.eventemitter.Listen(req.EventPath)

	for {
		select {
		case <-rpc.server.ctx.Done():
			return nil
		case <-stream.Context().Done():
			slog.Debug("rpc client disconnected")
			listener.Close()
			return nil

		case event := <-listener.Listen():
			data := make(map[string]string)
			for key, value := range event.Data {
				data[key] = fmt.Sprintf("%v", value)
			}

			msg := proto.ListenResponse{
				EventPath: event.Event,
				EventData: data,
			}

			if err := stream.Send(&msg); err != nil {
				slog.Error("failed to send event", "error", err)
			}
		}
	}
}
