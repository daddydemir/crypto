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
	Username        string `json:"username"`
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

func (r *Repository) GetAll(username string) []domain.Notification {
	var result []Result
	err := r.cacheService.GetZList("kripto:grafikler", &result)
	if err != nil {
		slog.Error("Notification::GetAll", "error", err)
		return []domain.Notification{}
	}
	response := make([]domain.Notification, 0)
	for _, v := range result {
		if v.Username == "" || v.Username != username {
			continue
		}
		response = append(response, domain.Notification{
			Type:       v.Tip,
			Coin:       v.Coin,
			Image:      v.Gorsel,
			CreateTime: v.OlusturmaZamani,
		})
	}
	return response
}
