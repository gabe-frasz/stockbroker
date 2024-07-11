package entity

type Order struct {
	ID            string
	Investor      *Investor
	Asset         *Asset
	Shares        int
	PendingShares int
	Price         float64
	Type          OrderType
	Status        OrderStatus
	Transactions  []*Transaction
}

type OrderType string

const (
	BuyOrder  OrderType = "BUY"
	SellOrder OrderType = "SELL"
)

type OrderStatus string

const (
	OpenOrder   OrderStatus = "OPEN"
	ClosedOrder OrderStatus = "CLOSED"
)

func NewOrder(id string, investor *Investor, asset *Asset, shares int, price float64, orderType OrderType) *Order {
	return &Order{
		ID:            id,
		Investor:      investor,
		Asset:         asset,
		Shares:        shares,
		PendingShares: shares,
		Price:         price,
		Type:          orderType,
		Status:        OpenOrder,
		Transactions:  []*Transaction{},
	}
}

func (o *Order) Close() (ok bool) {
	if o.PendingShares == 0 {
		o.Status = ClosedOrder
		return true
	} else {
		return false
	}
}

func (o *Order) UpdatePendingShares(shares int) {
	o.PendingShares += shares
}
