package domain

import "time"

const (
	LogicAND        = "AND"
	LogicOR         = "OR"
	RunBacktest     = "BACKTEST"
	RunForward      = "FORWARD_TEST"
	StatusDraft     = "DRAFT"
	StatusQueued    = "QUEUED"
	StatusRunning   = "RUNNING"
	StatusPaused    = "PAUSED"
	StatusCompleted = "COMPLETED"
	StatusStopped   = "STOPPED"
)

type Condition struct {
	Left       string   `json:"left"`
	Operator   string   `json:"operator"`
	RightField string   `json:"rightField,omitempty"`
	Value      *float64 `json:"value,omitempty"`
}
type RuleGroup struct {
	Logic      string      `json:"logic"`
	Conditions []Condition `json:"conditions"`
}
type RiskManagement struct {
	StopLossPercent   *float64 `json:"stopLossPercent,omitempty"`
	TakeProfitPercent *float64 `json:"takeProfitPercent,omitempty"`
}
type Configuration struct {
	Entry RuleGroup      `json:"entry"`
	Exit  RuleGroup      `json:"exit"`
	Risk  RiskManagement `json:"risk"`
}
type Strategy struct {
	ID            uint          `gorm:"primaryKey" json:"id"`
	Username      string        `gorm:"size:25;not null;index" json:"-"`
	Name          string        `gorm:"size:120;not null" json:"name"`
	Description   string        `gorm:"size:500" json:"description"`
	Version       int           `gorm:"not null;default:1" json:"version"`
	Configuration Configuration `gorm:"serializer:json;type:jsonb" json:"configuration"`
	RunCount      int64         `gorm:"not null;default:0" json:"runCount"`
	CreatedAt     time.Time     `json:"createdAt"`
	UpdatedAt     time.Time     `json:"updatedAt"`
}

func (Strategy) TableName() string { return "strategies" }

type StrategySnapshot struct {
	StrategyID    uint          `json:"strategyId"`
	Name          string        `json:"name"`
	Version       int           `json:"version"`
	Configuration Configuration `json:"configuration"`
}
type StrategyRun struct {
	ID                  uint             `gorm:"primaryKey" json:"id"`
	Username            string           `gorm:"size:25;not null;index" json:"-"`
	StrategyID          uint             `gorm:"index;not null" json:"strategyId"`
	StrategySnapshot    StrategySnapshot `gorm:"serializer:json;type:jsonb" json:"strategySnapshot"`
	CoinSymbol          string           `gorm:"size:20;not null;index" json:"coinSymbol"`
	Type                string           `gorm:"size:20;not null;index" json:"type"`
	Status              string           `gorm:"size:20;not null;index" json:"status"`
	StartDate           time.Time        `json:"startDate"`
	EndDate             *time.Time       `json:"endDate,omitempty"`
	InitialCapital      float64          `json:"initialCapital"`
	FeePercent          float64          `json:"feePercent"`
	SlippagePercent     float64          `json:"slippagePercent"`
	Cash                float64          `json:"cash"`
	Quantity            float64          `json:"quantity"`
	EntryPrice          float64          `json:"entryPrice"`
	EntryExecutionDate  *time.Time       `json:"entryExecutionDate,omitempty"`
	PendingAction       string           `gorm:"size:8" json:"pendingAction"`
	PendingExitReason   string           `gorm:"size:32" json:"pendingExitReason"`
	PendingSignalDate   *time.Time       `json:"pendingSignalDate,omitempty"`
	PendingSnapshot     SignalSnapshot   `gorm:"serializer:json;type:jsonb" json:"pendingSnapshot"`
	LastEvaluatedDate   *time.Time       `json:"lastEvaluatedDate,omitempty"`
	TotalFees           float64          `json:"totalFees"`
	BenchmarkQuantity   float64          `json:"benchmarkQuantity"`
	CreatedAt           time.Time        `json:"createdAt"`
	UpdatedAt           time.Time        `json:"updatedAt"`
	Metrics             Metrics          `gorm:"serializer:json;type:jsonb" json:"metrics"`
	MetricsCalculatedAt *time.Time       `json:"metricsCalculatedAt,omitempty"`
}

func (StrategyRun) TableName() string { return "strategy_runs" }

