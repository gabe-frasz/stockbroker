package entity

type Asset struct {
	Ticker       string
	MarketVolume int
}

func NewAsset(ticker string, marketVolume int) *Asset {
	return &Asset{
		Ticker:       ticker,
		MarketVolume: marketVolume,
	}
}
