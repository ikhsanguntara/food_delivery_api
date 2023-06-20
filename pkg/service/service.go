package service

import (
	"food_delivery_api/pkg/model"
	"food_delivery_api/pkg/storage/mysql"
)

type Service interface {
	// User
	AddUser(model.User) (model.User, error)
	GetUsers(model.QueryPagination) ([]model.User, int64, error)
	GetUser(model.User) (model.User, error)
	GetUserByEmailPassword(model.User) (model.User, error)
	EditUser(model.User) (model.User, error)
	RemoveUser(model.User) (model.User, error)

	// Category
	AddCategory(model.Category) (model.Category, error)
	GetCategories(model.QueryPagination) ([]model.Category, int64, error)
	GetCategory(model.Category) (model.Category, error)
	EditCategory(model.Category) (model.Category, error)
	RemoveCategory(model.Category) (model.Category, error)

	// UOM
	AddUOM(model.UOM) (model.UOM, error)
	GetUOMs(model.QueryPagination) ([]model.UOM, int64, error)
	GetUOM(model.UOM) (model.UOM, error)
	EditUOM(model.UOM) (model.UOM, error)
	RemoveUOM(model.UOM) (model.UOM, error)

	// Product
	AddProduct(model.ProductRequest) (model.Product, error)
	GetProducts(model.QueryPagination) ([]model.Product, int64, error)
	GetProduct(model.Product) (model.Product, error)
	EditProduct(model.ProductRequest) (model.Product, error)
	RemoveProduct(model.Product) (model.Product, error)

	// Transaction
	AddTransaction(model.Transaction, string) (model.Transaction, error)
	GetTransactions(model.QueryGetTransactions) ([]model.Transaction, int64, error)
	GetTransaction(model.Transaction) (model.Transaction, error)
	EditTransaction(model.Transaction) (model.Transaction, error)
	RemoveTransaction(model.Transaction) (model.Transaction, error)
	
	// Reports
	GetReportDashboard() (model.Dashboard, error)
	GetReportMetodologi(model.MetodologiFilter) ([]model.ResponseMetodologi, error)
}

type service struct {
	rmy mysql.RepositoryMySQL
}

func NewService(rmy mysql.RepositoryMySQL) Service {
	return &service{rmy}
}
