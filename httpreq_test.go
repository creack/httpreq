package httpreq

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"
)

func TestHandler(t *testing.T) {
	origTZ := os.Getenv("TZ")
	defer func() { _ = os.Setenv("TZ", origTZ) }()
	_ = os.Setenv("TZ", "UTC")

	type Req struct {
		Fields    []string
		Limit     int
		Page      int
		Timestamp time.Time
		F         float64
		B         bool
		Time      time.Time
	}
	hdlr := func(w http.ResponseWriter, req *http.Request) {
		if err := req.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		data := &Req{}
		if err := (ParsingMap{
			{Field: "limit", Fct: ToInt, Dest: &data.Limit},
			{Field: "page", Fct: ToInt, Dest: &data.Page},
			{Field: "fields", Fct: ToCommaList, Dest: &data.Fields},
			{Field: "timestamp", Fct: ToTSTime, Dest: &data.Timestamp},
			{Field: "f", Fct: ToFloat64, Dest: &data.F},
			{Field: "b", Fct: ToBool, Dest: &data.B},
			{Field: "t", Fct: ToRFC3339Time, Dest: &data.Time},
		}.Parse(req.Form)); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		_ = json.NewEncoder(w).Encode(data)
	}
	ts := httptest.NewServer(http.HandlerFunc(hdlr))
	defer ts.Close()

	resp, err := http.Get(ts.URL + "?timestamp=1437743020&limit=10&page=1&fields=a,b,c&f=1.5&b=true&t=2006-01-02T15:04:05Z")
	if err != nil {
		t.Fatal(err)
	}
	if expect, got := resp.StatusCode, http.StatusOK; expect != got {
		t.Fatalf("Unexpected result.\nExpect:\t%d\nGot:\t%d\n", expect, got)
	}
	buf, err := ioutil.ReadAll(resp.Body)
	_ = resp.Body.Close()
	if err != nil {
		t.Fatal(err)
	}
	if expect, got := `{"Fields":["a","b","c"],"Limit":10,"Page":1,"Timestamp":"2015-07-24T13:03:40Z","F":1.5,"B":true,"Time":"2006-01-02T15:04:05Z"}`+"\n", string(buf); expect != got {
		t.Fatalf("Unexpected result.\nExpect:\t%s\nGot:\t%s\n", expect, got)
	}
}

type mockForm map[string]string

func (m mockForm) Get(k string) string { return m[k] }

func TestParseFail(t *testing.T) {
	src, dest := "abc", 0
	form := mockForm{"limit": src}
	if err := (ParsingMap{
		{Field: "limit", Fct: ToInt, Dest: &dest},
	}).Parse(form); err == nil {
		t.Fatal("Invalid element in parseMap should yield an error")
	}
}

func TestToCommaList(t *testing.T) {
	// List
	src, dest := "a,b,c", []string{}
	if err := ToCommaList(src, &dest); err != nil {
		t.Fatal(err)
	}
	if expect, got := src, strings.Join(dest, ","); expect != got {
		t.Fatalf("Unexpected result.\nExpect:\t%s\nGot:\t%s\n", expect, got)
	}

	// Single element list
	src, dest = "a", []string{}
	if err := ToCommaList(src, &dest); err != nil {
		t.Fatal(err)
	}
	if expect, got := src, strings.Join(dest, ","); expect != got {
		t.Fatalf("Unexpected result.\nExpect:\t%s\nGot:\t%s\n", expect, got)
	}

	// Empty list
	src, dest = "", []string{}
	if err := ToCommaList(src, &dest); err != nil {
		t.Fatal(err)
	}
	if expect, got := src, strings.Join(dest, ","); expect != got {
		t.Fatalf("Unexpected result.\nExpect:\t%s\nGot:\t%s\n", expect, got)
	}

	// Error check
	if err := ToCommaList(src, dest); err != ErrWrongType {
		t.Fatalf("Wrong type didn't yield the proper error: %v\n", err)
	}
}

func TestToString(t *testing.T) {
	src, dest := "hello", ""
	if err := ToString(src, &dest); err != nil {
		t.Fatal(err)
	}
	if expect, got := src, dest; expect != got {
		t.Fatalf("Unexpected result.\nExpect:\t%s\nGot:\t%s\n", expect, got)
	}

	// Error check
	if err := ToString(src, dest); err != ErrWrongType {
		t.Fatalf("Wrong type didn't yield the proper error: %v\n", err)
	}
}

func TestToBool(t *testing.T) {
	// "true"
	src, dest := "true", false
	if err := ToBool(src, &dest); err != nil {
		t.Fatal(err)
	}
	if expect, got := true, dest; expect != got {
		t.Fatalf("Unexpected result.\nExpect:\t%t\nGot:\t%t\n", expect, got)
	}

	// "1" -> true
	src, dest = "1", false
	if err := ToBool(src, &dest); err != nil {
		t.Fatal(err)
	}
	if expect, got := true, dest; expect != got {
		t.Fatalf("Unexpected result.\nExpect:\t%t\nGot:\t%t\n", expect, got)
	}

	// "on" -> true
	src, dest = "on", false
	if err := ToBool(src, &dest); err != nil {
		t.Fatal(err)
	}
	if expect, got := true, dest; expect != got {
		t.Fatalf("Unexpected result.\nExpect:\t%t\nGot:\t%t\n", expect, got)
	}

	// "false" -> false
	src, dest = "false", true
	if err := ToBool(src, &dest); err != nil {
		t.Fatal(err)
	}
	if expect, got := false, dest; expect != got {
		t.Fatalf("Unexpected result.\nExpect:\t%t\nGot:\t%t\n", expect, got)
	}

	// emtpy -> false
	src, dest = "", true
	if err := ToBool(src, &dest); err != nil {
		t.Fatal(err)
	}
	if expect, got := false, dest; expect != got {
		t.Fatalf("Unexpected result.\nExpect:\t%t\nGot:\t%t\n", expect, got)
	}

	// random -> false
	src, dest = "abcd", true
	if err := ToBool(src, &dest); err != nil {
		t.Fatal(err)
	}
	if expect, got := false, dest; expect != got {
		t.Fatalf("Unexpected result.\nExpect:\t%t\nGot:\t%t\n", expect, got)
	}

	// Error check
	if err := ToBool(src, dest); err != ErrWrongType {
		t.Fatalf("Wrong type didn't yield the proper error: %v\n", err)
	}
}

