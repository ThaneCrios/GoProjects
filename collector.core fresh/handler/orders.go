package handler

import (
	"context"
	"github.com/pkg/errors"
	"gitlab.com/faemproject/backend/eda/eda.core/services/collector/models"
	"gitlab.com/faemproject/backend/eda/eda.core/services/collector/proto"
	orderModels "gitlab.com/faemproject/backend/eda/eda.core/services/orders/proto"
	"math/rand"
)

//GetOrderByUUID возвращает заказ по указанному UUID
func (h *Handler) GetOrderByUUID(ctx context.Context, uuid string) (models.Order, error) {
	order, err := h.DB.GetOrderByUUID(ctx, uuid)
	if err != nil {
		return order, errors.Wrap(err, "fail to get order by uuid")
	}

	return order, nil
}

func (h *Handler) GetDuplicateOrderByUUID(ctx context.Context, uuid string) (models.OrderForCollectorDuplicate, error) {
	order, err := h.DB.GetDuplicateOrderByUUID(ctx, uuid)
	if err != nil {
		return models.OrderForCollectorDuplicate{}, err
	}

	return order, nil
}

func (h *Handler) GetProductsFromOrderByFilter(ctx context.Context, uuid string, filter string) (models.OrderForCollectorDuplicate, error) {
	order, err := h.DB.GetDuplicateOrderByUUID(ctx, uuid)
	if err != nil {
		return order, errors.Wrap(err, "fail to get order by uuid")
	}

	order.CartItems = filterProducts(order.CartItems, filter)

	return order, nil
}

func (h *Handler) CancelOrder(ctx context.Context, uuid proto.OrderCancel) error {
	err := h.DB.CancelOrder(ctx, uuid)
	if err != nil {
		errors.Wrap(err, "fail to cancel order.")
	}
	return nil
}

func (h *Handler) FinishCollectOrder(ctx context.Context, uuid string) (models.Order, error) {
	orderDuplicate, err := h.DB.GetDuplicateOrderByUUID(ctx, uuid)
	if err != nil {
		return models.Order{}, errors.Wrap(err, "fail to get order by uuid")
	}

	orderOriginal, err := h.DB.GetOrderByUUID(ctx, uuid)
	if err != nil {
		return models.Order{}, errors.Wrap(err, "fail to ger original order by uuid")
	}

	orderOriginal, err = convertDuplicateToOriginal(orderOriginal, orderDuplicate)
	if err != nil {
		return models.Order{}, errors.Wrap(err, "failed to finish order")
	}

	orderUpdated := models.ConvertAllOrdersForSync(orderDuplicate, orderOriginal)

	orderOriginal, err = h.DB.FinishCollectOrder(ctx, &orderOriginal)
	if err != nil {
		return models.Order{}, errors.Wrap(err, "failed to finish order")
	}

	err = h.RPC.UpdateOrder(ctx, orderUpdated)
	if err != nil {
		return models.Order{}, errors.Wrap(err, "fail to update order")
	}

	var params orderModels.OrdersStateOptions
	params.State = proto.OrderStateReady
	params.InitiatorRole = "user"
	err = h.RPC.SetOrderState(ctx, params, orderOriginal.UUID)
	if err != nil {
		return models.Order{}, errors.Wrap(err, "fail to change order state")
	}

	return orderOriginal, nil
}

func (h *Handler) CreateOrder(ctx context.Context, order *models.Order) (models.Order, error) {
	order.UUID = h.Functions.IDs.GenUUID()
	order.ID = h.Functions.IDs.SliceUUID(order.UUID)
	order, err := h.DB.CreateOrder(ctx, order)
	if err != nil {
		return *order, errors.Wrap(err, "fail to create order")
	}
	return *order, nil
}

func (h *Handler) GetFreeOrders(ctx context.Context, uuid string) ([]models.Order, error) {
	collector, err := h.RPC.GetUserDataByUUID(ctx, uuid)
	if err != nil {
		return []models.Order{}, errors.Wrap(err, "rpc error")
	}

	orders, err := h.DB.GetFreeOrders(ctx, collector.Meta.StoresUUID)
	if err != nil {
		return orders, errors.Wrap(err, "fail to get all free orders")
	}

	return orders, nil
}

