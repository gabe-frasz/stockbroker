package dto

import "github.com/gabe-frasz/stockbroker/market/internal/app/entity"

type OrderInput struct {
	ID            string           `json:"id"`
	InvestorID    string           `json:"investor_id"`
	AssetTicker   string           `json:"asset_ticker"`
	CurrentShares int              `json:"current_shares"`
	Shares        int              `json:"shares"`
	Price         float64          `json:"price"`
	OrderType     entity.OrderType `json:"order_type"`
}

type OrderOutput struct {
	ID           string             `json:"id"`
	InvestorID   string             `json:"investor_id"`
	AssetTicker  string             `json:"asset_ticker"`
	OrderType    entity.OrderType   `json:"order_type"`
	Status       entity.OrderStatus `json:"status"`
	Shares       int                `json:"shares"`
	PeningShares int                `json:"pending_shares"`
	Transactions []*Transaction     `json:"transactions"`
}

type Transaction struct {
	ID          string  `json:"id"`
	BuyerID     string  `json:"buyer_id"`
	SellerID    string  `json:"seller_id"`
	AssetTicker string  `json:"asset_ticker"`
	Price       float64 `json:"price"`
	Shares      int     `json:"shares"`
}
