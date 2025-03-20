package main

import (
	"context"
	"database/sql"
	"fmt"
	"net"

	pb "grpc/protos/employee"

	_ "github.com/lib/pq"
	"google.golang.org/grpc"
)

const (
	dbConnectionStr = "user=postgres password=postgres@123 dbname=go_employee_db sslmode=disable"
)

type server struct {
	pb.UnimplementedEmployeeServiceServer
	db *sql.DB
}

func main() {

	db, err := sql.Open("postgres", dbConnectionStr)
	if err != nil {
		fmt.Println("An error occurred while connecting to the database:", err)
		return
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		fmt.Println("An error occurred while pinging the database:", err)
		return
	}
	fmt.Println("Connected to DB")

	// Start listening on port
	listen, err := net.Listen("tcp", ":50051")
	if err != nil {
		fmt.Println("Failed to listen on port:", err)
		return
	}

	s := grpc.NewServer()
	pb.RegisterEmployeeServiceServer(s, &server{db: db})
	fmt.Println("Server is running on port 50051")
	if err := s.Serve(listen); err != nil {
		fmt.Println("Error occurred while serving:", err)
	}
}

func (s *server) CreateEmployee(c context.Context, request *pb.Employee) (*pb.Employee, error) {
	fmt.Println("Recieved Create Request:", request)
	query := `INSERT INTO employees (name, role, department) VALUES ($1, $2, $3) RETURNING id`
	err := s.db.QueryRow(query, request.Name, request.Role, request.Department).Scan(&request.Id)
	if err != nil {
		return nil, err
	}
	return request, nil
}

func (s *server) GetEmployee(c context.Context, request *pb.EmployeeRequest) (*pb.Employee, error) {
	fmt.Println("Recieved Get Request:", request)
	emp := &pb.Employee{}
	query := `SELECT id, name, role, department FROM employees WHERE id=$1`
	err := s.db.QueryRow(query, request.Id).Scan(&emp.Id, &emp.Name, &emp.Role, &emp.Department)
	if err != nil {
		return nil, err
	}
	return emp, nil
}
func (s *server) UpdateEmployee(c context.Context, request *pb.Employee) (*pb.Employee, error) {
	fmt.Println("Recieved Update Request:", request)
	emp := &pb.Employee{}
	query := `UPDATE employees SET name=$1, role=$2, department=$3 WHERE id=$4 RETURNING id, name, role, department`
	err := s.db.QueryRow(query, request.Name, request.Role, request.Department, request.Id).Scan(&emp.Id, &emp.Name, &emp.Role, &emp.Department)
	if err != nil {
		return nil, err
	}
	return emp, nil
}

func (s *server) DeleteEmployee(c context.Context, request *pb.EmployeeRequest) (*pb.StringResponse, error) {
	fmt.Println("Recieved Delete Request:", request)
	query := `DELETE FROM employees WHERE id=$1`
	_, err := s.db.Exec(query, request.Id)
	if err != nil {
		return nil, err
	}
	return &pb.StringResponse{Response: "Deleted Successfully"}, nil
}
