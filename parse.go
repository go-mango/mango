package mango

import (
	"encoding/json"
	"errors"
	"net/http"
	"reflect"
	"strconv"
)

func ParseQuery[Q any](c *Context) *Q {
	q := new(Q)
	v := reflect.ValueOf(q).Elem()
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if !f.IsExported() {
			Abort(http.StatusInternalServerError, errors.New("query struct can not contain unexported fields"))
		}
		tag := f.Tag.Get("query")
		if tag == "" {
			tag = f.Name
		}
		str := c.Request().URL.Query().Get(tag)
		setFieldValue(v.Field(i), str)
	}
	if err := c.app.validate(q); err != nil {
		Abort(http.StatusUnprocessableEntity, err)
	}
	return q
}

func ParsePath[P any](c *Context) *P {
	p := new(P)
	v := reflect.ValueOf(p).Elem()
	t := v.Type()
	req := c.Request()
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if !f.IsExported() {
			Abort(http.StatusInternalServerError, errors.New("path struct can not contain unexported fields"))
		}
		tag := f.Tag.Get("path")
		if tag == "" {
			tag = f.Name
		}
		str := req.PathValue(tag)
		setFieldValue(v.Field(i), str)
	}
	if err := c.app.validate(p); err != nil {
		Abort(http.StatusUnprocessableEntity, err)
	}
	return p
}

func ParseBody[B any](c *Context) *B {
	b := new(B)
	if err := json.NewDecoder(c).Decode(b); err != nil {
		Abort(http.StatusBadRequest, err)
	}
	if err := c.app.validate(b); err != nil {
		Abort(http.StatusUnprocessableEntity, err)
	}
	return b
}

func setFieldValue(v reflect.Value, str string) {
	switch v.Kind() {
	case reflect.String:
		v.SetString(str)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		parsedVal, err := strconv.ParseInt(str, 10, 64)
		if err != nil {
			Abort(http.StatusBadRequest, err)
		}
		v.SetInt(parsedVal)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		parsedVal, err := strconv.ParseUint(str, 10, 64)
		if err != nil {
			Abort(http.StatusBadRequest, err)
		}
		v.SetUint(parsedVal)
	case reflect.Float32, reflect.Float64:
		parsedVal, err := strconv.ParseFloat(str, 64)
		if err != nil {
			Abort(http.StatusBadRequest, err)
		}
		v.SetFloat(parsedVal)
	default:
		Abort(http.StatusInternalServerError, errors.New("unsupported query param type"))
	}
}
