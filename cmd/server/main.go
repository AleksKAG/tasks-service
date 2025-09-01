package main

import (
	"log"

	"github.com/AleksKAG/tasks-service/internal/database"
	"github.com/AleksKAG/tasks-service/internal/task"
	transportgrpc "github.com/AleksKAG/tasks-service/internal/transport/grpc"
)

func main() {
	log.Println("Starting server initialization...")
	database.InitDB()
	log.Println("Database initialized successfully")
	repo := task.NewRepository(database.DB)
	log.Println("Repository created")
	svc := task.NewService(repo)
	log.Println("Service created")
	userClient, conn, err := transportgrpc.NewUserClient("localhost:50051")
	if err != nil {
		log.Fatalf("failed to connect to users: %v", err)
	}
	defer conn.Close()
	log.Println("Connected to Users service")
	log.Println("Starting gRPC server on :50052...")
	if err := transportgrpc.RunGRPC(svc, userClient); err != nil {
		log.Fatalf("Tasks gRPC server error: %v", err)
	}
}
