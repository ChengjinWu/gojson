// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Represents JSON data structure using native Go types: booleans, floats,
// strings, arrays, and maps.

package gojson

import (
	"encoding"
	"errors"
	"fmt"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf16"
	"unicode/utf8"
)

const INT32_MAX = int32(^uint32(0) >> 1)
const INT64_MAX = int64(^uint64(0) >> 1)

// Unmarshal parses the JSON-encoded data and stores the result
// in the value pointed to by v.
//
// Unmarshal uses the inverse of the encodings that
// Marshal uses, allocating maps, slices, and pointers as necessary,
// with the following additional rules:
//
// To unmarshal JSON into a pointer, Unmarshal first handles the case of
// the JSON being the JSON literal null. In that case, Unmarshal sets
// the pointer to nil. Otherwise, Unmarshal unmarshals the JSON into
// the value pointed at by the pointer. If the pointer is nil, Unmarshal
// allocates a new value for it to point to.
//
// To unmarshal JSON into a value implementing the Unmarshaler interface,
// Unmarshal calls that value's UnmarshalJSON method, including
// when the input is a JSON null.
// Otherwise, if the value implements encoding.TextUnmarshaler
// and the input is a JSON quoted string, Unmarshal calls that value's
// UnmarshalText method with the unquoted form of the string.
//
// To unmarshal JSON into a struct, Unmarshal matches incoming object
// keys to the keys used by Marshal (either the struct field name or its tag),
// preferring an exact match but also accepting a case-insensitive match.
// Unmarshal will only set exported fields of the struct.
//
// To unmarshal JSON into an interface value,
// Unmarshal stores one of these in the interface value:
//
//	bool, for JSON booleans
//	float64, for JSON numbers
//	string, for JSON strings
//	[]interface{}, for JSON arrays
//	map[string]interface{}, for JSON objects
//	nil for JSON null
//
// To unmarshal a JSON array into a slice, Unmarshal resets the slice length
// to zero and then appends each element to the slice.
// As a special case, to unmarshal an empty JSON array into a slice,
// Unmarshal replaces the slice with a new empty slice.
//
// To unmarshal a JSON array into a Go array, Unmarshal decodes
// JSON array elements into corresponding Go array elements.
// If the Go array is smaller than the JSON array,
// the additional JSON array elements are discarded.
// If the JSON array is smaller than the Go array,
// the additional Go array elements are set to zero values.
//
// To unmarshal a JSON object into a map, Unmarshal first establishes a map to
// use. If the map is nil, Unmarshal allocates a new map. Otherwise Unmarshal
// reuses the existing map, keeping existing entries. Unmarshal then stores
// key-value pairs from the JSON object into the map. The map's key type must
// either be a string, an integer, or implement encoding.TextUnmarshaler.
//
// If a JSON value is not appropriate for a given target type,
// or if a JSON number overflows the target type, Unmarshal
// skips that field and completes the unmarshaling as best it can.
// If no more serious errors are encountered, Unmarshal returns
// an UnmarshalTypeError describing the earliest such error.
//
// The JSON null value unmarshals into an interface, map, pointer, or slice
// by setting that Go value to nil. Because null is often used in JSON to mean
// ``not present,'' unmarshaling a JSON null into any other Go type has no effect
// on the value and produces no error.
//
// When unmarshaling quoted strings, invalid UTF-8 or
// invalid UTF-16 surrogate pairs are not treated as an error.
// Instead, they are replaced by the Unicode replacement
// character U+FFFD.
//
func FromBytes(data []byte) (*JsonObject, error) {
	// Check for well-formedness.
	// Avoids filling out half a data structure
	// before discovering a JSON syntax error.
	var d decodeState
	err := checkValid(data, &d.scan)
	if err != nil {
		return nil, err
	}

	d.init(data)
	return d.fromBytes()
}

// Unmarshaler is the interface implemented by types
// that can unmarshal a JSON description of themselves.
// The input can be assumed to be a valid encoding of
// a JSON value. UnmarshalJSON must copy the JSON data
// if it wishes to retain the data after returning.
//
// By convention, to approximate the behavior of Unmarshal itself,
// Unmarshalers implement UnmarshalJSON([]byte("null")) as a no-op.
type Unmarshaler interface {
	UnmarshalJSON([]byte) error
}

