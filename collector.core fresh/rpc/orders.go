package rpc

import (
	"context"
	"github.com/pkg/errors"

	"gitlab.com/faemproject/backend/core/shared/tools"
	collectorModels "gitlab.com/faemproject/backend/eda/eda.core/services/collector/models"
	orderModels "gitlab.com/faemproject/backend/eda/eda.core/services/orders/proto"
)

func (r *RPC) SetOrderState(ctx context.Context, params orderModels.OrdersStateOptions, uuid string) error {

	err := tools.RPC("PUT", r.Handler.Config.Services.Orders+"/orders/"+uuid+"/state", nil, params, nil)
	if err != nil {
		return errors.Wrap(err, "fail to change order state")
	}

	return nil
}

func (r *RPC) UpdateOrder(ctx context.Context, order collectorModels.OrderUpdated) error {

	err := tools.RPC("PUT", r.Handler.Config.Services.Orders+"/orders/"+order.UUID, nil, order, nil)
	if err != nil {
		return errors.Wrap(err, "fail to change order state")
	}

	return nil
}
