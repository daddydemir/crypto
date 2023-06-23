package handler

import (
	"encoding/json"
	"github.com/daddydemir/crypto/pkg/dao"
	"github.com/daddydemir/crypto/pkg/database/service"
	"io/ioutil"
	"log"
	"net/http"
)

func dailyStart(w http.ResponseWriter, r *http.Request) {
	service.CreateDaily(true)
}

func dailyEnd(w http.ResponseWriter, r *http.Request) {
	service.CreateDaily(false)
}

func daily(w http.ResponseWriter, r *http.Request) {
	err := json.NewEncoder(w).Encode(service.GetDaily())
	if err != nil {
		log.Println("::daily:: err:{}", err)
	}
}

func getDaily(w http.ResponseWriter, r *http.Request) {
	var request dao.Date
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("::getDaily::ReadAll err:{}", err)
	}
	err = json.Unmarshal(body, &request)
	if err != nil {
		log.Println("::getDaily::Unmarshal err:{}", err)
	}
	err = json.NewEncoder(w).Encode(service.GetDailyFromDb(request))
	if err != nil {
		log.Println("::getDaily:: err:{}", err)
	}
}
