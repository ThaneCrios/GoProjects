package repository

import (
	"context"
)

type VersionerRepository interface {
	GetCurrentVersion(context.Context) (string, error)
}
type DbRepository interface {

}
func (p *Pg) GetCurrentVersion(ctx context.Context) (string, error) {
	return "v3", nil
}
