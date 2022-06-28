package controller

import (
	"encoding/json"
	"github.com/goodsign/monday"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

func TestHomeHandler(t *testing.T) {
	req, err := http.NewRequest(GET.String(), "/", nil)

	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(HomeHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var want = map[string]bool{
		"ok": true,
	}
	var got map[string]bool

	errGot := json.Unmarshal(rr.Body.Bytes(), &got)
	if errGot != nil {
		log.Fatal(errGot)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("handler returned unexpected body: got %v want %v", got, want)
	}
}

func TestCurrentTime(t *testing.T) {
	req, err := http.NewRequest(GET.String(), "/time", nil)

	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(CurrentTime)

	handler.ServeHTTP(rr, req)
	timeNow := time.Now()

	var want = map[string]string{
		"date": monday.Format(timeNow, monday.FullFormatsByLocale[monday.LocaleEnUS], monday.LocaleEnUS),
		"time": timeNow.Format(time.Kitchen),
	}

	var got map[string]string

	errGot := json.Unmarshal(rr.Body.Bytes(), &got)

	if errGot != nil {
		t.Fatal(errGot)
	}

	if !reflect.DeepEqual(want, got) {
		t.Errorf("handler returned unexpected body: got %v want %v", got, want)
	}
}
