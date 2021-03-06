// Copyright 2016 The Xorm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package builder

import (
	"bytes"
	"io"
)

// Writer defines the interface
type Writer interface {
	io.Writer
	Append(...interface{})
	Key(string) string
}

var _ Writer = NewWriter()

// BytesWriter implments Writer and save SQL in bytes.Buffer
type BytesWriter struct {
	writer    *bytes.Buffer
	buffer    []byte
	args      []interface{}
	keyFilter func(string) string
}

// NewWriter creates a new string writer
func NewWriter() *BytesWriter {
	w := &BytesWriter{}
	w.writer = bytes.NewBuffer(w.buffer)
	return w
}

// Write writes data to Writer
func (s *BytesWriter) Write(buf []byte) (int, error) {
	return s.writer.Write(buf)
}

// Append appends args to Writer
func (s *BytesWriter) Append(args ...interface{}) {
	s.args = append(s.args, args...)
}

func (s *BytesWriter) Key(key string) string {
	if s.keyFilter != nil {
		return s.keyFilter(key)
	}
	return key
}

func (s *BytesWriter) String() string {
	return s.writer.String()
}

func (s *BytesWriter) Bytes() []byte {
	return s.writer.Bytes()
}

func (s *BytesWriter) Args() []interface{} {
	return s.args
}

// Cond defines an interface
type Cond interface {
	WriteTo(Writer) error
	And(...Cond) Cond
	Or(...Cond) Cond
	IsValid() bool
}

type condEmpty struct{}

var _ Cond = condEmpty{}

// NewCond creates an empty condition
func NewCond() Cond {
	return condEmpty{}
}

func (condEmpty) WriteTo(w Writer) error {
	return nil
}

func (condEmpty) And(conds ...Cond) Cond {
	return And(conds...)
}

func (condEmpty) Or(conds ...Cond) Cond {
	return Or(conds...)
}

func (condEmpty) IsValid() bool {
	return false
}

func condToSQL(cond Cond, keyFilters ...func(string) string) (string, []interface{}, error) {
	if cond == nil || !cond.IsValid() {
		return "", nil, nil
	}

	var keyFilter func(string) string
	if len(keyFilters) > 0 {
		keyFilter = keyFilters[0]
	}
	w := NewWriter()
	w.keyFilter = keyFilter
	if err := cond.WriteTo(w); err != nil {
		return "", nil, err
	}
	return w.writer.String(), w.args, nil
}
