package fastweb

import (
	"encoding/xml"
	"errors"
	"github.com/yebaoping/fastweb/pkg/json"
	"github.com/yebaoping/fastweb/pkg/str"
	"reflect"
	"strconv"
	"strings"
)

type Binder interface {
	Bind(ctx *HttpContext, val interface{}) error
}

type defaultBinder struct {
}

func DefaultBinder() *defaultBinder {
	return new(defaultBinder)
}

var ErrBindValueNoValid = errors.New("value is no valid")
var ErrBindValueNotPointer = errors.New("value is not pointer")
var ErrBindValueIsNil = errors.New("value is nil")
var ErrBindValueNotStruct = errors.New("value is not struct")
var ErrUnsupportedContentType = errors.New("unsupported content-type")

func (b *defaultBinder) Bind(ctx *HttpContext, val interface{}) error {
	v := reflect.ValueOf(val)
	if !v.IsValid() {
		return ErrBindValueNoValid
	}

	if v.Kind() != reflect.Ptr {
		return ErrBindValueNotPointer
	}

	if v.IsNil() {
		return ErrBindValueIsNil
	}

	v = v.Elem()
	if v.Kind() != reflect.Struct {
		return ErrBindValueNotStruct
	}

	contentType := str.String(ctx.Request.Header.ContentType())
	if strings.HasPrefix(contentType, ContentTypeFormUrlEncoded) {
		return b.bindFormUrlEncoded(ctx, v)
	} else if strings.HasPrefix(contentType, ContentTypeMultipartFormData) {
		return b.bindMultipartForm(ctx, v)
	} else if strings.HasPrefix(contentType, ContentTypeJson) {
		return b.bindJson(ctx, val)
	} else if strings.HasPrefix(contentType, ContentTypeXml) {
		return b.bindXml(ctx, val)
	} else {
		return ErrUnsupportedContentType
	}
}

func setFiledVal(field reflect.Value, vals [][]byte) error {
	if len(vals) == 0 || len(vals[0]) == 0 {
		return nil
	}

	val := str.String(vals[0])

	switch field.Kind() {
	case reflect.Bool:
		v, err := strconv.ParseBool(val)
		if err != nil {
			return err
		}

		field.SetBool(v)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return err
		}

		field.SetInt(v)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		v, err := strconv.ParseUint(val, 10, 64)
		if err != nil {
			return err
		}

		field.SetUint(v)
	case reflect.Float32, reflect.Float64:
		v, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return err
		}

		field.SetFloat(v)
	case reflect.String:
		field.SetString(val)
	case reflect.Array, reflect.Slice:
		t := field.Type().Elem()
		s := reflect.MakeSlice(field.Type(), 0, field.Cap())
		for _, bytes := range vals {
			v := reflect.New(t).Elem()

			if err := setFiledVal(v, [][]byte{bytes}); err != nil {
				return err
			}

			s = reflect.Append(s, v)
		}

		field.Set(s)
	}

	return nil
}

func bindForm(val reflect.Value, formVal map[string][][]byte) error {
	for i := 0; i < val.NumField(); i++ {
		f := val.Type().Field(i)
		fv := val.Field(i)

		if !fv.CanSet() {
			continue
		}

		name := f.Tag.Get("json")
		if name == "-" {
			continue
		}

		if name == "" {
			name = f.Name
		}

		vals, ext := formVal[name]
		if !ext {
			continue
		}

		if err := setFiledVal(fv, vals); err != nil {
			return err
		}
	}

	return nil
}

func (b *defaultBinder) bindFormUrlEncoded(ctx *HttpContext, val reflect.Value) error {
	vals := make(map[string][][]byte)
	ctx.PostArgs().VisitAll(func(key, value []byte) {
		_, ext := vals[string(key)]
		if ext {
			vals[string(key)] = append(vals[string(key)], value)
		} else {
			vals[string(key)] = [][]byte{value}
		}
	})

	return bindForm(val, vals)
}

func (b *defaultBinder) bindMultipartForm(ctx *HttpContext, val reflect.Value) error {
	form, err := ctx.MultipartForm()
	if err != nil {
		return err
	}

	formValue := make(map[string][][]byte)

	for k, vals := range form.Value {
		formValue[k] = make([][]byte, len(vals))
		for i, val := range vals {
			formValue[k][i] = str.Bytes(val)
		}
	}

	return bindForm(val, formValue)
}

func (b *defaultBinder) bindJson(ctx *HttpContext, val interface{}) error {
	body := ctx.PostBody()
	if len(body) == 0 {
		return nil
	}

	return json.Unmarshal(body, val)
}

func (b *defaultBinder) bindXml(ctx *HttpContext, val interface{}) error {
	body := ctx.PostBody()
	if len(body) == 0 {
		return nil
	}

	return xml.Unmarshal(body, val)
}
