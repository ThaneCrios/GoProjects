package rabbit

type (
	RoutingKey string
)

const (
	// channels
	OrderStatesChannel = "orderStatesChannel"

	// exchanges
	OrderExchange = "eda.events.order"

	// queues
	CreatedOrdersQueues = "collector.orders.created"

	CreatedOrdersConsumer = "collector.createdorders.consumer"
)
