package product

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/danilotadeu/products/app"
	mockAppProduct "github.com/danilotadeu/products/mock/app/product"
	productModel "github.com/danilotadeu/products/model/product"
	"github.com/gofiber/fiber/v2"
	"github.com/golang/mock/gomock"
	"gotest.tools/v3/assert"
)

func TestHandlerDelete(t *testing.T) {
	endpoint := "/products/:id"
	cases := map[string]struct {
		InputParamID       string
		ExpectedErr        error
		ExpectedStatusCode int
		PrepareMockApp     func(mockPlanetApp *mockAppProduct.MockApp)
	}{
		"should delete the planet": {
			InputParamID: "1",
			ExpectedErr:  nil,
			PrepareMockApp: func(mockPlanetApp *mockAppProduct.MockApp) {
				mockPlanetApp.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(nil)
			},
			ExpectedStatusCode: http.StatusNoContent,
		},
		"should throw error with parse int": {
			InputParamID: "xpto",
			ExpectedErr:  nil,
			PrepareMockApp: func(mockPlanetApp *mockAppProduct.MockApp) {
			},
			ExpectedStatusCode: http.StatusBadRequest,
		},
		"should return with planet not found": {
			InputParamID: "1",
			ExpectedErr:  nil,
			PrepareMockApp: func(mockPlanetApp *mockAppProduct.MockApp) {
				mockPlanetApp.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(productModel.ErrorProductNotFound)
			},
			ExpectedStatusCode: http.StatusNotFound,
		},
		"should throw error": {
			InputParamID: "1",
			ExpectedErr:  nil,
			PrepareMockApp: func(mockPlanetApp *mockAppProduct.MockApp) {
				mockPlanetApp.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(fmt.Errorf("error"))
			},
			ExpectedStatusCode: http.StatusInternalServerError,
		},
	}
	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			ctrl, ctx := gomock.WithContext(context.Background(), t)
			mockAppProduct := mockAppProduct.NewMockApp(ctrl)
			cs.PrepareMockApp(mockAppProduct)

			h := apiImpl{
				apps: &app.Container{
					Product: mockAppProduct,
				},
			}

			app := fiber.New()
			app.Delete(endpoint, h.productDelete)
			req := httptest.NewRequest(http.MethodDelete, strings.ReplaceAll(endpoint, ":id", cs.InputParamID), nil).WithContext(ctx)
			req.Header.Set("Content-Type", fiber.MIMEApplicationJSON)
			resp, err := app.Test(req, -1)
			if err != nil {
				t.Errorf("Error app.Test: %s", err.Error())
				return
			}

			assert.Equal(t, cs.ExpectedErr, err)
			assert.Equal(t, cs.ExpectedStatusCode, resp.StatusCode)
		})
	}
}

func TestHandlerGetPlanetByID(t *testing.T) {
	endpoint := "/products/:id"
	cases := map[string]struct {
		InputParamID       string
		ExpectedErr        error
		ExpectedStatusCode int
		PrepareMockApp     func(mockPlanetApp *mockAppProduct.MockApp)
	}{
		"should return success with planet": {
			InputParamID: "1",
			ExpectedErr:  nil,
			PrepareMockApp: func(mockPlanetApp *mockAppProduct.MockApp) {
				mockPlanetApp.EXPECT().GetOneByID(gomock.Any(), gomock.Any()).Return(&productModel.ProductDB{
					ID:       1,
					Name:     "Planet 1",
					Quantity: 1,
				}, nil)
			},
			ExpectedStatusCode: http.StatusOK,
		},
		"should throw error with parse int": {
			InputParamID: "xpto",
			ExpectedErr:  nil,
			PrepareMockApp: func(mockPlanetApp *mockAppProduct.MockApp) {
			},
			ExpectedStatusCode: http.StatusBadRequest,
		},
		"should return with planet not found": {
			InputParamID: "1",
			ExpectedErr:  nil,
			PrepareMockApp: func(mockPlanetApp *mockAppProduct.MockApp) {
				mockPlanetApp.EXPECT().GetOneByID(gomock.Any(), gomock.Any()).Return(nil, productModel.ErrorProductNotFound)
			},
			ExpectedStatusCode: http.StatusNotFound,
		},
		"should throw error": {
			InputParamID: "1",
			ExpectedErr:  nil,
			PrepareMockApp: func(mockPlanetApp *mockAppProduct.MockApp) {
				mockPlanetApp.EXPECT().GetOneByID(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("error"))
			},
			ExpectedStatusCode: http.StatusInternalServerError,
		},
	}
	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			ctrl, ctx := gomock.WithContext(context.Background(), t)
			mockAppProduct := mockAppProduct.NewMockApp(ctrl)
			cs.PrepareMockApp(mockAppProduct)

			h := apiImpl{
				apps: &app.Container{
					Product: mockAppProduct,
				},
			}

			app := fiber.New()
			app.Get(endpoint, h.product)
			req := httptest.NewRequest(http.MethodGet, strings.ReplaceAll(endpoint, ":id", cs.InputParamID), nil).WithContext(ctx)
			req.Header.Set("Content-Type", fiber.MIMEApplicationJSON)
			resp, err := app.Test(req, -1)
			if err != nil {
				t.Errorf("Error app.Test: %s", err.Error())
				return
			}

			assert.Equal(t, cs.ExpectedErr, err)
			assert.Equal(t, cs.ExpectedStatusCode, resp.StatusCode)
		})
	}
}

