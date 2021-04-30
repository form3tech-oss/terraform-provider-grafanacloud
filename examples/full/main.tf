provider "grafanacloud" {
  organisation = var.org_name
}

resource "grafanacloud_stack" "default" {
  name = var.stack_slug
  slug = var.stack_slug
}

resource "grafanacloud_portal_api_key" "prometheus_remote_write" {
  name = "prometheus-remote-write"
  role = "MetricsPublisher"
}

resource "grafanacloud_grafana_api_key" "api_client" {
  name  = "api_client"
  role  = "Editor"
  stack = grafanacloud_stack.default.slug
}

resource "grafanacloud_grafana_api_key" "expires" {
  name            = "expires"
  role            = "Viewer"
  stack           = grafanacloud_stack.default.slug
  seconds_to_live = 10
}

data "grafanacloud_stack" "demo" {
  slug = grafanacloud_stack.default.slug
}

data "grafanacloud_stacks" "all" {
  depends_on = [
    grafanacloud_stack.default
  ]
}
