// Licensed to Elasticsearch B.V. under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Elasticsearch B.V. licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package ratelimitprocessor // import "github.com/elastic/opentelemetry-collector-components/processor/ratelimitprocessor"

import (
	"context"
	"sync"
	"time"

	"github.com/elastic/opentelemetry-collector-components/processor/ratelimitprocessor/internal/metadata"
	"github.com/elastic/opentelemetry-collector-components/processor/ratelimitprocessor/internal/telemetry"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/processor"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"golang.org/x/time/rate"
)

var _ RateLimiter = (*localRateLimiter)(nil)

type localRateLimiter struct {
	cfg *Config
	set processor.Settings
	// TODO use an LRU to keep a cap on the number of limiters.
	// When the LRU capacity is exceeded, reuse the evicted limiter.
	limiters         sync.Map
	telemetryBuilder *metadata.TelemetryBuilder
}

func newLocalRateLimiter(cfg *Config, set processor.Settings, telemetryBuilder *metadata.TelemetryBuilder) (*localRateLimiter, error) {
	return &localRateLimiter{cfg: cfg, set: set, telemetryBuilder: telemetryBuilder}, nil
}

func (r *localRateLimiter) Start(ctx context.Context, host component.Host) error {
	return nil
}

func (r *localRateLimiter) Shutdown(ctx context.Context) error {
	r.limiters = sync.Map{}
	return nil
}

func (r *localRateLimiter) RateLimit(ctx context.Context, hits int) error {
	// Each (shared) processor gets its own rate limiter,
	// so it's enough to use client metadata-based unique key.
	key := getUniqueKey(ctx, r.cfg.MetadataKeys)
	v, _ := r.limiters.LoadOrStore(key, rate.NewLimiter(rate.Limit(r.cfg.Rate), r.cfg.Burst))
	limiter := v.(*rate.Limiter)
	switch r.cfg.ThrottleBehavior {
	case ThrottleBehaviorError:
		if ok := limiter.AllowN(time.Now(), hits); !ok {
			r.requestTelemetry(ctx, []attribute.KeyValue{
				telemetry.WithDecision("throttled"),
			})
			return errTooManyRequests
		}
	case ThrottleBehaviorDelay:
		lr := limiter.ReserveN(time.Now(), hits)
		if !lr.OK() {
			r.requestTelemetry(ctx, []attribute.KeyValue{
				telemetry.WithDecision("throttled"),
			})
			return errTooManyRequests
		}
		timer := time.NewTimer(lr.Delay())
		defer timer.Stop()
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-timer.C:
		}
		r.requestTelemetry(ctx, []attribute.KeyValue{
			telemetry.WithDecision("throttled"),
		})
		return nil
	}

	r.requestTelemetry(ctx, []attribute.KeyValue{
		telemetry.WithReason(telemetry.StatusUnderLimit),
		telemetry.WithDecision("accepted"),
	})
	return nil
}

func (r *localRateLimiter) requestTelemetry(ctx context.Context, baseAttrs []attribute.KeyValue) {
	attrs := attrsFromMetadata(ctx, r.cfg.MetadataKeys, baseAttrs)
	r.telemetryBuilder.RatelimitRequests.Add(ctx, 1, metric.WithAttributeSet(attribute.NewSet(attrs...)))
}
