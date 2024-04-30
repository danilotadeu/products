package product

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/danilotadeu/products/app"
	errorsP "github.com/danilotadeu/products/model/errors_handler"
	genericModel "github.com/danilotadeu/products/model/generic"
	productModel "github.com/danilotadeu/products/model/product"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type apiImpl struct {
	apps      *app.Container
	validator *validator.Validate
}

// NewAPI planet function..
func NewAPI(g fiber.Router, apps *app.Container, validate *validator.Validate) {
	api := apiImpl{
		apps:      apps,
		validator: validate,
	}

	g.Get("/", api.products)
	g.Get("/:id", api.product)
	g.Delete("/:id", api.productDelete)
	g.Post("/", api.productCreate)
	g.Put("/:id", api.productUpdate)
}

// CreateProduct godoc
// @Summary      Endpoint to create products
// @Description  Endpoint to create products
// @Tags         products
// @Accept       json
// @Produce      json
// @Param product   body productModel.ProductDB true "Request Product"
// @Success      200  {object}  productModel.ProductDB
// @Failure      400  {object}  errorsP.ErrorsResponse
// @Failure      404  {object}  errorsP.ErrorsResponse
// @Failure      500  {object}  errorsP.ErrorsResponse
// @Router       /api/products [post]
// productCreate is a handle to create products
func (p *apiImpl) productCreate(c *fiber.Ctx) error {
	ctx := c.Context()
	request := productModel.ProductDB{}
	if err := c.BodyParser(&request); err != nil {
		logrus.WithFields(logrus.Fields{"trace": "api.product.product.create.BodyParser"}).Error(err)
		return c.Status(http.StatusBadRequest).JSON(errorsP.ErrorsResponse{
			Message: err.Error(),
		})
	}

	err := p.validator.Struct(request)
	if err != nil {
		logrus.WithFields(logrus.Fields{"trace": "api.product.product.create.validator.Struct"}).Error(err)
		return c.Status(http.StatusBadRequest).JSON(errorsP.ErrorsResponse{
			Message: err.Error(),
		})
	}

	result, err := p.apps.Product.SaveProduct(ctx, request)
	if err != nil {
		logrus.WithFields(logrus.Fields{"trace": "api.product.product.create.Create"}).Error(err)
		return c.Status(http.StatusInternalServerError).JSON(errorsP.ErrorsResponse{
			Message: "aconteceu um erro interno",
		})
	}

	return c.Status(http.StatusOK).JSON(productModel.ProductDB{ID: *result})
}

// UpdateProduct godoc
// @Summary      Endpoint to update products
// @Description  Endpoint to update products
// @Tags         products
// @Accept       json
// @Produce      json
// @Param product   body productModel.ProductDB true "Request Product"
// @Success      200  {object}  productModel.ProductDB
// @Failure      400  {object}  errorsP.ErrorsResponse
// @Failure      404  {object}  errorsP.ErrorsResponse
// @Failure      500  {object}  errorsP.ErrorsResponse
// @Router       /api/products/{id} [put]
// productUpdate is a handle to update products
func (p *apiImpl) productUpdate(c *fiber.Ctx) error {
	ctx := c.Context()
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		logrus.WithFields(logrus.Fields{"trace": "api.product.product.productUpdate.ParseInt"}).Error(err)
		return c.Status(http.StatusBadRequest).JSON(errorsP.ErrorsResponse{
			Message: "aconteceu um erro interno",
		})
	}

	request := productModel.ProductDB{}
	if err := c.BodyParser(&request); err != nil {
		logrus.WithFields(logrus.Fields{"trace": "api.product.product.productUpdate.BodyParser"}).Error(err)
		return c.Status(http.StatusBadRequest).JSON(errorsP.ErrorsResponse{
			Message: err.Error(),
		})
	}

	err = p.validator.Struct(request)
	if err != nil {
		logrus.WithFields(logrus.Fields{"trace": "api.product.product.productUpdate.validator.Struct"}).Error(err)
		return c.Status(http.StatusBadRequest).JSON(errorsP.ErrorsResponse{
			Message: err.Error(),
		})
	}

	request.ID = id
	err = p.apps.Product.UpdateProduct(ctx, request)
	if err != nil {
		logrus.WithFields(logrus.Fields{"trace": "api.product.product.productUpdate.UpdateProduct"}).Error(err)
		return c.Status(http.StatusInternalServerError).JSON(errorsP.ErrorsResponse{
			Message: "aconteceu um erro interno",
		})
	}

	return c.Status(http.StatusOK).JSON(productModel.ProductDB{ID: request.ID})
}

// ShowProduct godoc
// @Summary      Show a product
// @Description  get product by ID
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Product ID"
// @Success      200  {object}  productModel.ProductDB
// @Failure      400  {object}  errorsP.ErrorsResponse
// @Failure      404  {object}  errorsP.ErrorsResponse
// @Failure      500  {object}  errorsP.ErrorsResponse
// @Router       /api/products/{id} [get]
func (p *apiImpl) product(c *fiber.Ctx) error {
	id := c.Params("id")
	iid, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		logrus.WithFields(logrus.Fields{"trace": "api.product.product.ParseInt"}).Error(err)
		return c.Status(http.StatusBadRequest).JSON(errorsP.ErrorsResponse{
			Message: "Por favor envie o id",
		})
	}

	ctx := c.Context()
	planet, err := p.apps.Product.GetOneByID(ctx, iid)
	if err != nil {
		logrus.WithFields(logrus.Fields{"trace": "api.product.product.GetOneByID"}).Error(err)
		if errors.Is(err, productModel.ErrorProductNotFound) {
			return c.Status(http.StatusNotFound).JSON(errorsP.ErrorsResponse{
				Message: fmt.Sprintf("Planeta (%d) não encontrado", iid),
			})
		}
		return c.Status(http.StatusInternalServerError).JSON(errorsP.ErrorsResponse{
			Message: "Aconteceu um erro interno..",
		})
	}

	return c.Status(http.StatusOK).JSON(planet)
}

