package rpc

import "gitlab.com/faemproject/backend/eda/eda.core/services/collector/handler"

type (
	RPC struct {
		Handler *handler.Handler
	}
)
