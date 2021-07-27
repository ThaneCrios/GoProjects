package subscriber

import (
	"context"
	"gitlab.com/faemproject/backend/delivery/collector.core/models"
	"gitlab.com/faemproject/backend/delivery/collector.core/pkg/rabbit"
	"gitlab.com/faemproject/backend/delivery/collector.core/proto"

	"github.com/korovkin/limiter"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"

	"gitlab.com/faemproject/backend/core/shared/lang"
	"gitlab.com/faemproject/backend/core/shared/logs"
	"gitlab.com/faemproject/backend/core/shared/structures/errpath"
)

const (
	maxNewOrderStatesAllowed = 100
)

func (s *Subscriber) HandleNewOrderState(ctx context.Context, msg amqp.Delivery) error {
	// Decode incoming message
	var order models.Order
	if err := s.Encoder.Decode(msg.Body, &order); err != nil {
		return errors.Wrap(err, "failed to decode order")
	}

	log := logs.Eloger.WithFields(logrus.Fields{
		"event":      "handling new order state from rabbit",
		"order-uuid": order.UUID,
	})
	//order.UUID = s.Handler.RAM.IDs.GenUUID()
	//order.ID = s.Handler.RAM.IDs.SliceUUID(order.UUID)
	_, err := s.Handler.DB.CreateOrder(ctx, &order)
	if err != nil {
		log.WithField("reason", "failed to create order and tasks in DB").Error(err)
		return errors.Wrap(err, "fail to create order")
	}
	log.Info(order)

	return nil
}

func (s *Subscriber) initOrderState() error {
	receiverOrderStatesChannel, err := s.Rabbit.GetReceiver(rabbit.OrderStatesChannel)
	if err != nil {
		return errors.Wrapf(err, "failed to get a receiver channel")
	}

	// Declare an exchange first
	err = receiverOrderStatesChannel.ExchangeDeclare(
		rabbit.OrderExchange, // name
		"topic",              // type
		true,                 // durable
		false,                // auto-deleted
		false,                // internal
		false,                // no-wait
		nil,                  // arguments
	)
	if err != nil {
		return errors.Wrap(err, "failed to declare order state exchange")
	}

	// объявляем очередь для получения статусов заказов
	queue, err := receiverOrderStatesChannel.QueueDeclare(
		rabbit.CreatedOrdersQueues, // name
		true,                       // durable
		false,                      // delete when unused
		false,                      // exclusive
		false,                      // no-wait
		nil,                        // arguments
	)
	if err != nil {
		return errors.Wrap(err, "failed to declare order state queue")
	}

	// биндим очередь для получения статусов заказов
	err = receiverOrderStatesChannel.QueueBind(
		queue.Name,                       // queue name
		"state."+proto.OrderStateCreated, // routing key
		rabbit.OrderExchange,             // exchange
		false,
		nil,
	)
	if err != nil {
		return errors.Wrap(err, "failed to bind order state queue")
	}

	msgs, err := receiverOrderStatesChannel.Consume(
		queue.Name,                   // queue
		rabbit.CreatedOrdersConsumer, // consumer
		true,                         // auto-ack
		false,                        // exclusive
		false,                        // no-local
		false,                        // no-wait
		nil,                          // args
	)
	if err != nil {
		return errors.Wrap(err, "failed to consume from a channel")
	}

	s.wg.Add(1)
	go s.handleNewOrderStates(msgs) // handle incoming messages
	return nil
}

func (s *Subscriber) handleNewOrderStates(messages <-chan amqp.Delivery) {
	defer s.wg.Done()

	limit := limiter.NewConcurrencyLimiter(maxNewOrderStatesAllowed)
	defer limit.Wait()

	for {
		select {
		case <-s.closed:
			return
		case msg := <-messages:
			// Start new goroutine to handle multiple requests at the same time
			limit.Execute(lang.Recover(
				func() {
					if err := s.HandleNewOrderState(context.Background(), msg); err != nil {
						logs.Eloger.Errorln(errpath.Err(err, "failed to handle new order state"))
					}
				},
			))
		}
	}
}