func (h *Handler) GetMyOrders(ctx context.Context, collectorUUID string) ([]models.OrderForCollectorDuplicate, error) {
	orders, err := h.DB.GetMyOrders(ctx, collectorUUID)
	if err != nil {
		return orders, errors.Wrap(err, "fail to get collector orders")
	}

	return orders, nil
}

func (h *Handler) MarkProduct(ctx context.Context, params proto.OrderProductMark) error {
	order, err := h.DB.GetDuplicateOrderByUUID(ctx, params.OrderUUID)
	if err != nil {
		return errors.Wrap(err, "fail to get order")
	}

	if len(params.Products) > 0 {
		for _, v := range params.Products {
			var param proto.OrderProductMark

			param.ProductUUID = v.ID
			param.ProductCount = v.Count

			order.CartItems, err = markItemAsCollected(order.CartItems, param)
			if err != nil {
				return errors.Wrap(err, "fail to mark product")
			}
		}
	} else {
		order.CartItems, err = markItemAsCollected(order.CartItems, params)
		if err != nil {
			return errors.Wrap(err, "fail to mark product")
		}
	}
	order.TotalPrice = CalcOrderPrice(order.CartItems) + order.DeliveryData.Price
	err = h.DB.UpdateDuplicateOrder(ctx, order)
	if err != nil {
		return errors.Wrap(err, "fail to update order")
	}

	return nil
}

func (h *Handler) RemoveProductFromOrder(ctx context.Context, params proto.OrderProductRemove) (models.OrderForCollectorDuplicate, error) {
	order, err := h.DB.GetDuplicateOrderByUUID(ctx, params.OrderUUID)
	if err != nil {
		return models.OrderForCollectorDuplicate{}, errors.Wrap(err, "fail to get order")
	}

	if len(params.Products) > 0 {
		for _, v := range params.Products {
			var param proto.OrderProductRemove
			param.ProductUUID = v.ProductID
			order.CartItems, err = removeProductFromOrder(order.CartItems, param)
			if err != nil {
				return models.OrderForCollectorDuplicate{}, errors.Wrap(err, "fail to remove product")
			}
		}
	} else {
		order.CartItems, err = removeProductFromOrder(order.CartItems, params)
		if err != nil {
			return models.OrderForCollectorDuplicate{}, errors.Wrap(err, "fail to remove product")
		}
	}
	order.TotalPrice = CalcOrderPrice(order.CartItems) + order.DeliveryData.Price
	err = h.DB.UpdateDuplicateOrder(ctx, order)
	if err != nil {
		return models.OrderForCollectorDuplicate{}, errors.Wrap(err, "fail to update order")
	}

	return order, nil
}

func (h *Handler) AddProductToOrder(ctx context.Context, params proto.ProductAdd) (models.OrderForCollectorDuplicate, error) {
	order, err := h.DB.GetDuplicateOrderByUUID(ctx, params.OrderUUID)
	if err != nil {
		return models.OrderForCollectorDuplicate{}, errors.Wrap(err, "fail to get order")
	}

	if len(params.Products) > 0 {
		for _, v := range params.Products {
			var param proto.ProductAdd
			param.ProductCount = v.ProductCount
			product, err := h.RPC.GetProductByUUID(ctx, v.ProductUUID)
			if err != nil {
				return models.OrderForCollectorDuplicate{}, errors.Wrap(err, "fail to get product by uuid")
			}
			order.CartItems, err = addProductToArray(order.CartItems, *product, param)
			if err != nil {
				return models.OrderForCollectorDuplicate{}, errors.Wrap(err, "fail to add product")
			}
		}
	} else {
		product, err := h.RPC.GetProductByUUID(ctx, params.ProductUUID)
		if err != nil {
			return models.OrderForCollectorDuplicate{}, errors.Wrap(err, "fail to get product by uuid")
		}

		order.CartItems, err = addProductToArray(order.CartItems, *product, params)
		if err != nil {
			return models.OrderForCollectorDuplicate{}, errors.Wrap(err, "fail to add product")
		}
	}
	order.TotalPrice = CalcOrderPrice(order.CartItems) + order.DeliveryData.Price
	err = h.DB.UpdateDuplicateOrder(ctx, order)
	if err != nil {
		return models.OrderForCollectorDuplicate{}, errors.Wrap(err, "fail to update order")
	}

	return order, nil
}

