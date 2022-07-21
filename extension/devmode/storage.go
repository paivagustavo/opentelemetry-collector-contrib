// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package devmode // import "github.com/open-telemetry/opentelemetry-collector-contrib/extension/storage/dbstorage"

import (
	"context"
	"database/sql"
	"fmt"

	// SQLite driver
	_ "github.com/mattn/go-sqlite3"
)

const (
	createTable = `
	create table if not exists %s (
		span_id text primary key, 
		trace_id text, 
		parent_id text, 
		start_time int,
		end_time  int, 
		attributes text,
		resource_attributes text
	)
`
	getQueryText = "select trace_id from %s where span_id=?"
	setQueryText = "insert into %s(span_id, trace_id) values(?,?) on conflict(span_id) do update set trace_id=?"
)

type dbStorageClient struct {
	driverName     string
	datasourceName string

	db       *sql.DB
	getQuery *sql.Stmt
	setQuery *sql.Stmt
}

func newClient(ctx context.Context, driverName, tableName string) (*dbStorageClient, error) {
	client := &dbStorageClient{}
	var err error

	client.db, err = sql.Open(driverName, ":memory:")
	if err != nil {
		return nil, err
	}

	if err := client.db.Ping(); err != nil {
		return nil, err
	}

	_, err = client.db.ExecContext(ctx, fmt.Sprintf(createTable, tableName))
	if err != nil {
		return nil, err
	}

	selectQuery, err := client.db.PrepareContext(ctx, fmt.Sprintf(getQueryText, tableName))
	if err != nil {
		return nil, err
	}
	setQuery, err := client.db.PrepareContext(ctx, fmt.Sprintf(setQueryText, tableName))
	if err != nil {
		return nil, err
	}

	client.getQuery = selectQuery
	client.setQuery = setQuery

	return client, nil
}

// Get will retrieve data from storage that corresponds to the specified key
func (c *dbStorageClient) Get(ctx context.Context, key string) ([]byte, error) {
	rows, err := c.getQuery.QueryContext(ctx, key)
	if err != nil {
		return nil, err
	}
	if !rows.Next() {
		return nil, nil
	}
	var result []byte
	err = rows.Scan(&result)
	if err != nil {
		return result, err
	}
	err = rows.Close()
	return result, err
}

// Set will store data. The data can be retrieved using the same key
func (c *dbStorageClient) Set(ctx context.Context, key string, value []byte) error {
	_, err := c.setQuery.ExecContext(ctx, key, value, value)
	return err
}

// Close will close the database
func (c *dbStorageClient) Close(_ context.Context) error {
	if err := c.setQuery.Close(); err != nil {
		return err
	}
	if err := c.getQuery.Close(); err != nil {
		return err
	}
	return nil
}
