package entity

type Investor struct {
	ID            string
	Name          string
	AssetPosition []*AssetPosition
}

type AssetPosition struct {
	Ticker string
	Shares int
}

func NewInvestor(id, name string) *Investor {
	return &Investor{
		ID:            id,
		Name:          name,
		AssetPosition: []*AssetPosition{},
	}
}

func (i *Investor) AddAssetPosition(ticker string, shares int) {
	i.AssetPosition = append(i.AssetPosition, &AssetPosition{ticker, shares})
}

func (i *Investor) GetAssetPosition(ticker string) *AssetPosition {
	for _, asset := range i.AssetPosition {
		if asset.Ticker == ticker {
			return asset
		}
	}
	return nil
}

func (i *Investor) UpdateAssetPosition(ticker string, shares int) {
	asset := i.GetAssetPosition(ticker)
	if asset != nil {
		asset.Shares += shares
	} else {
		i.AddAssetPosition(ticker, shares)
	}
}
