package main

import (
    "context"
    "database/sql"
    "testing"

    _ "github.com/lib/pq"
    pb "grpc/protos/employee"
)

func TestGrpcServer(t *testing.T) {
    db, err := sql.Open("postgres", dbConnectionStr)
    if err != nil {
        t.Fatalf("Failed to connect to database: %v", err)
    }
    defer db.Close()

    server := &server{db: db}
    ctx := context.Background()

    // Test Create Employee
    createReq := &pb.Employee{
        Name:       "Sanjay A U",
        Role:       "Developer",
        Department: "Engineering",
    }
    createResp, err := server.CreateEmployee(ctx, createReq)
    if err != nil {
        t.Fatalf("CreateEmployee() error = %v", err)
    }
    if createResp.Id == 0 {
        t.Fatalf("CreateEmployee() returned invalid ID")
    }
    employeeID := createResp.Id

    // Test Get Employee
    getReq := &pb.EmployeeRequest{Id: employeeID}
    getResp, err := server.GetEmployee(ctx, getReq)
    if err != nil {
        t.Errorf("GetEmployee() error = %v", err)
    }
    if getResp.Name != createReq.Name {
        t.Errorf("GetEmployee() returned incorrect name, got %v, want %v", getResp.Name, createReq.Name)
    }

    // Test Update Employee
    updateReq := &pb.Employee{
        Id:         employeeID,
        Name:       "Jane Doe",
        Role:       "Senior Developer",
        Department: "Engineering",
    }
    updateResp, err := server.UpdateEmployee(ctx, updateReq)
    if err != nil {
        t.Errorf("UpdateEmployee() error = %v", err)
    }
    if updateResp.Role != "Senior Developer" {
        t.Errorf("UpdateEmployee() failed to update role, got %v", updateResp.Role)
    }

    // Test Delete Employee
    deleteReq := &pb.EmployeeRequest{Id: employeeID}
    deleteResp, err := server.DeleteEmployee(ctx, deleteReq)
    if err != nil {
        t.Errorf("DeleteEmployee() error = %v", err)
    }
    if deleteResp.Response != "Deleted Successfully" {
        t.Errorf("DeleteEmployee() returned incorrect response, got %v", deleteResp.Response)
    }
}
