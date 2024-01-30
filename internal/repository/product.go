package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"storage/internal"

	"github.com/go-sql-driver/mysql"
)

// NewSellerMysql creates a new instance of the seller repository
func NewProductMysql(db *sql.DB) *ProductMysql {
	return &ProductMysql{db}
}

type ProductMysql struct {
	db *sql.DB
}

func (p *ProductMysql) FindAll() (products []internal.Product, err error) {
	// query

	rows, err := p.db.Query("SELECT p.`id`, p.`name`, p.`quantity`, p.`code_value`, p.`is_published`, p.`expiration`, p.`price` FROM `products` AS  `p`")
	if err != nil {
		return
	}
	defer rows.Close()

	// serialize the products
	for rows.Next() {
		var product internal.Product
		err = rows.Scan(&product.ID, &product.Name, &product.Quantity, &product.CodeValue, &product.IsPublished, &product.Expiration, &product.Price)
		if err != nil {
			return
		}
		products = append(products, product)
	}
	err = rows.Err()
	if err != nil {
		return
	}
	return
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
			return
		}
		return
	}
	return
}

func (p *ProductMysql) Delete(id int) (err error) {
	// query
	_, err = p.db.Exec("DELETE FROM `products` WHERE `id` = ?", id)
	if err != nil {
		return
	}
	return
}

func (p *ProductMysql) Create(product *internal.Product) (err error) {
	// execute the query
	result, err := p.db.Exec("INSERT INTO `products` (`name`, `quantity`, `code_value`, `is_published`, `expiration`, `price`) VALUES (?, ?, ?, ?, ?, ?)", (*product).Name, (*product).Quantity, (*product).CodeValue, (*product).IsPublished, (*product).Expiration, (*product).Price)

	if err != nil {
		var mySqlErr *mysql.MySQLError
		if errors.As(err, &mySqlErr) {
			switch mySqlErr.Number {
			case 1062:
				err = internal.ErrProductRepositoryDuplicated
			default:
				err = fmt.Errorf("internal server error")
			}
			return
		}
		return
	}

	// get the last inserted id
	id, err := result.LastInsertId()
	if err != nil {
		return
	}

	// set the id of the warehouse
	(*product).ID = int(id)

	return

	return
}

func (p *ProductMysql) Update(product *internal.Product) (err error) {
	// execute the query

	//_, err = p.db.Exec("UPDATE `products` AS `p` SET p.`name` = ?, p.`quantity` = ?, p.`code_value` = ?, p.`is_published` = ?, p.`expiration` = ?, p.`price` = ? WHERE p.`id` = ?", (*product).Name, (*product).Quantity, (*product).CodeValue, (*product).IsPublished, (*product).Expiration, (*product).Price, (*product).ID)

	_, err = p.db.Exec("UPDATE `products` AS `p` SET p.`name` = ?, p.`quantity` = ?, p.`code_value` = ?, p.`is_published` = ?, p.`expiration` = ?, p.`price` = ? WHERE p.`id` = ?", (*product).Name, (*product).Quantity, (*product).CodeValue, (*product).IsPublished, (*product).Expiration, (*product).Price, (*product).ID)

	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) {
			switch mysqlErr.Number {
			case 1062:
				err = internal.ErrProductRepositoryDuplicated
			default:
				err = internal.ErrInternalServerError
			}
			return
		}
		return
	}
	return
}
