package grpc

import (
	"log"
	"net"

	taskpb "github.com/AleksKAG/project-protos/proto/task"
	"github.com/AleksKAG/tasks-service/internal/task"
	"google.golang.org/grpc"
)

func RunGRPC(svc *task.Service) error {
	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		return err
	}
	grpcSrv := grpc.NewServer()
	handler := NewHandler(svc)
	taskpb.RegisterTaskServiceServer(grpcSrv, handler)
	log.Println("Server running at :50052")
	return grpcSrv.Serve(lis)
}