// An UnmarshalTypeError describes a JSON value that was
// not appropriate for a value of a specific Go type.
type UnmarshalTypeError struct {
	Value  string       // description of JSON value - "bool", "array", "number -5"
	Type   reflect.Type // type of Go value it could not be assigned to
	Offset int64        // error occurred after reading Offset bytes
	Struct string       // name of the struct type containing the field
	Field  string       // name of the field holding the Go value
}

func (e *UnmarshalTypeError) Error() string {
	if e.Struct != "" || e.Field != "" {
		return "json: cannot unmarshal " + e.Value + " into Go struct field " + e.Struct + "." + e.Field + " of type " + e.Type.String()
	}
	return "json: cannot unmarshal " + e.Value + " into Go value of type " + e.Type.String()
}

// An UnmarshalFieldError describes a JSON object key that
// led to an unexported (and therefore unwritable) struct field.
// (No longer used; kept for compatibility.)
type UnmarshalFieldError struct {
	Key   string
	Type  reflect.Type
	Field reflect.StructField
}

func (e *UnmarshalFieldError) Error() string {
	return "json: cannot unmarshal object key " + strconv.Quote(e.Key) + " into unexported field " + e.Field.Name + " of type " + e.Type.String()
}

// An InvalidUnmarshalError describes an invalid argument passed to Unmarshal.
// (The argument to Unmarshal must be a non-nil pointer.)
type InvalidUnmarshalError struct {
	Type reflect.Type
}

func (e *InvalidUnmarshalError) Error() string {
	if e.Type == nil {
		return "json: Unmarshal(nil)"
	}

	if e.Type.Kind() != reflect.Ptr {
		return "json: Unmarshal(non-pointer " + e.Type.String() + ")"
	}
	return "json: Unmarshal(nil " + e.Type.String() + ")"
}

func (d *decodeState) fromBytes() (object *JsonObject, err error) {
	defer func() {
		if r := recover(); r != nil {
			if _, ok := r.(runtime.Error); ok {
				panic(r)
			}
			err = r.(error)
		}
	}()

	d.scan.reset()
	// We decode rv not rv.Elem because the Unmarshaler interface
	// test must be applied at the top level of the value.
	object = d.value()
	return object, d.savedError
}

// A Number represents a JSON number literal.
type Number string

// String returns the literal text of the number.
func (n Number) String() string { return string(n) }

// Float64 returns the number as a float64.
func (n Number) Float64() (float64, error) {
	return strconv.ParseFloat(string(n), 64)
}

// Int64 returns the number as an int64.
func (n Number) Int64() (int64, error) {
	return strconv.ParseInt(string(n), 10, 64)
}

// isValidNumber reports whether s is a valid JSON number literal.
func isValidNumber(s string) bool {
	// This function implements the JSON numbers grammar.
	// See https://tools.ietf.org/html/rfc7159#section-6
	// and http://json.org/number.gif

	if s == "" {
		return false
	}

	// Optional -
	if s[0] == '-' {
		s = s[1:]
		if s == "" {
			return false
		}
	}

	// Digits
	switch {
	default:
		return false

	case s[0] == '0':
		s = s[1:]

	case '1' <= s[0] && s[0] <= '9':
		s = s[1:]
		for len(s) > 0 && '0' <= s[0] && s[0] <= '9' {
			s = s[1:]
		}
	}

	// . followed by 1 or more digits.
	if len(s) >= 2 && s[0] == '.' && '0' <= s[1] && s[1] <= '9' {
		s = s[2:]
		for len(s) > 0 && '0' <= s[0] && s[0] <= '9' {
			s = s[1:]
		}
	}

	// e or E followed by an optional - or + and
	// 1 or more digits.
	if len(s) >= 2 && (s[0] == 'e' || s[0] == 'E') {
		s = s[1:]
		if s[0] == '+' || s[0] == '-' {
			s = s[1:]
			if s == "" {
				return false
			}
		}
		for len(s) > 0 && '0' <= s[0] && s[0] <= '9' {
			s = s[1:]
		}
	}

	// Make sure we are at the end.
	return s == ""
}

// decodeState represents the state while decoding a JSON value.
type decodeState struct {
	data         []byte
	off          int // read offset in data
	scan         scanner
	nextscan     scanner  // for calls to nextValue
	errorContext struct { // provides context for type errors
		Struct string
		Field  string
	}
	savedError error
	useNumber  bool
}

