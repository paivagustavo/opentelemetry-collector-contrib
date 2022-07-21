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

	"go.opentelemetry.io/collector/pdata/ptrace"
	"go.uber.org/zap"
)

type devmodeProcessor struct {
	logger *zap.Logger
}

func (dev *devmodeProcessor) processTraces(ctx context.Context, td ptrace.Traces) (ptrace.Traces, error) {
	rss := td.ResourceSpans()
	for i := 0; i < rss.Len(); i++ {
		ss := rss.At(i).ScopeSpans()
		// TODO: rss.At(i).Resource() used this for getting the resource attributes
		for j := 0; j < ss.Len(); j++ {
			spans := ss.At(j).Spans()
			for k := 0; k < spans.Len(); k++ {
				span := spans.At(k)
				// TODO: store this span on the devmode storage.
				// 	remember to save the resource attributes as well.

				// TODO: Need to define how the processor will write to the extension storage.

				// Log each trace_id/span_id
				dev.logger.Info(span.TraceID().HexString() + " - " + span.SpanID().HexString())
			}
		}
	}
	return td, nil
}
