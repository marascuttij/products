package handler

import (
	"encoding/json"
	"net/http"
	"storage/internal"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type ProductJSON struct {
	Id          int     `json:"id"`
	Name        string  `json:"name"`
	Quantity    int     `json:"quantity"`
	CodeValue   int     `json:"codeValue"`
	IsPublished bool    `json:"isPublished"`
	Expiration  string  `json:"expiration"`
	Price       float64 `json:"price"`
}

type BodyRequestProductJSON struct {
	Quantity    int     `json:"quantity"`
	CodeValue   int     `json:"codeValue"`
	IsPublished bool    `json:"isPublished"`
	Expiration  string  `json:"expiration"`
	Price       float64 `json:"price"`
}

type ResponseProduct struct {
	Data []ProductJSON `json:"data"`
}

type ResponseProductJSON struct {
	Data ProductJSON `json:"data"`
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

// GetByID returns a product
func (h *ProductDefault) GetByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// get id from url and convert to int

		id, err := strconv.Atoi(chi.URLParam(r, "id"))

		// check for errors
		if err != nil {
			JSON(w, http.StatusBadRequest, map[string]interface{}{
				"message": "failed to convert id to int",
				"data":    nil,
			})
			return
		}

		// check for errors
		if err != nil {
			// Validation failed, return the error
			JSON(w, http.StatusBadRequest, map[string]interface{}{
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
			//	response.Error(w, http.StatusNotFound, "product not found")
			default:
				//	response.Error(w, http.StatusInternalServerError, "internal server error")
			}
			return
		}

		// serialize the product to json
		data := ProductJSON{
			Id: product.ID,
		}

		// create the response
		JSON(w, http.StatusOK, map[string]any{
			"data": data,
		})
	}
}

// JSON writes json response
func JSON(w http.ResponseWriter, code int, body any) {
	// check body
	if body == nil {
		w.WriteHeader(code)
		return
	}

	// marshal body
	bytes, err := json.Marshal(body)
	if err != nil {
		// default error
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// set header (before code due to it sets by default "text/plain")
	w.Header().Set("Content-Type", "application/json")

	// set status code
	w.WriteHeader(code)

	// write body
	w.Write(bytes)
}
