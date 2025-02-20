// Copyright 2021 Woodpecker Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package datastore

import (
	"context"

	"github.com/rs/zerolog"
	"xorm.io/xorm"
	xlog "xorm.io/xorm/log"

	"go.woodpecker-ci.org/woodpecker/v3/server/store"
	"go.woodpecker-ci.org/woodpecker/v3/server/store/datastore/migration"
)

type storage struct {
	engine *xorm.Engine
}

const perPage = 50

func NewEngine(opts *store.Opts) (store.Store, error) {
	engine, err := xorm.NewEngine(opts.Driver, opts.Config)
	if err != nil {
		return nil, err
	}

	level := xlog.LogLevel(zerolog.GlobalLevel())
	if !opts.XORM.Log {
		level = xlog.LOG_OFF
	}

	logger := newXORMLogger(level)
	engine.SetLogger(logger)
	engine.ShowSQL(opts.XORM.ShowSQL)
	engine.SetMaxOpenConns(opts.XORM.MaxOpenConns)
	engine.SetMaxIdleConns(opts.XORM.MaxIdleConns)
	engine.SetConnMaxLifetime(opts.XORM.ConnMaxLifetime)

	return &storage{
		engine: engine,
	}, nil
}

func (s storage) Ping() error {
	return s.engine.Ping()
}

// Migrate old storage or init new one.
func (s storage) Migrate(ctx context.Context, allowLong bool) error {
	return migration.Migrate(ctx, s.engine, allowLong)
}

func (s storage) Close() error {
	return s.engine.Close()
}
