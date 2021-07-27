package repository

import (
	"context"
	"fmt"
	"github.com/go-pg/pg"
	"time"
)

const (
	queryTimeout = time.Second * 10
)

type Pg struct {
	Db *pg.DB
}


func (p *Pg) CheckExistsByUUID(ctx context.Context, model interface{}, uuid string) error {
	timeout, cancelFunc := context.WithTimeout(ctx, queryTimeout)
	defer cancelFunc()

	exists, _ := p.Db.ModelContext(timeout, model).
		Column("uuid").
		Where("uuid = ?", uuid).
		Exists()

	if !exists {
		return fmt.Errorf("record with uuid=%s not exists", uuid)
	}

	return nil
}
