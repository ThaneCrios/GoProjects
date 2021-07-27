package rpc

import (
	"context"
	"github.com/pkg/errors"
	"gitlab.com/faemproject/backend/core/shared/tools"
	"gitlab.com/faemproject/backend/eda/eda.core/services/collector/models"
)

func (r *RPC) GetUserDataByUUID(ctx context.Context, uuid string) (*models.User, error) {
	userData := new(models.User)

	err := tools.RPC("GET", r.Handler.Config.Services.Orders+"/users/"+uuid, nil, nil, userData)
	if err != nil {
		return nil, errors.Wrap(err, "fail to get user data by uuid")
	}

	return userData, nil
}
