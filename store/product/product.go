package product

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	productModel "github.com/danilotadeu/products/model/product"
	"github.com/sirupsen/logrus"
)

// Store is a contract to Product..
//
//go:generate mockgen -destination ../../mock/store/product/product_store_mock.go -package mockStoreProduct . Store
type Store interface {
	SaveProduct(ctx context.Context, Product productModel.ProductDB) (*int64, error)
	Update(ctx context.Context, product productModel.ProductDB) error
	GetOne(ctx context.Context, name string) (*productModel.ProductDB, error)
	GetOneByID(ctx context.Context, id int64) (*productModel.ProductDB, error)
	GetAll(ctx context.Context, page, limit int64, name string) ([]*productModel.ProductDB, error)
	Delete(ctx context.Context, id int64) error
	GetTotalProducts(ctx context.Context) (*int64, error)
}

type storeImpl struct {
	db *sql.DB
}

// NewApp init a Product
func NewStore(db *sql.DB) Store {
	return &storeImpl{
		db: db,
	}
}

func (a *storeImpl) SaveProduct(ctx context.Context, product productModel.ProductDB) (*int64, error) {
	query := fmt.Sprintf("INSERT INTO products(name, quantity) VALUES ('%s','%d')",
		product.Name, product.Quantity)
	res, err := a.db.Exec(query)

	if err != nil {
		logrus.WithFields(logrus.Fields{"trace": "store.product.SaveProduct.Exec"}).Error(err)
		return nil, err
	}

	lastId, err := res.LastInsertId()
	if err != nil {
		logrus.WithFields(logrus.Fields{"trace": "store.product.SaveProduct.LastInsertId"}).Error(err)
		return nil, err
	}

	return &lastId, nil
}

func (a *storeImpl) Update(ctx context.Context, product productModel.ProductDB) error {
	query := fmt.Sprintf("UPDATE products SET name = '%s', quantity = '%d' WHERE id = '%d'", product.Name, product.Quantity, product.ID)
	res, err := a.db.Exec(query)
	if err != nil {
		logrus.WithFields(logrus.Fields{"trace": "store.product.Delete.Exec_1"}).Error(err)
		return err

	}
	_, err = res.RowsAffected()
	if err != nil {
		logrus.WithFields(logrus.Fields{"trace": "store.product.Delete.RowsAffected_1"}).Error(err)
		return err
	}

	return nil
}

func (a *storeImpl) GetOne(ctx context.Context, name string) (*productModel.ProductDB, error) {
	res, err := a.db.Query("SELECT * FROM products WHERE deleted_at IS NULL and name = ?", name)
	if err != nil {
		logrus.WithFields(logrus.Fields{"trace": "store.product.GetOne.Query"}).Error(err)
		return nil, err
	}
	defer res.Close()

	if res.Next() {
		var Product productModel.ProductDB
		err := res.Scan(
			&Product.ID,
			&Product.Name,
			&Product.Quantity,
			&Product.CreatedAt,
			&Product.DeletedAt,
		)
		if err != nil {
			logrus.WithFields(logrus.Fields{"trace": "store.product.GetOne.Scan"}).Error(err)
			return nil, err
		}

		return &Product, nil
	} else {
		return nil, nil
	}
}

func (a *storeImpl) GetOneByID(ctx context.Context, id int64) (*productModel.ProductDB, error) {
	res, err := a.db.Query("SELECT * FROM products WHERE deleted_at IS NULL and id = ?", id)
	if err != nil {
		logrus.WithFields(logrus.Fields{"trace": "store.product.GetOneByID.Query"}).Error(err)
		return nil, err
	}
	defer res.Close()

	if res.Next() {
		var Product productModel.ProductDB
		err := res.Scan(
			&Product.ID,
			&Product.Name,
			&Product.Quantity,
			&Product.CreatedAt,
			&Product.DeletedAt,
		)
		if err != nil {
			logrus.WithFields(logrus.Fields{"trace": "store.product.GetOneByID.Scan"}).Error(err)
			return nil, err
		}

		return &Product, nil
	} else {
		return nil, productModel.ErrorProductNotFound
	}
}

func (a *storeImpl) GetAll(ctx context.Context, page, limit int64, name string) ([]*productModel.ProductDB, error) {
	query := `SELECT * FROM products WHERE deleted_at IS NULL`
	params := []interface{}{}
	if len(name) > 0 {
		params = append(params, "%"+name+"%")
		query += ` AND name LIKE ? `
	}

	query += ` LIMIT ? OFFSET ?`
	params = append(params, limit, page)

	res, err := a.db.Query(query, params...)
	if err != nil {
		logrus.WithFields(logrus.Fields{"trace": "store.product.GetAll.Query"}).Error(err)
		return nil, err
	}
	defer res.Close()

	var results []*productModel.ProductDB
	for res.Next() {
		var Product productModel.ProductDB
		err := res.Scan(
			&Product.ID,
			&Product.Name,
			&Product.Quantity,
			&Product.CreatedAt,
			&Product.DeletedAt,
		)
		if err != nil {
			logrus.WithFields(logrus.Fields{"trace": "store.product.GetAll.Scan"}).Error(err)
			return nil, err
		}
		results = append(results, &Product)
	}

	return results, nil
}

func (a *storeImpl) Delete(ctx context.Context, id int64) error {
	query := fmt.Sprintf("UPDATE products SET deleted_at = '%s' WHERE id = '%d'", time.Now().Format("2006-01-02 15:04:05"), id)
	res, err := a.db.Exec(query)
	if err != nil {
		logrus.WithFields(logrus.Fields{"trace": "store.product.Delete.Exec_1"}).Error(err)
		return err

	}
	_, err = res.RowsAffected()
	if err != nil {
		logrus.WithFields(logrus.Fields{"trace": "store.product.Delete.RowsAffected_1"}).Error(err)
		return err
	}

	return nil
}

func (a *storeImpl) GetTotalProducts(ctx context.Context) (*int64, error) {
	res, err := a.db.Query("SELECT COUNT(*) FROM products WHERE deleted_at IS NULL")
	if err != nil {
		logrus.WithFields(logrus.Fields{"trace": "store.product.getTotalProducts.Query"}).Error(err)
		return nil, err
	}
	defer res.Close()

	if res.Next() {
		var Product productModel.ProductsTotal
		err := res.Scan(
			&Product.Total,
		)
		if err != nil {
			logrus.WithFields(logrus.Fields{"trace": "store.product.getTotalProducts.Scan"}).Error(err)
			return nil, err
		}

		return &Product.Total, nil
	} else {
		return nil, productModel.ErrorProductNotFound
	}
}
