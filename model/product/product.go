package product

import (
	"errors"
	"time"

	genericModel "github.com/danilotadeu/products/model/generic"
)

var ErrorProductNotFound = errors.New("product not found")

type ProductDB struct {
	ID        int64      `json:"id"`
	Name      string     `json:"name" validate:"required"`
	Quantity  int64      `json:"quantity" validate:"required"`
	CreatedAt time.Time  `json:"created_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

type ProductsTotal struct {
	Total int64 `json:"total"`
}

type ResponseProducts struct {
	Data               []*ProductDB            `json:"data"`
	ResponsePagination genericModel.Pagination `json:"pagination"`
}
