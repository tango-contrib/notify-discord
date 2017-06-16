package discord

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/lunny/tango"
)

func TestNotify(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	tg := tango.Classic()
	tg.Use(Discord(Options{
		WebhookID:    "",
		WebhookToken: "",
		Source:       "",
	}))

	b := bytes.NewBufferString("testest")
	req, err := http.NewRequest("POST", "http://localhost:8000/", b)
	if err != nil {
		t.Error(err)
	}

	tg.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusNotFound)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "You are a writer")
}

/* Test Helpers */
func expect(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Errorf("Expected %v (type %v) - Got %v (type %v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
	}
}

func refute(t *testing.T, a interface{}, b interface{}) {
	if a == b {
		t.Errorf("Did not expect %v (type %v) - Got %v (type %v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
	}
}
