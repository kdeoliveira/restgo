package endpoints

import (
	"encoding/json"
	"github.com/goodsign/monday"
	"github.com/gorilla/websocket"
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
		errorResponse(err, response)
	}
}

func AllLocales(response http.ResponseWriter, _ *http.Request) {
	out := monday.ListLocales()
	err := json.NewEncoder(response).Encode(
		map[string]interface{}{
			"locale":  out,
			"default": out[0],
		})

	if err != nil {
		errorResponse(err, response)
	}
}

func WsSendDateTime(conn *websocket.Conn) {

	sendTicker := time.NewTicker(60 * time.Second)
	defer func() {
		sendTicker.Stop()
		conn.Close()
	}()

	for {
		select {
		case <-sendTicker.C:
			timeNow := time.Now()
			var val, err = json.Marshal(map[string]string{
				"date": monday.Format(timeNow, monday.FullFormatsByLocale[monday.LocaleEnUS], monday.LocaleEnUS),
				"time": timeNow.Format(time.Kitchen),
			})

			if err != nil {
				return
			}

			err = conn.SetWriteDeadline(time.Now().Add(10 * time.Second))

			if err != nil {
				return
			}

			if err := conn.WriteMessage(websocket.TextMessage, val); err != nil {
				return
			}
		}
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

func errorResponse(error error, response http.ResponseWriter) {
	response.WriteHeader(http.StatusBadRequest)
	response.Write([]byte(error.Error()))
}

func Pipeline(response http.ResponseWriter, request *http.Request) {

}
