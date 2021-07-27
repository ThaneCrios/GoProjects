package handler

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gitlab.com/faemproject/backend/core/shared/logs"
	"gitlab.com/faemproject/backend/delivery/collector.core/models"
	"gitlab.com/faemproject/backend/delivery/collector.core/proto"
)

//GetOrderByUUID возвращает заказ по указанному UUID
func (h *Handler) GetOrderByUUID(ctx context.Context, uuid string) (models.Order, error) {
	log := logs.Eloger.WithFields(logrus.Fields{
		"event":        "getting order",
		"from address": uuid,
	})
	order, err := h.DB.GetOrderByUUID(ctx, uuid)
	if err != nil {
		log.WithField("reason", "failed to getting order in DB").Error(err)
		return order, errors.Wrap(err, "fail to get order by uuid")
	}
	return order, nil
}

func (h *Handler) GetProductsFromOrderByFilter(ctx context.Context, uuid string, filter string) (models.OrderForCollectorDuplicate, error) {
	log := logs.Eloger.WithFields(logrus.Fields{
		"event": "getting order",
	})
	order, err := h.DB.GetDuplicateOrderByUUID(ctx, uuid)
	if err != nil {
		log.WithField("reason", "failed to getting order in DB").Error(err)
		return order, errors.Wrap(err, "fail to get order by uuid")
	}

	order.CartItems = filterProducts(order.CartItems, filter)

	return order, nil
}

func (h *Handler) FinishCollectOrder(ctx context.Context, uuid string) (models.Order, error) {
	log := logs.Eloger.WithFields(logrus.Fields{
		"event":           "finish collect order",
		"with order uuid": uuid,
	})
	orderDuplicate, err := h.DB.GetDuplicateOrderByUUID(ctx, uuid)
	if err != nil {
		log.WithField("reason", "failed to getting order in DB").Error(err)
		return models.Order{}, errors.Wrap(err, "fail to get order by uuid")
	}
	orderOriginal, err := h.DB.GetOrderByUUID(ctx, uuid)
	if err != nil {
		log.WithField("reason", "failed to get original order from DB")
		return models.Order{}, errors.Wrap(err, "fail to ger original order by uuid")
	}
	orderOriginal, err = convertDuplicateToOriginal(orderOriginal, orderDuplicate)
	if err != nil {
		log.WithField("reason", "failed to finish order")
		return models.Order{}, errors.Wrap(err, "failed to finish order")
	}
	orderOriginal, err = h.DB.FinishCollectOrder(ctx, &orderOriginal)
	if err != nil {
		log.WithField("reason", "failed to finish order")
		return models.Order{}, errors.Wrap(err, "failed to finish order")
	}

	return orderOriginal, nil
}

func (h *Handler) CreateOrder(ctx context.Context, order *models.Order) (models.Order, error) {
	log := logs.Eloger.WithFields(logrus.Fields{
		"event": "order creating",
	})

	order.UUID = h.RAM.IDs.GenUUID()
	order.ID = h.RAM.IDs.SliceUUID(order.UUID)
	order, err := h.DB.CreateOrder(ctx, order)
	if err != nil {
		log.WithField("reason", "failed to create order and tasks in DB").Error(err)
		return *order, errors.Wrap(err, "fail to create order")
	}
	return *order, nil
}

func (h *Handler) GetFreeOrders(ctx context.Context) ([]models.Order, error) {
	log := logs.Eloger.WithFields(logrus.Fields{
		"event": "get free orders",
	})

	orders, err := h.DB.GetFreeOrders(ctx)
	if err != nil {
		log.WithField("reason", "fail to get all free orders").Error(err)
		return orders, errors.Wrap(err, "fail to get all free orders")
	}
	return orders, nil
}

func (h *Handler) MarkProduct(ctx context.Context, params proto.OrderProductMark) error {
	log := logs.Eloger.WithFields(logrus.Fields{
		"event": "mark product as collected",
	})
	order, err := h.DB.GetDuplicateOrderByUUID(ctx, params.OrderUUID)
	if err != nil {
		log.WithField("reason", "fail to get order").Error(err)
		return errors.Wrap(err, "fail to get order")
	}
	order.CartItems, err = markItemAsCollected(order.CartItems, params)
	if err != nil {
		log.WithField("reason", "fail to mark product").Error(err)
		return errors.Wrap(err, "fail to mark product")
	}
	err = h.DB.UpdateDuplicateOrder(ctx, order)
	if err != nil {
		log.WithField("reason", "fail to update order").Error(err)
		return errors.Wrap(err, "fail to update order")
	}
	return nil
}

