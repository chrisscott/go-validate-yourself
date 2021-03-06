package validator

import (
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"time"
)

// BakedInValidators is the map of ValidationFunc used internally
// but can be used with any new Validator if desired
var BakedInValidators = map[string]ValidationFunc{
	"required":    hasValue,
	"len":         hasLengthOf,
	"min":         hasMinOf,
	"max":         hasMaxOf,
	"lt":          isLt,
	"lte":         isLte,
	"gt":          isGt,
	"gte":         isGte,
	"gtefield":    isGteField,
	"gtfield":     isGtField,
	"ltefield":    isLteField,
	"ltfield":     isLtField,
	"alpha":       isAlpha,
	"alphanum":    isAlphanum,
	"numeric":     isNumeric,
	"number":      isNumber,
	"hexadecimal": isHexadecimal,
	"hexcolor":    isHexcolor,
	"rgb":         isRgb,
	"rgba":        isRgba,
	"hsl":         isHsl,
	"hsla":        isHsla,
	"email":       isEmail,
	"url":         isURL,
	"uri":         isURI,
}

func isURI(top interface{}, current interface{}, field interface{}, param string) bool {

	st := reflect.ValueOf(field)

	switch st.Kind() {

	case reflect.String:
		_, err := url.ParseRequestURI(field.(string))

		return err == nil
	}

	panic(fmt.Sprintf("Bad field type %T", field))
}

func isURL(top interface{}, current interface{}, field interface{}, param string) bool {

	st := reflect.ValueOf(field)

	switch st.Kind() {

	case reflect.String:
		url, err := url.ParseRequestURI(field.(string))

		if err != nil {
			return false
		}

		if len(url.Scheme) == 0 {
			return false
		}

		return err == nil
	}

	panic(fmt.Sprintf("Bad field type %T", field))
}

func isEmail(top interface{}, current interface{}, field interface{}, param string) bool {

	st := reflect.ValueOf(field)

	switch st.Kind() {

	case reflect.String:
		return emailRegex.MatchString(field.(string))
	}

	panic(fmt.Sprintf("Bad field type %T", field))
}

func isHsla(top interface{}, current interface{}, field interface{}, param string) bool {

	st := reflect.ValueOf(field)

	switch st.Kind() {

	case reflect.String:
		return hslaRegex.MatchString(field.(string))
	}

	panic(fmt.Sprintf("Bad field type %T", field))
}

func isHsl(top interface{}, current interface{}, field interface{}, param string) bool {

	st := reflect.ValueOf(field)

	switch st.Kind() {

	case reflect.String:
		return hslRegex.MatchString(field.(string))
	}

	panic(fmt.Sprintf("Bad field type %T", field))
}

func isRgba(top interface{}, current interface{}, field interface{}, param string) bool {

	st := reflect.ValueOf(field)

	switch st.Kind() {

	case reflect.String:
		return rgbaRegex.MatchString(field.(string))
	}

	panic(fmt.Sprintf("Bad field type %T", field))
}

func isRgb(top interface{}, current interface{}, field interface{}, param string) bool {

	st := reflect.ValueOf(field)

	switch st.Kind() {

	case reflect.String:
		return rgbRegex.MatchString(field.(string))
	}

	panic(fmt.Sprintf("Bad field type %T", field))
}

func isHexcolor(top interface{}, current interface{}, field interface{}, param string) bool {

	st := reflect.ValueOf(field)

	switch st.Kind() {

	case reflect.String:
		return hexcolorRegex.MatchString(field.(string))
	}

	panic(fmt.Sprintf("Bad field type %T", field))
}

func isHexadecimal(top interface{}, current interface{}, field interface{}, param string) bool {

	st := reflect.ValueOf(field)

	switch st.Kind() {

	case reflect.String:
		return hexadecimalRegex.MatchString(field.(string))
	}

	panic(fmt.Sprintf("Bad field type %T", field))
}

func isNumber(top interface{}, current interface{}, field interface{}, param string) bool {

	st := reflect.ValueOf(field)

	switch st.Kind() {

	case reflect.String:
		return numberRegex.MatchString(field.(string))
	}

	panic(fmt.Sprintf("Bad field type %T", field))
}

func isNumeric(top interface{}, current interface{}, field interface{}, param string) bool {

	st := reflect.ValueOf(field)

	switch st.Kind() {

	case reflect.String:
		return numericRegex.MatchString(field.(string))
	}

	panic(fmt.Sprintf("Bad field type %T", field))
}

