package handler

var Api = "/api"

// Clean
var (
	GetAllCleaning     = "/api/cleanings"
	PostNewCleaning    = "/api/cleanings"
	DeleteCleaningById = "/api/cleanings"
)

// Employee
var (
	GetAllEmployee     = "/api/employees"
	PostNewEmployee    = "/api/employees"
	DeleteEmployeeById = "/api/employees"
)

// Contact with employee
var (
	GetAllEmployeeContacts     = "/api/contacts"
	PostNewEmployeeContacts    = "/api/contacts"
	DeleteEmployeeContactsById = "/api/contacts"
)

//STORE

// Dishes
var (
	GetAllDishes     = "/api/dishes"
	PostNewDishes    = "/api/dishes"
	DeleteDishesById = "/api/dishes"
)

// Orders
var (
	GetAllOrders = "/api/orders"
	GetOrderByID = "/api/orders/id"
	PostNewOrder = "/api/orders"
	DeleteOrder  = "/api/orders"
)

// OrderItem
var (
	GetOrderItems   = "/api/order-items"
	UpdateOrderItem = "/api/order-items"
	DeleteOrderItem = "/api/order-items"
)
