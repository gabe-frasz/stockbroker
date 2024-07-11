package entity

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuyAsset(t *testing.T) {
	asset1 := NewAsset("ASST1", 100)

	investor1 := NewInvestor("1", "Investor 1")
	investor2 := NewInvestor("2", "Investor 2")

	investor1.AddAssetPosition("ASST1", 10)

	wg := sync.WaitGroup{}
	orderChan := make(chan *Order)
	orderChanOut := make(chan *Order)

	book := NewBook(orderChan, orderChanOut, &wg)
	go book.Trade()

	wg.Add(1)
	order1 := NewOrder(investor1, asset1, 5, 5, SellOrder)
	orderChan <- order1

	order2 := NewOrder(investor2, asset1, 5, 5, BuyOrder)
	orderChan <- order2
	wg.Wait()

	assert := assert.New(t)
	assert.Equal(ClosedOrder, order1.Status, "Order 1 should be closed")
	assert.Equal(0, order1.PendingShares, "Order 1 should have 0 PendingShares")
	assert.Equal(ClosedOrder, order2.Status, "Order 2 should be closed")
	assert.Equal(0, order2.PendingShares, "Order 2 should have 0 PendingShares")

	assert.Equal(5, investor1.GetAssetPosition("ASST1").Shares, "Investor 1 should have 5 shares of asset 1")
	assert.Equal(5, investor2.GetAssetPosition("ASST1").Shares, "Investor 2 should have 5 shares of asset 1")
}

func TestBuyAssetWithDifferentAssets(t *testing.T) {
	asset1 := NewAsset("ASST1", 100)
	asset2 := NewAsset("ASST2", 100)

	investor1 := NewInvestor("1", "Investor 1")
	investor2 := NewInvestor("2", "Investor 2")

	investor1.AddAssetPosition("ASST1", 10)

	investor2.AddAssetPosition("ASST2", 10)

	wg := sync.WaitGroup{}
	orderChan := make(chan *Order)
	orderChanOut := make(chan *Order)

	book := NewBook(orderChan, orderChanOut, &wg)
	go book.Trade()

	order1 := NewOrder(investor1, asset1, 5, 5, SellOrder)
	orderChan <- order1

	order2 := NewOrder(investor2, asset2, 5, 5, BuyOrder)
	orderChan <- order2

	assert := assert.New(t)
	assert.Equal(OpenOrder, order1.Status, "Order 1 should be open")
	assert.Equal(5, order1.PendingShares, "Order 1 should have 5 PendingShares")
	assert.Equal(OpenOrder, order2.Status, "Order 2 should be open")
	assert.Equal(5, order2.PendingShares, "Order 2 should have 5 PendingShares")
}

func TestBuyPartialAsset(t *testing.T) {
	asset1 := NewAsset("ASST1", 100)

	investor1 := NewInvestor("1", "Investor 1")
	investor2 := NewInvestor("2", "Investor 2")
	investor3 := NewInvestor("3", "Investor 3")

	investor1.AddAssetPosition("ASST1", 3)
	investor3.AddAssetPosition("ASST1", 5)

	wg := sync.WaitGroup{}
	orderChan := make(chan *Order)
	orderChanOut := make(chan *Order)

	book := NewBook(orderChan, orderChanOut, &wg)
	go book.Trade()

	wg.Add(1)
	// investidor 2 quer comprar 5 shares
	order2 := NewOrder(investor2, asset1, 5, 5.0, BuyOrder)
	orderChan <- order2

	// investidor 1 quer vender 3 shares
	order1 := NewOrder(investor1, asset1, 3, 5.0, SellOrder)
	orderChan <- order1

	go func() {
		for range orderChanOut {
		}
	}()
	wg.Wait()

	assert := assert.New(t)
	assert.Equal(ClosedOrder, order1.Status, "Order 1 should be closed")
	assert.Equal(0, order1.PendingShares, "Order 1 should have 0 PendingShares")

	assert.Equal(OpenOrder, order2.Status, "Order 2 should be OpenOrder")
	assert.Equal(2, order2.PendingShares, "Order 2 should have 2 PendingShares")

	assert.Equal(0, investor1.GetAssetPosition("ASST1").Shares, "Investor 1 should have 0 shares of asset 1")
	assert.Equal(3, investor2.GetAssetPosition("ASST1").Shares, "Investor 2 should have 3 shares of asset 1")

	wg.Add(1)
	order3 := NewOrder(investor3, asset1, 2, 5.0, SellOrder)
	orderChan <- order3
	wg.Wait()

	assert.Equal(ClosedOrder, order3.Status, "Order 3 should be closed")
	assert.Equal(0, order3.PendingShares, "Order 3 should have 0 PendingShares")

	assert.Equal(ClosedOrder, order2.Status, "Order 2 should be closed")
	assert.Equal(0, order2.PendingShares, "Order 2 should have 0 PendingShares")

	assert.Equal(2, len(book.Transactions), "Should have 2 transactions")
	assert.Equal(15.0, float64(book.Transactions[0].Total), "Transaction should have price 15")
	assert.Equal(10.0, float64(book.Transactions[1].Total), "Transaction should have price 10")
}

