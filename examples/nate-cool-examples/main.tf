terraform {
  required_providers {
    apollo = {
      source  = "terraform.local/local/apollo"
      version = "1.0.0"
    }
  }
}

provider "apollo" {
  personal_api_key = "user:po.proctor-and-gamble.EN9763:BQHs7LrtV_B9f358ZqenqQ"
}

resource "apollo_graph" "graph" {
  org_id     = "procter-gamble"
  graph_name = "testing-graph-123"
}