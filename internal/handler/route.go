package handler

import "fmt"

var Api = "/api"

// Clean
var (
	GetAllCleaning     = fmt.Sprintf("GET %s/cleaning", Api)
	PostNewCleaning    = fmt.Sprintf("POST %s/cleaning", Api)
	DeleteCleaningById = fmt.Sprintf("DELETE %s/cleaning", Api)
)

// Empolyee
var (
	GetAllEmployee     = fmt.Sprintf("GET %s/employee", Api)
	PostNewEmployee    = fmt.Sprintf("POST %s/employee", Api)
	DeleteEmployeeById = fmt.Sprintf("DELETE %s/employee", Api)
)

// Contact with employee
var (
	GetAllEmployeeContacts     = fmt.Sprintf("GET %s/contacts", Api)
	PostNewEmployeeContacts    = fmt.Sprintf("POST %s/contacts", Api)
	DeleteEmployeeContactsById = fmt.Sprintf("DELETE %s/contacts", Api)
)
