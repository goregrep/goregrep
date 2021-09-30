// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package ggrep provides functionality for replace (regenerate) generated code.
package ggrep

import (
	"fmt"
	"io"
	"path/filepath"
	"regexp"

	"github.com/goccy/go-yaml"
	"github.com/goregrep/goregrep/regenerate"
	"golang.org/x/tools/imports"
)

func New(yml io.Reader, opts ...Option) error {
	var cfg Configuration

	for _, opt := range opts {
		opt(&cfg)
	}

	var yamlOptions []yaml.DecodeOption

	if len(cfg.references) != 0 {
		yamlOptions = append(yamlOptions, yaml.ReferenceReaders(cfg.references...))
	}

	dec := yaml.NewDecoder(yml, yamlOptions...)

	err := dec.Decode(&cfg)
	if err != nil {
		return fmt.Errorf("decode yaml: %w", err)
	}

	for _, reg := range cfg.Regenerates {
		var regenerates []regenerate.Option

		for _, rep := range reg.Replace.Strings {
			regenerates = append(regenerates,
				regenerate.ReplaceString(rep.Match, rep.Replacement),
			)
		}

		for _, rep := range reg.Replace.Regexps {
			re, err := regexp.Compile(rep.Match)
			if err != nil {
				return fmt.Errorf("regexp compile: %w", err)
			}

			regenerates = append(regenerates,
				regenerate.ReplaceRegexp(re, rep.Replacement),
			)
		}

		if cfg.gofmt != nil {
			regenerates = append(regenerates, regenerate.WithGofmt(cfg.gofmt))
		}

		err = regenerate.Glob(filepath.Join(cfg.directory, reg.File), regenerates...)
		if err != nil {
			return fmt.Errorf("regenerate pipe: %s", err)
		}
	}

	return nil
}

type Configuration struct {
	Regenerates []Regenerate `yaml:"regenerates"`
	directory   string
	references  []io.Reader
	gofmt       *imports.Options
}

type Regenerate struct {
	File    string  `yaml:"file"`
	Replace Replace `yaml:"replace"`
}

type Replace struct {
	Strings []String `yaml:"strings"`
	Regexps []Regexp `yaml:"regexps"`
}

type String struct {
	Match       string `yaml:"match"`
	Replacement string `yaml:"replacement"`
}

type Regexp struct {
	Match       string `yaml:"match"`
	Replacement string `yaml:"replacement"`
}

// Option changes configuration.
type Option func(*Configuration)

// WithDirectory sets directory.
func WithDirectory(dir string) Option {
	return func(c *Configuration) {
		c.directory = dir
	}
}

// WithReferences setes reference to anchor defined by passed readers.
func WithReferences(refs ...io.Reader) Option {
	return func(c *Configuration) {
		c.references = append(c.references, refs...)
	}
}

func WithGofmt(opts *imports.Options) Option {
	return func(c *Configuration) {
		c.gofmt = opts
	}
}
