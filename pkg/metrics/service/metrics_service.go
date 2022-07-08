/*
Copyright 2021 The Kubernetes Authors.

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

package service

import (
	"context"
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"k8s.io/klog"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

const metricServicePort = 8381

// Service listener implementation
type MetricServiceListener struct{}

// Start listener and link to prometheus http handler
func (listener *MetricServiceListener) Start(context.Context) error {
	http.HandleFunc("/metrics", promhttp.Handler().ServeHTTP)
	klog.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", metricServicePort), nil))

	return nil
}

// Add the metrics service to the manager
func Add(mgr manager.Manager) error {
	if err := mgr.Add(&MetricServiceListener{}); err != nil {
		return err
	}

	return nil
}