func isAlphanum(top interface{}, current interface{}, field interface{}, param string) bool {

	st := reflect.ValueOf(field)

	switch st.Kind() {

	case reflect.String:
		return alphaNumericRegex.MatchString(field.(string))
	}

	panic(fmt.Sprintf("Bad field type %T", field))
}

func isAlpha(top interface{}, current interface{}, field interface{}, param string) bool {

	st := reflect.ValueOf(field)

	switch st.Kind() {

	case reflect.String:
		return alphaRegex.MatchString(field.(string))
	}

	panic(fmt.Sprintf("Bad field type %T", field))
}

func hasValue(top interface{}, current interface{}, field interface{}, param string) bool {

	st := reflect.ValueOf(field)

	switch st.Kind() {

	case reflect.Slice, reflect.Map, reflect.Array:
		return field != nil && int64(st.Len()) > 0

	default:
		return field != nil && field != reflect.Zero(reflect.TypeOf(field)).Interface()
	}
}

func isGteField(top interface{}, current interface{}, field interface{}, param string) bool {

	if current == nil {
		panic("struct not passed for cross validation")
	}

	currentVal := reflect.ValueOf(current)

	if currentVal.Kind() == reflect.Ptr && !currentVal.IsNil() {
		currentVal = reflect.ValueOf(currentVal.Elem().Interface())
	}

	var currentFielVal reflect.Value

	switch currentVal.Kind() {

	case reflect.Struct:

		if currentVal.Type() == reflect.TypeOf(time.Time{}) {
			currentFielVal = currentVal
			break
		}

		f := currentVal.FieldByName(param)

		if f.Kind() == reflect.Invalid {
			panic(fmt.Sprintf("Field \"%s\" not found in struct", param))
		}

		currentFielVal = f

	default:

		currentFielVal = currentVal
	}

	if currentFielVal.Kind() == reflect.Ptr && !currentFielVal.IsNil() {

		currentFielVal = reflect.ValueOf(currentFielVal.Elem().Interface())
	}

	fv := reflect.ValueOf(field)

	switch fv.Kind() {

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:

		return fv.Int() >= currentFielVal.Int()

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:

		return fv.Uint() >= currentFielVal.Uint()

	case reflect.Float32, reflect.Float64:

		return fv.Float() >= currentFielVal.Float()

	case reflect.Struct:

		if fv.Type() == reflect.TypeOf(time.Time{}) {

			if currentFielVal.Type() != reflect.TypeOf(time.Time{}) {
				panic("Bad Top Level field type")
			}

			t := currentFielVal.Interface().(time.Time)
			fieldTime := field.(time.Time)

			return fieldTime.After(t) || fieldTime.Equal(t)
		}
	}

	panic(fmt.Sprintf("Bad field type %T", field))
}

func isGtField(top interface{}, current interface{}, field interface{}, param string) bool {

	if current == nil {
		panic("struct not passed for cross validation")
	}

	currentVal := reflect.ValueOf(current)

	if currentVal.Kind() == reflect.Ptr && !currentVal.IsNil() {
		currentVal = reflect.ValueOf(currentVal.Elem().Interface())
	}

	var currentFielVal reflect.Value

	switch currentVal.Kind() {

	case reflect.Struct:

		if currentVal.Type() == reflect.TypeOf(time.Time{}) {
			currentFielVal = currentVal
			break
		}

		f := currentVal.FieldByName(param)

		if f.Kind() == reflect.Invalid {
			panic(fmt.Sprintf("Field \"%s\" not found in struct", param))
		}

		currentFielVal = f

	default:

		currentFielVal = currentVal
	}

	if currentFielVal.Kind() == reflect.Ptr && !currentFielVal.IsNil() {

		currentFielVal = reflect.ValueOf(currentFielVal.Elem().Interface())
	}

	fv := reflect.ValueOf(field)

	switch fv.Kind() {

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:

		return fv.Int() > currentFielVal.Int()

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:

		return fv.Uint() > currentFielVal.Uint()

	case reflect.Float32, reflect.Float64:

		return fv.Float() > currentFielVal.Float()

	case reflect.Struct:

		if fv.Type() == reflect.TypeOf(time.Time{}) {

			if currentFielVal.Type() != reflect.TypeOf(time.Time{}) {
				panic("Bad Top Level field type")
			}

			t := currentFielVal.Interface().(time.Time)
			fieldTime := field.(time.Time)

			return fieldTime.After(t)
		}
	}

	panic(fmt.Sprintf("Bad field type %T", field))
}

