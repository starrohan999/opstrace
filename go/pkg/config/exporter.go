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

type ExporterAccess struct {
	access *graphql.GraphqlAccess
}

func NewExporterAccess(graphqlURL *url.URL, graphqlSecret string) ExporterAccess {
	return ExporterAccess{
		graphql.NewGraphqlAccess(graphqlURL, graphqlSecret),
	}
}

func (c *ExporterAccess) ListID(tenantID string) (*graphql.GetExportersResponse, error) {
	req, err := graphql.NewGetExportersRequest(
		c.access.URL,
		&graphql.GetExportersVariables{TenantId: graphql.UUID(tenantID)},
	)
	if err != nil {
		return nil, err
	}

	var result graphql.GetExportersResponse
	if err := c.access.Execute(req.Request, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *ExporterAccess) GetIDByName(tenantID string, name string) (*string, error) {
	req, err := graphql.NewGetExporterIdByNameRequest(
		c.access.URL,
		&graphql.GetExporterIdByNameVariables{TenantId: graphql.UUID(tenantID), Name: graphql.String(name)},
	)
	if err != nil {
		return nil, err
	}

	var result graphql.GetExporterIdByNameResponse
	if err := c.access.Execute(req.Request, &result); err != nil {
		return nil, err
	}
	if len(result.Exporter) != 1 {
		// Not found
		return nil, nil
	}
	return &result.Exporter[0].ID, nil
}

func (c *ExporterAccess) GetID(tenantID string, id string) (*graphql.GetExporterResponse, error) {
	req, err := graphql.NewGetExporterRequest(
		c.access.URL,
		&graphql.GetExporterVariables{TenantId: graphql.UUID(tenantID), ID: graphql.UUID(id)},
	)
	if err != nil {
		return nil, err
	}

	var result graphql.GetExporterResponse
	if err := c.access.Execute(req.Request, &result); err != nil {
		return nil, err
	}
	if len(result.Exporter) != 1 {
		// Not found
		return nil, nil
	}
	return &result, nil
}

func (c *ExporterAccess) DeleteID(tenantID string, id string) (*string, error) {
	req, err := graphql.NewDeleteExporterRequest(
		c.access.URL,
		&graphql.DeleteExporterVariables{TenantId: graphql.UUID(tenantID), ID: graphql.UUID(id)},
	)
	if err != nil {
		return nil, err
	}

	var result graphql.DeleteExporterResponse
	if err := c.access.Execute(req.Request, &result); err != nil {
		return nil, err
	}
	if len(result.DeleteExporter.Returning) != 1 {
		// Not found
		return nil, nil
	}
	return &result.DeleteExporter.Returning[0].ID, nil
}

// Insert inserts one or more exporters, returns an error if any already exists.
func (c *ExporterAccess) InsertID(tenantID string, inserts []graphql.ExporterInsertInput) error {
	// Ensure the inserts each have the correct tenant
	gtenantID := graphql.UUID(tenantID)
	insertsWithTenant := make([]graphql.ExporterInsertInput, 0)
	for _, insert := range inserts {
		insert.TenantId = &gtenantID
		insertsWithTenant = append(insertsWithTenant, insert)
	}

	req, err := graphql.NewCreateExportersRequest(
		c.access.URL,
		&graphql.CreateExportersVariables{Exporters: &insertsWithTenant},
	)
	if err != nil {
		return err
	}

	var result graphql.CreateExportersResponse
	return c.access.Execute(req.Request, &result)
}

// Update updates an existing exporter, returns an error if a exporter of the same tenant/name doesn't exist.
func (c *ExporterAccess) UpdateID(tenantID string, update graphql.UpdateExporterVariables) error {
	// Ensure the update has the correct tenant
	update.TenantId = graphql.UUID(tenantID)

	req, err := graphql.NewUpdateExporterRequest(c.access.URL, &update)
	if err != nil {
		return err
	}

	var result graphql.UpdateExporterResponse
	return c.access.Execute(req.Request, &result)
}