type ConditionResult struct {
	Condition  Condition `json:"condition"`
	LeftValue  float64   `json:"leftValue"`
	RightValue float64   `json:"rightValue"`
	Matched    bool      `json:"matched"`
}
type SignalSnapshot struct {
	Date       time.Time          `json:"date"`
	Values     map[string]float64 `json:"values"`
	Conditions []ConditionResult  `json:"conditions"`
}
type StrategyTrade struct {
	ID                 uint           `gorm:"primaryKey" json:"id"`
	RunID              uint           `gorm:"index;not null" json:"runId"`
	EntrySignalDate    time.Time      `json:"entrySignalDate"`
	EntryExecutionDate time.Time      `json:"entryExecutionDate"`
	EntryPrice         float64        `json:"entryPrice"`
	EntryFee           float64        `json:"entryFee"`
	EntrySnapshot      SignalSnapshot `gorm:"serializer:json;type:jsonb" json:"entrySnapshot"`
	ExitSignalDate     *time.Time     `json:"exitSignalDate,omitempty"`
	ExitExecutionDate  *time.Time     `json:"exitExecutionDate,omitempty"`
	ExitPrice          float64        `json:"exitPrice"`
	ExitFee            float64        `json:"exitFee"`
	ExitSnapshot       SignalSnapshot `gorm:"serializer:json;type:jsonb" json:"exitSnapshot"`
	Quantity           float64        `json:"quantity"`
	PnL                float64        `json:"pnl"`
	PnLPercent         float64        `json:"pnlPercent"`
	HoldingDays        int            `json:"holdingDays"`
	ExitReason         string         `gorm:"size:32;index" json:"exitReason"`
	Note               string         `gorm:"size:2000" json:"note"`
	CreatedAt          time.Time      `json:"createdAt"`
	UpdatedAt          time.Time      `json:"updatedAt"`
}

func (StrategyTrade) TableName() string { return "strategy_trades" }

type EquitySnapshot struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	RunID           uint      `gorm:"uniqueIndex:ux_run_equity_date;not null" json:"runId"`
	Date            time.Time `gorm:"uniqueIndex:ux_run_equity_date" json:"date"`
	Close           float64   `json:"close"`
	Equity          float64   `json:"equity"`
	BenchmarkEquity float64   `json:"benchmarkEquity"`
	StrategyReturn  float64   `json:"strategyReturn"`
	BenchmarkReturn float64   `json:"benchmarkReturn"`
	RunningPeak     float64   `json:"runningPeak"`
	Drawdown        float64   `json:"drawdown"`
	CreatedAt       time.Time `json:"createdAt"`
}

func (EquitySnapshot) TableName() string { return "strategy_equity_snapshots" }

type Evaluation struct {
	ID        uint      `gorm:"primaryKey"`
	RunID     uint      `gorm:"uniqueIndex:ux_run_evaluation_date"`
	Date      time.Time `gorm:"uniqueIndex:ux_run_evaluation_date"`
	CreatedAt time.Time
}

func (Evaluation) TableName() string { return "strategy_run_evaluations" }

type Metrics struct {
	CurrentEquity      float64  `json:"currentEquity"`
	TotalReturn        float64  `json:"totalReturn"`
	RealizedPnL        float64  `json:"realizedPnl"`
	UnrealizedPnL      float64  `json:"unrealizedPnl"`
	MaxDrawdown        float64  `json:"maxDrawdown"`
	WinRate            float64  `json:"winRate"`
	NumberOfTrades     int      `json:"numberOfTrades"`
	AverageProfit      float64  `json:"averageProfit"`
	AverageLoss        float64  `json:"averageLoss"`
	BestTrade          float64  `json:"bestTrade"`
	WorstTrade         float64  `json:"worstTrade"`
	AverageHoldingDays float64  `json:"averageHoldingDays"`
	TotalFees          float64  `json:"totalFees"`
	ProfitFactor       *float64 `json:"profitFactor"`
	BenchmarkReturn    float64  `json:"benchmarkReturn"`
	Difference         float64  `json:"difference"`
}
type PricePoint struct {
	Date  time.Time
	Close float64
	High  float64
	Low   float64
}
