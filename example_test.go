package httpreq

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"time"
)

func Example() {
	defer os.Setenv("TZ", os.Getenv("TZ"))
	os.Setenv("TZ", "UTC")

	type Req struct {
		Fields    []string
		Limit     int
		Page      int
		Timestamp time.Time
	}
	hdlr := func(w http.ResponseWriter, req *http.Request) {
		_ = req.ParseForm()
		data := &Req{}
		if err := (ParsingMap{
			{Field: "limit", Fct: ToInt, Dest: &data.Limit},
			{Field: "page", Fct: ToInt, Dest: &data.Page},
			{Field: "fields", Fct: ToCommaList, Dest: &data.Fields},
			{Field: "timestamp", Fct: ToTSTime, Dest: &data.Timestamp},
		}.Parse(req.Form)); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		_ = json.NewEncoder(w).Encode(data)
	}
	ts := httptest.NewServer(http.HandlerFunc(hdlr))
	defer ts.Close()

	resp, err := http.Get(ts.URL + "?timestamp=1437743020&limit=10&page=1&fields=a,b,c")
	if err != nil {
		log.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Error: %d", resp.StatusCode)
	}
	_, _ = io.Copy(os.Stdout, resp.Body)
	// output:
	// {"Fields":["a","b","c"],"Limit":10,"Page":1,"Timestamp":"2015-07-24T13:03:40Z"}
}

func Example_chain() {
	type Req struct {
		Fields []string
		Limit  int
		Page   int
	}
	hdlr := func(w http.ResponseWriter, req *http.Request) {
		_ = req.ParseForm()
		data := &Req{}
		if err := NewParsingMap().
			Add("limit", ToInt, &data.Limit).
			Add("page", ToInt, &data.Page).
			Add("fields", ToCommaList, &data.Fields).
			Parse(req.Form); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		_ = json.NewEncoder(w).Encode(data)
	}
	ts := httptest.NewServer(http.HandlerFunc(hdlr))
	defer ts.Close()

	resp, err := http.Get(ts.URL + "?timestamp=1437743020&limit=10&page=1&fields=a,b,c")
	if err != nil {
		log.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Error: %d", resp.StatusCode)
	}
	_, _ = io.Copy(os.Stdout, resp.Body)
	// output:
	// {"Fields":["a","b","c"],"Limit":10,"Page":1}
}
