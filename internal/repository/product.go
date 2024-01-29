package repository

import (
	"database/sql"
	"storage/internal"
)

// NewSellerMysql creates a new instance of the seller repository
func NewProductMysql(db *sql.DB) *ProductMysql {
	return &ProductMysql{db}
}

type ProductMysql struct {
	db *sql.DB
}

func (p *ProductMysql) FindByID(id int) (product internal.Product, err error) {
	// query
	query := "SELECT id, name, quantity, code_value, is_published, expiration, price FROM product WHERE id = ?"

	row, err := p.db.Query(query, id)

	// check errors
	if row.Err() != nil {
		err = row.Err()
	}

	// serialize the product
	err = row.Scan(&product.ID, &product.Quantity, &product.CodeValue, &product.IsPublished, &product.Expiration, &product.Price)

	// check errors
	if err != nil {
		if err == sql.ErrNoRows {
			err = internal.ErrProductRepositoryNotFound
			return
		}
	}
	return
}