func (h *Handler) ChangeProduct(ctx context.Context, params proto.OrderProductChange) error {
	order, err := h.DB.GetDuplicateOrderByUUID(ctx, params.OrderUUID)
	if err != nil {
		return errors.Wrap(err, "fail to get order")
	}

	newProduct, err := h.RPC.GetProductByUUID(ctx, params.NewProductUUID)
	if err != nil {
		return errors.Wrap(err, "rpc error")
	}

	order.CartItems, err = changeItem(order.CartItems, params, *newProduct)
	if err != nil {
		return errors.Wrap(err, "fail to change product")
	}

	order.TotalPrice = CalcOrderPrice(order.CartItems) + order.DeliveryData.Price
	err = h.DB.UpdateDuplicateOrder(ctx, order)
	if err != nil {
		return errors.Wrap(err, "fail to update order")
	}

	return nil
}

func (h *Handler) SetCollectorToOrder(ctx context.Context, ids *proto.OrderCollector) (proto.OrderForCollector, error) {
	collectorUUID := ids.CollectorUUID
	orderUUID := ids.OrderUUID

	order, err := h.DB.GetOrderByUUID(ctx, orderUUID)
	if err != nil {
		return proto.OrderForCollector{}, errors.Wrap(err, "failed to get order by uuid")
	}

	if order.CollectorUUID != "" {
		return proto.OrderForCollector{}, errors.New("На этот заказ уже назначен сборщик.")
	}

	collector, err := h.RPC.GetUserDataByUUID(ctx, collectorUUID)
	if err != nil {
		return proto.OrderForCollector{}, errors.Wrap(err, "failed to getting collector")
	}

	order = collectorDataConvert(order, *collector)
	err = h.DB.SetCollectorToOrder(ctx, order)
	if err != nil {
		return proto.OrderForCollector{}, errors.Wrap(err, "fail to set collector to order")
	}

	//var params orderModels.OrdersStateOptions
	//params.State = proto.OrderStateCooking
	//params.InitiatorRole = "user"
	//if ids.CookingTime != 0 {
	//	params.CookingTime = ids.CookingTime
	//}
	//err = h.RPC.SetOrderState(ctx, params, order.UUID)
	//if err != nil {
	//	return proto.OrderForCollector{}, errors.Wrap(err, "fail to change order state")
	//}

	orderDuplicate := models.CreateOrderDuplicate(order)
	orderDuplicate.UUID = h.Functions.IDs.GenUUID()
	err = h.DB.DuplicateOrder(ctx, orderDuplicate)
	if err != nil {
		return proto.OrderForCollector{}, errors.Wrap(err, "fail to duplicate order")
	}

	orderOutput := createOrderForCollector(orderDuplicate)
	return orderOutput, nil
}

func collectorDataConvert(order models.Order, collector models.User) models.Order {
	order.CollectorUUID = collector.UUID
	order.CollectorData = models.CollectorInfo{
		UUID: collector.UUID,
	}
	return order
}

func createOrderForCollector(order models.OrderForCollectorDuplicate) (orderOutput proto.OrderForCollector) {
	for _, v := range order.CartItems {
		orderOutput.CartItems = append(orderOutput.CartItems, convertFromCoreToCollector(v))
	}
	orderOutput.OrderUUID = order.OrderUUID
	orderOutput.CollectorUUID = order.CollectorUUID
	orderOutput.CallbackPhone = order.CallbackPhone
	orderOutput.ClientData.UUID = order.ClientData.UUID
	orderOutput.ClientData.Name = order.ClientData.Name
	orderOutput.ClientData.Application = order.ClientData.Application
	orderOutput.ClientData.MainPhone = order.ClientData.MainPhone
	orderOutput.ClientData.DevicesID = order.ClientData.DevicesID
	orderOutput.ClientData.Deleted = order.ClientData.Deleted
	orderOutput.ClientData.Blocked = order.ClientData.Blocked
	orderOutput.ClientData.Addresses = order.ClientData.Addresses
	orderOutput.ClientData.Meta = order.ClientData.Meta
	orderOutput.CollectTime = 1
	return orderOutput
}

