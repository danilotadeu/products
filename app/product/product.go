package product

import (
	"context"

	productModel "github.com/danilotadeu/products/model/product"
	"github.com/danilotadeu/products/store"
	"github.com/sirupsen/logrus"
)

//go:generate mockgen -destination ../../mock/app/product/product_app_mock.go -package mockAppProduct . App
type App interface {
	SaveProduct(ctx context.Context, product productModel.ProductDB) (*int64, error)
	UpdateProduct(ctx context.Context, product productModel.ProductDB) error
	GetOneByID(ctx context.Context, id int64) (*productModel.ProductDB, error)
	GetAllProducts(ctx context.Context, page, offset int64, name string) ([]*productModel.ProductDB, error)
	Delete(ctx context.Context, productID int64) error
	GetTotalProducts(ctx context.Context) (*int64, error)
}

type appImpl struct {
	store *store.Container
}

// NewApp init a planet
func NewApp(store *store.Container) App {
	return &appImpl{
		store: store,
	}
}

func (a *appImpl) SaveProduct(ctx context.Context, product productModel.ProductDB) (*int64, error) {
	id, err := a.store.Product.SaveProduct(ctx, product)
	if err != nil {
		logrus.WithFields(logrus.Fields{"trace": "app.product.SaveProduct.Store.Product.SaveProduct"}).Error(err)
		return nil, err
	}

	return id, nil
}

func (a *appImpl) UpdateProduct(ctx context.Context, product productModel.ProductDB) error {
	err := a.store.Product.Update(ctx, product)
	if err != nil {
		logrus.WithFields(logrus.Fields{"trace": "app.product.UpdateProduct.Store.Product.Update"}).Error(err)
		return err
	}

	return nil
}

func (a *appImpl) GetOneByID(ctx context.Context, id int64) (*productModel.ProductDB, error) {
	product, err := a.store.Product.GetOneByID(ctx, id)
	if err != nil {
		logrus.WithFields(logrus.Fields{"trace": "app.product.GetOneByID.Store.Product.GetOneByID"}).Error(err)
		return nil, err
	}

	return product, nil
}

func (a *appImpl) GetAllProducts(ctx context.Context, page, offset int64, name string) ([]*productModel.ProductDB, error) {
	planets, err := a.store.Product.GetAll(ctx, page, offset, name)
	if err != nil {
		logrus.WithFields(logrus.Fields{"trace": "app.product.GetAllPlanets.Store.Planet.GetAll"}).Error(err)
		return nil, err
	}

	if len(planets) == 0 {
		return nil, productModel.ErrorProductNotFound
	}

	planetIDs := make([]int64, len(planets))
	for idx, planet := range planets {
		planetIDs[idx] = planet.ID
	}

	return planets, nil
}

func (a *appImpl) Delete(ctx context.Context, productID int64) error {
	product, err := a.store.Product.GetOneByID(ctx, productID)
	if err != nil {
		logrus.WithFields(logrus.Fields{"trace": "app.product.Delete.Store.Product.GetOneByID"}).Error(err)
		return err
	}

	err = a.store.Product.Delete(ctx, product.ID)
	if err != nil {
		logrus.WithFields(logrus.Fields{"trace": "app.product.Delete.Store.Product.Delete"}).Error(err)
		return err
	}
	return nil
}

func (a *appImpl) GetTotalProducts(ctx context.Context) (*int64, error) {
	total, err := a.store.Product.GetTotalProducts(ctx)
	if err != nil {
		logrus.WithFields(logrus.Fields{"trace": "app.product.GetTotalProducts.Store.Planet.GetTotalProducts"}).Error(err)
		return nil, err
	}
	return total, nil
}
