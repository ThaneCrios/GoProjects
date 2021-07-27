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

var Localisation map[string]string

func LocalisationInit() {
	localisation := make(map[string]string)

	localisation[Variables.Tasks.TaskType.PickUp] = "забрать"
	localisation[Variables.Tasks.TaskType.Deliver] = "доставить"
	localisation[Variables.Tasks.TaskType.DeliverWaitDeliverBack] = "доставить, подождать, доставить обратно"
	localisation[Variables.Tasks.TaskType.PickUpAndPay] = "забрать и оплатить"

	localisation[Variables.Order.OrderType.PickUpPayAndDeliverWaitDeliverBack] = "забрать, оплатить, доставить, подождать, доставить обратнл"
	localisation[Variables.Order.OrderType.PickUpAndDeliver] = "забрать и доставить"
	localisation[Variables.Order.OrderType.PickUpPayAndDeliver] = "забрать, оплатить и доставить"

	localisation[Variables.Order.OrderState.Created] = "создан"
	localisation[Variables.Order.OrderState.InProgress] = "в процессе выполнения"
	localisation[Variables.Order.OrderState.AwaitForExecution] = "ждёт выполнения"

	localisation[Variables.Payment.PaymentType.Cash] = "наличные"
	localisation[Variables.Payment.PaymentType.BankCard] = "банковская карта"

	localisation[Variables.Payment.PaymentStatus.AwaitsPayment] = "ожидает оплаты"
	localisation[Variables.Payment.PaymentStatus.PaymentNotRequired] = "оплата не требуется"
	localisation[Variables.Payment.PaymentStatus.Paid] = "оплачен"

	localisation[Variables.Courier.CourierType.OnFoot] = "пешком"
	localisation[Variables.Courier.CourierType.Machine] = "на машине"

	localisation[Variables.Courier.CourierState.DayOff] = "выходной"
	localisation[Variables.Courier.CourierState.Moderated] = "промодерирован"
	localisation[Variables.Courier.CourierState.OnModeration] = "на модерации"
	localisation[Variables.Courier.CourierState.LookingForOrder] = "ожидает следующего заказа"
	localisation[Variables.Courier.CourierState.ReadyToTakeOrder] = "готов принимать заказы"
	localisation[Variables.Courier.CourierState.OnOrder] = "на заказе"

}
