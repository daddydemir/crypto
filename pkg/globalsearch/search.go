package globalsearch

import (
	"sort"
	"strconv"
	"strings"

	"github.com/daddydemir/crypto/pkg/infrastructure"
	"gorm.io/gorm"
)

type Item struct {
	ID       string `json:"id"`
	Type     string `json:"type"`
	Title    string `json:"title"`
	Subtitle string `json:"subtitle"`
	Meta     string `json:"meta"`
}
type Results struct {
	Coins      []Item `json:"coins"`
	Strategies []Item `json:"strategies"`
	Runs       []Item `json:"strategyRuns"`
	Alerts     []Item `json:"alerts"`
	Trades     []Item `json:"trades"`
}
type Service struct {
	db    *gorm.DB
	coins *infrastructure.CoinCatalog
}

func New(db *gorm.DB, coins *infrastructure.CoinCatalog) *Service {
	return &Service{db: db, coins: coins}
}

func rank(q string, values ...string) int {
	best := 3
	for _, value := range values {
		v := strings.ToLower(value)
		if v == q {
			return 0
		}
		if strings.HasPrefix(v, q) && best > 1 {
			best = 1
		}
		if strings.Contains(v, q) && best > 2 {
			best = 2
		}
	}
	return best
}
func (s *Service) Search(username, query string, limit int) (Results, error) {
	q := strings.ToLower(strings.TrimSpace(query))
	if limit < 1 || limit > 10 {
		limit = 5
	}
	out := Results{Coins: []Item{}, Strategies: []Item{}, Runs: []Item{}, Alerts: []Item{}, Trades: []Item{}}
	coins, err := s.coins.List()
	if err != nil {
		return out, err
	}
	type ranked struct {
		item  Item
		rank  int
		order int
	}
	matches := []ranked{}
	for _, c := range coins {
		r := rank(q, c.Symbol, c.Name)
		if r < 3 {
			matches = append(matches, ranked{Item{c.ID, "COIN", c.Symbol, c.Name, "Coin"}, r, c.MarketCapRank})
		}
	}
	sort.SliceStable(matches, func(i, j int) bool {
		if matches[i].rank != matches[j].rank {
			return matches[i].rank < matches[j].rank
		}
		return matches[i].order < matches[j].order
	})
	for i := 0; i < len(matches) && i < limit; i++ {
		out.Coins = append(out.Coins, matches[i].item)
	}
	pattern := "%" + q + "%"
	prefix := q + "%"
	err = s.db.Raw(`select id::text id,'STRATEGY' type,name title,coalesce(description,'') subtitle,concat(run_count,' runs') meta from strategies where username=? and (name ilike ? or description ilike ?) order by case when lower(name)=? then 0 when lower(name) like ? then 1 else 2 end,updated_at desc limit ?`, username, pattern, pattern, q, prefix, limit).Scan(&out.Strategies).Error
	if err != nil {
		return out, err
	}
	err = s.db.Raw(`select id::text id,'STRATEGY_RUN' type,strategy_snapshot->>'name' title,concat(coin_symbol,' · ',replace(type,'_',' ')) subtitle,status meta from strategy_runs where username=? and (coin_symbol ilike ? or strategy_snapshot->>'name' ilike ? or type ilike ? or status ilike ?) order by case when lower(coin_symbol)=? then 0 when lower(coin_symbol) like ? then 1 else 2 end,created_at desc limit ?`, username, pattern, pattern, pattern, pattern, q, prefix, limit).Scan(&out.Runs).Error
	if err != nil {
		return out, err
	}
	err = s.db.Raw(`select id::text id,'ALERT' type,concat(coin,' ',case when is_above then 'Above' else 'Below' end,' Alert') title,coin subtitle,case when is_active then 'ACTIVE' else 'INACTIVE' end meta from alerts where coin ilike ? order by case when lower(coin)=? then 0 else 1 end,create_date desc limit ?`, pattern, q, limit).Scan(&out.Alerts).Error
	if err != nil {
		return out, err
	}
	err = s.db.Raw(`select t.id::text id,'TRADE' type,r.coin_symbol title,r.strategy_snapshot->>'name' subtitle,case when t.exit_execution_date is null then 'OPEN' else concat(case when t.pn_l_percent>=0 then '+' else '' end,round(t.pn_l_percent::numeric,2),'%') end meta from strategy_trades t join strategy_runs r on r.id=t.run_id where r.username=? and (r.coin_symbol ilike ? or r.strategy_snapshot->>'name' ilike ?) order by case when lower(r.coin_symbol)=? then 0 else 1 end,t.entry_execution_date desc limit ?`, username, pattern, pattern, q, limit).Scan(&out.Trades).Error
	return out, err
}
func ID(v uint) string { return strconv.FormatUint(uint64(v), 10) }
