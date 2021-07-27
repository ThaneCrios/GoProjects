package proto

var Variables struct {
	Courier CourierVars
	Order   OrderVars
	Tasks   TaskVars
	Payment PaymentVars
}

var VariablesHuman struct {
	Courier CourierVars
	Order   OrderVars
	Tasks   TaskVars
	Payment PaymentVars
}

type CourierVars struct {
	CourierType      CourierTypes
	CourierTypeHuman CourierTypes

	CourierState      CourierStates
	CourierStateHuman CourierStates
}

type OrderVars struct {
	OrderType      OrderTypes
	OrderTypeHuman OrderTypes

	OrderState      OrderStates
	OrderStateHuman OrderStates
}

type TaskVars struct {
	TaskType      TaskTypes
	TaskTypeHuman TaskTypes

	TaskState      TaskStates
	TaskStateHuman TaskStates
}

type TaskTypes struct {
	PickUp                 string
	Deliver                string
	PickUpAndPay           string
	DeliverWaitDeliverBack string
}

type TaskStates struct {
	Created           string
	InProgress        string
	AwaitForExecution string
}

type PaymentVars struct {
	PaymentType      PaymentTypes
	PaymentTypeHuman PaymentTypes

	PaymentStatus      PaymentStatuses
	PaymentStatusHuman PaymentStatuses
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

type ConstsArray struct {
	ConstsBack  string `json:"consts_back"`
	ConstsHuman string `json:"consts_human"`
	Description string `json:"description"`
}

func InitCourierStatesVariables() {
	Variables.Courier.CourierState.OnOrder = "on order"
	Variables.Courier.CourierState.OnModeration = "on moderation"
	Variables.Courier.CourierState.Moderated = "moderated"
	Variables.Courier.CourierState.DayOff = "day off"
	Variables.Courier.CourierState.ReadyToTakeOrder = "ready to take order"
	Variables.Courier.CourierState.LookingForOrder = "looking for order"

	Variables.Courier.CourierStateHuman.OnOrder = "на заказе"
	Variables.Courier.CourierStateHuman.OnModeration = "на модерации"
	Variables.Courier.CourierStateHuman.Moderated = "промодерирован"
	Variables.Courier.CourierStateHuman.DayOff = "выходной"
	Variables.Courier.CourierStateHuman.ReadyToTakeOrder = "готов принимать заказы"
	Variables.Courier.CourierStateHuman.LookingForOrder = "ожидает следующего заказа"
}

func InitCourierTypesVariables() {
	Variables.Courier.CourierType.Machine = "machine"
	Variables.Courier.CourierType.OnFoot = "on foot"

	Variables.Courier.CourierTypeHuman.Machine = "на машине"
	Variables.Courier.CourierTypeHuman.OnFoot = "пешком"
}

func InitOrderStatesVariables() {
	Variables.Order.OrderState.Created = "created"
	Variables.Order.OrderState.InProgress = "in progress"
	Variables.Order.OrderState.AwaitForExecution = "await for execution"

	Variables.Order.OrderStateHuman.Created = "создан"
	Variables.Order.OrderStateHuman.InProgress = "выполняется"
	Variables.Order.OrderStateHuman.AwaitForExecution = "ожидает выполнения"
}

func InitOrderTypesVariables() {
	Variables.Order.OrderType.PickUpPayAndDeliverWaitDeliverBack = "pickup, pay, deliver, wait and deliver back"
	Variables.Order.OrderType.PickUpAndDeliver = "pickup and deliver"
	Variables.Order.OrderType.PickUpPayAndDeliver = "pickup, pay and deliver"

	Variables.Order.OrderTypeHuman.PickUpPayAndDeliverWaitDeliverBack = "забрать, оплатить, доставить, подождать и вернуть"
	Variables.Order.OrderTypeHuman.PickUpAndDeliver = "забрать и доставить"
	Variables.Order.OrderTypeHuman.PickUpPayAndDeliver = "забрать, оплатить и доставить"
}

func InitTaskTypesVariables() {
	Variables.Tasks.TaskType.PickUp = "pick up"
	Variables.Tasks.TaskType.Deliver = "deliver"
	Variables.Tasks.TaskType.PickUpAndPay = "pick up and pay"
	Variables.Tasks.TaskType.DeliverWaitDeliverBack = "deliver, wait, deliver back"

	Variables.Tasks.TaskTypeHuman.PickUp = "забрать"
	Variables.Tasks.TaskTypeHuman.Deliver = "доставить"
	Variables.Tasks.TaskTypeHuman.PickUpAndPay = "забрать и оплатить"
	Variables.Tasks.TaskTypeHuman.DeliverWaitDeliverBack = "доставить, подождать, доставить обратно"
}

func InitTaskStatesVariables() {
	Variables.Tasks.TaskState.Created = "created"
	Variables.Tasks.TaskState.InProgress = "in progress"
	Variables.Tasks.TaskState.AwaitForExecution = "await for execution"

	Variables.Tasks.TaskStateHuman.Created = "создан"
	Variables.Tasks.TaskStateHuman.InProgress = "выполняется"
	Variables.Tasks.TaskStateHuman.AwaitForExecution = "ожидает выполнения"
}

func InitPaymentTypesVariables() {
	Variables.Payment.PaymentType.BankCard = "bank card"
	Variables.Payment.PaymentType.Cash = "cash"

	Variables.Payment.PaymentTypeHuman.BankCard = "банковская карта"
	Variables.Payment.PaymentTypeHuman.Cash = "наличные"
}

func InitPaymentStatusesVariables() {
	Variables.Payment.PaymentStatus.AwaitsPayment = "awaits payment"
	Variables.Payment.PaymentStatus.Paid = "paid"
	Variables.Payment.PaymentStatus.PaymentNotRequired = "payment not required"

	Variables.Payment.PaymentStatusHuman.AwaitsPayment = "ожидает оплаты"
	Variables.Payment.PaymentStatusHuman.Paid = "оплачен"
	Variables.Payment.PaymentStatusHuman.PaymentNotRequired = "оплата не нужна"
}

var CourierStateArray []ConstsArray
var CourierTypeArray []ConstsArray

var OrderTypeArray []ConstsArray
var OrderStateArray []ConstsArray

var TaskTypeArray []ConstsArray
var TaskStateArray []ConstsArray

var PaymentTypeArray []ConstsArray
var PaymentStatusArray []ConstsArray

func InitCourierStateArray() {
	CourierStateArray = append(CourierStateArray, ConstsArray{
		ConstsBack:  Variables.Courier.CourierState.OnOrder,
		ConstsHuman: Variables.Courier.CourierStateHuman.OnOrder,
		Description: "",
	})
	CourierStateArray = append(CourierStateArray, ConstsArray{
		ConstsBack:  Variables.Courier.CourierState.ReadyToTakeOrder,
		ConstsHuman: Variables.Courier.CourierStateHuman.ReadyToTakeOrder,
		Description: "",
	})
	CourierStateArray = append(CourierStateArray, ConstsArray{
		ConstsBack:  Variables.Courier.CourierState.DayOff,
		ConstsHuman: Variables.Courier.CourierStateHuman.DayOff,
		Description: "",
	})
	CourierStateArray = append(CourierStateArray, ConstsArray{
		ConstsBack:  Variables.Courier.CourierState.Moderated,
		ConstsHuman: Variables.Courier.CourierStateHuman.Moderated,
		Description: "",
	})
	CourierStateArray = append(CourierStateArray, ConstsArray{
		ConstsBack:  Variables.Courier.CourierState.OnModeration,
		ConstsHuman: Variables.Courier.CourierStateHuman.OnModeration,
		Description: "",
	})
	CourierStateArray = append(CourierStateArray, ConstsArray{
		ConstsBack:  Variables.Courier.CourierState.LookingForOrder,
		ConstsHuman: Variables.Courier.CourierStateHuman.LookingForOrder,
		Description: "",
	})
}

func InitCourierTypeArray() {
	CourierTypeArray = append(CourierTypeArray, ConstsArray{
		ConstsBack:  Variables.Courier.CourierType.Machine,
		ConstsHuman: Variables.Courier.CourierTypeHuman.Machine,
		Description: "",
	})
	CourierTypeArray = append(CourierTypeArray, ConstsArray{
		ConstsBack:  Variables.Courier.CourierType.OnFoot,
		ConstsHuman: Variables.Courier.CourierTypeHuman.OnFoot,
		Description: "",
	})
}

func InitOrderStateArray() {
	OrderStateArray = append(OrderStateArray, ConstsArray{
		ConstsBack:  Variables.Order.OrderState.Created,
		ConstsHuman: Variables.Order.OrderStateHuman.Created,
		Description: "",
	})
	OrderStateArray = append(OrderStateArray, ConstsArray{
		ConstsBack:  Variables.Order.OrderState.InProgress,
		ConstsHuman: Variables.Order.OrderStateHuman.InProgress,
		Description: "",
	})
	OrderStateArray = append(OrderStateArray, ConstsArray{
		ConstsBack:  Variables.Order.OrderState.AwaitForExecution,
		ConstsHuman: Variables.Order.OrderStateHuman.AwaitForExecution,
		Description: "",
	})
}

func InitOrderTypeArray() {
	OrderTypeArray = append(OrderTypeArray, ConstsArray{
		ConstsBack:  Variables.Order.OrderType.PickUpAndDeliver,
		ConstsHuman: Variables.Order.OrderTypeHuman.PickUpAndDeliver,
		Description: "",
	})
	OrderTypeArray = append(OrderTypeArray, ConstsArray{
		ConstsBack:  Variables.Order.OrderType.PickUpPayAndDeliver,
		ConstsHuman: Variables.Order.OrderTypeHuman.PickUpPayAndDeliver,
		Description: "",
	})
	OrderTypeArray = append(OrderTypeArray, ConstsArray{
		ConstsBack:  Variables.Order.OrderType.PickUpPayAndDeliverWaitDeliverBack,
		ConstsHuman: Variables.Order.OrderTypeHuman.PickUpPayAndDeliverWaitDeliverBack,
		Description: "",
	})
}

func InitTaskTypeArray() {
	TaskTypeArray = append(TaskTypeArray, ConstsArray{
		ConstsBack:  Variables.Tasks.TaskType.Deliver,
		ConstsHuman: Variables.Tasks.TaskTypeHuman.Deliver,
		Description: "",
	})
	TaskTypeArray = append(TaskTypeArray, ConstsArray{
		ConstsBack:  Variables.Tasks.TaskType.PickUp,
		ConstsHuman: Variables.Tasks.TaskTypeHuman.PickUp,
		Description: "",
	})
	TaskTypeArray = append(TaskTypeArray, ConstsArray{
		ConstsBack:  Variables.Tasks.TaskType.PickUpAndPay,
		ConstsHuman: Variables.Tasks.TaskTypeHuman.PickUpAndPay,
		Description: "",
	})
	TaskTypeArray = append(TaskTypeArray, ConstsArray{
		ConstsBack:  Variables.Tasks.TaskType.DeliverWaitDeliverBack,
		ConstsHuman: Variables.Tasks.TaskTypeHuman.DeliverWaitDeliverBack,
		Description: "",
	})
}

func InitTaskStateArray() {
	TaskStateArray = append(TaskStateArray, ConstsArray{
		ConstsBack:  Variables.Tasks.TaskState.Created,
		ConstsHuman: Variables.Tasks.TaskStateHuman.Created,
		Description: "",
	})
	TaskStateArray = append(TaskStateArray, ConstsArray{
		ConstsBack:  Variables.Tasks.TaskState.InProgress,
		ConstsHuman: Variables.Tasks.TaskStateHuman.InProgress,
		Description: "",
	})
	TaskStateArray = append(TaskStateArray, ConstsArray{
		ConstsBack:  Variables.Tasks.TaskState.AwaitForExecution,
		ConstsHuman: Variables.Tasks.TaskStateHuman.AwaitForExecution,
		Description: "",
	})
}

func InitPaymentTypeArray() {
	PaymentTypeArray = append(PaymentTypeArray, ConstsArray{
		ConstsBack:  Variables.Payment.PaymentType.Cash,
		ConstsHuman: Variables.Payment.PaymentTypeHuman.Cash,
		Description: "",
	})
	PaymentTypeArray = append(PaymentTypeArray, ConstsArray{
		ConstsBack:  Variables.Payment.PaymentType.BankCard,
		ConstsHuman: Variables.Payment.PaymentTypeHuman.BankCard,
		Description: "",
	})
}

func InitPaymentStatusArray() {
	PaymentStatusArray = append(PaymentStatusArray, ConstsArray{
		ConstsBack:  Variables.Payment.PaymentStatus.AwaitsPayment,
		ConstsHuman: Variables.Payment.PaymentStatusHuman.AwaitsPayment,
		Description: "",
	})
	PaymentStatusArray = append(PaymentStatusArray, ConstsArray{
		ConstsBack:  Variables.Payment.PaymentStatus.Paid,
		ConstsHuman: Variables.Payment.PaymentStatusHuman.Paid,
		Description: "",
	})
	PaymentStatusArray = append(PaymentStatusArray, ConstsArray{
		ConstsBack:  Variables.Payment.PaymentStatus.PaymentNotRequired,
		ConstsHuman: Variables.Payment.PaymentStatusHuman.PaymentNotRequired,
		Description: "",
	})
}