func isGte(top interface{}, current interface{}, field interface{}, param string) bool {

	st := reflect.ValueOf(field)

	switch st.Kind() {

	case reflect.String:
		p := asInt(param)

		return int64(len(st.String())) >= p

	case reflect.Slice, reflect.Map, reflect.Array:
		p := asInt(param)

		return int64(st.Len()) >= p

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		p := asInt(param)

		return st.Int() >= p

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		p := asUint(param)

		return st.Uint() >= p

	case reflect.Float32, reflect.Float64:
		p := asFloat(param)

		return st.Float() >= p

	case reflect.Struct:

		if st.Type() == reflect.TypeOf(time.Time{}) {

			now := time.Now().UTC()
			t := field.(time.Time)

			return t.After(now) || t.Equal(now)
		}
	}

	panic(fmt.Sprintf("Bad field type %T", field))
}

func isGt(top interface{}, current interface{}, field interface{}, param string) bool {

	st := reflect.ValueOf(field)

	switch st.Kind() {

	case reflect.String:
		p := asInt(param)

		return int64(len(st.String())) > p

	case reflect.Slice, reflect.Map, reflect.Array:
		p := asInt(param)

		return int64(st.Len()) > p

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		p := asInt(param)

		return st.Int() > p

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		p := asUint(param)

		return st.Uint() > p

	case reflect.Float32, reflect.Float64:
		p := asFloat(param)

		return st.Float() > p
	case reflect.Struct:

		if st.Type() == reflect.TypeOf(time.Time{}) {

			return field.(time.Time).After(time.Now().UTC())
		}
	}

	panic(fmt.Sprintf("Bad field type %T", field))
}

// length tests whether a variable's length is equal to a given
// value. For strings it tests the number of characters whereas
// for maps and slices it tests the number of items.
func hasLengthOf(top interface{}, current interface{}, field interface{}, param string) bool {

	st := reflect.ValueOf(field)

	switch st.Kind() {

	case reflect.String:
		p := asInt(param)

		return int64(len(st.String())) == p

	case reflect.Slice, reflect.Map, reflect.Array:
		p := asInt(param)

		return int64(st.Len()) == p

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		p := asInt(param)

		return st.Int() == p

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		p := asUint(param)

		return st.Uint() == p

	case reflect.Float32, reflect.Float64:
		p := asFloat(param)

		return st.Float() == p
	}

	panic(fmt.Sprintf("Bad field type %T", field))
}

// min tests whether a variable value is larger or equal to a given
// number. For number types, it's a simple lesser-than test; for
// strings it tests the number of characters whereas for maps
// and slices it tests the number of items.
func hasMinOf(top interface{}, current interface{}, field interface{}, param string) bool {

	return isGte(top, current, field, param)
}

func isLteField(top interface{}, current interface{}, field interface{}, param string) bool {

	if current == nil {
		panic("struct not passed for cross validation")
	}

	currentVal := reflect.ValueOf(current)

	if currentVal.Kind() == reflect.Ptr && !currentVal.IsNil() {
		currentVal = reflect.ValueOf(currentVal.Elem().Interface())
	}

	var currentFielVal reflect.Value

	switch currentVal.Kind() {

	case reflect.Struct:

		if currentVal.Type() == reflect.TypeOf(time.Time{}) {
			currentFielVal = currentVal
			break
		}

		f := currentVal.FieldByName(param)

		if f.Kind() == reflect.Invalid {
			panic(fmt.Sprintf("Field \"%s\" not found in struct", param))
		}

		currentFielVal = f

	default:

		currentFielVal = currentVal
	}

	if currentFielVal.Kind() == reflect.Ptr && !currentFielVal.IsNil() {

		currentFielVal = reflect.ValueOf(currentFielVal.Elem().Interface())
	}

	fv := reflect.ValueOf(field)

	switch fv.Kind() {

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:

		return fv.Int() <= currentFielVal.Int()

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:

		return fv.Uint() <= currentFielVal.Uint()

	case reflect.Float32, reflect.Float64:

		return fv.Float() <= currentFielVal.Float()

	case reflect.Struct:

		if fv.Type() == reflect.TypeOf(time.Time{}) {

			if currentFielVal.Type() != reflect.TypeOf(time.Time{}) {
				panic("Bad Top Level field type")
			}

			t := currentFielVal.Interface().(time.Time)
			fieldTime := field.(time.Time)

			return fieldTime.Before(t) || fieldTime.Equal(t)
		}
	}

	panic(fmt.Sprintf("Bad field type %T", field))
}

