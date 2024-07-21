package main

import (
	"context"
	"log"
	"net"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "github.com/mashinolol/auth-service-grpc/proto/"
)

const (
	port = ":50051"
)

type server struct {
	db *sqlx.DB
	pb.UnimplementedUserServiceServer
}

func (s *server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	// Валидация данных
	if req.Name == "" || req.Email == "" || req.Password == "" {
		return &pb.CreateUserResponse{
			Success: false,
			Message: "Invalid input",
		}, nil
	}

	// Сохранение данных пользователя в БД
	_, err := s.db.Exec("INSERT INTO users (name, email, password) VALUES ($1, $2, $3)", req.Name, req.Email, req.Password)
	if err != nil {
		return &pb.CreateUserResponse{
			Success: false,
			Message: "Failed to create user",
		}, err
	}

	return &pb.CreateUserResponse{
		Success: true,
		Message: "User created successfully",
	}, nil
}

func main() {
	db, err := sqlx.Connect("postgres", "user=postgres dbname=mydb sslmode=disable")
	if err != nil {
		log.Fatalln(err)
	}

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterUserServiceServer(s, &server{db: db})
	reflection.Register(s)

	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
