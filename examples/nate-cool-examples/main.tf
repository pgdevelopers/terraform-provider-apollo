terraform {
  required_providers {
    apollo = {
      source  = "terraform.local/local/apollo"
      version = "1.0.0"
    }
  }
}

provider "apollo" {
  personal_api_key = "fake-api-key"
}

# resource "apollo_graph" "graph" {
#   org_id     = "fake-org"
#   graph_name = "fake-graph"
# }

# resource "apollo_apikey" "apikey" {
#   graph_id = apollo_graph.graph.id
#   key_name = "fake-key"
# }