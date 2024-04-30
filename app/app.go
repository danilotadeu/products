package app

import (
	"github.com/danilotadeu/products/app/product"
	"github.com/danilotadeu/products/store"
	"github.com/sirupsen/logrus"
)

// Container ...
type Container struct {
	Product product.App
}

// Register app container
func Register(store *store.Container) *Container {
	container := &Container{
		Product: product.NewApp(store),
	}

	logrus.WithFields(logrus.Fields{"trace": "app"}).Infof("Registered - App")
	return container
}
