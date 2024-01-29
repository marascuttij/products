package service

import (
	"errors"
	"fmt"
	"storage/internal"
)

// NewProductDefault creates a new instance of the product service
func NewProductDefault(rp internal.ProductRepository) *ProductDefault {
	return &ProductDefault{
		rp: rp,
	}
}

// NewProductDefault is the default implementation of the product service
type ProductDefault struct {
	// rp is the repository used by the service
	rp internal.ProductRepository
}

// FindByID returns a product
func (s *ProductDefault) FindByID(id int) (product internal.Product, err error) {

	// get the product from the repository
	product, err = s.rp.FindByID(id)

	// check for errors
	if err != nil {
		switch {
		case errors.Is(err, internal.ErrProductRepositoryNotFound):
			err = internal.ErrProductRepositoryNotFound
		default:
			err = fmt.Errorf("internal server error")
		}
		return
	}

	// return the product
	return product, nil
}
