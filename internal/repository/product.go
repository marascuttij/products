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

	row := p.db.QueryRow("SELECT p.`id`, p.`name`, p.`quantity`, p.`code_value`, p.`is_published`, p.`expiration`, p.`price` FROM `products` AS  `p` WHERE p.`id` = ?", id)

	// serialize the product
	err = row.Scan(&product.ID, &product.Name, &product.Quantity, &product.CodeValue, &product.IsPublished, &product.Expiration, &product.Price)

	// check errors
	if err != nil {
		if err == sql.ErrNoRows {
			err = internal.ErrProductRepositoryNotFound
		}
		return
	}
	return
}
