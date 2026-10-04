package app

import (
	"fmt"
	"github.com/daddydemir/crypto/pkg/purchasegoals/domain"
	"github.com/daddydemir/crypto/pkg/purchasegoals/infra"
	"time"
)

type App struct{ repo *infra.Repository }

func New(repo *infra.Repository) *App { return &App{repo: repo} }
func view(g domain.Goal) domain.GoalView {
	target, purchased := 0.0, 0.0
	for _, i := range g.Items {
		v := domain.ItemTarget(i)
		if i.Status == domain.StatusPurchased {
			if i.ActualPurchasePrice != nil {
				purchased += *i.ActualPurchasePrice
			}
		} else if i.Status != domain.StatusCancelled {
			target += v
		}
	}
	saved := 0.0
	for _, s := range g.Savings {
		saved += s.Amount
	}
	remaining := target - saved
	if remaining < 0 {
		remaining = 0
	}
	progress := 0.0
	if target > 0 {
		progress = saved / target * 100
	}
	status := g.Status
	if status != domain.StatusPurchased && status != domain.StatusCancelled {
		if target > 0 && saved >= target {
			status = "READY"
		} else if saved != 0 {
			status = domain.StatusSaving
		} else {
			status = domain.StatusWanted
		}
	}
	var months *float64
	if remaining > 0 && g.ExpectedMonthlySaving > 0 {
		m := remaining / g.ExpectedMonthlySaving
		months = &m
	}
	return domain.GoalView{Goal: g, TargetAmount: target, SavedAmount: saved, RemainingAmount: remaining, Progress: progress, DisplayStatus: status, PurchasedValue: purchased, RemainingItemValue: target, EstimatedMonths: months}
}
func (a *App) List(user string) ([]domain.GoalView, error) {
	gs, e := a.repo.List(user)
	if e != nil {
		return nil, e
	}
	out := make([]domain.GoalView, len(gs))
	for i := range gs {
		out[i] = view(gs[i])
	}
	return out, nil
}
func (a *App) Get(user string, id uint) (*domain.GoalView, error) {
	g, e := a.repo.Find(user, id)
	if e != nil {
		return nil, e
	}
	v := view(*g)
	return &v, nil
}
func (a *App) Create(user string, g domain.Goal) (*domain.GoalView, error) {
	g.Username = user
	if e := g.NormalizeAndValidate(); e != nil {
		return nil, e
	}
	for i := range g.Items {
		if e := g.Items[i].NormalizeAndValidate(g.Currency); e != nil {
			return nil, fmt.Errorf("item: %w", e)
		}
	}
	if e := a.repo.Create(&g); e != nil {
		return nil, e
	}
	return a.Get(user, g.ID)
}
func (a *App) Update(user string, id uint, in domain.Goal) (*domain.GoalView, error) {
	g, e := a.repo.Find(user, id)
	if e != nil {
		return nil, e
	}
	g.Name, g.Notes, g.Priority, g.TargetDate, g.ExpectedMonthlySaving = in.Name, in.Notes, in.Priority, in.TargetDate, in.ExpectedMonthlySaving
	if e = g.NormalizeAndValidate(); e != nil {
		return nil, e
	}
	if e = a.repo.SaveGoal(g); e != nil {
		return nil, e
	}
	return a.Get(user, id)
}
func (a *App) Delete(user string, id uint) error { return a.repo.DeleteGoal(user, id) }
func (a *App) AddItem(user string, goalID uint, i domain.Item) (*domain.GoalView, error) {
	g, e := a.repo.Find(user, goalID)
	if e != nil {
		return nil, e
	}
	i.GoalID = goalID
	i.ID = 0
	i.Status = domain.StatusWanted
	i.ActualPurchasePrice = nil
	i.PurchasedAt = nil
	if e = i.NormalizeAndValidate(g.Currency); e != nil {
		return nil, e
	}
	if e = a.repo.CreateItem(&i); e != nil {
		return nil, e
	}
	return a.Get(user, goalID)
}
func (a *App) UpdateItem(user string, id uint, in domain.Item) (*domain.GoalView, error) {
	i, g, e := a.repo.FindItem(user, id)
	if e != nil {
		return nil, e
	}
	oldPrice := i.CurrentPrice
	in.ID = i.ID
	in.GoalID = i.GoalID
	in.Currency = g.Currency
	in.CreatedAt = i.CreatedAt
	in.ActualPurchasePrice = i.ActualPurchasePrice
	in.PurchasedAt = i.PurchasedAt
	in.PriceHistory = i.PriceHistory
	if e = in.NormalizeAndValidate(g.Currency); e != nil {
		return nil, e
	}
	if e = a.repo.SaveItem(&in); e != nil {
		return nil, e
	}
	if oldPrice != in.CurrentPrice {
		if e = a.repo.UpdatePrice(&in, in.CurrentPrice, in.Currency, time.Now()); e != nil {
			return nil, e
		}
	}
	return a.Get(user, g.ID)
}
func (a *App) DeleteItem(user string, id uint) error { return a.repo.DeleteItem(user, id) }
func (a *App) UpdatePrice(user string, id uint, price float64, date time.Time) (*domain.GoalView, error) {
	i, g, e := a.repo.FindItem(user, id)
	if e != nil {
		return nil, e
	}
	if price < 0 {
		return nil, fmt.Errorf("price cannot be negative")
	}
	if date.IsZero() {
		date = time.Now()
	}
	if e = a.repo.UpdatePrice(i, price, g.Currency, date); e != nil {
		return nil, e
	}
	return a.Get(user, g.ID)
}
func (a *App) Saving(user string, goalID, id uint, in domain.SavingTransaction) (*domain.GoalView, error) {
	if in.Amount == 0 {
		return nil, fmt.Errorf("saving amount cannot be zero")
	}
	if in.Date.IsZero() {
		in.Date = time.Now()
	}
	if id == 0 {
		if _, e := a.repo.Find(user, goalID); e != nil {
			return nil, e
		}
		total, e := a.repo.SavingTotal(goalID, 0)
		if e != nil {
			return nil, e
		}
		if total+in.Amount < 0 {
			return nil, infra.ErrNegativeSavings
		}
		in.GoalID = goalID
		e = a.repo.AddSaving(&in)
		if e != nil {
			return nil, e
		}
	} else {
		s, e := a.repo.FindSaving(user, goalID, id)
		if e != nil {
			return nil, e
		}
		total, e := a.repo.SavingTotal(goalID, id)
		if e != nil {
			return nil, e
		}
		if total+in.Amount < 0 {
			return nil, infra.ErrNegativeSavings
		}
		s.Amount, s.Date, s.Note = in.Amount, in.Date, in.Note
		if e = a.repo.SaveSaving(s); e != nil {
			return nil, e
		}
	}
	return a.Get(user, goalID)
}
func (a *App) DeleteSaving(user string, goalID, id uint) (*domain.GoalView, error) {
	if e := a.repo.DeleteSaving(user, goalID, id); e != nil {
		return nil, e
	}
	return a.Get(user, goalID)
}
func (a *App) Purchase(user string, id uint, price float64, date time.Time) (*domain.GoalView, error) {
	i, g, e := a.repo.FindItem(user, id)
	if e != nil {
		return nil, e
	}
	if price < 0 {
		return nil, fmt.Errorf("purchase price cannot be negative")
	}
	if date.IsZero() {
		date = time.Now()
	}
	i.Status = domain.StatusPurchased
	i.ActualPurchasePrice = &price
	i.PurchasedAt = &date
	if e = a.repo.SaveItem(i); e != nil {
		return nil, e
	}
	updated, e := a.repo.Find(user, g.ID)
	if e != nil {
		return nil, e
	}
	allPurchased, hasPurchased := true, false
	for _, candidate := range updated.Items {
		if candidate.Status == domain.StatusCancelled {
			continue
		}
		hasPurchased = true
		if candidate.Status != domain.StatusPurchased {
			allPurchased = false
		}
	}
	if hasPurchased && allPurchased {
		updated.Status = domain.StatusPurchased
		if e = a.repo.SaveGoal(updated); e != nil {
			return nil, e
		}
	}
	return a.Get(user, g.ID)
}
func (a *App) CancelItem(user string, id uint) (*domain.GoalView, error) {
	i, g, e := a.repo.FindItem(user, id)
	if e != nil {
		return nil, e
	}
	i.Status = domain.StatusCancelled
	if e = a.repo.SaveItem(i); e != nil {
		return nil, e
	}
	return a.Get(user, g.ID)
}
func (a *App) CancelGoal(user string, id uint) (*domain.GoalView, error) {
	g, e := a.repo.Find(user, id)
	if e != nil {
		return nil, e
	}
	g.Status = domain.StatusCancelled
	if e = a.repo.SaveGoal(g); e != nil {
		return nil, e
	}
	return a.Get(user, id)
}