// errPhase is used for errors that should not happen unless
// there is a bug in the JSON decoder or something is editing
// the data slice while the decoder executes.
var errPhase = errors.New("JSON decoder out of sync - data changing underfoot?")

func (d *decodeState) init(data []byte) *decodeState {
	d.data = data
	d.off = 0
	d.savedError = nil
	d.errorContext.Struct = ""
	d.errorContext.Field = ""
	return d
}

// error aborts the decoding by panicking with err.
func (d *decodeState) error(err error) {
	panic(d.addErrorContext(err))
}

// saveError saves the first err it is called with,
// for reporting at the end of the unmarshal.
func (d *decodeState) saveError(err error) {
	if d.savedError == nil {
		d.savedError = d.addErrorContext(err)
	}
}

// addErrorContext returns a new error enhanced with information from d.errorContext
func (d *decodeState) addErrorContext(err error) error {
	if d.errorContext.Struct != "" || d.errorContext.Field != "" {
		switch err := err.(type) {
		case *UnmarshalTypeError:
			err.Struct = d.errorContext.Struct
			err.Field = d.errorContext.Field
			return err
		}
	}
	return err
}

// next cuts off and returns the next full JSON value in d.data[d.off:].
// The next value is known to be an object or array, not a literal.
func (d *decodeState) next() []byte {
	c := d.data[d.off]
	item, rest, err := nextValue(d.data[d.off:], &d.nextscan)
	if err != nil {
		d.error(err)
	}
	d.off = len(d.data) - len(rest)

	// Our scanner has seen the opening brace/bracket
	// and thinks we're still in the middle of the object.
	// invent a closing brace/bracket to get it out.
	if c == '{' {
		d.scan.step(&d.scan, '}')
	} else {
		d.scan.step(&d.scan, ']')
	}

	return item
}

// scanWhile processes bytes in d.data[d.off:] until it
// receives a scan code not equal to op.
// It updates d.off and returns the new scan code.
func (d *decodeState) scanWhile(op int) int {
	var newOp int
	for {
		if d.off >= len(d.data) {
			newOp = d.scan.eof()
			d.off = len(d.data) + 1 // mark processed EOF with len+1
		} else {
			c := d.data[d.off]
			d.off++
			newOp = d.scan.step(&d.scan, c)
		}
		if newOp != op {
			break
		}
	}
	return newOp
}

// value decodes a JSON value from d.data[d.off:] into the value.
// it updates d.off to point past the decoded value.
func (d *decodeState) value() *JsonObject {
	switch op := d.scanWhile(scanSkipSpace); op {
	default:
		d.error(errPhase)
		return nil

	case scanBeginArray:
		return d.array()

	case scanBeginObject:
		return d.object()

	case scanBeginLiteral:
		return d.literal()
	}
}

type unquotedValue struct{}

// valueQuoted is like value but decodes a
// quoted string literal or literal null into an interface value.
// If it finds anything other than a quoted string literal or null,
// valueQuoted returns unquotedValue{}.
func (d *decodeState) valueQuoted() interface{} {
	switch op := d.scanWhile(scanSkipSpace); op {
	default:
		d.error(errPhase)

	case scanBeginArray:
		d.array()

	case scanBeginObject:
		d.object()

	case scanBeginLiteral:
		switch v := d.literalInterface().(type) {
		case nil, string:
			return v
		}
	}
	return unquotedValue{}
}

// indirect walks down v allocating pointers as needed,
// until it gets to a non-pointer.
// if it encounters an Unmarshaler, indirect stops and returns that.
// if decodingNull is true, indirect stops at the last pointer so it can be set to nil.
func (d *decodeState) indirect(v reflect.Value, decodingNull bool) (Unmarshaler, encoding.TextUnmarshaler, reflect.Value) {
	// If v is a named type and is addressable,
	// start with its address, so that if the type has pointer methods,
	// we find them.
	if v.Kind() != reflect.Ptr && v.Type().Name() != "" && v.CanAddr() {
		v = v.Addr()
	}
	for {
		// Load value from interface, but only if the result will be
		// usefully addressable.
		if v.Kind() == reflect.Interface && !v.IsNil() {
			e := v.Elem()
			if e.Kind() == reflect.Ptr && !e.IsNil() && (!decodingNull || e.Elem().Kind() == reflect.Ptr) {
				v = e
				continue
			}
		}

		if v.Kind() != reflect.Ptr {
			break
		}

		if v.Elem().Kind() != reflect.Ptr && decodingNull && v.CanSet() {
			break
		}
		if v.IsNil() {
			v.Set(reflect.New(v.Type().Elem()))
		}
		if v.Type().NumMethod() > 0 {
			if u, ok := v.Interface().(Unmarshaler); ok {
				return u, nil, reflect.Value{}
			}
			if !decodingNull {
				if u, ok := v.Interface().(encoding.TextUnmarshaler); ok {
					return nil, u, reflect.Value{}
				}
			}
		}
		v = v.Elem()
	}
	return nil, nil, v
}

