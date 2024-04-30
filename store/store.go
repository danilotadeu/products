package store

import (
	"database/sql"

	"github.com/danilotadeu/products/store/product"
	"github.com/sirupsen/logrus"

	_ "github.com/go-sql-driver/mysql"
)

// Container ...
type Container struct {
	Product product.Store
}

// Register store container
func Register(db *sql.DB) *Container {
	container := &Container{
		Product: product.NewStore(db),
	}

	logrus.WithFields(logrus.Fields{"trace": "store"}).Infof("Registered - Store")
	return container
}
