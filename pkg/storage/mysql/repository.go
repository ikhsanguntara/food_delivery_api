package mysql

import (
	"food_delivery_api/cfg"
	"food_delivery_api/pkg/model"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type RepositoryMySQL interface {
	// Users
	CreateUser(model.User) (model.User, error)
	ReadUsers(model.QueryPagination) ([]model.User, int64, error)
	ReadUser(model.User) (model.User, error)
	ReadUserByEmailPassword(model.User) (model.User, error)
	UpdateUser(model.User) (model.User, error)
	DeleteUser(model.User) (model.User, error)

	// Categories
	CreateCategory(model.Category) (model.Category, error)
	ReadCategories(model.QueryPagination) ([]model.Category, int64, error)
	ReadCategory(model.Category) (model.Category, error)
	UpdateCategory(model.Category) (model.Category, error)
	DeleteCategory(model.Category) (model.Category, error)

	// UOMs
	CreateUOM(model.UOM) (model.UOM, error)
	ReadUOMs(model.QueryPagination) ([]model.UOM, int64, error)
	ReadUOM(model.UOM) (model.UOM, error)
	UpdateUOM(model.UOM) (model.UOM, error)
	DeleteUOM(model.UOM) (model.UOM, error)

	// Products
	CreateProduct(model.Product) (model.Product, error)
	ReadProducts(model.QueryPagination) ([]model.Product, int64, error)
	ReadProduct(model.Product) (model.Product, error)
	UpdateProduct(model.Product) (model.Product, error)
	DeleteProduct(model.Product) (model.Product, error)

	// Transactions
	CreateTransaction(model.Transaction) (model.Transaction, error)
	CreateTransactions([]model.Transaction) ([]model.Transaction, error)
	ReadTransactions(model.QueryGetTransactions) ([]model.Transaction, int64, error)
	ReadTransactionsBetweenDate(model.QueryGetTransactions) ([]model.Transaction, int64, error)
	ReadTransaction(model.Transaction) (model.Transaction, error)
	UpdateTransaction(model.Transaction) (model.Transaction, error)
	DeleteTransaction(model.Transaction) (model.Transaction, error)

	// Transaction Lines
	CreateTransactionLine(model.TransactionLine) (model.TransactionLine, error)
	CreateTransactionLines([]model.TransactionLine) ([]model.TransactionLine, error)
	ReadTransactionLines() ([]model.TransactionLine, error)
	ReadTransactionLine(model.TransactionLine) (model.TransactionLine, error)
	UpdateTransactionLine(model.TransactionLine) (model.TransactionLine, error)
	DeleteTransactionLine(model.TransactionLine) (model.TransactionLine, error)
	TotalData(model.MetodologiFilter) ([]model.ResponseMetodologi, error)

	// Reports
	ReadReportDashboard() (model.Dashboard, error)
}

type Storage struct {
	db *gorm.DB
}

func NewStorage(c cfg.MySQL, goEnv string) (*Storage, error) {
	var err error

	s := new(Storage)

	var ll logger.LogLevel
	if goEnv == "local" {
		ll = logger.Info
	} else {
		ll = logger.Silent
	}

	s.db, err = gorm.Open(mysql.Open(c.DSN), &gorm.Config{
		Logger: logger.Default.LogMode(ll),
	})
	if err != nil {
		return s, err
	}

	// if err = autoMigrateDB(s); err != nil {
	// 	return nil, err
	// }

	if err = seedDB(s); err != nil {
		return nil, err
	}

	log.Println("MySQL connected")

	return s, nil
}

func autoMigrateDB(s *Storage) error {
	// Migrate the schema
	err := s.db.AutoMigrate(
		&model.User{},
		&model.Category{},
		&model.UOM{},
		&model.Product{},
		&model.Transaction{},
		&model.TransactionLine{},
	)

	return err
}
