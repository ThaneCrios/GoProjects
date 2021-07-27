package proto

var States struct {
	Order struct {
		Preparing          Constant `json:"preparing"`
		WaitingForDelivery Constant `json:"waiting_for_delivery"`
		OnTheWay           Constant `json:"on_the_way"`
		Delivered          Constant `json:"delivered"`
	}
	Courier struct {
		DoOrder       Constant `json:"do_order"`
		OnModeration  Constant `json:"on_moderation"`
		Moderated     Constant `json:"moderated"`
		DayOff        Constant `json:"day_off"`
		ReadyForOrder Constant `json:"ready_for_order"`
		OnOrder       Constant `json:"on_order"`
		WaitForOrder  Constant `json:"wait_for_order"`
	}
}

func initStates() {
	States.Courier.DoOrder = "do order"
	States.Courier.OnModeration = "on moderation"
	States.Courier.Moderated = "moderated"
	States.Courier.DayOff = "day off"
	States.Courier.ReadyForOrder = "ready for order"
	States.Courier.OnOrder = "on order"
	States.Courier.WaitForOrder = "wait for order"

	States.Order.Preparing = "preparing"
	States.Order.WaitingForDelivery = "waiting for delivery"
	States.Order.OnTheWay = "on the way"
	States.Order.Delivered = "delivered"
}

type AllVariables struct {
	Couriers map[string]string
	Tasks    map[string]string
	Orders   map[string]string
	Payments map[string]string
}

var Variable AllVariables

func InitLocalisationMaps() {
	Variable.Couriers = make(map[string]string)
	Variable.Orders = make(map[string]string)
	Variable.Tasks = make(map[string]string)
	Variable.Payments = make(map[string]string)
}

func LocalisationInit() {
	Variable.Tasks[Variables.Tasks.TaskType.PickUp] = "забрать"
	Variable.Tasks[Variables.Tasks.TaskType.Deliver] = "доставить"
	Variable.Tasks[Variables.Tasks.TaskType.DeliverWaitDeliverBack] = "доставить, подождать, доставить обратно"
	Variable.Tasks[Variables.Tasks.TaskType.PickUpAndPay] = "забрать и оплатить"

	Variable.Tasks[Variables.Tasks.TaskState.Created] = "создан"
	Variable.Tasks[Variables.Tasks.TaskState.InProgress] = "выполняется"
	Variable.Tasks[Variables.Tasks.TaskState.AwaitForExecution] = "ожидает выполнения"

	Variable.Orders[Variables.Order.OrderType.PickUpPayAndDeliverWaitDeliverBack] = "забрать, оплатить, доставить, подождать и вернуть"
	Variable.Orders[Variables.Order.OrderType.PickUpAndDeliver] = "забрать и доставить"
	Variable.Orders[Variables.Order.OrderType.PickUpPayAndDeliver] = "забрать, оплатить и доставить"

	Variable.Orders[Variables.Order.OrderState.Created] = "создан"
	Variable.Orders[Variables.Order.OrderState.InProgress] = "выполняется"
	Variable.Orders[Variables.Order.OrderState.AwaitForExecution] = "ожидает выполнения"

	Variable.Payments[Variables.Payment.PaymentType.Cash] = "наличные"
	Variable.Payments[Variables.Payment.PaymentType.BankCard] = "банковская карта"

	Variable.Payments[Variables.Payment.PaymentStatus.AwaitsPayment] = "ожидает оплаты"
	Variable.Payments[Variables.Payment.PaymentStatus.PaymentNotRequired] = "оплата не требуется"
	Variable.Payments[Variables.Payment.PaymentStatus.Paid] = "оплачен"

	Variable.Couriers[Variables.Courier.CourierType.OnFoot] = "пешком"
	Variable.Couriers[Variables.Courier.CourierType.Machine] = "на машине"

	Variable.Couriers[Variables.Courier.CourierState.DayOff] = "выходной"
	Variable.Couriers[Variables.Courier.CourierState.Moderated] = "промодерирован"
	Variable.Couriers[Variables.Courier.CourierState.OnModeration] = "на модерации"
	Variable.Couriers[Variables.Courier.CourierState.LookingForOrder] = "ожидает следующего заказа"
	Variable.Couriers[Variables.Courier.CourierState.ReadyToTakeOrder] = "готов принимать заказы"
	Variable.Couriers[Variables.Courier.CourierState.OnOrder] = "на заказе"
}