func convertFromCoreToCollector(items models.CartItemDuplicate) proto.CartItem {
	return proto.CartItem{
		ID:              items.ID,
		Product:         convertProductToCollector(items.Product),
		VariantGroups:   nil,
		SingleItemPrice: items.SingleItemPrice,
		TotalItemPrice:  items.TotalItemPrice,
		Count:           items.Count,
		Hash:            items.Hash,
	}
}

func convertProductToCollector(product models.CartProductDuplicate) proto.CartProduct {
	return proto.CartProduct{
		UUID:              product.UUID,
		Name:              product.Name,
		StoreUUID:         product.StoreUUID,
		Type:              proto.ProductType(product.Type),
		Price:             product.Price,
		Weight:            product.Weight,
		WeightMeasurement: product.WeightMeasurement,
		Meta: proto.ProductMeta{
			ShortDescription:  product.Meta.ShortDescription,
			Description:       product.Meta.Description,
			Composition:       product.Meta.Composition,
			Weight:            product.Meta.Weight,
			WeightMeasurement: product.Meta.WeightMeasurement,
			Discount:          product.Meta.Discount,
			Images:            product.Meta.Images,
		},
	}
}

func markItemAsCollected(array []models.CartItemDuplicate, params proto.OrderProductMark) ([]models.CartItemDuplicate, error) {
	for i, _ := range array {
		if array[i].ID == params.ProductUUID {
			if array[i].Count != params.ProductCount {
				array[i].WasChanged = true
				array[i].Count = params.ProductCount
				array[i].TotalItemPrice = array[i].SingleItemPrice * float64(array[i].Count)
				return array, nil
			}
			array[i].CollectionSign = "collected"
			return array, nil
		}
	}
	return array, errors.New("Не могу найти продукт с указанным id. Попробуйте ещё раз.")
}

func removeProductFromOrder(array []models.CartItemDuplicate, params proto.OrderProductRemove) ([]models.CartItemDuplicate, error) {
	for i, _ := range array {
		if array[i].ID == params.ProductUUID {
			array = append(array[:i], array[i+1:]...)
			return array, nil
		}
	}
	return array, errors.New("Не могу найти продукт с указанным id. Попробуйте ещё раз.")
}

func changeItem(array []models.CartItemDuplicate, params proto.OrderProductChange, newProduct models.Product) ([]models.CartItemDuplicate, error) {
	for i, _ := range array {
		if array[i].ID == params.ProductUUID {
			array[i] = changeProducts(array[i], newProduct, params.NewProductCount)
			return array, nil
		}
	}
	return array, errors.New("Не могу найти продукт с указанным id. Попробуйте ещё раз.")
}

func addProductToArray(cartItems []models.CartItemDuplicate, array models.Product, params proto.ProductAdd) ([]models.CartItemDuplicate, error) {
	newProduct := convertModelProductToDuplicate(array, params.ProductCount)
	cartItems = append(cartItems, newProduct)
	return cartItems, nil
}

func convertModelProductToDuplicate(array models.Product, productCount int) models.CartItemDuplicate {
	return models.CartItemDuplicate{
		ID: GenerateID(),
		Product: models.CartProductDuplicate{
			UUID:              array.UUID,
			Name:              array.Name,
			StoreUUID:         array.StoreUUID,
			Type:              array.Type,
			Price:             array.Price,
			Leftover:          array.Leftover,
			Weight:            0,
			WeightMeasurement: "",
			Meta: models.ProductMeta{
				ShortDescription:  array.Meta.ShortDescription,
				Description:       array.Meta.Description,
				Composition:       array.Meta.Composition,
				Weight:            array.Meta.Weight,
				WeightMeasurement: array.Meta.WeightMeasurement,
				Discount:          array.Meta.Discount,
				Images:            array.Meta.Images,
			},
		},
		VariantGroups:       nil,
		SingleItemPrice:     array.Price,
		TotalItemPrice:      array.Price * float64(productCount),
		CollectionSign:      "uncollected",
		WasChanged:          false,
		PreviousProductUUID: "",
		Count:               productCount,
		Hash:                "",
	}
}

