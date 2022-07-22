// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package devmodeprocessor // import "github.com/open-telemetry/opentelemetry-collector-contrib/processor/devmodeprocessor"

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/open-telemetry/opentelemetry-collector-contrib/extension/devmodeextension"
	"go.opentelemetry.io/collector/pdata/ptrace"
	"go.uber.org/zap"
	"strings"
)

type devmodeProcessor struct {
	logger *zap.Logger
}

func (dev *devmodeProcessor) processTraces(ctx context.Context, td ptrace.Traces) (ptrace.Traces, error) {
	rss := td.ResourceSpans()
	for i := 0; i < rss.Len(); i++ {
		ss := rss.At(i).ScopeSpans()
		resource := rss.At(i).Resource()
		// TODO: rss.At(i).Resource() used this for getting the resource attributes
		for j := 0; j < ss.Len(); j++ {
			spans := ss.At(j).Spans()
			for k := 0; k < spans.Len(); k++ {
				span := spans.At(k)

				newString := func(str string) sql.NullString {
					return sql.NullString{
						String: str,
						Valid:  true,
					}
				}

				ds := devmodeextension.Span{
					SpanID:    newString(span.SpanID().HexString()),
					TraceID:   newString(span.TraceID().HexString()),
					StartTime: sql.NullInt64{Int64: span.StartTimestamp().AsTime().UnixMilli(), Valid: true},
					EndTime:   sql.NullInt64{Int64: span.EndTimestamp().AsTime().UnixMilli(), Valid: true},
				}

				if !span.ParentSpanID().IsEmpty() {
					ds.TraceID = newString(span.ParentSpanID().HexString())
				}
				var attrs []string
				for key := range span.Attributes().AsRaw() {
					value, ok := span.Attributes().Get(key)

					if ok {
						attrs = append(attrs, fmt.Sprintf("%s=%s", key, value.AsString()))
					}
				}
				ds.Attributes = newString(strings.Join(attrs, ","))

				var rscAttrs []string
				for key := range resource.Attributes().AsRaw() {
					value, ok := resource.Attributes().Get(key)
					if ok {
						rscAttrs = append(rscAttrs, fmt.Sprintf("%s=%s", key, value.AsString()))
					}
				}
				ds.ResourceAttributes = newString(strings.Join(rscAttrs, ","))

				devmodeextension.Storage.StoreTrace(ds)
			}
		}
	}
	return td, nil
}