// array consumes an array from d.data[d.off-1:], decoding into the value v.
// the first byte of the array ('[') has been read already.
func (d *decodeState) array() *JsonObject {
	object := &JsonObject{
		VType: reflect.Slice,
	}
	values := []*JsonObject{}
	for {
		// Look ahead for ] - can only happen on first iteration.
		op := d.scanWhile(scanSkipSpace)
		if op == scanEndArray {
			break
		}

		// Back up so d.value can have the byte we just read.
		d.off--
		d.scan.undo(op)

		value := d.value()
		values = append(values, value)

		// Next token must be , or ].
		op = d.scanWhile(scanSkipSpace)
		if op == scanEndArray {
			break
		}
		if op != scanArrayValue {
			d.error(errPhase)
		}
	}
	object.Value = values
	return object
}

var nullLiteral = []byte("null")
var textUnmarshalerType = reflect.TypeOf(new(encoding.TextUnmarshaler)).Elem()

// object consumes an object from d.data[d.off-1:], decoding into the value v.
// the first byte ('{') of the object has been read already.
func (d *decodeState) object() *JsonObject {
	object := &JsonObject{
		Attributes: make(map[string]*JsonObject, 1),
		VType:      reflect.Struct,
	}
	for {
		// Read opening " of string key or closing }.
		op := d.scanWhile(scanSkipSpace)
		if op == scanEndObject {
			// closing } - can only happen on first iteration.
			break
		}
		if op != scanBeginLiteral {
			d.error(errPhase)
		}

		// Read key.
		start := d.off - 1
		op = d.scanWhile(scanContinue)
		item := d.data[start : d.off-1]
		key, ok := unquoteBytes(item)
		if !ok {
			d.error(errPhase)
		}

		// Read : before value.
		if op == scanSkipSpace {
			op = d.scanWhile(scanSkipSpace)
		}
		if op != scanObjectKey {
			d.error(errPhase)
		}

		value := d.value()
		object.Attributes[string(key)] = value

		// Next token must be , or }.
		op = d.scanWhile(scanSkipSpace)
		if op == scanEndObject {
			break
		}
		if op != scanObjectValue {
			d.error(errPhase)
		}

		d.errorContext.Struct = ""
		d.errorContext.Field = ""
	}
	return object
}

// literal consumes a literal from d.data[d.off-1:], decoding into the value v.
// The first byte of the literal has been read already
// (that's how the caller knows it's a literal).
func (d *decodeState) literal() *JsonObject {
	// All bytes inside literal return scanContinue op code.
	start := d.off - 1
	op := d.scanWhile(scanContinue)

	// Scan read one byte too far; back up.
	d.off--
	d.scan.undo(op)

	return d.literalStore(d.data[start:d.off])
}

// convertNumber converts the number literal s to a float64 or a Number
// depending on the setting of d.useNumber.
func (d *decodeState) convertNumber(s string) (interface{}, error) {
	if d.useNumber {
		return Number(s), nil
	}
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return nil, &UnmarshalTypeError{Value: "number " + s, Type: reflect.TypeOf(0.0), Offset: int64(d.off)}
	}
	return f, nil
}

var numberType = reflect.TypeOf(Number(""))

