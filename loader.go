package env

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	defaultTagName = "env"
)

// tagInfo maintains information about the struct tags
type tagInfo struct {
	key    string
	defVal string
}

// Loader populates the specified struct based on environment variables
type Loader struct {
	once sync.Once
	opts *options
}

func New(options ...Option) *Loader {
	p := new(Loader)
	p.lazyInit()
	for _, o := range options {
		o(p.opts)
	}
	return p
}

func (p *Loader) lazyInit() {
	p.once.Do(func() {
		p.opts = &options{
			prefix:   "",
			tagName:  defaultTagName,
			override: false,
			getter:   &getter{},
		}
	})
}

func (p *Loader) Load(i interface{}) error {
	p.lazyInit()
	return p.loadInterface(i, p.opts.prefix)
}

func (p *Loader) loadInterface(i interface{}, prefix string) error {
	refVal := reflect.ValueOf(i)
	if refVal.Kind() != reflect.Ptr {
		return ErrNotStructPtr
	}

	refVal = refVal.Elem()
	if refVal.Kind() != reflect.Struct {
		return ErrNotStructPtr
	}

	return p.loadValue(refVal, prefix)
}

func (p *Loader) loadValue(refVal reflect.Value, prefix string) error {
	refType := refVal.Type()

	for i := 0; i < refType.NumField(); i++ {
		field := refVal.Field(i)
		if !field.CanSet() {
			continue
		}

		structField := refType.Field(i)
		tag, err := p.parseTags(structField)
		if err != nil {
			return err
		}
		if tag == nil {
			continue
		}

		// create a new object if nil pointer for struct-type
		if field.Kind() == reflect.Ptr && field.Type().Elem().Kind() == reflect.Struct {
			if field.IsNil() {
				field.Set(reflect.New(field.Type().Elem()))
			}
			field = field.Elem()
		}

		if field.Kind() == reflect.Struct && field.CanAddr() {
			if len(getSetters(field)) == 0 {
				if err := p.loadValue(field.Addr().Elem(), p.opts.getter.Merge(prefix, tag.key)); err != nil {
					return err
				}
				continue
			}
		}

		if !field.IsZero() && !p.opts.override {
			continue
		}

		// Get value by specified key
		key := p.opts.getter.Merge(prefix, tag.key)
		value, found, err := p.opts.getter.Get(key)
		if err != nil {
			return err
		}
		// Use default value if the key not be set and field value is zero.
		if !found && field.IsZero() {
			value = tag.defVal
		}
		if value == "" {
			continue
		}

		if err := setField(field, value); err != nil {
			return &ParseError{
				KeyName:   key,
				FieldName: refType.Name() + "." + structField.Name,
				TypeName:  field.Type().String(),
				Value:     value,
				Err:       err,
			}
		}
	}
	return nil
}

// parseTags split the struct tag's into the expected key and desired option, if any.
// return nil if no tags set.
func (p *Loader) parseTags(structField reflect.StructField) (*tagInfo, error) {
	structTag := structField.Tag
	value, ok := structTag.Lookup(p.opts.tagName)
	if !ok {
		return nil, nil
	}

	values := strings.SplitN(value, ",", -1)
	key, args := values[0], values[1:]
	if key == "" || key == "-" {
		return nil, nil
	}

	if strings.Contains(key, " ") {
		return nil, fmt.Errorf("env: assigning '%s': invalid key in tag '%s', cannot contain white space characters", structField.Name, structField.Tag)
	}

	tags := &tagInfo{
		key:    key,
		defVal: "",
	}

	for _, arg := range args {
		x := strings.SplitN(arg, "=", 2)
		k := x[0]
		switch k {
		case "default":
			if len(x) != 2 {
				return nil, fmt.Errorf("env: assigning '%s': cannot parse keyword 'default' from tag '%s', format sample: 'default=xxx'", structField.Name, structField.Tag)
			}
			tags.defVal = x[1]
		default:
			//
		}
	}
	return tags, nil
}

// setField set value to the struct field
func setField(field reflect.Value, value string) error {
	refType := field.Type()
	// create a new object if nil pointer
	if refType.Kind() == reflect.Ptr {
		refType = refType.Elem()
		if field.IsNil() {
			field.Set(reflect.New(refType))
		}
		field = field.Elem()
	}

	if setters := getSetters(field); len(setters) != 0 {
		var errs []error
		for _, setter := range setters {
			if err := setter.Set(value); err == nil {
				return nil
			} else {
				errs = append(errs, err)
			}
		}
		return fmt.Errorf("%v", errs)
	}

	switch refType.Kind() {
	case reflect.String:
		field.SetString(value)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		var err error
		var v int64
		if _, ok := field.Interface().(time.Duration); ok {
			var d time.Duration
			d, err = time.ParseDuration(value)
			v = int64(d)
		} else {
			v, err = strconv.ParseInt(value, 0, refType.Bits())
		}

		if err != nil {
			return err
		}

		field.SetInt(v)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		v, err := strconv.ParseUint(value, 0, refType.Bits())
		if err != nil {
			return err
		}
		field.SetUint(v)
	case reflect.Float32, reflect.Float64:
		v, err := strconv.ParseFloat(value, refType.Bits())
		if err != nil {
			return err
		}
		field.SetFloat(v)
	case reflect.Bool:
		v, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		field.SetBool(v)
	case reflect.Slice:
		parts := strings.Split(value, " ")
		sl := reflect.MakeSlice(refType, len(parts), len(parts))
		for i, p := range parts {
			if err := setField(sl.Index(i), p); err != nil {
				return err
			}
		}
		field.Set(sl)
	case reflect.Array:
		parts := strings.Split(value, " ")
		if len(parts) != field.Len() {
			return fmt.Errorf("not enough elements for set %s", refType.String())
		}
		for i, p := range parts {
			if err := setField(field.Index(i), p); err != nil {
				return err
			}
		}
	case reflect.Map:
		mp := reflect.MakeMap(refType)
		pairs := strings.Split(value, " ")

		for _, pair := range pairs {
			kv := strings.Split(pair, ":")
			if len(kv) < 2 {
				return errors.New("invalid map items")
			}
			k := reflect.New(refType.Key()).Elem()
			if err := setField(k, kv[0]); err != nil {
				return err
			}
			v := reflect.New(refType.Elem()).Elem()
			if err := setField(v, strings.Join(kv[1:], ":")); err != nil {
				return err
			}
			mp.SetMapIndex(k, v)
		}
		field.Set(mp)
	default:
		return fmt.Errorf("type '%s' is not supported in array/slice/map", refType.String())
	}
	return nil
}
