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

				ds := devmodeextension.Span{
					SpanID:    span.SpanID().HexString(),
					Name:      span.Name(),
					TraceID:   span.TraceID().HexString(),
					StartTime: span.StartTimestamp().AsTime().UnixMilli(),
					EndTime:   span.EndTimestamp().AsTime().UnixMilli(),
				}

				if !span.ParentSpanID().IsEmpty() {
					ds.ParentID = span.ParentSpanID().HexString()
				}
				var attrs []string
				for key := range span.Attributes().AsRaw() {
					value, ok := span.Attributes().Get(key)

					if ok {
						attrs = append(attrs, fmt.Sprintf("%s=%s", key, value.AsString()))
					}
				}
				ds.Attributes = strings.Join(attrs, ",")

				var rscAttrs []string
				for key := range resource.Attributes().AsRaw() {
					value, ok := resource.Attributes().Get(key)
					if ok {
						rscAttrs = append(rscAttrs, fmt.Sprintf("%s=%s", key, value.AsString()))
					}
				}
				ds.ResourceAttributes = strings.Join(rscAttrs, ",")

				devmodeextension.Storage.StoreSpan(ds)
			}
		}
	}
	return td, nil
}
