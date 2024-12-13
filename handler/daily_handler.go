package handler

import (
	"encoding/json"
	"io/ioutil"
	"log/slog"
	"net/http"

	"github.com/daddydemir/crypto/pkg/dao"
	"github.com/daddydemir/crypto/pkg/service"
)

func dailyStart(_ http.ResponseWriter, _ *http.Request) {
	//service.CreateDaily(true)
}

func dailyEnd(_ http.ResponseWriter, _ *http.Request) {
	//service.CreateDaily(false)
}

func daily(w http.ResponseWriter, _ *http.Request) {
	response := service.GetDaily()
	slog.Info("daily:service.GetDaily", "dailies", response)
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		slog.Error("daily:json.Encode", "error", err)
	}
}

func getDaily(w http.ResponseWriter, r *http.Request) {
	var request dao.Date
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		slog.Error("getDaily:ioutil.ReadAll", "error", err)
	}
	err = json.Unmarshal(body, &request)
	if err != nil {
		slog.Error("getDaily:json.Unmarshal", "error", err)
	}
	dailyService := serviceFactory.NewDailyService()
	response, err := dailyService.FindByDateRange(request.StartDate, request.EndDate)
	if err != nil {
		slog.Error("daily:FindByDateRange", "error", err)
	}
	slog.Info("daily:service.GetDailyFromDb", "request", request, "response", response)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		slog.Error("daily:json.Encode", "error", err)
	}
}

func getDailyWithId(w http.ResponseWriter, r *http.Request) {
	var request dao.Date
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		slog.Error("getDailyWithId:ioutil.ReadAll", "error", err)
	}
	err = json.Unmarshal(body, &request)
	if err != nil {
		slog.Error("getDailyWithId:json.Unmarshal", "error", err)
	}
	dailyService := serviceFactory.NewDailyService()
	response, err := dailyService.FindByIdAndDateRange(request.Id, request.StartDate, request.EndDate)
	if err != nil {
		slog.Error("daily:FindByIdAndDateRange", "error", err)
	}
	slog.Info("daily:service.GetDailyWithId", "request", request, "response", response)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		slog.Error("daily:json.Encode", "error", err)
	}
}