func TestToInt(t *testing.T) {
	src, dest := "42", 0
	if err := ToInt(src, &dest); err != nil {
		t.Fatal(err)
	}
	if expect, got := 42, dest; expect != got {
		t.Fatalf("Unexpected result.\nExpect:\t%d\nGot:\t%d\n", expect, got)
	}

	// Error check
	src, dest = "abc", -1
	if err := ToInt(src, &dest); err == nil {
		t.Fatal("Invalid integer should yield an error")
	}

	if err := ToInt(src, dest); err != ErrWrongType {
		t.Fatalf("Wrong type didn't yield the proper error: %v\n", err)
	}
}

func TestToFloat64(t *testing.T) {
	src, dest := "42.", 0.
	if err := ToFloat64(src, &dest); err != nil {
		t.Fatal(err)
	}
	if expect, got := 42., dest; expect != got {
		t.Fatalf("Unexpected result.\nExpect:\t%f\nGot:\t%f\n", expect, got)
	}

	// Error check
	src, dest = "abc", -1.
	if err := ToFloat64(src, &dest); err == nil {
		t.Fatal("Invalid integer should yield an error")
	}

	if err := ToFloat64(src, dest); err != ErrWrongType {
		t.Fatalf("Wrong type didn't yield the proper error: %v\n", err)
	}
}

func TestToTSTimeChain(t *testing.T) {
	src, dest := "1437743020", time.Time{}
	if err := NewParsingMap().ToTSTime("ts", &dest).Parse(mockForm{"ts": src}); err != nil {
		t.Fatal(err)
	}
	if expect, got := time.Unix(1437743020, 0), dest; expect.Sub(got) != 0 {
		t.Fatalf("Unexpected result.\nExpect:\t%s\nGot:\t%s\n", expect, got)
	}
}

func TestToTSTime(t *testing.T) {
	src, dest := "1437743020", time.Time{}
	if err := ToTSTime(src, &dest); err != nil {
		t.Fatal(err)
	}
	if expect, got := time.Unix(1437743020, 0), dest; expect.Sub(got) != 0 {
		t.Fatalf("Unexpected result.\nExpect:\t%s\nGot:\t%s\n", expect, got)
	}

	src, dest2 := "1437743020", &time.Time{}
	if err := ToTSTime(src, &dest2); err != nil {
		t.Fatal(err)
	}
	tt := time.Unix(1437743020, 0)
	if expect, got := &tt, dest2; expect.Sub(*got) != 0 {
		t.Fatalf("Unexpected result.\nExpect:\t%s\nGot:\t%s\n", expect, got)
	}

	// Error check
	if err := ToTSTime(src, dest); err != ErrWrongType {
		t.Fatalf("Wrong type didn't yield the proper error: %v\n", err)
	}
	src, dest = "abc", time.Time{}
	if err := ToTSTime(src, &dest); err == nil {
		t.Fatal("Invalid timestamp should yield an error")
	}
}

func TestToRFC3339TimeChain(t *testing.T) {
	src, dest := "2006-01-02T15:04:05Z", time.Time{}
	if err := NewParsingMap().ToRFC3339Time("date", &dest).Parse(mockForm{"date": src}); err != nil {
		t.Fatal(err)
	}
	tt, err := time.Parse(time.RFC3339, src)
	if err != nil {
		t.Fatal(err)
	}
	if expect, got := tt, dest; expect.Sub(got) != 0 {
		t.Fatalf("Unexpected result.\nExpect:\t%s\nGot:\t%s\n", expect, got)
	}
}

func TestToRFC3339Time(t *testing.T) {
	src, dest := "2006-01-02T15:04:05Z", time.Time{}
	if err := ToRFC3339Time(src, &dest); err != nil {
		t.Fatal(err)
	}
	tt, err := time.Parse(time.RFC3339, src)
	if err != nil {
		t.Fatal(err)
	}
	if expect, got := tt, dest; expect.Sub(got) != 0 {
		t.Fatalf("Unexpected result.\nExpect:\t%s\nGot:\t%s\n", expect, got)
	}

	src, dest2 := "2006-01-02T15:04:05Z", &time.Time{}
	if err := ToRFC3339Time(src, &dest2); err != nil {
		t.Fatal(err)
	}
	if expect, got := &tt, dest2; expect.Sub(*got) != 0 {
		t.Fatalf("Unexpected result.\nExpect:\t%s\nGot:\t%s\n", expect, got)
	}

	// Error check
	if err := ToRFC3339Time(src, dest); err != ErrWrongType {
		t.Fatalf("Wrong type didn't yield the proper error: %v\n", err)
	}
	src, dest = "abc", time.Time{}
	if err := ToRFC3339Time(src, &dest); err == nil {
		t.Fatal("Invalid timestamp should yield an error")
	}
}
