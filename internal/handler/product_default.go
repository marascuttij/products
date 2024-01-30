package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"storage/internal"

	"strconv"

	"github.com/bootcamp-go/web/request"
	"github.com/bootcamp-go/web/response"

	"github.com/go-chi/chi/v5"
)

type ProductJSON struct {
	// ID is the unique identifier of the product
	Id int `json:"id"`
	// Name is the name of the product
	Name string `json:"name"`
	// Quantity is the unique identifier of the product
	Quantity int `json:"quantity"`
	// CodeValue is the universal code of the product
	CodeValue string `json:"code_value"`
	// IsPublished is the status of the product
	IsPublished string `json:"is_published"`
	// Expiration is the date of expiration of the product
	Expiration string `json:"expiration"`
	// Price is the price of the product
	Price float64 `json:"price"`
}

type BodyRequestProductJSON struct {
	Name        string  `json:"name"`
	Quantity    int     `json:"quantity"`
	CodeValue   string  `json:"code_value"`
	IsPublished string  `json:"is_published"`
	Expiration  string  `json:"expiration"`
	Price       float64 `json:"price"`
}

type BodyResponseProductJSON struct {
	Id          int     `json:"id"`
	Name        string  `json:"name"`
	Quantity    int     `json:"quantity"`
	CodeValue   string  `json:"code_value"`
	IsPublished string  `json:"is_published"`
	Expiration  string  `json:"expiration"`
	Price       float64 `json:"price"`
}

type ResponseProduct struct {
	Data ProductJSON `json:"data"`
}

type ResponseProductJSON struct {
	Data []ProductJSON `json:"data"`
}

// NewProductDefault creates a new instance of the product handler

func NewProductDefault(sv internal.ProductService) *ProductDefault {
	return &ProductDefault{
		sv: sv,
	}
}

type ProductDefault struct {
	// sv is the service used by the handler
	sv internal.ProductService
}

// GetAll returns all products
func (h *ProductDefault) GetAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		//process
		products, err := h.sv.FindAll()
		if err != nil {
			response.Error(w, http.StatusInternalServerError, "internal server error")
			return
		}

		// serealize to json
		productsJSON := make([]ProductJSON, 0)
		for _, product := range products {
			productsJSON = append(productsJSON, ProductJSON{
				Id:          product.ID,
				Name:        product.Name,
				Quantity:    product.Quantity,
				CodeValue:   product.CodeValue,
				IsPublished: product.IsPublished,
				Expiration:  product.Expiration,
				Price:       product.Price,
			})
		}

		//return response
		response.JSON(w, http.StatusOK, ResponseProductJSON{
			Data: productsJSON,
		})
	}
}

// GetByID returns a product
func (h *ProductDefault) GetByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// get id from url and convert to int
		id, err := strconv.Atoi(chi.URLParam(r, "id"))

		// check for errors
		if err != nil {
			response.JSON(w, http.StatusBadRequest, map[string]interface{}{
				"message": "failed to convert id to int",
				"data":    nil,
			})
			return
		}

		// get the product from the service
		product, err := h.sv.FindByID(id)

		// check for errors
		if err != nil {
			switch err {
			case internal.ErrProductRepositoryNotFound:
				response.Error(w, http.StatusNotFound, "product not found")
			default:
				response.Error(w, http.StatusInternalServerError, "internal server error")
			}
			return
		}

		// return response
		response.JSON(w, http.StatusOK, product)
	}
}

// Delete a product
func (h *ProductDefault) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// get id from url and convert to int
		id, err := strconv.Atoi(chi.URLParam(r, "id"))

		// check for errors
		if err != nil {
			response.JSON(w, http.StatusBadRequest, map[string]interface{}{
				"message": "failed to convert id to int",
				"data":    nil,
			})
			return
		}

		// delete the product from the service
		err = h.sv.Delete(id)

		// check for errors
		if err != nil {
			switch err {
			case internal.ErrProductRepositoryNotFound:
				response.Error(w, http.StatusNotFound, "product not found")
			default:
				response.Error(w, http.StatusInternalServerError, "internal server error")
			}
			return
		}

		// return response
		response.JSON(w, http.StatusOK, map[string]interface{}{
			"message": "product deleted successfully",
			"data":    nil,
		})
	}
}

