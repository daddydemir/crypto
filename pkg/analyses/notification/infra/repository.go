package infra

import (
	"log/slog"

	"github.com/daddydemir/crypto/pkg/analyses/notification/domain"
	"github.com/daddydemir/crypto/pkg/cache"
)

type Repository struct {
	cacheService cache.Cache
}

type Result struct {
	Tip             string `json:"tip"`
	Coin            string `json:"coin"`
	Gorsel          string `json:"gorsel"`
	OlusturmaZamani int64  `json:"olusturma_zamani"`
}

func NewRepository(cacheService cache.Cache) *Repository {
	return &Repository{
		cacheService: cacheService,
	}
}

func (r *Repository) GetAll() []domain.Notification {
	var result []Result
	err := r.cacheService.GetZList("kripto:grafikler", &result)
	if err != nil {
		slog.Error("Notification::GetAll", "error", err)
		return nil
	}
	var response []domain.Notification
	for _, v := range result {
		response = append(response, domain.Notification{
			Type:       v.Tip,
			Coin:       v.Coin,
			Image:      v.Gorsel,
			CreateTime: v.OlusturmaZamani,
		})
	}
	return response
}
