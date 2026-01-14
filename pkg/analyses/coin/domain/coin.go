package domain

type StatsDTO struct {
	ID                  string  `json:"id"`
	Name                string  `json:"name"`
	Symbol              string  `json:"symbol"`
	Price               float64 `json:"price"`
	Change24h           float64 `json:"change24h"`
	Change7d            float64 `json:"change7d"`
	Change30d           float64 `json:"change30d"`
	ArithmeticChange7d  float64 `json:"arithmeticChange7d"`
	ArithmeticChange30d float64 `json:"arithmeticChange30d"`
}

type Coin struct {
	ID       string
	Name     string
	Symbol   string
	PriceUSD float64
}

type PriceResult struct {
	ExchangeID          string  `gorm:"column:exchange_id"`
	CurrentPrice        float64 `gorm:"column:current_price"`
	DayAgoPrice         float64 `gorm:"column:day_ago_price"`
	WeekAgoPrice        float64 `gorm:"column:week_ago_price"`
	MonthAgoPrice       float64 `gorm:"column:month_ago_price"`
	AvgWeekAgoPrice     float64 `gorm:"column:avg_7_days_price"`
	AvgMonthAgoPrice    float64 `gorm:"column:avg_30_days_price"`
	Change24h           float64 `gorm:"column:change_24h"`
	Change7d            float64 `gorm:"column:change_7d"`
	Change30d           float64 `gorm:"column:change_30d"`
	ChangeArithmetic7d  float64 `gorm:"column:change_arithmetic_7d"`
	ChangeArithmetic30d float64 `gorm:"column:change_arithmetic_30d"`
}
