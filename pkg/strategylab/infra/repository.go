package infra

import (
	"errors"
	"github.com/daddydemir/crypto/pkg/strategylab/domain"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type Repository struct{ DB *gorm.DB }

func NewRepository(db *gorm.DB) *Repository { return &Repository{db} }
func (r *Repository) Migrate() error {
	if err := r.DB.Exec(`create table if not exists app_schema_migrations (name varchar(120) primary key, applied_at timestamptz not null default current_timestamp)`).Error; err != nil {
		return err
	}
	var count int64
	if err := r.DB.Raw(`select count(*) from app_schema_migrations where name = ?`, "strategy_lab_v5_persisted_charts").Scan(&count).Error; err != nil || count > 0 {
		return err
	}
	if err := r.DB.AutoMigrate(&domain.Strategy{}, &domain.StrategyRun{}, &domain.StrategyTrade{}, &domain.EquitySnapshot{}, &domain.Evaluation{}); err != nil {
		return err
	}
	if err := r.DB.Exec(`with values as (
		select e.id,
			case when r.initial_capital > 0 then round((((e.equity / r.initial_capital) - 1) * 100)::numeric, 2) else 0 end strategy_return,
			case when r.initial_capital > 0 then round((((e.benchmark_equity / r.initial_capital) - 1) * 100)::numeric, 2) else 0 end benchmark_return,
			max(e.equity) over(partition by e.run_id order by e.date rows between unbounded preceding and current row) running_peak
		from strategy_equity_snapshots e join strategy_runs r on r.id=e.run_id
	), calculated as (
		select id,strategy_return,benchmark_return,running_peak,
			case when running_peak > 0 then round((((select equity from strategy_equity_snapshots where id=values.id) - running_peak) / running_peak * 100)::numeric, 2) else 0 end drawdown
		from values
	)
	update strategy_equity_snapshots e set strategy_return=c.strategy_return,benchmark_return=c.benchmark_return,running_peak=c.running_peak,drawdown=c.drawdown from calculated c where c.id=e.id`).Error; err != nil {
		return err
	}
	if err := r.DB.Exec(`update strategies s set run_count = c.total from (select strategy_id, count(*) as total from strategy_runs group by strategy_id) c where s.id = c.strategy_id`).Error; err != nil {
		return err
	}
	return r.DB.Exec(`insert into app_schema_migrations(name) values (?) on conflict (name) do nothing`, "strategy_lab_v5_persisted_charts").Error
}
func (r *Repository) Prices(symbol string) ([]domain.PricePoint, error) {
	var v []domain.PricePoint
	err := r.DB.Raw(`select candle_date as date, close_price as close, high_price as high, low_price as low
		from yahoo_candles where upper(symbol)=upper(?) and close_price is not null
		and high_price is not null and low_price is not null order by candle_date`, symbol).Scan(&v).Error
	return v, err
}
func (r *Repository) ListStrategies(u string) ([]domain.Strategy, error) {
	var v []domain.Strategy
	err := r.DB.Where("username=?", u).Order("updated_at desc").Find(&v).Error
	return v, err
}
func (r *Repository) Strategy(u string, id uint) (domain.Strategy, error) {
	var v domain.Strategy
	err := r.DB.Where("id=? and username=?", id, u).First(&v).Error
	return v, err
}
func (r *Repository) SaveStrategy(v *domain.Strategy) error { return r.DB.Save(v).Error }
func (r *Repository) DeleteStrategy(u string, id uint) error {
	res := r.DB.Where("id=? and username=?", id, u).Delete(&domain.Strategy{})
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return res.Error
}
func (r *Repository) SaveRun(v *domain.StrategyRun) error { return r.DB.Save(v).Error }
func (r *Repository) IncrementRunCount(strategyID uint) error {
	return r.DB.Model(&domain.Strategy{}).Where("id=?", strategyID).UpdateColumn("run_count", gorm.Expr("run_count + 1")).Error
}
func (r *Repository) DeleteRun(username string, id uint) error {
	return r.DB.Transaction(func(tx *gorm.DB) error {
		var run domain.StrategyRun
		if err := tx.Where("id=? and username=?", id, username).First(&run).Error; err != nil {
			return err
		}
		if run.Status == domain.StatusQueued || run.Status == domain.StatusRunning {
			return errors.New("active strategy runs cannot be deleted")
		}
		for _, model := range []any{&domain.StrategyTrade{}, &domain.EquitySnapshot{}, &domain.Evaluation{}} {
			if err := tx.Where("run_id=?", id).Delete(model).Error; err != nil {
				return err
			}
		}
		if err := tx.Delete(&run).Error; err != nil {
			return err
		}
		return tx.Model(&domain.Strategy{}).Where("id=?", run.StrategyID).UpdateColumn("run_count", gorm.Expr("case when run_count > 0 then run_count - 1 else 0 end")).Error
	})
}
func (r *Repository) Run(u string, id uint) (domain.StrategyRun, error) {
	var v domain.StrategyRun
	err := r.DB.Where("id=? and username=?", id, u).First(&v).Error
	return v, err
}
func (r *Repository) RunByID(id uint) (domain.StrategyRun, error) {
	var v domain.StrategyRun
	err := r.DB.First(&v, id).Error
	return v, err
}
func (r *Repository) Runs(u string) ([]domain.StrategyRun, error) {
	var v []domain.StrategyRun
	err := r.DB.Where("username=?", u).Order("created_at desc").Limit(100).Find(&v).Error
	return v, err
}
func (r *Repository) Running() (v []domain.StrategyRun, err error) {
	err = r.DB.Where("type=? and status=?", domain.RunForward, domain.StatusRunning).Find(&v).Error
	return
}
func (r *Repository) Trades(id uint) (v []domain.StrategyTrade, err error) {
	err = r.DB.Where("run_id=?", id).Order("entry_execution_date").Find(&v).Error
	return
}
func (r *Repository) Equity(id uint) (v []domain.EquitySnapshot, err error) {
	err = r.DB.Where("run_id=?", id).Order("date").Find(&v).Error
	return
}
func (r *Repository) PeakEquity(id uint) (float64, error) {
	var peak float64
	err := r.DB.Model(&domain.EquitySnapshot{}).Where("run_id=?", id).Select("coalesce(max(equity),0)").Scan(&peak).Error
	return peak, err
}
func (r *Repository) CreateTrade(v *domain.StrategyTrade) error { return r.DB.Create(v).Error }
func (r *Repository) SaveTrade(v *domain.StrategyTrade) error   { return r.DB.Save(v).Error }
func (r *Repository) OpenTrade(id uint) (domain.StrategyTrade, error) {
	var v domain.StrategyTrade
	err := r.DB.Where("run_id=? and exit_execution_date is null", id).First(&v).Error
	return v, err
}
func (r *Repository) Snapshot(v domain.EquitySnapshot) error {
	return r.DB.Clauses(clause.OnConflict{DoNothing: true}).Create(&v).Error
}
func (r *Repository) Snapshots(v []domain.EquitySnapshot) error {
	if len(v) == 0 {
		return nil
	}
	return r.DB.Clauses(clause.OnConflict{DoNothing: true}).CreateInBatches(v, 500).Error
}
func (r *Repository) Mark(id uint, d time.Time) bool {
	res := r.DB.Clauses(clause.OnConflict{DoNothing: true}).Create(&domain.Evaluation{RunID: id, Date: d})
	return res.Error == nil && res.RowsAffected == 1
}
func IsNotFound(e error) bool { return errors.Is(e, gorm.ErrRecordNotFound) }
