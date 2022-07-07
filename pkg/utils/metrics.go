// Copyright 2021 The Kubernetes Authors.
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

package utils

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Collector for a checkout process
type CheckoutSummary struct {
	SuccessfulCount     int
	FailedCount         int
	SuccessfulLatencyMS int
	FailedLatencyMS     int
}

// Subscriber type enum
type MetricSubscriberType int

//nolint // underscores in var names are required for readability
const (
	// Subscriber types
	MetricSubscriberType_Git          MetricSubscriberType = 0
	MetricSubscriberType_HelmRepo     MetricSubscriberType = 1
	MetricSubscriberType_ObjectBucket MetricSubscriberType = 2
)

// Return the string value of the various subscribers
func (subscriberType MetricSubscriberType) ToString() string {
	switch subscriberType {
	case MetricSubscriberType_Git:
		return "Git"
	case MetricSubscriberType_HelmRepo:
		return "HelmRepo"
	case MetricSubscriberType_ObjectBucket:
		return "ObjectBucket"
	default:
		return "Unknown"
	}
}

var subscriberVectorLabels = []string{"subscriber_type", "subscriber_namespace", "subscriber_name"}

// #################
// #### Metrics ####
// #################
var SuccessfulCheckoutCount = *promauto.NewCounterVec(prometheus.CounterOpts{
	Name: "subscriber_successful_checkout_count",
	Help: "Counter for successful checkout process count",
}, subscriberVectorLabels)

var SuccessfulCheckoutLatency = *promauto.NewHistogramVec(prometheus.HistogramOpts{
	Name: "subscriber_successful_checkout_latency",
	Help: "Histogram of successful checkout process latency",
}, subscriberVectorLabels)

var FailedCheckoutCount = *promauto.NewCounterVec(prometheus.CounterOpts{
	Name: "subscriber_failed_checkout_count",
	Help: "Counter for failed checkout process count",
}, subscriberVectorLabels)

var FailedCheckoutLatency = *promauto.NewHistogramVec(prometheus.HistogramOpts{
	Name: "subscriber_failed_checkout_latency",
	Help: "Histogram of failed checkout process latency",
}, subscriberVectorLabels)

// ###########################
// #### Utility functions ####
// ###########################

// Update the checkout various checkout metrics
func UpdateCheckoutMetrics(subType MetricSubscriberType, subNamespace string, subName string, checkoutSummary CheckoutSummary) {
	if checkoutSummary.SuccessfulCount > 0 {
		// Update the successful checkout count
		SuccessfulCheckoutCount.WithLabelValues(subType.ToString(), subNamespace, subName).
			Add(float64(checkoutSummary.SuccessfulCount))
		// Update the successful checkout latency
		SuccessfulCheckoutLatency.WithLabelValues(subType.ToString(), subNamespace, subName).
			Observe(float64(checkoutSummary.SuccessfulLatencyMS))
	}

	if checkoutSummary.FailedCount > 0 {
		// Update the failed checkout count
		FailedCheckoutCount.WithLabelValues(subType.ToString(), subNamespace, subName).
			Add(float64(checkoutSummary.FailedCount))
		// Update the failed checkout latency
		FailedCheckoutLatency.WithLabelValues(subType.ToString(), subNamespace, subName).
			Observe(float64(checkoutSummary.FailedLatencyMS))
	}
}