// Create a product
func (h *ProductDefault) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// read the request body to []bytes
		bytes, err := io.ReadAll(r.Body)

		// check for errors
		if err != nil {
			response.JSON(w, http.StatusBadRequest, map[string]any{
				"message": "failed to read request body to []bytes",
			})
			return
		}

		// transform the []bytes to map[string]any
		var bodyMap map[string]any

		// check for errors
		if err := json.Unmarshal(bytes, &bodyMap); err != nil {
			response.JSON(w, http.StatusBadRequest, map[string]any{
				"message": "failed to transform []bytes to map[string]any",
			})
			return
		}

		// validate key exist in bodyMap
		if err := validateExistsKeys(bodyMap, "name", "quantity", "code_value", "is_published", "expiration", "price"); err != nil {
			response.JSON(w, http.StatusUnprocessableEntity, map[string]any{
				"message": "one or more keys have not been sent in the request body",
			})
			return
		}

		var body BodyResponseProductJSON

		// check for errors
		if err := json.Unmarshal(bytes, &body); err != nil {
			response.JSON(w, http.StatusBadRequest, map[string]any{
				"message": "failed to transform []bytes to BodyResponseProductJSON",
			})
			return
		}

		// serialize the body to a product
		product := internal.Product{
			Name:        body.Name,
			Quantity:    body.Quantity,
			CodeValue:   body.CodeValue,
			IsPublished: body.IsPublished,
			Expiration:  body.Expiration,
			Price:       body.Price,
		}

		// create the product in the service
		err = h.sv.Create(&product)

		// check for errors
		if err != nil {
			switch {
			case errors.Is(err, internal.ErrProductServiceInvalidField):
				response.Error(w, http.StatusBadGateway, err.Error())
				return
			case errors.Is(err, internal.ErrProductRepositoryDuplicated):
				response.Error(w, http.StatusConflict, "product_code already exists")
				return
			default:
				response.Error(w, http.StatusInternalServerError, "internal server error")
			}
		}

		// parsing the product to ProductJSON
		data := ProductJSON{
			Id:          product.ID,
			Name:        product.Name,
			Quantity:    product.Quantity,
			CodeValue:   product.CodeValue,
			IsPublished: product.IsPublished,
			Expiration:  product.Expiration,
			Price:       product.Price,
		}

		// create the response
		ProductJSON := ResponseProductJSON{
			Data: []ProductJSON{data},
		}

		// return response
		response.JSON(w, http.StatusCreated, ProductJSON)
	}
}

// Update a product
func (h *ProductDefault) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// get id from url and convert to int
		id, err := strconv.Atoi(chi.URLParam(r, "id"))

		// check for errors
		if err != nil {
			response.JSON(w, http.StatusBadRequest, map[string]any{
				"message": "failed to convert id to int",
				"data":    nil,
			})
			return
		}

		//get the body of the request
		var bodyJSON ProductJSON
		err = request.JSON(r, &bodyJSON)
		if err != nil {
			response.Error(w, http.StatusBadRequest, "invalid body")
			return
		}

		if bodyJSON.Id != id {
			response.JSON(w, http.StatusBadRequest, map[string]any{
				"message": "id in url and body are different",
			})
			return
		}

		// validate previus product exists
		product, err := h.sv.FindByID(id)

		// check for errors
		if err != nil {
			switch {
			case errors.Is(err, internal.ErrProductRepositoryNotFound):
				response.Error(w, http.StatusNotFound, "product not found")
				return
			default:
				response.Error(w, http.StatusInternalServerError, "internal server error")
				return
			}
		}

		// create the reqBody
		reqBody := ProductJSON{
			Id:          product.ID,
			Name:        product.Name,
			Quantity:    product.Quantity,
			CodeValue:   product.CodeValue,
			IsPublished: product.IsPublished,
			Expiration:  product.Expiration,
			Price:       product.Price,
		}

		updateProduct(&reqBody, bodyJSON)

		// create the productUpdate
		productUpdate := internal.Product{
			ID:          reqBody.Id,
			Name:        reqBody.Name,
			Quantity:    reqBody.Quantity,
			CodeValue:   reqBody.CodeValue,
			IsPublished: reqBody.IsPublished,
			Expiration:  reqBody.Expiration,
			Price:       reqBody.Price,
		}

		// validate id in url and body are different

		// check for errors
		err = h.sv.Update(&productUpdate)
		if err != nil {
			switch {
			case errors.Is(err, internal.ErrProductRepositoryNotFound):
				response.Error(w, http.StatusNotFound, "product not found")
				return
			case errors.Is(err, internal.ErrProductRepositoryDuplicated):
				response.Error(w, http.StatusConflict, "product code already exists")
				return
			default:
				response.Error(w, http.StatusInternalServerError, "internal server error")
				return
			}
		}

		//response
		ProductListJSON := ResponseProduct{
			Data: reqBody,
		}

		// return response
		response.JSON(w, http.StatusCreated, ProductListJSON)
	}
}

// function validateExistKeys validates if the keys exist in the map
func validateExistsKeys(mp map[string]any, keys ...string) (err error) {
	for _, key := range keys {
		if _, ok := mp[key]; !ok {
			err = fmt.Errorf("key %s not found", key)
			return
		}
	}
	return
}

func updateProduct(productPersisted *ProductJSON, productUpdated ProductJSON) {
	if productUpdated.Name != "" {
		productPersisted.Name = productUpdated.Name
	}
	if productUpdated.Quantity != 0 {
		productPersisted.Quantity = productUpdated.Quantity
	}
	if productUpdated.CodeValue != "" {
		productPersisted.CodeValue = productUpdated.CodeValue
	}
	if productUpdated.IsPublished != "" {
		productPersisted.IsPublished = productUpdated.IsPublished
	}

	if productUpdated.Expiration != "" {
		productPersisted.Expiration = productUpdated.Expiration
	}
	if productUpdated.Price != 0 {
		productPersisted.Price = productUpdated.Price
	}

}
