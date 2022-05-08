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

type MetricSubscriberType int

const (
	// Subscriber types
	MetricSubscriberType_Git          MetricSubscriberType = 0
	MetricSubscriberType_HelmRepo     MetricSubscriberType = 1
	MetricSubscriberType_ObjectBucket MetricSubscriberType = 2
	// Checkout metrics info
	metricInfo_SuccessfulCheckoutCount_Name   = "subscriber_successful_checkout_count"
	metricInfo_SuccessfulCheckoutCount_Help   = "Counter for successful checkout process count"
	metricInfo_SuccessfulCheckoutLatency_Name = "subscriber_successful_checkout_latency"
	metricInfo_SuccessfulCheckoutLatency_Help = "Histogram of successful checkout process latency"
	metricInfo_FailedCheckoutCount_Name       = "subscriber_failed_checkout_count"
	metricInfo_FailedCheckoutCount_Help       = "Counter for failed checkout process count"
	metricInfo_FailedCheckoutLatency_Name     = "subscriber_failed_checkout_latency"
	metricInfo_FailedCheckoutLatency_Help     = "Histogram of failed checkout process latency"
)

// Vector labels
var metricLabels_Subscriber = []string{"sub_type", "sub_namespace", "sub_name"}

// Metrics store
type SubscriberMetricStore struct {
	successfulCheckoutCount   *prometheus.CounterVec
	successfulCheckoutLatency *prometheus.HistogramVec
	failedCount               *prometheus.CounterVec
	failedLatency             *prometheus.HistogramVec
}

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

// Get a new instance of the metrics store
func NewSubscriberMetricStore() SubscriberMetricStore {
	return SubscriberMetricStore{
		successfulCheckoutCount:   newSuccessfulCheckoutCountCounter(),
		successfulCheckoutLatency: newSuccessfulCheckoutLatencyHistogram(),
		failedCount:               newFailedCheckoutCountCounter(),
		failedLatency:             newFailedCheckoutLatencyHistogram(),
	}
}

// Create a new counter for incrementing checkout successful count
func newSuccessfulCheckoutCountCounter() *prometheus.CounterVec {
	return promauto.NewCounterVec(prometheus.CounterOpts{
		Name: metricInfo_SuccessfulCheckoutCount_Name,
		Help: metricInfo_SuccessfulCheckoutCount_Help,
	}, metricLabels_Subscriber)
}

// Create a new histogram for aggregating checkout successful latency
func newSuccessfulCheckoutLatencyHistogram() *prometheus.HistogramVec {
	return promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: metricInfo_SuccessfulCheckoutLatency_Name,
		Help: metricInfo_SuccessfulCheckoutLatency_Help,
	}, metricLabels_Subscriber)
}

// Create a new counter for incrementing checkout failed count
func newFailedCheckoutCountCounter() *prometheus.CounterVec {
	return promauto.NewCounterVec(prometheus.CounterOpts{
		Name: metricInfo_FailedCheckoutCount_Name,
		Help: metricInfo_FailedCheckoutCount_Help,
	}, metricLabels_Subscriber)
}

// Create a new histogram for aggregating checkout failed latency
func newFailedCheckoutLatencyHistogram() *prometheus.HistogramVec {
	return promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: metricInfo_FailedCheckoutLatency_Name,
		Help: metricInfo_FailedCheckoutLatency_Help,
	}, metricLabels_Subscriber)
}

// Update the checkout various checkout metrics
func UpdateCheckoutMetrics(
	subType MetricSubscriberType, subNamespace string, subName string,
	checkoutSummary CheckoutSummary, metricStore SubscriberMetricStore) {

	if checkoutSummary.SuccessfulCount > 0 { // Checkout was successful
		// Update the successful checkout count
		metricStore.successfulCheckoutCount.
			WithLabelValues(subType.ToString(), subNamespace, subName).
			Add(float64(checkoutSummary.SuccessfulCount))
		// Update the successful checkout latency
		metricStore.successfulCheckoutLatency.
			WithLabelValues(subType.ToString(), subNamespace, subName).
			Observe(float64(checkoutSummary.SuccessfulLatencyMS))
	} else if checkoutSummary.FailedCount > 0 { // Checkout has failed
		// Update the failed checkout count
		metricStore.failedCount.
			WithLabelValues(subType.ToString(), subNamespace, subName).
			Add(float64(checkoutSummary.FailedCount))
		// Update the failed checkout latency
		metricStore.failedLatency.
			WithLabelValues(subType.ToString(), subNamespace, subName).
			Observe(float64(checkoutSummary.FailedLatencyMS))
	}
}
