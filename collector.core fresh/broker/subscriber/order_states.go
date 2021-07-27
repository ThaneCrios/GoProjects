package subscriber

import (
	"context"
	"gitlab.com/faemproject/backend/eda/eda.core/pkg/rabbit"
	"gitlab.com/faemproject/backend/eda/eda.core/services/collector/models"

	"github.com/korovkin/limiter"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"

	"gitlab.com/faemproject/backend/core/shared/lang"
	"gitlab.com/faemproject/backend/core/shared/logs"
	"gitlab.com/faemproject/backend/core/shared/structures/errpath"
)

func (s *Subscriber) HandleChangeOrderState(ctx context.Context, msg amqp.Delivery) error {
	// Decode incoming message
	var order models.Order
	if err := s.Encoder.Decode(msg.Body, &order); err != nil {
		return errors.Wrap(err, "failed to decode order")
	}

	log := logs.Eloger.WithFields(logrus.Fields{
		"event":      "handling new order state from rabbit",
		"order-uuid": order.UUID,
	})

	_, err := s.Handler.DB.UpdateOrderState(ctx, &order)
	if err != nil {
		log.WithField("reason", "failed to change order state in DB").Error(err)
		return errors.Wrap(err, "failed to change order state")
	}
	log.Info(order)

	return nil
}

func (s *Subscriber) initOrderChangeState() error {
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
		rabbit.ChangeStatesOrdersQueue, // name
		true,                           // durable
		false,                          // delete when unused
		false,                          // exclusive
		false,                          // no-wait
		nil,                            // arguments
	)
	if err != nil {
		return errors.Wrap(err, "failed to declare order state queue")
	}

	// биндим очередь для получения статусов заказов
	err = receiverOrderStatesChannel.QueueBind(
		queue.Name,           // queue name
		"state.*",            // routing key
		rabbit.OrderExchange, // exchange
		false,
		nil,
	)
	if err != nil {
		return errors.Wrap(err, "failed to bind order state queue")
	}

	msgs, err := receiverOrderStatesChannel.Consume(
		queue.Name,                        // queue
		rabbit.ChangeStatesOrdersConsumer, // consumer
		true,                              // auto-ack
		false,                             // exclusive
		false,                             // no-local
		false,                             // no-wait
		nil,                               // args
	)
	if err != nil {
		return errors.Wrap(err, "failed to consume from a channel")
	}

	s.wg.Add(1)
	go s.handleChangeOrderStates(msgs) // handle incoming messages
	return nil
}

func (s *Subscriber) handleChangeOrderStates(messages <-chan amqp.Delivery) {
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
					if err := s.HandleChangeOrderState(context.Background(), msg); err != nil {
						logs.Eloger.Errorln(errpath.Err(err, "failed to handle new order state"))
					}
				},
			))
		}
	}
}
