package grpc

import (
	"log"
	"net"

	taskpb "github.com/AleksKAG/project-protos/proto/task"
	userpb "github.com/AleksKAG/project-protos/proto/user"
	"github.com/AleksKAG/tasks-service/internal/task"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func RunGRPC(svc *task.Service, uc userpb.UserServiceClient) error {
	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		return err
	}
	grpcSrv := grpc.NewServer()
	handler := NewHandler(svc, uc)
	taskpb.RegisterTaskServiceServer(grpcSrv, handler)

	// включаем server reflection
	reflection.Register(grpcSrv)

	log.Println("Server running at :50052")
	return grpcSrv.Serve(lis)
}
