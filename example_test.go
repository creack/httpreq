package httpreq_test

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"time"

	"github.com/creack/httpreq"
)

func Example() {
	origTZ := os.Getenv("TZ")
	defer func() { _ = os.Setenv("TZ", origTZ) }()
	_ = os.Setenv("TZ", "UTC")

	type Req struct {
		Fields    []string
		Limit     int
		Page      int
		Timestamp time.Time
	}
	hdlr := func(w http.ResponseWriter, req *http.Request) {
		_ = req.ParseForm()
		data := &Req{}
		if err := (httpreq.ParsingMap{
			0: {Field: "limit", Fct: httpreq.ToInt, Dest: &data.Limit},
			1: {Field: "page", Fct: httpreq.ToInt, Dest: &data.Page},
			2: {Field: "fields", Fct: httpreq.ToCommaList, Dest: &data.Fields},
			3: {Field: "timestamp", Fct: httpreq.ToTSTime, Dest: &data.Timestamp},
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

func Example_add() {
	type Req struct {
		Fields []string
		Limit  int
		Page   int
	}
	hdlr := func(w http.ResponseWriter, req *http.Request) {
		data := &Req{}
		if err := httpreq.NewParsingMap().
			Add("limit", httpreq.ToInt, &data.Limit).
			Add("page", httpreq.ToInt, &data.Page).
			Add("fields", httpreq.ToCommaList, &data.Fields).
			Parse(req.URL.Query()); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		_ = json.NewEncoder(w).Encode(data)
	}
	ts := httptest.NewServer(http.HandlerFunc(hdlr))
	defer ts.Close()

	resp, err := http.Get(ts.URL + "?limit=10&page=1&fields=a,b,c")
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

func Example_chain() {
	type Req struct {
		Fields    []string
		Limit     int
		DryRun    bool
		Threshold float64
		Name      string
	}
	hdlr := func(w http.ResponseWriter, req *http.Request) {
		data := &Req{}
		if err := httpreq.NewParsingMapPre(5).
			ToInt("limit", &data.Limit).
			ToBool("dryrun", &data.DryRun).
			ToFloat64("threshold", &data.Threshold).
			ToCommaList("fields", &data.Fields).
			ToString("name", &data.Name).
			Parse(req.URL.Query()); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		_ = json.NewEncoder(w).Encode(data)
	}
	ts := httptest.NewServer(http.HandlerFunc(hdlr))
	defer ts.Close()

	resp, err := http.Get(ts.URL + "?limit=10&dryrun=true&fields=a,b,c&threshold=42.5&name=creack")
	if err != nil {
		log.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Error: %d", resp.StatusCode)
	}
	_, _ = io.Copy(os.Stdout, resp.Body)
	// output:
	// {"Fields":["a","b","c"],"Limit":10,"DryRun":true,"Threshold":42.5,"Name":"creack"}
}
