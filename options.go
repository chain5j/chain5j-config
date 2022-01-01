// Package config
//
// @author: xwc1125
package config

import (
	"fmt"
	"github.com/chain5j/chain5j-protocol/protocol"
)

type option func(f *config) error

func apply(f *config, opts ...option) error {
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		if err := opt(f); err != nil {
			return fmt.Errorf("option apply err:%v", err)
		}
	}
	return nil
}

func WithDB(db protocol.DatabaseReader) option {
	return func(f *config) error {
		f.db = db
		return nil
	}
}
