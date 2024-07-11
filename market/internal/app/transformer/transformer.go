package transformer

import (
	"github.com/gabe-frasz/stockbroker/market/internal/app/dto"
	"github.com/gabe-frasz/stockbroker/market/internal/app/entity"
)

func ToDomainOrder(input *dto.OrderInput) *entity.Order {
	asset := entity.NewAsset(input.AssetTicker, 1000)
	investor := entity.NewInvestor(input.InvestorID, "John Doe")
	order := entity.NewOrder(input.ID, investor, asset, input.Shares, input.Price, input.OrderType)
	if input.CurrentShares > 0 {
		investor.AddAssetPosition(asset.Ticker, input.CurrentShares)
	}
	return order
}

func ToDtoOrder(order *entity.Order) *dto.OrderOutput {
	output := &dto.OrderOutput{
		ID:           order.ID,
		InvestorID:   order.Investor.ID,
		AssetTicker:  order.Asset.Ticker,
		OrderType:    order.Type,
		Status:       order.Status,
		Shares:       order.Shares,
		PeningShares: order.PendingShares,
	}

	var transactions []*dto.Transaction
	for _, transaction := range order.Transactions {
		transactions = append(transactions, &dto.Transaction{
			ID:          transaction.ID,
			BuyerID:     transaction.BuyingOrder.Investor.ID,
			SellerID:    transaction.SellingOrder.Investor.ID,
			AssetTicker: order.Asset.Ticker,
			Price:       transaction.Price,
			Shares:      transaction.SellingOrder.Shares - transaction.SellingOrder.PendingShares,
		})
	}
	output.Transactions = transactions

	return output
}