func TestHandlerGetPlanets(t *testing.T) {
	cases := map[string]struct {
		InputPage          string
		InputLimit         string
		ExpectedErr        error
		ExpectedStatusCode int
		PrepareMockApp     func(mockPlanetApp *mockAppProduct.MockApp)
	}{
		"should return success with planet": {
			InputPage:   "1",
			InputLimit:  "10",
			ExpectedErr: nil,
			PrepareMockApp: func(mockPlanetApp *mockAppProduct.MockApp) {
				mockPlanetApp.EXPECT().GetAllProducts(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return([]*productModel.ProductDB{
					{
						ID:       1,
						Name:     "Planet 1",
						Quantity: 1,
					},
				}, nil)
				mockPlanetApp.EXPECT().GetAllProducts(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, productModel.ErrorProductNotFound)
				var total int64 = 1
				mockPlanetApp.EXPECT().GetTotalProducts(gomock.Any()).Return(&total, nil)
			},
			ExpectedStatusCode: http.StatusOK,
		},
		"should throw error when get planets to paginate": {
			InputPage:   "1",
			InputLimit:  "10",
			ExpectedErr: nil,
			PrepareMockApp: func(mockPlanetApp *mockAppProduct.MockApp) {
				mockPlanetApp.EXPECT().GetAllProducts(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return([]*productModel.ProductDB{
					{
						ID:       1,
						Name:     "Planet 1",
						Quantity: 1,
					},
				}, nil)
				mockPlanetApp.EXPECT().GetAllProducts(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("error"))
			},
			ExpectedStatusCode: http.StatusInternalServerError,
		},
		"should throw error when get total planets": {
			InputPage:   "1",
			InputLimit:  "10",
			ExpectedErr: nil,
			PrepareMockApp: func(mockPlanetApp *mockAppProduct.MockApp) {
				mockPlanetApp.EXPECT().GetAllProducts(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return([]*productModel.ProductDB{
					{
						ID:       1,
						Name:     "Planet 1",
						Quantity: 1,
					},
				}, nil)
				mockPlanetApp.EXPECT().GetAllProducts(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return([]*productModel.ProductDB{
					{
						ID:       1,
						Name:     "Planet 1",
						Quantity: 1,
					},
				}, nil)
				mockPlanetApp.EXPECT().GetTotalProducts(gomock.Any()).Return(nil, fmt.Errorf("error"))
			},
			ExpectedStatusCode: http.StatusInternalServerError,
		},
		"should throw error with parse int page": {
			ExpectedErr: nil,
			InputPage:   "xpto",
			PrepareMockApp: func(mockPlanetApp *mockAppProduct.MockApp) {
			},
			ExpectedStatusCode: http.StatusBadRequest,
		},
		"should throw error with parse int limit": {
			ExpectedErr: nil,
			InputLimit:  "xpto",
			PrepareMockApp: func(mockPlanetApp *mockAppProduct.MockApp) {
			},
			ExpectedStatusCode: http.StatusBadRequest,
		},
		"should return with planet not found": {
			ExpectedErr: nil,
			PrepareMockApp: func(mockPlanetApp *mockAppProduct.MockApp) {
				mockPlanetApp.EXPECT().GetAllProducts(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, productModel.ErrorProductNotFound)
			},
			ExpectedStatusCode: http.StatusNotFound,
		},
		"should throw error": {
			ExpectedErr: nil,
			PrepareMockApp: func(mockPlanetApp *mockAppProduct.MockApp) {
				mockPlanetApp.EXPECT().GetAllProducts(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("error"))
			},
			ExpectedStatusCode: http.StatusInternalServerError,
		},
	}
	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			ctrl, ctx := gomock.WithContext(context.Background(), t)
			mockAppProduct := mockAppProduct.NewMockApp(ctrl)
			cs.PrepareMockApp(mockAppProduct)

			h := apiImpl{
				apps: &app.Container{
					Product: mockAppProduct,
				},
			}
			endpoint := "/planets"
			app := fiber.New()
			app.Get(endpoint, h.products)

			if len(cs.InputPage) > 0 && len(cs.InputLimit) > 0 {
				endpoint += "?page=" + cs.InputPage + "&limit=" + cs.InputLimit
			} else if len(cs.InputPage) > 0 {
				endpoint += "?page=" + cs.InputPage
			} else if len(cs.InputLimit) > 0 {
				endpoint += "?limit=" + cs.InputLimit
			}

			req := httptest.NewRequest(http.MethodGet, endpoint, nil).WithContext(ctx)
			req.Header.Set("Content-Type", fiber.MIMEApplicationJSON)
			resp, err := app.Test(req, -1)
			if err != nil {
				t.Errorf("Error app.Test: %s", err.Error())
				return
			}

			assert.Equal(t, cs.ExpectedErr, err)
			assert.Equal(t, cs.ExpectedStatusCode, resp.StatusCode)
		})
	}
}
