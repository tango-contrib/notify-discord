package discord

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

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
		Source:       "discord-test",
	}))
	tg.Post("/", func(ctx *tango.Context) {
		ctx.Abort(500, "test2 discord error")
	})

	b := bytes.NewBufferString("testest")
	req, err := http.NewRequest("POST", "http://localhost:8000/", b)
	if err != nil {
		t.Error(err)
	}

	tg.ServeHTTP(recorder, req)

	time.Sleep(time.Second * 5)
	expect(t, recorder.Code, 500)
	refute(t, len(buff.String()), 0)
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
