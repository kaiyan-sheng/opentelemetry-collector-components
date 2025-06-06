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

// Code generated by mdatagen. DO NOT EDIT.

package metadatatest

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/component/componenttest"
	"go.opentelemetry.io/collector/processor"
	"go.opentelemetry.io/collector/processor/processortest"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"go.opentelemetry.io/otel/sdk/metric/metricdata/metricdatatest"
)

func NewSettings(tt *componenttest.Telemetry) processor.Settings {
	set := processortest.NewNopSettings(processortest.NopType)
	set.ID = component.NewID(component.MustNewType("ratelimit"))
	set.TelemetrySettings = tt.NewTelemetrySettings()
	return set
}

func AssertEqualRatelimitConcurrentRequests(t *testing.T, tt *componenttest.Telemetry, dps []metricdata.DataPoint[int64], opts ...metricdatatest.Option) {
	want := metricdata.Metrics{
		Name:        "otelcol_ratelimit.concurrent_requests",
		Description: "Number of in-flight requests at any given time",
		Unit:        "{requests}",
		Data: metricdata.Gauge[int64]{
			DataPoints: dps,
		},
	}
	got, err := tt.GetMetric("otelcol_ratelimit.concurrent_requests")
	require.NoError(t, err)
	metricdatatest.AssertEqual(t, want, got, opts...)
}

func AssertEqualRatelimitRequestDuration(t *testing.T, tt *componenttest.Telemetry, dps []metricdata.HistogramDataPoint[float64], opts ...metricdatatest.Option) {
	want := metricdata.Metrics{
		Name:        "otelcol_ratelimit.request_duration",
		Description: "Time(in seconds) taken to process a rate limit request",
		Unit:        "{seconds}",
		Data: metricdata.Histogram[float64]{
			Temporality: metricdata.CumulativeTemporality,
			DataPoints:  dps,
		},
	}
	got, err := tt.GetMetric("otelcol_ratelimit.request_duration")
	require.NoError(t, err)
	metricdatatest.AssertEqual(t, want, got, opts...)
}

func AssertEqualRatelimitRequests(t *testing.T, tt *componenttest.Telemetry, dps []metricdata.DataPoint[int64], opts ...metricdatatest.Option) {
	want := metricdata.Metrics{
		Name:        "otelcol_ratelimit.requests",
		Description: "Number of rate-limiting requests",
		Unit:        "{requests}",
		Data: metricdata.Sum[int64]{
			Temporality: metricdata.CumulativeTemporality,
			IsMonotonic: true,
			DataPoints:  dps,
		},
	}
	got, err := tt.GetMetric("otelcol_ratelimit.requests")
	require.NoError(t, err)
	metricdatatest.AssertEqual(t, want, got, opts...)
}
