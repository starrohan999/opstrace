// Copyright 2021 Opstrace, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package config

import (
	"net/url"

	"github.com/opstrace/opstrace/go/pkg/graphql"
)

type TenantAccess struct {
	access *graphql.GraphqlAccess
}

func NewTenantAccess(graphqlURL *url.URL, graphqlSecret string) TenantAccess {
	return TenantAccess{
		graphql.NewGraphqlAccess(graphqlURL, graphqlSecret),
	}
}

func (c *TenantAccess) GetByName(tenantName string) (*graphql.GetTenantByNameResponse, error) {
	req, err := graphql.NewGetTenantByNameRequest(
		c.access.URL,
		&graphql.GetTenantByNameVariables{TenantName: graphql.String(tenantName)},
	)
	if err != nil {
		return nil, err
	}

	var result graphql.GetTenantByNameResponse
	if err := c.access.Execute(req.Request, &result); err != nil {
		return nil, err
	}
	if len(result.Tenant) != 1 {
		// Not found
		return nil, nil
	}
	return &result, nil
}

func (c *TenantAccess) GetByID(tenantID string) (*graphql.GetTenantByIdResponse, error) {
	req, err := graphql.NewGetTenantByIdRequest(
		c.access.URL,
		&graphql.GetTenantByIdVariables{TenantId: graphql.UUID(tenantID)},
	)
	if err != nil {
		return nil, err
	}

	var result graphql.GetTenantByIdResponse
	if err := c.access.Execute(req.Request, &result); err != nil {
		return nil, err
	}
	if result.TenantByPk.ID == "" {
		// Not found
		return nil, nil
	}
	return &result, nil
}
