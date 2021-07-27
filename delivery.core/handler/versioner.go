package handler

import (
	"context"
)

func (h *Handler) GetCurrentVersion(ctx context.Context) (string, error) {
	return h.DB.GetCurrentVersion(ctx)
}
