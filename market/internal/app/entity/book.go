package entity

import (
	"container/heap"
	"sync"
)

type Book struct {
	Order         []*Order
	Transactions  []*Transaction
	OrdersChan    <-chan *Order // input
	OrdersChanOut chan *Order
	Wg            *sync.WaitGroup
}

func NewBook(orderChan <-chan *Order, orderChanOut chan *Order, wg *sync.WaitGroup) *Book {
	return &Book{
		Order:         []*Order{},
		Transactions:  []*Transaction{},
		OrdersChan:    orderChan,
		OrdersChanOut: orderChanOut,
		Wg:            wg,
	}
}

func (b *Book) Trade() {
	buyOrders := make(map[string]*OrderQueue)
	sellOrders := make(map[string]*OrderQueue)

	for order := range b.OrdersChan {
		asset := order.Asset.Ticker

		if buyOrders[asset] == nil {
			buyOrders[asset] = NewOrderQueue()
			heap.Init(buyOrders[asset])
		}

		if sellOrders[asset] == nil {
			sellOrders[asset] = NewOrderQueue()
			heap.Init(sellOrders[asset])
		}

		buyOrdersQueue := buyOrders[asset]
		sellOrdersQueue := sellOrders[asset]

		if order.Type == BuyOrder {
			buyOrdersQueue.Push(order)
			if sellOrdersQueue.Len() > 0 && order.Price >= sellOrdersQueue.Orders[0].Price {
				sellOrder := sellOrdersQueue.Pop().(*Order)
				if sellOrder.PendingShares <= 0 {
					sellOrder.Close()
					continue
				}
				transaction := NewTransaction(sellOrder, order, min(sellOrder.Shares, order.Shares), sellOrder.Price)
				b.AddTransaction(transaction, b.Wg)
				sellOrder.Transactions = append(sellOrder.Transactions, transaction)
				order.Transactions = append(order.Transactions, transaction)
				b.OrdersChanOut <- sellOrder
				b.OrdersChanOut <- order
				if sellOrder.PendingShares > 0 {
					sellOrdersQueue.Push(sellOrder)
				}
			}
		} else if order.Type == SellOrder {
			sellOrdersQueue.Push(order)
			if buyOrdersQueue.Len() > 0 && order.Price <= buyOrdersQueue.Orders[0].Price {
				buyOrder := buyOrdersQueue.Pop().(*Order)
				if buyOrder.PendingShares <= 0 {
					buyOrder.Close()
					continue

				}
				transaction := NewTransaction(order, buyOrder, min(order.Shares, buyOrder.Shares), order.Price)
				b.AddTransaction(transaction, b.Wg)
				order.Transactions = append(order.Transactions, transaction)
				buyOrder.Transactions = append(buyOrder.Transactions, transaction)
				b.OrdersChanOut <- order
				b.OrdersChanOut <- buyOrder
				if buyOrder.PendingShares > 0 {
					buyOrdersQueue.Push(buyOrder)
				}
			}
		}
	}
}

func (b *Book) AddTransaction(transaction *Transaction, wg *sync.WaitGroup) {
	defer wg.Done()

	sellingShares := transaction.SellingOrder.PendingShares
	buyingShares := transaction.BuyingOrder.PendingShares

	minShares := min(sellingShares, buyingShares)

	transaction.SellingOrder.Investor.UpdateAssetPosition(transaction.SellingOrder.Asset.Ticker, -minShares)
	transaction.SellingOrder.UpdatePendingShares(-minShares)

	transaction.BuyingOrder.Investor.UpdateAssetPosition(transaction.BuyingOrder.Asset.Ticker, minShares)
	transaction.BuyingOrder.UpdatePendingShares(-minShares)

	transaction.SellingOrder.Close()
	transaction.BuyingOrder.Close()
	b.Transactions = append(b.Transactions, transaction)
}