// literalStore decodes a literal stored in item into v.
//
// fromQuoted indicates whether this literal came from unwrapping a
// string from the ",string" struct tag option. this is used only to
// produce more helpful error messages.
func (d *decodeState) literalStore(item []byte) *JsonObject {
	// Check for unmarshaler.
	if len(item) == 0 {
		//Empty string given
		d.saveError(fmt.Errorf("json: invalid use of ,string struct tag, trying to unmarshal %q into 自变量", item))
		return nil
	}

	switch c := item[0]; c {
	case 'n': // null
		if string(item) != "null" {
			d.saveError(fmt.Errorf("json: invalid use of ,string struct tag, trying to unmarshal %q into null", item))
			break
		}
		return &JsonObject{
			Value: nil,
			VType: reflect.Interface,
		}
	case 't', 'f': // true, false
		value := item[0] == 't'
		// The main parser checks that only true and false can reach here,
		// but if this was a quoted string input, it could be anything.
		if string(item) != "true" && string(item) != "false" {
			d.saveError(fmt.Errorf("json: invalid use of ,string struct tag, trying to unmarshal %q into bool", item))
			break
		}

		return &JsonObject{
			Value: value,
			VType: reflect.Bool,
		}
	case '"': // string
		s, ok := unquoteBytes(item)
		if !ok {
			d.error(fmt.Errorf("json: invalid use of ,string struct tag, trying to unmarshal %q into string", item))
		}
		return &JsonObject{
			Value: string(s),
			VType: reflect.String,
		}

	default: // number
		s := string(item)
		numberLength := strings.Index(s, ".")
		if numberLength == -1 {
			numberLength = len(s)
		}
		if (c != '-' && (c < '0' || c > '9')) || (numberLength > 18) {
			d.error(fmt.Errorf("json: invalid use of ,string struct tag, trying to unmarshal %q into number", item))
			return &JsonObject{}
		}
		if strings.Contains(s, ".") {
			n, err := strconv.ParseFloat(s, 64)
			if err != nil {
				d.saveError(&UnmarshalTypeError{Value: "number " + s, Offset: int64(d.off)})
				break
			}
			return &JsonObject{
				Value: n,
				VType: reflect.Float64,
			}
		} else {
			n, err := strconv.ParseInt(s, 10, 64)
			if err != nil {
				d.saveError(&UnmarshalTypeError{Value: "number " + s, Offset: int64(d.off)})
				break
			}
			if n > int64(INT32_MAX){
				return &JsonObject{
					Value: n,
					VType: reflect.Int64,
				}
			}else {
				return &JsonObject{
					Value: n,
					VType: reflect.Int32,
				}
			}
		}

	}
	return &JsonObject{}
}

// The xxxInterface routines build up a value to be stored
// in an empty interface. They are not strictly necessary,
// but they avoid the weight of reflection in this common case.

// valueInterface is like value but returns interface{}
func (d *decodeState) valueInterface() interface{} {
	switch d.scanWhile(scanSkipSpace) {
	default:
		d.error(errPhase)
		panic("unreachable")
	case scanBeginArray:
		return d.arrayInterface()
	case scanBeginObject:
		return d.objectInterface()
	case scanBeginLiteral:
		return d.literalInterface()
	}
}

// arrayInterface is like array but returns []interface{}.
func (d *decodeState) arrayInterface() []interface{} {
	var v = make([]interface{}, 0)
	for {
		// Look ahead for ] - can only happen on first iteration.
		op := d.scanWhile(scanSkipSpace)
		if op == scanEndArray {
			break
		}

		// Back up so d.value can have the byte we just read.
		d.off--
		d.scan.undo(op)

		v = append(v, d.valueInterface())

		// Next token must be , or ].
		op = d.scanWhile(scanSkipSpace)
		if op == scanEndArray {
			break
		}
		if op != scanArrayValue {
			d.error(errPhase)
		}
	}
	return v
}

// objectInterface is like object but returns map[string]interface{}.
func (d *decodeState) objectInterface() map[string]interface{} {
	m := make(map[string]interface{})
	for {
		// Read opening " of string key or closing }.
		op := d.scanWhile(scanSkipSpace)
		if op == scanEndObject {
			// closing } - can only happen on first iteration.
			break
		}
		if op != scanBeginLiteral {
			d.error(errPhase)
		}

		// Read string key.
		start := d.off - 1
		op = d.scanWhile(scanContinue)
		item := d.data[start : d.off-1]
		key, ok := unquote(item)
		if !ok {
			d.error(errPhase)
		}

		// Read : before value.
		if op == scanSkipSpace {
			op = d.scanWhile(scanSkipSpace)
		}
		if op != scanObjectKey {
			d.error(errPhase)
		}

		// Read value.
		m[key] = d.valueInterface()

		// Next token must be , or }.
		op = d.scanWhile(scanSkipSpace)
		if op == scanEndObject {
			break
		}
		if op != scanObjectValue {
			d.error(errPhase)
		}
	}
	return m
}

