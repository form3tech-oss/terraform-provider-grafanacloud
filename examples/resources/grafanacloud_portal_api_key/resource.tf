resource "grafanacloud_portal_api_key" "prometheus_remote_write" {
  name = "prometheus_remote_write"
  role = "MetricsPublisher"
}
