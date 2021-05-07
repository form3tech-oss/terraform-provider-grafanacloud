resource "grafanacloud_grafana_api_key" "api_client" {
  name  = "api_client"
  role  = "Editor"
  stack = "demo"
}
