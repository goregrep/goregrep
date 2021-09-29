// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package regenerate regenerates some generated code.
package regenerate

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"regexp"

	"golang.org/x/tools/imports"
)

func Pipe(in io.Reader, out file, opts ...Option) error {
	var cfg Configuration

	for _, opt := range opts {
		opt(&cfg)
	}

	b, err := ioutil.ReadAll(in)
	if err != nil {
		return fmt.Errorf("io util read all: %w", err)
	}

	for _, replace := range cfg.strings {
		b = bytes.ReplaceAll(b, []byte(replace.Match), []byte(replace.Replacement))
	}

	for _, replace := range cfg.regexps {
		b = replace.Match.ReplaceAll(b, []byte(replace.Replacement))
	}

	err = out.Truncate(0)
	if err != nil {
		return fmt.Errorf("io file truncate: %w", err)
	}

	if cfg.gofmt != nil {
		b, err = imports.Process("", b, cfg.gofmt)
		if err != nil {
			return fmt.Errorf("imports process: %w", err)
		}
	}

	_, err = out.WriteAt(b, 0)
	if err != nil {
		return fmt.Errorf("os file write at 0: %w", err)
	}

	return nil
}

type file interface {
	io.WriterAt
	Truncate(size int64) error
}

type String struct {
	Match       string
	Replacement string
}

type Regexp struct {
	Match       *regexp.Regexp
	Replacement string
}

// Option changes configuration.
type Option func(*Configuration)

// Configuration holds values changeable by options.
type Configuration struct {
	strings []String
	regexps []Regexp
	gofmt   *imports.Options
}

// ReplaceString add replacement.
func ReplaceString(match string, replacement string) Option {
	return func(c *Configuration) {
		c.strings = append(c.strings, String{
			Match:       match,
			Replacement: replacement,
		})
	}
}

// ReplaceRegexp add replacement.
func ReplaceRegexp(match *regexp.Regexp, replacement string) Option {
	return func(c *Configuration) {
		c.regexps = append(c.regexps, Regexp{
			Match:       match,
			Replacement: replacement,
		})
	}
}

func WithGofmt(opts *imports.Options) Option {
	return func(c *Configuration) { c.gofmt = opts }
}