// DeleteProduct godoc
// @Summary      Delete a products
// @Description  delete products by ID
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Product ID"
// @Success      204
// @Failure      400  {object}  errorsP.ErrorsResponse
// @Failure      404  {object}  errorsP.ErrorsResponse
// @Failure      500  {object}  errorsP.ErrorsResponse
// @Router       /api/products/{id} [delete]
func (p *apiImpl) productDelete(c *fiber.Ctx) error {
	id := c.Params("id")
	iid, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		logrus.WithFields(logrus.Fields{"trace": "api.product.productDelete.ParseInt"}).Error(err)
		return c.Status(http.StatusBadRequest).JSON(errorsP.ErrorsResponse{
			Message: "Por favor envie o id",
		})
	}

	ctx := c.Context()
	err = p.apps.Product.Delete(ctx, iid)
	if err != nil {
		logrus.WithFields(logrus.Fields{"trace": "api.product.productDelete.Delete"}).Error(err)
		if errors.Is(err, productModel.ErrorProductNotFound) {
			return c.Status(http.StatusNotFound).JSON(errorsP.ErrorsResponse{
				Message: fmt.Sprintf("Planeta (%d) não encontrado", iid),
			})
		}
		return c.Status(http.StatusInternalServerError).JSON(errorsP.ErrorsResponse{
			Message: "Aconteceu um erro interno..",
		})
	}

	return c.Status(http.StatusNoContent).JSON(true)
}

// ListProducts godoc
// @Summary      List products
// @Description  get products
// @Tags         products
// @Accept       json
// @Produce      json
// @Param page query int false "page"
// @Param limit query int false "limit"
// @Param name query string false "name"
// @Success      200  {object}  productModel.ResponseProducts
// @Failure      400  {object}  errorsP.ErrorsResponse
// @Failure      404  {object}  errorsP.ErrorsResponse
// @Failure      500  {object}  errorsP.ErrorsResponse
// @Router       /api/products [get]
func (p *apiImpl) products(c *fiber.Ctx) error {
	ctx := c.Context()

	limit := c.Query("limit")
	var ilimit int64 = 10
	if len(limit) > 0 {
		limitConv, err := strconv.ParseInt(limit, 10, 64)
		if err != nil {
			logrus.WithFields(logrus.Fields{"trace": "api.product.products.ParseInt.limit"}).Error(err)
			return c.Status(http.StatusBadRequest).JSON(errorsP.ErrorsResponse{
				Message: "Por favor envie o limit corretamente.",
			})
		}
		ilimit = limitConv
	}

	page := c.Query("page")
	var ipage int64

	if len(page) > 0 {
		pageConv, err := strconv.ParseInt(page, 10, 64)
		if err != nil {
			logrus.WithFields(logrus.Fields{"trace": "api.product.products.ParseInt.page"}).Error(err)
			return c.Status(http.StatusBadRequest).JSON(errorsP.ErrorsResponse{
				Message: "Por favor envie o page corretamente.",
			})
		}
		ipage = pageConv
	}

	name := c.Query("name")

	planets, err := p.apps.Product.GetAllProducts(ctx, ipage, ilimit, name)
	if err != nil {
		logrus.WithFields(logrus.Fields{"trace": "api.product.products.GetAllPlanets"}).Error(err)
		if errors.Is(err, productModel.ErrorProductNotFound) {
			return c.Status(http.StatusNotFound).JSON(errorsP.ErrorsResponse{
				Message: "Dados nao encontrados",
			})
		}

		return c.Status(http.StatusInternalServerError).JSON(errorsP.ErrorsResponse{
			Message: "Aconteceu um erro interno..",
		})
	}

	nextPage, previousPage := genericModel.MakePagination(ipage)

	_, err = p.apps.Product.GetAllProducts(ctx, *nextPage, ilimit, name)
	if err != nil {
		if !errors.Is(err, productModel.ErrorProductNotFound) {
			logrus.WithFields(logrus.Fields{"trace": "api.product.products.GetAllPlanets_1"}).Error(err)
			return c.Status(http.StatusInternalServerError).JSON(errorsP.ErrorsResponse{
				Message: "Aconteceu um erro interno..",
			})
		}
		nextPage = nil
	}

	total, err := p.apps.Product.GetTotalProducts(ctx)
	if err != nil {
		logrus.WithFields(logrus.Fields{"trace": "api.product.products.GetTotalPlanets"}).Error(err)
		return c.Status(http.StatusInternalServerError).JSON(errorsP.ErrorsResponse{
			Message: "Aconteceu um erro interno..",
		})
	}

	return c.Status(http.StatusOK).JSON(productModel.ResponseProducts{
		Data: planets,
		ResponsePagination: genericModel.Pagination{
			Count:        *total,
			NextPage:     nextPage,
			PreviousPage: previousPage,
		},
	})
}
