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

// FindAll returns all products
func (s *ProductDefault) FindAll() (products []internal.Product, err error) {

	// get the products from the repository
	products, err = s.rp.FindAll()

	// check for errors
	if err != nil {
		err = fmt.Errorf("internal server error")
		return
	}
	// return the products
	return

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
	return
}

// Delete deletes a product
func (s *ProductDefault) Delete(id int) (err error) {

	// delete the product from the repository
	err = s.rp.Delete(id)

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
	return
}

// Create creates a new product
func (s *ProductDefault) Create(product *internal.Product) (err error) {

	// validate the warehouse fields
	err = validateProductFields(product)

	// check for errors
	if err != nil {
		return err
	}

	// create the product in the repository
	err = s.rp.Create(product)

	// check for errors
	if err != nil {
		switch {
		case errors.Is(err, internal.ErrProductRepositoryDuplicated):
			err = internal.ErrProductRepositoryDuplicated
		default:
			err = fmt.Errorf("internal server error")
		}
		return
	}
	// return the product
	return
}

// Update updates a product
func (p *ProductDefault) Update(product *internal.Product) (err error) {

	err = p.rp.Update(product)

	// check for errors
	if err != nil {
		switch {
		case errors.Is(err, internal.ErrProductRepositoryDuplicated):
			err = internal.ErrProductRepositoryDuplicated
		default:
			err = internal.ErrInternalServerError

		}
		return
	}
	return
}

// validateProductFields validates the warehouse fields
func validateProductFields(product *internal.Product) (err error) {

	// validate the product name
	if product.Name == "" {
		return fmt.Errorf("%w: name", internal.ErrProductServiceInvalidField)
	}

	// validate the product quantity
	if product.Quantity == 0 {
		return fmt.Errorf("%w: quantity", internal.ErrProductServiceInvalidField)
	}

	// validate the product code_value
	if product.CodeValue == "" {
		return fmt.Errorf("%w: code_value", internal.ErrProductServiceInvalidField)
	}

	// validate the product is_published
	if product.IsPublished == "" {
		return fmt.Errorf("%w: is_published", internal.ErrProductServiceInvalidField)
	}

	// validate the product expiration
	if product.Expiration == "" {
		return fmt.Errorf("%w: expiration", internal.ErrProductServiceInvalidField)
	}

	// validate the product price
	if product.Price == 0 {
		return fmt.Errorf("%w: price", internal.ErrProductServiceInvalidField)
	}

	return nil
}