func isLtField(top interface{}, current interface{}, field interface{}, param string) bool {

	if current == nil {
		panic("struct not passed for cross validation")
	}

	currentVal := reflect.ValueOf(current)

	if currentVal.Kind() == reflect.Ptr && !currentVal.IsNil() {
		currentVal = reflect.ValueOf(currentVal.Elem().Interface())
	}

	var currentFielVal reflect.Value

	switch currentVal.Kind() {

	case reflect.Struct:

		if currentVal.Type() == reflect.TypeOf(time.Time{}) {
			currentFielVal = currentVal
			break
		}

		f := currentVal.FieldByName(param)

		if f.Kind() == reflect.Invalid {
			panic(fmt.Sprintf("Field \"%s\" not found in struct", param))
		}

		currentFielVal = f

	default:

		currentFielVal = currentVal
	}

	if currentFielVal.Kind() == reflect.Ptr && !currentFielVal.IsNil() {

		currentFielVal = reflect.ValueOf(currentFielVal.Elem().Interface())
	}

	fv := reflect.ValueOf(field)

	switch fv.Kind() {

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:

		return fv.Int() < currentFielVal.Int()

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:

		return fv.Uint() < currentFielVal.Uint()

	case reflect.Float32, reflect.Float64:

		return fv.Float() < currentFielVal.Float()

	case reflect.Struct:

		if fv.Type() == reflect.TypeOf(time.Time{}) {

			if currentFielVal.Type() != reflect.TypeOf(time.Time{}) {
				panic("Bad Top Level field type")
			}

			t := currentFielVal.Interface().(time.Time)
			fieldTime := field.(time.Time)

			return fieldTime.Before(t)
		}
	}

	panic(fmt.Sprintf("Bad field type %T", field))
}

func isLte(top interface{}, current interface{}, field interface{}, param string) bool {

	st := reflect.ValueOf(field)

	switch st.Kind() {

	case reflect.String:
		p := asInt(param)

		return int64(len(st.String())) <= p

	case reflect.Slice, reflect.Map, reflect.Array:
		p := asInt(param)

		return int64(st.Len()) <= p

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		p := asInt(param)

		return st.Int() <= p

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		p := asUint(param)

		return st.Uint() <= p

	case reflect.Float32, reflect.Float64:
		p := asFloat(param)

		return st.Float() <= p

	case reflect.Struct:

		if st.Type() == reflect.TypeOf(time.Time{}) {

			now := time.Now().UTC()
			t := field.(time.Time)

			return t.Before(now) || t.Equal(now)
		}
	}

	panic(fmt.Sprintf("Bad field type %T", field))
}

func isLt(top interface{}, current interface{}, field interface{}, param string) bool {

	st := reflect.ValueOf(field)

	switch st.Kind() {

	case reflect.String:
		p := asInt(param)

		return int64(len(st.String())) < p

	case reflect.Slice, reflect.Map, reflect.Array:
		p := asInt(param)

		return int64(st.Len()) < p

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		p := asInt(param)

		return st.Int() < p

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		p := asUint(param)

		return st.Uint() < p

	case reflect.Float32, reflect.Float64:
		p := asFloat(param)

		return st.Float() < p

	case reflect.Struct:

		if st.Type() == reflect.TypeOf(time.Time{}) {

			return field.(time.Time).Before(time.Now().UTC())
		}
	}

	panic(fmt.Sprintf("Bad field type %T", field))
}

// max tests whether a variable value is lesser than a given
// value. For numbers, it's a simple lesser-than test; for
// strings it tests the number of characters whereas for maps
// and slices it tests the number of items.
func hasMaxOf(top interface{}, current interface{}, field interface{}, param string) bool {

	return isLte(top, current, field, param)
}

// asInt retuns the parameter as a int64
// or panics if it can't convert
func asInt(param string) int64 {

	i, err := strconv.ParseInt(param, 0, 64)

	if err != nil {
		panic(err.Error())
	}

	return i
}

// asUint returns the parameter as a uint64
// or panics if it can't convert
func asUint(param string) uint64 {

	i, err := strconv.ParseUint(param, 0, 64)

	if err != nil {
		panic(err.Error())
	}

	return i
}

// asFloat returns the parameter as a float64
// or panics if it can't convert
func asFloat(param string) float64 {

	i, err := strconv.ParseFloat(param, 64)

	if err != nil {
		panic(err.Error())
	}

	return i
}
