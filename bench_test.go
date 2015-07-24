package httpreq

import (
	"encoding/json"
	"testing"
)

// Straight benchmark literal.
func BenchmarkRawLiteral(b *testing.B) {
	mock := mockForm{
		"limit":  "10",
		"page":   "1",
		"fields": "a,b,c",
	}
	type Req struct {
		Fields []string
		Limit  int
		Page   int
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		data := &Req{}
		if err := (ParsingMap{
			{Field: "limit", Fct: ToInt, Dest: &data.Limit},
			{Field: "page", Fct: ToInt, Dest: &data.Page},
			{Field: "fields", Fct: ToCommaList, Dest: &data.Fields},
		}.Parse(mock)); err != nil {
			b.Fatal(err)
		}
	}
}

// Straight benchmark Add chain.
func BenchmarkRawAdd(b *testing.B) {
	mock := mockForm{
		"limit":  "10",
		"page":   "1",
		"fields": "a,b,c",
	}
	type Req struct {
		Fields []string
		Limit  int
		Page   int
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		data := &Req{}
		if err := NewParsingMap().
			Add("limit", ToInt, &data.Limit).
			Add("page", ToInt, &data.Page).
			Add("fields", ToCommaList, &data.Fields).
			Parse(mock); err != nil {
			b.Fatal(err)
		}
	}
}

// Straight benchmark json.Unmarshal for comparison.
func BenchmarkRawJSONUnmarshal(b *testing.B) {
	mock := []byte(`{"Fields":["a","b","c"],"Limit":10,"Page":1}`)
	type Req struct {
		Fields []string
		Limit  int
		Page   int
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		data := &Req{}
		if err := json.Unmarshal(mock, data); err != nil {
			b.Fatal(err)
		}
	}
}

// Parallel benchmark literal.
func BenchmarkRawPLiteral(b *testing.B) {
	mock := mockForm{
		"limit":  "10",
		"page":   "1",
		"fields": "a,b,c",
	}
	type Req struct {
		Fields []string
		Limit  int
		Page   int
	}
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			data := &Req{}
			if err := (ParsingMap{
				{Field: "limit", Fct: ToInt, Dest: &data.Limit},
				{Field: "page", Fct: ToInt, Dest: &data.Page},
				{Field: "fields", Fct: ToCommaList, Dest: &data.Fields},
			}.Parse(mock)); err != nil {
				b.Fatal(err)
			}
		}
	})
}

// Parallel benchmark Add chain.
func BenchmarkRawPAdd(b *testing.B) {
	mock := mockForm{
		"limit":  "10",
		"page":   "1",
		"fields": "a,b,c",
	}
	type Req struct {
		Fields []string
		Limit  int
		Page   int
	}
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			data := &Req{}
			if err := NewParsingMap().
				Add("limit", ToInt, &data.Limit).
				Add("page", ToInt, &data.Page).
				Add("fields", ToCommaList, &data.Fields).
				Parse(mock); err != nil {
				b.Fatal(err)
			}
		}
	})
}

// Parallel benchmark json.Unmarshal for comparison.
func BenchmarkRawPJSONUnmarshal(b *testing.B) {
	mock := []byte(`{"Fields":["a","b","c"],"Limit":10,"Page":1}`)
	type Req struct {
		Fields []string
		Limit  int
		Page   int
	}
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			data := &Req{}
			if err := json.Unmarshal(mock, data); err != nil {
				b.Fatal(err)
			}
		}
	})
}
