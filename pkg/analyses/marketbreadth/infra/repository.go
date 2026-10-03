package infra

import (
	"fmt"
	"strings"
	"time"

	"github.com/daddydemir/crypto/pkg/analyses/marketbreadth/domain"
	"gorm.io/gorm"
)

type Repository struct{ db *gorm.DB }

func NewRepository(db *gorm.DB) *Repository { return &Repository{db: db} }

func (r *Repository) Candles(symbols []string, maxPoints int) (map[string][]domain.Candle, error) {
	bucket := `date_trunc('day', candle_date)`
	query := fmt.Sprintf(`with bucketed as (
		select distinct on (symbol, bucket) symbol, bucket as time, close_price as close
		from (select symbol, candle_date, close_price, %s as bucket from yahoo_candles where symbol in (?) and close_price is not null) c
		order by symbol, bucket, candle_date desc
	), ranked as (select *, row_number() over(partition by symbol order by time desc) rn from bucketed)
	select symbol,time,close from ranked where rn <= ? order by symbol,time`, bucket)
	var rows []domain.Candle
	if err := r.db.Raw(query, symbols, maxPoints).Scan(&rows).Error; err != nil {
		return nil, err
	}
	result := map[string][]domain.Candle{}
	for _, row := range rows {
		result[strings.ToUpper(row.Symbol)] = append(result[strings.ToUpper(row.Symbol)], row)
	}
	return result, nil
}

func LatestTime(series map[string][]domain.Candle) time.Time {
	var latest time.Time
	for _, rows := range series {
		if len(rows) > 0 && rows[len(rows)-1].Time.After(latest) {
			latest = rows[len(rows)-1].Time
		}
	}
	return latest
}
