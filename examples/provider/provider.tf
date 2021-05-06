provider "grafanacloud" {
  # Can also be provided via GRAFANA_CLOUD_API_KEY
  api_key = var.your_secret_api_key

  # Can also be provided via GRAFANA_CLOUD_ORGANISATION
  organisation = "org-slug"
}
