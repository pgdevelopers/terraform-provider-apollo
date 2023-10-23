package client

import (
	"context"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/machinebox/graphql"
)

type Client struct {
	ApiKey            string
	EnterPriseEnabled bool
	GraphClient       *graphql.Client
}

func (cl *Client) Init() {
	cl.GraphClient = graphql.NewClient("https://graphql.api.apollographql.com/api/graphql")
}

func (cl *Client) Query(c context.Context, q string, response interface{}) error {

	c = tflog.SetField(c, "query", q)
	graphqlRequest := graphql.NewRequest(q)
	graphqlRequest.Header.Add("X-API-Key", cl.ApiKey)
	c = tflog.SetField(c, "response", response)

	return cl.GraphClient.Run(c, graphqlRequest, &response)

}
