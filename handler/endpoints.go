package handler

import (
	"encoding/json"
	"github.com/goodsign/monday"
	"net/http"
	"time"
)

func HomeHandler(response http.ResponseWriter, _ *http.Request) {
	response.WriteHeader(http.StatusOK)
	err := json.NewEncoder(response).Encode(
		map[string]bool{
			"ok": true,
		})
	if err != nil {
		return
	}

}

func CurrentTime(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	response.Header().Set("Content-Type", "application/json")
	response.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	response.Header().Set("Access-Control-Allow-Credentials", "true")
	queries := request.URL.Query()
	timeNow := time.Now()
	var (
		output string
	)
	if val, ok := queries["locale"]; ok {
		if containLocale(monday.Locale(val[0])) {
			output = monday.Format(timeNow, monday.FullFormatsByLocale[monday.Locale(val[0])], monday.Locale(val[0]))
		} else {
			output = monday.Format(timeNow, monday.FullFormatsByLocale[monday.LocaleEnUS], monday.LocaleEnUS)
		}
	} else {
		output = monday.Format(timeNow, monday.FullFormatsByLocale[monday.LocaleEnUS], monday.LocaleEnUS)
	}

	response.WriteHeader(http.StatusOK)

	err := json.NewEncoder(response).Encode(
		map[string]string{
			"date": output,
			"time": timeNow.Format(time.Kitchen),
		})
	if err != nil {
		return
	}

}

func containLocale(locale monday.Locale) bool {
	for _, v := range monday.ListLocales() {
		if v == locale {
			return true

		}
	}
	return false
}
