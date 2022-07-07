<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [Prometheus Metrics](#prometheus-metrics)
  - [Subscriber Metrics](#subscriber-metrics)
    - [Scarping Metrics with Prometheus](#scarping-metrics-with-prometheus)
    - [Collecting Metrics for Observability](#collecting-metrics-for-observability)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# Prometheus Metrics

## Subscriber Metrics

The following metrics can be scrapped from *Managed Clusters*:

| Name                                   | Description                                      |
| -------------------------------------- | ------------------------------------------------ |
| subscriber_successful_checkout_count   | Counter for successful checkout process count    |
| subscriber_successful_checkout_latency | Histogram of successful checkout process latency |
| subscriber_failed_checkout_count       | Counter for failed checkout process count        |
| subscriber_failed_checkout_latency     | Histogram of failed checkout process latency     |

With every metric recorded, you can find the following custom vector labels for identifying its source:

- `subscriber_type` (*Git*/*HelmRepo*/*ObjectBucket*)
- `subscriber_namespace`
- `subscriber_name`

### Scarping Metrics with Prometheus

The metrics service, `multicluster-subscriber-metrics` resides in the `open-cluster-management-agent-addon` namespace in the *Managed Clusters*.</br>
For scraping with Prometheus, you need to create 3 resources in each *Managed Cluster*:

- a `Role` in the `open-cluster-management-agent-addon` namespace for allowing access to the metrics data.
- a `RoleBinding` in Prometheus' namespace** binding Prometheus's existing `ServiceAccount` named `prometheus-k8s` to the newly created `Role`.
- a `ServiceMonitor` in Prometheus' namespace** for explicitly configuring Prometheus to scrap our metric service.

** In *OpenShift* clusters, Prometheus' namespace is `openshift-monitoring`, for none-*OpenShift*_* clusters, it's `monitoring`.

Here's an example of the 3 required resources, please modify these to fit your environment:

```yaml
---
# create a service monitor for collecting services exposing metrics
apiVersion: monitoring.coreos.com/v1
kind: ServiceMonitor
metadata:
  name: application-manager-addon-service-monitor
  namespace: openshift-monitoring  # the namespace of the prometheus operator
spec:
  endpoints:
  - port: metrics # the name of port object exposing the metrics
  namespaceSelector:
    matchNames:
    - open-cluster-management-agent-addon # the namespace to look for the exposed metrics service
  selector:
    matchLabels:
      app: multicluster-subscriber-metrics # the label to match the services exposing metrics

---
# bind the designated role for monitoring to the prometheus-k8s service account
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: prometheus-k8s-monitoring-binding
  namespace: open-cluster-management-agent-addon # the rolebinding goes in the namespace of the metrics service
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: Role
  name: prometheus-k8s-monitoring
subjects:
- kind: ServiceAccount
  name: prometheus-k8s
  namespace: openshift-monitoring # the namespace of the prometheus operator

---
# permissions for the prometheus-k8s service account for monitoring
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: prometheus-k8s-monitoring
  namespace: open-cluster-management-agent-addon # the role goes in the namespace of the metrics service
rules:
- apiGroups:
  - ""
  resources:
  - services
  - endpoints
  - pods
  verbs:
  - get
  - list
  - watch
- apiGroups:
  - extensions
  resources:
  - ingresses
  verbs:
  - get
  - list
  - watch
- apiGroups:
  - networking.k8s.io
  resources:
  - ingresses
  verbs:
  - get
  - list
  - watch
```

### Collecting Metrics for Observability

For the [Observability Operator](https://github.com/stolostron/multicluster-observability-operator) to collect the aforementioned metrics from the *Managed Clusters* and display them on the *Hub Cluster*, you need to configure the `observability-metrics-custom-allowlist` *ConfigMap* in the `open-cluster-management-observability` namespace on the *Hub Cluster*.</br>
Note that for *Histograms* type metrics, there are actually 3 metrics created per each *Histogram*, they are identified by the *bucket*, *count*, and *sum* prefixes to the metric name. Here's an example of a working *ConfigMap*:

```yaml
apiVersion: v1
kind: ConfigMap
data:
  metrics_list.yaml: |
    names:
    - subscriber_successful_checkout_count
    - subscriber_successful_checkout_latency_bucket
    - subscriber_successful_checkout_latency_count
    - subscriber_successful_checkout_latency_sum
    - subscriber_failed_checkout_count
    - subscriber_failed_checkout_latency_bucket
    - subscriber_failed_checkout_latency_count
    - subscriber_failed_checkout_latency_sum
```
