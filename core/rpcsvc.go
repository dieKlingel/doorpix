package core

//go:generate

import (
	"log/slog"
	"sync"

	"github.com/dieklingel/doorpix/core/internal/doorpix"
	"github.com/dieklingel/doorpix/core/internal/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type connection struct {
	write chan *proto.RunRequest
	read  chan *proto.RunResponse
}

type RPCService struct {
	proto.UnimplementedCoreServer

	Config doorpix.Config

	server      *grpc.Server
	connections []*connection
	mutex       sync.Mutex
}

func (service *RPCService) Name() string {
	return "rpc-service"
}

func (service *RPCService) Init() error {
	slog.Debug("init rpc service")

	service.connections = make([]*connection, 0)
	service.server = grpc.NewServer()
	reflection.Register(service.server)

	proto.RegisterCoreServer(service.server, service)
	slog.Debug("successfully initialized rpc service")
	return nil
}

/*
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
		service.server.Stop()
		slog.Debug("successfully shut down rpc service")
	}()

	return nil
}

func (service *RPCService) Run(act doorpix.Action, event *doorpix.ActionHook) bool {
	action, ok := act.(doorpix.RPCAction)
	if !ok {
		return false
	}

	data := make(map[string]string)
	for k, v := range event.Data {
		data[k] = fmt.Sprintf("%v", v)
	}

	success := false
	for _, conn := range service.connections {
		conn.write <- &proto.RunRequest{
			Spec: action.Spec,
			Data: data,
		}

		select {
		case resp, ok := <-conn.read:
			if !ok {
				continue
			}
			slog.Debug("received response", "response", resp)

			if resp.Success {
				for k, v := range resp.Data {
					event.Data[k] = v
				}
				success = true
			}
		case <-time.After(action.Timeout):
			slog.Warn("timeout")
		}

		if success && !action.Multiple {
			break
		}
	}

	return success
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

func (service *RPCService) Listen(stream grpc.BidiStreamingServer[proto.RunResponse, proto.RunRequest]) error {
	conn := &connection{
		write: make(chan *proto.RunRequest),
		read:  make(chan *proto.RunResponse),
	}
	service.mutex.Lock()
	service.connections = append(service.connections, conn)
	index := len(service.connections) - 1
	service.mutex.Unlock()

	go func() {
		for {
			msg, ok := <-conn.write
			if !ok {
				return
			}

			if err := stream.Send(msg); err != nil {
				slog.Warn("failed to send message", "error", err)
			}
		}
	}()

	for {
		msg, err := stream.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) || errors.Is(err, context.Canceled) {
				close(conn.write)
				close(conn.read)
				service.mutex.Lock()
				service.connections = append(service.connections[:index], service.connections[index+1:]...)
				service.mutex.Unlock()
				return nil
			}

			return err
		}

		select {
		case conn.read <- msg:
		default:
			slog.Warn("could not use rpc message, due to actiuon run is pending", "message", msg)
		}
	}
}
*/