func TestBuyWithDifferentPrice(t *testing.T) {
	asset1 := NewAsset("ASST1", 100)

	investor1 := NewInvestor("1", "Investor 1")
	investor2 := NewInvestor("2", "Investor 2")
	investor3 := NewInvestor("3", "Investor 3")

	investor1.AddAssetPosition("ASST1", 3)
	investor3.AddAssetPosition("ASST1", 5)

	wg := sync.WaitGroup{}
	orderChan := make(chan *Order)
	orderChanOut := make(chan *Order)

	book := NewBook(orderChan, orderChanOut, &wg)
	go book.Trade()

	wg.Add(1)
	// investidor 2 quer comprar 5 shares
	order2 := NewOrder(investor2, asset1, 5, 5.0, BuyOrder)
	orderChan <- order2

	// investidor 1 quer vender 3 shares
	order1 := NewOrder(investor1, asset1, 3, 4.0, SellOrder)
	orderChan <- order1

	go func() {
		for range orderChanOut {
		}
	}()
	wg.Wait()

	assert := assert.New(t)
	assert.Equal(ClosedOrder, order1.Status, "Order 1 should be closed")
	assert.Equal(0, order1.PendingShares, "Order 1 should have 0 PendingShares")

	assert.Equal(OpenOrder, order2.Status, "Order 2 should be open")
	assert.Equal(2, order2.PendingShares, "Order 2 should have 2 PendingShares")

	assert.Equal(0, investor1.GetAssetPosition("ASST1").Shares, "Investor 1 should have 0 shares of asset 1")
	assert.Equal(3, investor2.GetAssetPosition("ASST1").Shares, "Investor 2 should have 3 shares of asset 1")

	wg.Add(1)
	order3 := NewOrder(investor3, asset1, 3, 4.5, SellOrder)
	orderChan <- order3

	wg.Wait()

	assert.Equal(OpenOrder, order3.Status, "Order 3 should be open")
	assert.Equal(1, order3.PendingShares, "Order 3 should have 1 PendingShares")

	assert.Equal(ClosedOrder, order2.Status, "Order 2 should be closed")
	assert.Equal(0, order2.PendingShares, "Order 2 should have 0 PendingShares")

	assert.Equal(2, len(book.Transactions), "Should have 2 transactions")
	assert.Equal(12.0, float64(book.Transactions[0].Total), "Transaction should have price 12")
	assert.Equal(13.5, float64(book.Transactions[1].Total), "Transaction should have price 13.5")
}

func TestNoMatch(t *testing.T) {
	asset1 := NewAsset("ASST1", 100)

	investor1 := NewInvestor("1", "Investor 1")
	investor2 := NewInvestor("2", "Investor 2")

	investor1.AddAssetPosition("ASST1", 3)

	wg := sync.WaitGroup{}
	orderChan := make(chan *Order)
	orderChanOut := make(chan *Order)

	book := NewBook(orderChan, orderChanOut, &wg)
	go book.Trade()

	wg.Add(0)
	// investidor 1 quer vender 3 shares
	order1 := NewOrder(investor1, asset1, 3, 6.0, SellOrder)
	orderChan <- order1

	// investidor 2 quer comprar 5 shares
	order2 := NewOrder(investor2, asset1, 5, 5.0, BuyOrder)
	orderChan <- order2

	go func() {
		for range orderChanOut {
		}
	}()
	wg.Wait()

	assert := assert.New(t)
	assert.Equal(OpenOrder, order1.Status, "Order 1 should be open")
	assert.Equal(OpenOrder, order2.Status, "Order 2 should be open")
	assert.Equal(3, order1.PendingShares, "Order 1 should have 3 PendingShares")
	assert.Equal(5, order2.PendingShares, "Order 2 should have 5 PendingShares")
}
