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
	CodeValue int
	// IsPublished is the status of the product
	IsPublished bool
	// Expiration is the date of expiration of the product
	Expiration string
	// Price is the price of the product
	Price float64
}

var (
	ErrProductRepositoryNotFound = errors.New("repository: product not found")

	ErrProductRepositoryDuplicated = errors.New("repository: product already exists")
)

// ProductRepository is an interface that contains the methods that the product repository should support
type ProductRepository interface {
	// FindByID returns the product with the given ID
	FindByID(id int) (Product, error)
}

// ProductService is an interface that contains the methods that the product service should support
type ProductService interface {
	// FindByID returns the product with the given ID
	FindByID(id int) (Product, error)
}