func (h *Handler) ChangeProduct(ctx context.Context, params proto.OrderProductChange) error {
	log := logs.Eloger.WithFields(logrus.Fields{
		"event": "change product",
	})
	order, err := h.DB.GetDuplicateOrderByUUID(ctx, params.OrderUUID)
	if err != nil {
		log.WithField("reason", "fail to get order").Error(err)
		return errors.Wrap(err, "fail to get order")
	}

	newProduct, err := h.DB.GetProductByUUID(ctx, params.NewProductUUID)
	if err != nil {
		log.WithField("reason", "fail to get product").Error(err)
		return errors.Wrap(err, "fail to get product")
	}
	order.CartItems, err = changeItem(order.CartItems, params, newProduct)
	if err != nil {
		log.WithField("reason", "fail to change product").Error(err)
		return errors.Wrap(err, "fail to change product")
	}
	err = h.DB.UpdateDuplicateOrder(ctx, order)
	if err != nil {
		log.WithField("reason", "fail to update order").Error(err)
		return errors.Wrap(err, "fail to update order")
	}
	return nil
}

func (h *Handler) SetCollectorToOrder(ctx context.Context, ids *proto.OrderCollector) (proto.OrderForCollector, error) {
	log := logs.Eloger.WithFields(logrus.Fields{
		"event": "set courier",
	})
	collectorUUID := ids.CollectorUUID
	orderUUID := ids.OrderUUID

	order, err := h.DB.GetOrderByUUID(ctx, orderUUID)
	if err != nil {
		log.WithField("reason", "failed to getting tasks").Error(err)
		return proto.OrderForCollector{}, errors.Wrap(err, "failed to getting tasks")
	}
	if order.CollectorUUID != "" {
		log.WithField("reason", "failed set collector").Error(err)
		return proto.OrderForCollector{}, errors.New("collector has already been assigned")
	}

	collector, err := h.DB.GetCollectorByUUID(ctx, collectorUUID)
	if err != nil {
		log.WithField("reason", "failed to getting collector").Error(err)
		return proto.OrderForCollector{}, errors.Wrap(err, "failed to getting collector")
	}
	order = collectorDataConvert(order, *collector)
	err = h.DB.SetCollectorToOrder(ctx, order)
	if err != nil {
		log.WithField("reason", "failed to set collector to order in DB").Error(err)
		return proto.OrderForCollector{}, errors.Wrap(err, "fail to set collector to order")
	}
	orderDuplicate := models.CreateOrderDuplicate(order)
	orderDuplicate.UUID = h.RAM.IDs.GenUUID()
	err = h.DB.DuplicateOrder(ctx, orderDuplicate)
	if err != nil {
		log.WithField("reason", "failed to duplicate order").Error(err)
		return proto.OrderForCollector{}, errors.Wrap(err, "fail to duplicate order")
	}
	logs.Eloger.Info(fmt.Sprintf("order with uuid=%s updated", ids.OrderUUID))

	orderOutput := createOrderForCollector(orderDuplicate)
	return orderOutput, nil
}

func collectorDataConvert(order models.Order, collector models.Collector) models.Order {
	order.CollectorUUID = collector.UUID
	order.CollectorData = models.CollectorInfo{
		UUID:        collector.UUID,
		Name:        collector.CollectorMeta.FullName,
		PhoneNumber: collector.PhoneNumber,
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
	orderOutput.ClientData.Name = order.ClientData.Name
	orderOutput.ClientData.Phone = order.ClientData.Phone
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
			array[i].CollectionSign = "collected"
			if array[i].Count != params.ProductCount {
				array[i].WasChanged = true
				array[i].Count = params.ProductCount
			}
			return array, nil
		}
	}
	return array, errors.New("cant find product with this product_id and count")
}

func changeItem(array []models.CartItemDuplicate, params proto.OrderProductChange, newProduct models.Product) ([]models.CartItemDuplicate, error) {
	for i, _ := range array {
		if array[i].ID == params.ProductUUID {
			array[i] = changeProducts(array[i], newProduct, params.NewProductCount)
			return array, nil
		}
	}
	return array, errors.New("cant find product with this product_id")
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
	for i, _ := range orderOriginal.CartItems {
		if orderDuplicate.CartItems[i].CollectionSign == "uncollected" {
			return models.Order{}, errors.New("You must collect all of products to finish order.")
		}
		orderOriginal.CartItems[i] = convertItemsToOriginal(orderOriginal.CartItems[i], orderDuplicate.CartItems[i])
	}
	return orderOriginal, nil
}

func convertItemsToOriginal(originalItem models.CartItem, duplicateItem models.CartItemDuplicate) models.CartItem {
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
