package internal

import "errors"

// Product is a struct that contains the product's information
type Product struct {
	// ID is the unique identifier of the product
	ID int
	// Name is the name of the product
	Name string
	// Quantity is the unique identifier of the product
	Quantity int
	// CodeValue is the universal code of the product
	CodeValue string
	// IsPublished is the status of the product
	IsPublished string
	// Expiration is the date of expiration of the product
	Expiration string
	// Price is the price of the product
	Price float64
}

var (
	// ErrProductRepositoryNotFound is the error returned when the product is not found
	ErrProductRepositoryNotFound = errors.New("repository: product not found")
	// ErrProductRepositoryDuplicated is the error returned when the product already exists
	ErrProductRepositoryDuplicated = errors.New("repository: product already exists")
	// ErrProductRepositoryInvalidField is the error returned when the product has an invalid field
	ErrProductServiceInvalidField = errors.New("service: invalid field")
	// ErrInternalServerError is the error returned when an internal server error occurs
	ErrInternalServerError = errors.New("internal server error")
)

// ProductRepository is an interface that contains the methods that the product repository should support
type ProductRepository interface {
	// FindByID returns the product with the given ID
	FindByID(id int) (Product, error)
	// FindAll returns all the products
	FindAll() ([]Product, error)
	// Delete deletes the product with the given ID
	Delete(id int) error
	// Create creates a new product
	Create(product *Product) error
	// Update updates the product with the given ID
	Update(product *Product) error
}

// ProductService is an interface that contains the methods that the product service should support
type ProductService interface {
	// FindByID returns the product with the given ID
	FindByID(id int) (Product, error)
	// FindAll returns all the products
	FindAll() ([]Product, error)
	// Delete deletes the product with the given ID
	Delete(id int) error
	// Create creates a new product
	Create(product *Product) error
	// Update updates the product with the given ID
	Update(product *Product) error
}
