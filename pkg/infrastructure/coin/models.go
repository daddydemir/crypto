package coin

type Result struct {
	ExchangeId string
	Date       string
	Price      float64 `gorm:"column:first_price"`
}