// literalInterface is like literal but returns an interface value.
func (d *decodeState) literalInterface() interface{} {
	// All bytes inside literal return scanContinue op code.
	start := d.off - 1
	op := d.scanWhile(scanContinue)

	// Scan read one byte too far; back up.
	d.off--
	d.scan.undo(op)
	item := d.data[start:d.off]

	switch c := item[0]; c {
	case 'n': // null
		return nil

	case 't', 'f': // true, false
		return c == 't'

	case '"': // string
		s, ok := unquote(item)
		if !ok {
			d.error(errPhase)
		}
		return s

	default: // number
		if c != '-' && (c < '0' || c > '9') {
			d.error(errPhase)
		}
		n, err := d.convertNumber(string(item))
		if err != nil {
			d.saveError(err)
		}
		return n
	}
}

// getu4 decodes \uXXXX from the beginning of s, returning the hex value,
// or it returns -1.
func getu4(s []byte) rune {
	if len(s) < 6 || s[0] != '\\' || s[1] != 'u' {
		return -1
	}
	r, err := strconv.ParseUint(string(s[2:6]), 16, 64)
	if err != nil {
		return -1
	}
	return rune(r)
}

// unquote converts a quoted JSON string literal s into an actual string t.
// The rules are different than for Go, so cannot use strconv.Unquote.
func unquote(s []byte) (t string, ok bool) {
	s, ok = unquoteBytes(s)
	t = string(s)
	return
}

func unquoteBytes(s []byte) (t []byte, ok bool) {
	if len(s) < 2 || s[0] != '"' || s[len(s)-1] != '"' {
		return
	}
	s = s[1 : len(s)-1]

	// Check for unusual characters. If there are none,
	// then no unquoting is needed, so return a slice of the
	// original bytes.
	r := 0
	for r < len(s) {
		c := s[r]
		if c == '\\' || c == '"' || c < ' ' {
			break
		}
		if c < utf8.RuneSelf {
			r++
			continue
		}
		rr, size := utf8.DecodeRune(s[r:])
		if rr == utf8.RuneError && size == 1 {
			break
		}
		r += size
	}
	if r == len(s) {
		return s, true
	}

	b := make([]byte, len(s)+2*utf8.UTFMax)
	w := copy(b, s[0:r])
	for r < len(s) {
		// Out of room?  Can only happen if s is full of
		// malformed UTF-8 and we're replacing each
		// byte with RuneError.
		if w >= len(b)-2*utf8.UTFMax {
			nb := make([]byte, (len(b)+utf8.UTFMax)*2)
			copy(nb, b[0:w])
			b = nb
		}
		switch c := s[r]; {
		case c == '\\':
			r++
			if r >= len(s) {
				return
			}
			switch s[r] {
			default:
				return
			case '"', '\\', '/', '\'':
				b[w] = s[r]
				r++
				w++
			case 'b':
				b[w] = '\b'
				r++
				w++
			case 'f':
				b[w] = '\f'
				r++
				w++
			case 'n':
				b[w] = '\n'
				r++
				w++
			case 'r':
				b[w] = '\r'
				r++
				w++
			case 't':
				b[w] = '\t'
				r++
				w++
			case 'u':
				r--
				rr := getu4(s[r:])
				if rr < 0 {
					return
				}
				r += 6
				if utf16.IsSurrogate(rr) {
					rr1 := getu4(s[r:])
					if dec := utf16.DecodeRune(rr, rr1); dec != unicode.ReplacementChar {
						// A valid pair; consume.
						r += 6
						w += utf8.EncodeRune(b[w:], dec)
						break
					}
					// Invalid surrogate; fall back to replacement rune.
					rr = unicode.ReplacementChar
				}
				w += utf8.EncodeRune(b[w:], rr)
			}

		// Quote, control characters are invalid.
		case c == '"', c < ' ':
			return

		// ASCII
		case c < utf8.RuneSelf:
			b[w] = c
			r++
			w++

		// Coerce to well-formed UTF-8.
		default:
			rr, size := utf8.DecodeRune(s[r:])
			r += size
			w += utf8.EncodeRune(b[w:], rr)
		}
	}
	return b[0:w], true
}
