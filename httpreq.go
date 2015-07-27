// Package httpreq is a set of helper to extract data from HTTP Request.
package httpreq

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

// Common errors.
var (
	ErrWrongType = errors.New("wrong type for the given convertion function")
)

// ParsingMapElem represent the needed elements to parse a given element.
// - `Field` to be pulled from the `given Getter() interface`
// - `Fct`   the transform function betweend `Getter(Field)` and `Dest`
// - `Dest`   where to store the result.
type ParsingMapElem struct {
	Field string
	Fct   func(string, interface{}) error
	Dest  interface{}
}

// ParsingMap is a list of ParsingMapElem.
type ParsingMap []ParsingMapElem

// Getter is the basic interface to extract the intput data.
// Commonly used with http.Request.Form.
type Getter interface {
	// Get key return value.
	Get(string) string
}

// NewParsingMap create a new parsing map and returns a pointer to
// be able to call Add directly.
func NewParsingMap() *ParsingMap {
	return &ParsingMap{}
}

// NewParsingMapPre create a new preallocated parsing map and returns a pointer to
// be able to call Add directly.
func NewParsingMapPre(n int) *ParsingMap {
	p := make(ParsingMap, 0, n)
	return &p
}

// Add inserts a new field definition in the ParsingMap.
func (p *ParsingMap) Add(field string, fct func(string, interface{}) error, dest interface{}) *ParsingMap {
	*p = append(*p, ParsingMapElem{Field: field, Fct: fct, Dest: dest})
	return p
}

// Parse walks through the ParsingMap and executes it.
func (p ParsingMap) Parse(in Getter) error {
	for _, elem := range p {
		if e := in.Get(elem.Field); e != "" {
			if err := elem.Fct(e, elem.Dest); err != nil {
				return err
			}
		}
	}
	return nil
}

// ToCommaList takes the given string and splits it on `,` then set the resulting slice to `dest`.
func ToCommaList(src string, dest interface{}) error {
	d, ok := dest.(*[]string)
	if !ok {
		return ErrWrongType
	}
	*d = strings.Split(src, ",")
	return nil
}

// ToCommaList is a helper for ToCommaList.
func (p *ParsingMap) ToCommaList(field string, dest interface{}) *ParsingMap {
	*p = append(*p, ParsingMapElem{Field: field, Fct: ToCommaList, Dest: dest})
	return p
}

// ToString takes the given string and sets it to `dest`.
func ToString(src string, dest interface{}) error {
	d, ok := dest.(*string)
	if !ok {
		return ErrWrongType
	}
	*d = src
	return nil
}

// ToString is a helper for ToString.
func (p *ParsingMap) ToString(field string, dest interface{}) *ParsingMap {
	*p = append(*p, ParsingMapElem{Field: field, Fct: ToString, Dest: dest})
	return p
}

// ToBool takes the given string, parses it as bool and sets it to `dest`.
// NOTE: considers empty/invalid value as false
func ToBool(src string, dest interface{}) error {
	d, ok := dest.(*bool)
	if !ok {
		return ErrWrongType
	}
	if src == "on" {
		*d = true
		return nil
	}
	b, _ := strconv.ParseBool(src)
	*d = b
	return nil
}

// ToBool is a helper for ToBool.
func (p *ParsingMap) ToBool(field string, dest interface{}) *ParsingMap {
	*p = append(*p, ParsingMapElem{Field: field, Fct: ToBool, Dest: dest})
	return p
}

// ToInt takes the given string, parses it as int and sets it to `dest`.
func ToInt(src string, dest interface{}) error {
	d, ok := dest.(*int)
	if !ok {
		return ErrWrongType
	}
	i, err := strconv.Atoi(src)
	if err != nil {
		return err
	}
	*d = i
	return nil
}

// ToInt is a helper for ToInt.
func (p *ParsingMap) ToInt(field string, dest interface{}) *ParsingMap {
	*p = append(*p, ParsingMapElem{Field: field, Fct: ToInt, Dest: dest})
	return p
}

// ToFloat64 takes the given string, parses it as float64 and sets it to `dest`.
func ToFloat64(src string, dest interface{}) error {
	d, ok := dest.(*float64)
	if !ok {
		return ErrWrongType
	}
	f, err := strconv.ParseFloat(src, 64)
	if err != nil {
		return err
	}
	*d = f
	return nil
}

// ToFloat64 is a helper for ToFloat64.
func (p *ParsingMap) ToFloat64(field string, dest interface{}) *ParsingMap {
	*p = append(*p, ParsingMapElem{Field: field, Fct: ToFloat64, Dest: dest})
	return p
}

// ToTSTime takes the given string, parses it as timestamp and sets it to `dest`.
func ToTSTime(src string, dest interface{}) error {
	ts, err := strconv.ParseInt(src, 10, 64)
	if err != nil {
		return err
	}
	t := time.Unix(ts, 0)

	d, ok := dest.(**time.Time)
	if !ok {
		d, ok := dest.(*time.Time)
		if !ok {
			return ErrWrongType
		}
		*d = t
		return nil
	}
	*d = &t
	return nil
}

// ToTSTime is a helper for ToTSTime.
func (p *ParsingMap) ToTSTime(field string, dest interface{}) *ParsingMap {
	*p = append(*p, ParsingMapElem{Field: field, Fct: ToTSTime, Dest: dest})
	return p
}

// ToRFC3339Time takes the given string, parses it as timestamp and sets it to `dest`.
func ToRFC3339Time(src string, dest interface{}) error {
	t, err := time.Parse(time.RFC3339, src)
	if err != nil {
		return err
	}

	d, ok := dest.(**time.Time)
	if !ok {
		d, ok := dest.(*time.Time)
		if !ok {
			return ErrWrongType
		}
		*d = t
		return nil
	}
	*d = &t
	return nil
}

// ToRFC3339Time is a helper for ToRFC3339Time.
func (p *ParsingMap) ToRFC3339Time(field string, dest interface{}) *ParsingMap {
	*p = append(*p, ParsingMapElem{Field: field, Fct: ToRFC3339Time, Dest: dest})
	return p
}
