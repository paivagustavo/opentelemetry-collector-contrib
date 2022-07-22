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

package devmodeextension // import "github.com/open-telemetry/opentelemetry-collector-contrib/extension/storage/dbstorage"

import (
	"context"
	"database/sql"
	"fmt"
	"go.uber.org/zap"

	// SQLite driver
	_ "github.com/mattn/go-sqlite3"
)

const (
	createTable = `
	create table if not exists %s (
		name text,
		span_id text primary key, 
		trace_id text, 
		parent_id text, 
		start_time int,
		end_time  int, 
		attributes text,
		resource_attributes text
	)
`
	getQueryText    = "select trace_id from %s where span_id=?"
	getAllQueryText = "select * from %s order by start_time desc limit 50"
	setQueryText    = `
	insert into %s(
		name,
		span_id, 
		trace_id, 
		parent_id, 
		start_time,
		end_time, 
		attributes,
		resource_attributes
	) 
	values(?,?,?,?,?,?,?,?) 
`
)

var _ Storer = (*dbStorageClient)(nil)

type dbStorageClient struct {
	driverName     string
	datasourceName string

	db          *sql.DB
	getQuery    *sql.Stmt
	setQuery    *sql.Stmt
	getAllQuery *sql.Stmt
	logger      *zap.Logger
}

func (c *dbStorageClient) StoreSpan(span Span) {
	c.Set(context.Background(), span)
}

func newClient(ctx context.Context, driverName, tableName string, logger *zap.Logger) (*dbStorageClient, error) {
	client := &dbStorageClient{
		logger: logger,
	}
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
	selectAllQuery, err := client.db.PrepareContext(ctx, fmt.Sprintf(getAllQueryText, tableName))
	if err != nil {
		return nil, err
	}
	setQuery, err := client.db.PrepareContext(ctx, fmt.Sprintf(setQueryText, tableName))
	if err != nil {
		return nil, err
	}

	client.getQuery = selectQuery
	client.getAllQuery = selectAllQuery
	client.setQuery = setQuery

	return client, nil
}

// Get will retrieve data from storage that corresponds to the specified key
func (c *dbStorageClient) Get(ctx context.Context, key string) (Span, error) {
	rows, err := c.getQuery.QueryContext(ctx, key)
	if err != nil {
		return Span{}, err
	}
	defer rows.Close()

	if !rows.Next() {
		return Span{}, nil
	}
	span := Span{}
	err = rows.Scan(&span.Name, &span.SpanID, &span.TraceID, &span.ParentID, &span.StartTime, &span.EndTime, &span.Attributes, &span.ResourceAttributes)
	if err != nil {
		return Span{}, err
	}
	err = rows.Close()
	return span, nil
}

// GetAll will retrieve all data from storage
func (c *dbStorageClient) GetAll(ctx context.Context) ([]Span, error) {
	rows, err := c.getAllQuery.QueryContext(ctx)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var spans []Span
	for rows.Next() {
		span := Span{}
		err = rows.Scan(&span.Name, &span.SpanID, &span.TraceID, &span.ParentID, &span.StartTime, &span.EndTime, &span.Attributes, &span.ResourceAttributes)
		if err != nil {
			return nil, err
		}

		spans = append(spans, span)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return spans, nil
}

// Set will store data. The data can be retrieved using the same key
func (c *dbStorageClient) Set(ctx context.Context, span Span) error {
	_, err := c.setQuery.ExecContext(ctx, span.Name, span.SpanID, span.TraceID, span.ParentID, span.StartTime, span.EndTime, span.Attributes, span.ResourceAttributes)
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
