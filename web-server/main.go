package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type Employee struct {
	ID       int
	Name     string
	Age      int
	Division string
}

var employees = []Employee{
	{ID: 1, Name: "Airel", Age: 23, Division: "IT"},
	{ID: 2, Name: "Juwita", Age: 20, Division: "Finance"},
	{ID: 3, Name: "Bagus", Age: 30, Division: "IT"},
}

var PORT = ":8080"

func main() {
	http.HandleFunc("/employees", getEmployees)
	http.HandleFunc("/employee", createEmployee)

	fmt.Println("Application is listening on port", PORT)
	http.ListenAndServe(PORT, nil)
}

func getEmployees(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method == "GET" {
		json.NewEncoder(w).Encode(employees)
		return
	}

	http.Error(w, "Invalid method", http.StatusBadRequest)
	return
}

func createEmployee(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method == "POST" {
		name := r.FormValue("name")
		age := r.FormValue("age")
		division := r.FormValue("division")

		convertAge, err := strconv.Atoi(age)

		if err != nil {
			http.Error(w, "Invalid age", http.StatusBadRequest)
			return
		}

		newEmployee := Employee{
			ID:       len(employees) * 1,
			Name:     name,
			Age:      convertAge,
			Division: division,
		}

		employees = append(employees, newEmployee)

		json.NewEncoder(w).Encode(newEmployee)
		return
	}

	http.Error(w, "invalid method", http.StatusBadRequest)
}

/*
func greet(w http.ResponseWriter, r *http.Request) {
	msg := "Hello World"

	fmt.Fprint(w, msg)
}
*/
