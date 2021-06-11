terraform {
  required_providers {
    po = {
      source = "github.com/feniix/po"
    }
  }
}

provider "po" {
  config_path    = "../../config/kubeconfig"
  config_context = "kind-kind"
}

resource "po_service_monitor" "example" {
  metadata {
    name = "example"
    namespace = "default"
    labels = {
      "k8s-app" = "label1"
    }
  }
  spec {
    endpoints {
      port = "http-metrics"
      interval = "30s"
    }
    job_label = "myapp"
    namespace_selector {
      match_names = ["default"]
    }
    selector {
      match_labels = {
        "k8s-app" = "myapplabel"
      }
    }
  }
}