func changeProducts(old models.CartItemDuplicate, new models.Product, newProductCount int) models.CartItemDuplicate {
	old.PreviousProductUUID = old.Product.UUID
	old.SingleItemPrice = new.Price
	old.Count = newProductCount
	old.TotalItemPrice = old.SingleItemPrice * float64(old.Count)
	old.WasChanged = true
	old.Product.UUID = new.UUID
	old.Product.Name = new.Name
	old.Product.Type = new.Type
	old.Product.Meta = new.Meta

	return old
}

func convertDuplicateToOriginal(orderOriginal models.Order, orderDuplicate models.OrderForCollectorDuplicate) (models.Order, error) {
	orderOriginal.TotalPrice = orderDuplicate.TotalPrice
	if len(orderDuplicate.CartItems) == 0 {
		return models.Order{}, errors.New("В заказе нет продуктов. Невозможно завершить подобный заказ.")
	}
	orderOriginal.TotalPrice = 0
	orderOriginal.CartItems = nil
	for i, _ := range orderDuplicate.CartItems {
		if orderDuplicate.CartItems[i].CollectionSign == "uncollected" {
			return models.Order{}, errors.New("Вы должны собрать все продукты для завершения заказа.")
		}
		orderOriginal.CartItems = append(orderOriginal.CartItems, convertItemsToOriginal(orderDuplicate.CartItems[i]))
		orderOriginal.TotalPrice += orderDuplicate.CartItems[i].TotalItemPrice
	}
	return orderOriginal, nil
}

func convertItemsToOriginal(duplicateItem models.CartItemDuplicate) models.CartItem {
	var originalItem models.CartItem

	originalItem.TotalItemPrice = duplicateItem.TotalItemPrice
	originalItem.SingleItemPrice = duplicateItem.SingleItemPrice
	originalItem.Product = convertProductToOriginal(duplicateItem.Product)
	originalItem.Hash = duplicateItem.Hash
	originalItem.Count = duplicateItem.Count
	originalItem.ID = duplicateItem.ID
	originalItem.VariantGroups = duplicateItem.VariantGroups

	return originalItem
}

func convertProductToOriginal(duplicate models.CartProductDuplicate) (original models.CartProduct) {
	original.UUID = duplicate.UUID
	original.Meta = duplicate.Meta
	original.Type = duplicate.Type
	original.Name = duplicate.Name
	original.WeightMeasurement = duplicate.WeightMeasurement
	original.Weight = duplicate.Weight
	original.Price = duplicate.Price
	original.StoreUUID = duplicate.StoreUUID

	return original
}

func filterProducts(items []models.CartItemDuplicate, filter string) (itemsOut []models.CartItemDuplicate) {
	var array []models.CartItemDuplicate
	for i, _ := range items {
		if filter != "replaced" {
			switch filter {
			case "original":
				{
					array = append(array, items[i])
				}
			case "collected":
				{
					if items[i].CollectionSign == "collected" {
						array = append(array, items[i])
					}
				}
			case "uncollected":
				{
					if items[i].CollectionSign == "uncollected" {
						array = append(array, items[i])
					}
				}
			}
		} else {
			if items[i].WasChanged {
				array = append(array, items[i])
			}
		}
	}
	return array
}

func CalcOrderPrice(items []models.CartItemDuplicate) float64 {
	var total float64
	for _, v := range items {
		total += v.TotalItemPrice
	}
	return total
}

func GenerateID() int {
	minLimit := 1
	maxLimit := 100000

	ID := minLimit + rand.Intn(maxLimit-minLimit) // generate random int in range
	return ID
}
