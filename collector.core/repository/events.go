package repository

import (
	"context"
	"github.com/go-pg/pg"
	"gitlab.com/faemproject/backend/delivery/collector.core/models"
	"gitlab.com/faemproject/backend/delivery/collector.core/proto"
)

type EventsRepository interface {
	CreateEvent(ctx context.Context, event *models.Event) error
	EventsFilter(ctx context.Context, filter *proto.EventFilter) ([]models.Event, error)
	GetCourierEvents(ctx context.Context, orderUUID string) ([]models.Event, error)
	GetOrderEvents(ctx context.Context, orderUUID string) ([]models.Event, error)
}

func (p *Pg) CreateEvent(ctx context.Context, event *models.Event) error {
	timeout, cancelFunc := context.WithTimeout(ctx, queryTimeout)
	defer cancelFunc()

	_, err := p.Db.ModelContext(timeout, event).Insert()

	return err
}

func (p *Pg) GetCourierEvents(ctx context.Context, courierUUID string) ([]models.Event, error) {
	timeout, cancelFunc := context.WithTimeout(ctx, queryTimeout)
	defer cancelFunc()

	events := new([]models.Event)

	err := p.Db.ModelContext(timeout, events).
		Where("courier_uuid = ?", courierUUID).
		Select()

	return *events, err
}

func (p *Pg) GetOrderEvents(ctx context.Context, orderUUID string) ([]models.Event, error) {
	timeout, cancelFunc := context.WithTimeout(ctx, queryTimeout)
	defer cancelFunc()

	events := new([]models.Event)

	err := p.Db.ModelContext(timeout, events).
		Where("order_uuid = ?", orderUUID).
		Select()

	return *events, err
}

func (p *Pg) EventsFilter(ctx context.Context, filter *proto.EventFilter) ([]models.Event, error) {
	timeout, cancelFunc := context.WithTimeout(ctx, queryTimeout)
	defer cancelFunc()

	event := new([]models.Event)
	query := p.Db.ModelContext(timeout, event)
	if filter.UUID != "" {
		query.Where("(?) <@ event_uuid", pg.Array([]string{filter.UUID}))
	}

	if filter.EventType != "" {
		query.Where("event_type = ?", filter.EventType)
	}
	if err := query.Select(); err != nil {
		return nil, err
	}

	return *event, nil
}
