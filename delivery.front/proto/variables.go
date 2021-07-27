package proto

var Variables struct {
	Courier CourierVars
	Order   OrderVars
	Tasks   TaskVars
	Payment PaymentVars
}

type CourierVars struct {
	CourierType  CourierTypes
	CourierState CourierStates
}

type OrderVars struct {
	OrderType  OrderTypes
	OrderState OrderStates
}

type TaskVars struct {
	TaskType TaskTypes
}

type TaskTypes struct {
	PickUp                 string
	Deliver                string
	PickUpAndPay           string
	DeliverWaitDeliverBack string
}

type PaymentVars struct {
	PaymentType   PaymentTypes
	PaymentStatus PaymentStatuses
}

type PaymentTypes struct {
	Cash     string
	BankCard string
}

type PaymentStatuses struct {
	Paid               string
	AwaitsPayment      string
	PaymentNotRequired string
}

type CourierTypes struct {
	Machine string
	OnFoot  string
}

type CourierStates struct {
	OnOrder          string
	OnModeration     string
	Moderated        string
	DayOff           string
	ReadyToTakeOrder string
	LookingForOrder  string
}

type CancelReasons struct {
	CourierWasLate string
	AnotherReason  string
}

type OrderTypes struct {
	PickUpAndDeliver                   string
	PickUpPayAndDeliver                string
	PickUpPayAndDeliverWaitDeliverBack string
}

type OrderStates struct {
	Created           string
	InProgress        string
	AwaitForExecution string
}

type Services struct {
}

func InitCourierStatesVariables() {
	Variables.Courier.CourierState.OnOrder = "on order"
	Variables.Courier.CourierState.OnModeration = "on moderation"
	Variables.Courier.CourierState.Moderated = "moderated"
	Variables.Courier.CourierState.DayOff = "day off"
	Variables.Courier.CourierState.ReadyToTakeOrder = "ready to take order"
	Variables.Courier.CourierState.LookingForOrder = "looking for order"
}

func InitCourierTypesVariables() {
	Variables.Courier.CourierType.Machine = "machine"
	Variables.Courier.CourierType.OnFoot = "on foot"
}

func InitOrderStatesVariables() {
	Variables.Order.OrderState.Created = "created"
	Variables.Order.OrderState.InProgress = "in progress"
	Variables.Order.OrderState.AwaitForExecution = "await for execution"
}

func InitOrderTypesVariables() {
	Variables.Order.OrderType.PickUpAndDeliver = "pickup and deliver"
	Variables.Order.OrderType.PickUpPayAndDeliver = "pickup, pay and deliver"
	Variables.Order.OrderType.PickUpPayAndDeliverWaitDeliverBack = "pickup, pay, deliver, wait and deliver back"
}

func InitTaskTypesVariables() {
	Variables.Tasks.TaskType.PickUp = "pick up"
	Variables.Tasks.TaskType.Deliver = "deliver"
	Variables.Tasks.TaskType.PickUpAndPay = "pick up and pay"
	Variables.Tasks.TaskType.DeliverWaitDeliverBack = "deliver, wait, deliver back"
}

func InitPaymentTypesVariables() {
	Variables.Payment.PaymentType.BankCard = "bank card"
	Variables.Payment.PaymentType.Cash = "cash"
}

func InitPaymentStatusesVariables() {
	Variables.Payment.PaymentStatus.AwaitsPayment = "awaits payment"
	Variables.Payment.PaymentStatus.Paid = "paid"
	Variables.Payment.PaymentStatus.PaymentNotRequired = "payment not required"
}
