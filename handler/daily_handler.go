package handler

import (
	"encoding/json"
	"github.com/daddydemir/crypto/config/log"
	"github.com/daddydemir/crypto/pkg/dao"
	"github.com/daddydemir/crypto/pkg/database/service"
	"io/ioutil"
	"net/http"
)

func dailyStart(w http.ResponseWriter, r *http.Request) {
	service.CreateDaily(true)
}

func dailyEnd(w http.ResponseWriter, r *http.Request) {
	service.CreateDaily(false)
}

func daily(w http.ResponseWriter, r *http.Request) {
	response := service.GetDaily()
	log.Infoln("::daily:: response:{}", response)
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Errorln("::daily:: err:{}", err)
	}
}

func getDaily(w http.ResponseWriter, r *http.Request) {
	var request dao.Date
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Errorln("::getDaily::ReadAll err:{}", err)
	}
	err = json.Unmarshal(body, &request)
	if err != nil {
		log.Errorln("::getDaily::Unmarshal err:{}", err)
	}
	response := service.GetDailyFromDb(request)
	log.Infoln("::getDaily:: request:{} response:{}", request, response)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Errorln("::getDaily:: err:{}", err)
	}
}

func getDailyWithId(w http.ResponseWriter, r *http.Request) {
	var request dao.Date
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Errorln("::getDailyWithId::ReadAll err:{}", err)
	}
	err = json.Unmarshal(body, &request)
	if err != nil {
		log.Errorln("::getDailyWithId::Unmarshal err:{}", err)
	}
	response := service.GetDailyWithId(request)
	log.Infoln("::getDailyWithId:: request:{} response:{}", request, response)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Errorln("::getDailyWithId:: err:{}", err)
	}
}
