package main

import (
	"context"
	"fmt"
	"net"
	"testing"

	"log"

	pb "github.com/PierreDougnac/Todo-gRPC-Service/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024

// helper pour créer un client et un serveur en mémoire
func setupTestServer() (*grpc.ClientConn, pb.TodoServiceClient, func()) {
	lis := bufconn.Listen(bufSize)
	s := grpc.NewServer()
	pb.RegisterTodoServiceServer(s, &mockTodoServer{})

	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()

	conn, err := grpc.DialContext(context.Background(), "bufnet",
		grpc.WithContextDialer(func(ctx context.Context, s string) (net.Conn, error) {
			return lis.Dial()
		}),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("Failed to dial bufnet: %v", err)
	}

	client := pb.NewTodoServiceClient(conn)
	return conn, client, func() {
		conn.Close()
		s.Stop()
	}
}

// Mock server
type mockTodoServer struct {
	pb.UnimplementedTodoServiceServer
	todos  map[string]*pb.Todo
	nextID int
}

func (s *mockTodoServer) CreateTodo(ctx context.Context, req *pb.CreateTodoRequest) (*pb.CreateTodoResponse, error) {
	if s.todos == nil {
		s.todos = make(map[string]*pb.Todo)
		s.nextID = 1
	}
	id := fmt.Sprintf("%d", s.nextID)
	s.nextID++
	todo := &pb.Todo{Id: id, Title: req.Title}
	s.todos[id] = todo
	return &pb.CreateTodoResponse{Todo: todo}, nil
}

func (s *mockTodoServer) ListTodos(ctx context.Context, req *pb.ListTodosRequest) (*pb.ListTodosResponse, error) {
	todos := []*pb.Todo{}
	for _, t := range s.todos {
		todos = append(todos, t)
	}
	return &pb.ListTodosResponse{Todos: todos}, nil
}

func TestClientCreateListTodo(t *testing.T) {
	_, client, cleanup := setupTestServer()
	defer cleanup()

	// Test create
	createResp, err := client.CreateTodo(context.Background(), &pb.CreateTodoRequest{Title: "Test todo"})
	if err != nil {
		t.Fatalf("CreateTodo failed: %v", err)
	}
	if createResp.Todo.Title != "Test todo" {
		t.Errorf("expected title 'Test todo', got '%s'", createResp.Todo.Title)
	}

	// Test list
	listResp, err := client.ListTodos(context.Background(), &pb.ListTodosRequest{})
	if err != nil {
		t.Fatalf("ListTodos failed: %v", err)
	}
	if len(listResp.Todos) != 1 {
		t.Errorf("expected 1 todo, got %d", len(listResp.Todos))
	}
}
