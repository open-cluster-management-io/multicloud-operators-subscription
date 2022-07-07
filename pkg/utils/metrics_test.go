/*
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package utils

import (
	"testing"

	"github.com/onsi/gomega"
	. "github.com/onsi/gomega"
	"github.com/prometheus/client_golang/prometheus/testutil"
)

type subscriberSut struct {
	Type            MetricSubscriberType
	Namespace       string
	Name            string
	CheckoutSummary CheckoutSummary
}

var subscriber1Git = subscriberSut{
	Type:      MetricSubscriberType_Git,
	Name:      "fake-subscriber-git",
	Namespace: "fake-namespace",
	CheckoutSummary: CheckoutSummary{ // will update successful and failed metrics
		SuccessfulCount:     1,
		SuccessfulLatencyMS: 699,
		FailedCount:         1,
		FailedLatencyMS:     1566,
	},
}

var subscriber2Helm = subscriberSut{
	Type:      MetricSubscriberType_HelmRepo,
	Name:      "fake-subscriber-helm",
	Namespace: "another-fake-namespace",
	CheckoutSummary: CheckoutSummary{ // will only update failed metrics
		SuccessfulCount:     0,
		SuccessfulLatencyMS: 0,
		FailedCount:         2,
		FailedLatencyMS:     3544,
	},
}

var subscriber3Bucket = subscriberSut{ // will only update successful metrics
	Type:      MetricSubscriberType_ObjectBucket,
	Name:      "fake-subscriber-bucket",
	Namespace: "one-more-fake-namespace",
	CheckoutSummary: CheckoutSummary{
		SuccessfulCount:     1,
		SuccessfulLatencyMS: 321,
		FailedCount:         0,
		FailedLatencyMS:     0,
	},
}

func TestUpdateCheckoutMetrics(t *testing.T) {
	g := gomega.NewGomegaWithT(t)

	// when all subscribers invoke the update checkout metrics utility function
	UpdateCheckoutMetrics(
		subscriber1Git.Type, subscriber1Git.Namespace, subscriber1Git.Name, subscriber1Git.CheckoutSummary)
	UpdateCheckoutMetrics(
		subscriber2Helm.Type, subscriber2Helm.Namespace, subscriber2Helm.Name, subscriber2Helm.CheckoutSummary)
	UpdateCheckoutMetrics(
		subscriber3Bucket.Type, subscriber3Bucket.Namespace, subscriber3Bucket.Name, subscriber3Bucket.CheckoutSummary)

	// then verify metrics data based on the fixture data (when success, failure is 0 - will not update respectively)
	g.Expect(2).To(Equal(testutil.CollectAndCount(SuccessfulCheckoutCount)))
	g.Expect(2).To(Equal(testutil.CollectAndCount(SuccessfulCheckoutLatency)))
	g.Expect(2).To(Equal(testutil.CollectAndCount(FailedCheckoutCount)))
	g.Expect(2).To(Equal(testutil.CollectAndCount(FailedCheckoutLatency)))

	// and verify subscriber 1 (git) metrics
	g.Expect(float64(subscriber1Git.CheckoutSummary.SuccessfulCount)).
		To(Equal(testutil.ToFloat64(SuccessfulCheckoutCount.
			WithLabelValues(subscriber1Git.Type.ToString(), subscriber1Git.Namespace, subscriber1Git.Name))))
	g.Expect(float64(subscriber1Git.CheckoutSummary.FailedCount)).
		To(Equal(testutil.ToFloat64(FailedCheckoutCount.
			WithLabelValues(subscriber1Git.Type.ToString(), subscriber1Git.Namespace, subscriber1Git.Name))))

	// and verify subscriber 2 (helm) metrics
	g.Expect(float64(subscriber2Helm.CheckoutSummary.SuccessfulCount)).
		To(Equal(testutil.ToFloat64(SuccessfulCheckoutCount.
			WithLabelValues(subscriber2Helm.Type.ToString(), subscriber2Helm.Namespace, subscriber2Helm.Name))))
	g.Expect(float64(subscriber2Helm.CheckoutSummary.FailedCount)).
		To(Equal(testutil.ToFloat64(FailedCheckoutCount.
			WithLabelValues(subscriber2Helm.Type.ToString(), subscriber2Helm.Namespace, subscriber2Helm.Name))))

	// and verify subscriber 3 (bucket) metrics
	g.Expect(float64(subscriber3Bucket.CheckoutSummary.SuccessfulCount)).
		To(Equal(testutil.ToFloat64(SuccessfulCheckoutCount.
			WithLabelValues(subscriber3Bucket.Type.ToString(), subscriber3Bucket.Namespace, subscriber3Bucket.Name))))
	g.Expect(float64(subscriber3Bucket.CheckoutSummary.FailedCount)).
		To(Equal(testutil.ToFloat64(FailedCheckoutCount.
			WithLabelValues(subscriber3Bucket.Type.ToString(), subscriber3Bucket.Namespace, subscriber3Bucket.Name))))
}
