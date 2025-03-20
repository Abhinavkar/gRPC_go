package main

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	pb "grpc/protos/employee"
)

const (
	serverAddress = "localhost:50051"
)

func main() {
	conn, err := grpc.NewClient(serverAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Did not connect to server: %v", err)
	}
	defer conn.Close()

	client := pb.NewEmployeeServiceClient(conn)

	// Create Employee
	createdEmployee, err := createEmployee(client)
	if err != nil {
		log.Fatalf("Error creating employee: %v", err)
	}
	fmt.Printf("Created Employee: %v\n", createdEmployee)

	// Get Employee
	fetchedEmployee, err := getEmployee(client, createdEmployee.Id) // WE CAN CHANGE ACCORDING ID OF CREATED EMPLOYEE cretedEmployee.Id TO 1 OR 2
	if err != nil {
		log.Fatalf("Error getting employee: %v", err)
	}
	fmt.Printf("Fetched Employee: %v\n", fetchedEmployee)

	// Update Employee
	updatedEmployee, err := updateEmployee(client, fetchedEmployee)
	if err != nil {
		log.Fatalf("Error updating employee: %v", err)
	}
	fmt.Printf("Updated Employee: %v\n", updatedEmployee)

	// Delete Employee
	deleteResponse, err := deleteEmployee(client, updatedEmployee.Id)
	if err != nil {
		log.Fatalf("Error deleting employee: %v", err)
	}
	fmt.Printf("Delete Response: %v\n", deleteResponse.Response)
}

func createEmployee(client pb.EmployeeServiceClient) (*pb.Employee, error) {
	employee := &pb.Employee{
		Name:       "Abhinav Kar",
		Role:       "Developer",
		Department: "Engineering",
	}
	return client.CreateEmployee(context.Background(), employee)
}

func getEmployee(client pb.EmployeeServiceClient, employeeID int32) (*pb.Employee, error) {
	employeeRequest := &pb.EmployeeRequest{Id: employeeID}
	return client.GetEmployee(context.Background(), employeeRequest)
}

func updateEmployee(client pb.EmployeeServiceClient, employee *pb.Employee) (*pb.Employee, error) {
	updatedEmployee := &pb.Employee{
		Id:         employee.Id,
		Name:       "Abhinav Kar",
		Role:       "Senior Developer",
		Department: "Engineering",
	}
	return client.UpdateEmployee(context.Background(), updatedEmployee)
}

func deleteEmployee(client pb.EmployeeServiceClient, employeeID int32) (*pb.StringResponse, error) {
	employeeRequest := &pb.EmployeeRequest{Id: employeeID}
	return client.DeleteEmployee(context.Background(), employeeRequest)
}
