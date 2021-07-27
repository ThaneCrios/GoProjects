package rpc

import (
	"context"
	"github.com/pkg/errors"
	"gitlab.com/faemproject/backend/core/shared/tools"
	"gitlab.com/faemproject/backend/eda/eda.core/services/collector/models"
)

func (r *RPC) GetProductByUUID(ctx context.Context, uuid string) (*models.Product, error) {
	product := new(models.Product)

	err := tools.RPC("GET", r.Handler.Config.Services.Products+"/products/"+uuid, nil, nil, product)
	if err != nil {
		return nil, errors.Wrap(err, "fail to get user data by uuid")
	}

	return product, nil
}

func (r *RPC) GetProductUUIDByBarCode(ctx context.Context, params map[string]string) ([]models.ProductBarcode, error) {
	products := new([]models.ProductBarcode)

	err := tools.RPC("GET", r.Handler.Config.Services.Products+"/products/barcodes/filter", params, nil, &products)
	if err != nil {
		return nil, errors.Wrap(err, "fail to get product by bar code")
	}

	return *products, nil
}
