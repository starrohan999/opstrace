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

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"

	"github.com/opstrace/opstrace/go/pkg/config"
	"github.com/opstrace/opstrace/go/pkg/graphql"
)

// Information about an exporter. Custom type which omits the tenant field.
type ExporterInfo struct {
	ID         string      `yaml:"id"`
	Name       string      `yaml:"name"`
	Type       string      `yaml:"type,omitempty"`
	Credential string      `yaml:"credential,omitempty"`
	Config     interface{} `yaml:"config,omitempty"`
	CreatedAt  string      `yaml:"created_at,omitempty"`
	UpdatedAt  string      `yaml:"updated_at,omitempty"`
}

// Exporter entry for validation and interaction with GraphQL.
type Exporter struct {
	ID         string
	Name       string
	Type       string
	Credential string
	ConfigJSON string
}

// Raw exporter entry received from a POST request, converted to an Exporter.
type yamlExporter struct {
	ID         string      `yaml:"id"`
	Name       string      `yaml:"name"`
	Type       string      `yaml:"type"`
	Credential string      `yaml:"credential,omitempty"`
	Config     interface{} `yaml:"config"` // nested yaml
}

type exporterAPI struct {
	credentialAccess *config.CredentialAccess
	exporterAccess   *config.ExporterAccess
	tenantAccess     *config.TenantAccess
}

func newExporterAPI(
	credentialAccess *config.CredentialAccess,
	exporterAccess *config.ExporterAccess,
	tenantAccess *config.TenantAccess,
) *exporterAPI {
	return &exporterAPI{
		credentialAccess,
		exporterAccess,
		tenantAccess,
	}
}

func (e *exporterAPI) listExporters(tenant string, w http.ResponseWriter, r *http.Request) {
	tenantInfo, err := e.tenantAccess.GetByName(tenant)
	if err != nil {
		log.Warnf("Fetching tenant %s failed: %s", tenant, err)
		http.Error(w, fmt.Sprintf("Fetching tenant %s failed: %s", tenant, err), http.StatusInternalServerError)
		return
	}
	if tenantInfo == nil {
		log.Warnf("Tenant %s not found", tenant)
		http.Error(w, fmt.Sprintf("Tenant %s not found", tenant), http.StatusNotFound)
		return
	}
	tenantID := tenantInfo.Tenant[0].ID

	resp, err := e.exporterAccess.ListID(tenantID)
	if err != nil {
		log.Warnf("Listing exporters failed: %s", err)
		http.Error(w, fmt.Sprintf("Listing exporters failed: %s", err), http.StatusInternalServerError)
		return
	}

	log.Debugf("Listing %d exporters", len(resp.Exporter))

	// Create list payload to respond with.
	// Avoid passing entries individually to encoder since that won't consistently produce a list.
	entries := make([]ExporterInfo, len(resp.Exporter))
	for i, exporter := range resp.Exporter {
		configJSON := make(map[string]interface{})
		err := json.Unmarshal([]byte(exporter.Config), &configJSON)
		if err != nil {
			// give up and pass-through the json
			log.Warnf("Failed to decode JSON config for exporter %s (err: %s): %s", exporter.Name, err, exporter.Config)
			configJSON["json"] = exporter.Config
		}
		entries[i] = ExporterInfo{
			Name:       exporter.Name,
			Type:       exporter.Type,
			Credential: exporter.Credential.ID,
			Config:     configJSON,
			CreatedAt:  exporter.CreatedAt,
			UpdatedAt:  exporter.UpdatedAt,
		}
	}

	encoder := yaml.NewEncoder(w)
	encoder.Encode(entries)
}

func (e *exporterAPI) writeExporters(tenant string, w http.ResponseWriter, r *http.Request) {
	tenantInfo, err := e.tenantAccess.GetByName(tenant)
	if err != nil {
		log.Warnf("Fetching tenant %s failed: %s", tenant, err)
		http.Error(w, fmt.Sprintf("Fetching tenant %s failed: %s", tenant, err), http.StatusInternalServerError)
		return
	}
	if tenantInfo == nil {
		log.Warnf("Tenant %s not found", tenant)
		http.Error(w, fmt.Sprintf("Tenant %s not found", tenant), http.StatusNotFound)
		return
	}
	tenantID := tenantInfo.Tenant[0].ID

	// Collect map of existing name->type so that we can decide between insert vs update
	existingTypes, err := e.listExporterTypesID(tenantID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Listing exporters failed: %s", err), http.StatusInternalServerError)
		return
	}

	now := nowTimestamp()

	decoder := yaml.NewDecoder(r.Body)
	// Return error for unrecognized or duplicate fields in the input
	decoder.SetStrict(true)

	var inserts []graphql.ExporterInsertInput
	var updates []graphql.UpdateExporterVariables
	for {
		var yamlExporter yamlExporter
		if err := decoder.Decode(&yamlExporter); err != nil {
			if err != io.EOF {
				msg := fmt.Sprintf("Decoding exporter input at index=%d failed: %s", len(inserts)+len(updates), err)
				log.Debug(msg)
				http.Error(w, msg, http.StatusBadRequest)
				return
			}
			break
		}

		var credential *graphql.UUID
		if yamlExporter.Credential == "" {
			credential = nil
		} else {
			// TODO name -> uuid
			gcredential := graphql.UUID(yamlExporter.Credential)
			credential = &gcredential
		}

		config, err := convertYAMLExporterConfig(yamlExporter.Name, yamlExporter.Config)
		if err != nil {
			log.Debugf("Parsing exporter input at index=%d failed: %s", len(inserts)+len(updates), err)
			http.Error(w, fmt.Sprintf(
				"Parsing exporter input at index=%d failed: %s", len(inserts)+len(updates), err,
			), http.StatusBadRequest)
			return
		}

		exists, err := e.validateExporterID(
			tenantID,
			existingTypes,
			Exporter{
				Name:       yamlExporter.Name,
				Type:       yamlExporter.Type,
				Credential: yamlExporter.Credential,
				ConfigJSON: *config,
			},
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// TODO in insert case, generate ID if not provided
		id := graphql.String(yamlExporter.ID)
		gid := graphql.UUID(id)
		name := graphql.String(yamlExporter.Name)
		gconfig := graphql.Json(*config)
		if exists {
			// TODO check for no-op updates and skip them (and avoid unnecessary changes to UpdatedAt)
			updates = append(updates, graphql.UpdateExporterVariables{
				ID:           gid,
				Name:         name,
				CredentialId: credential,
				Config:       gconfig,
				UpdatedAt:    now,
			})
		} else {
			expType := graphql.String(yamlExporter.Type)
			inserts = append(inserts, graphql.ExporterInsertInput{
				ID:           &gid,
				Name:         &name,
				Type:         &expType,
				CredentialId: credential,
				Config:       &gconfig,
				CreatedAt:    &now,
				UpdatedAt:    &now,
			})
		}
	}

	if len(inserts)+len(updates) == 0 {
		log.Debugf("Writing exporters: No data provided")
		http.Error(w, "Missing exporter YAML data in request body", http.StatusBadRequest)
		return
	}

	log.Debugf("Writing exporters: %d insert, %d update", len(inserts), len(updates))

	if len(inserts) != 0 {
		if err := e.exporterAccess.InsertID(tenantID, inserts); err != nil {
			log.Warnf("Insert: %d exporters failed: %s", len(inserts), err)
			http.Error(w, fmt.Sprintf("Creating %d exporters failed: %s", len(inserts), err), http.StatusInternalServerError)
			return
		}
	}
	if len(updates) != 0 {
		for _, update := range updates {
			if err := e.exporterAccess.UpdateID(tenantID, update); err != nil {
				log.Warnf("Update: Exporter %s/%s failed: %s", tenant, update.Name, err)
				http.Error(w, fmt.Sprintf("Updating exporter %s failed: %s", update.Name, err), http.StatusInternalServerError)
				return
			}
		}
	}
}

func (e *exporterAPI) getExporter(tenant string, w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["name"]
	log.Debugf("Getting exporter: %s/%s", tenant, name)

	tenantInfo, err := e.tenantAccess.GetByName(tenant)
	if err != nil {
		log.Warnf("Fetching tenant %s failed: %s", tenant, err)
		http.Error(w, fmt.Sprintf("Fetching tenant %s failed: %s", tenant, err), http.StatusInternalServerError)
		return
	}
	if tenantInfo == nil {
		log.Warnf("Tenant %s not found", tenant)
		http.Error(w, fmt.Sprintf("Tenant %s not found", tenant), http.StatusNotFound)
		return
	}
	tenantID := tenantInfo.Tenant[0].ID

	id, err := e.exporterAccess.GetIDByName(tenantID, name)
	if err != nil {
		log.Warnf("Get: Credential %s ID failed: %s", name, err)
		http.Error(w, fmt.Sprintf("Getting credential ID failed: %s", err), http.StatusInternalServerError)
		return
	}
	if id == nil {
		log.Debugf("Get: Credential %s ID not found", name)
		http.Error(w, fmt.Sprintf("Credential ID not found: %s", name), http.StatusNotFound)
		return
	}

	resp, err := e.exporterAccess.GetID(tenantID, *id)
	if err != nil {
		log.Warnf("Get: Exporter %s/%s failed: %s", tenant, name, err)
		http.Error(w, fmt.Sprintf("Getting exporter failed: %s", err), http.StatusInternalServerError)
		return
	}
	if resp == nil {
		log.Debugf("Get: Exporter %s/%s not found", tenant, name)
		http.Error(w, fmt.Sprintf("Exporter not found: %s", name), http.StatusNotFound)
		return
	}

	configJSON := make(map[string]interface{})
	if err = json.Unmarshal([]byte(resp.Exporter[0].Config), &configJSON); err != nil {
		// give up and pass-through the json
		log.Warnf(
			"Failed to decode JSON config for exporter %s (err: %s): %s",
			resp.Exporter[0].Name, err, resp.Exporter[0].Config,
		)
		configJSON["json"] = resp.Exporter[0].Config
	}

	encoder := yaml.NewEncoder(w)
	encoder.Encode(ExporterInfo{
		Name:       resp.Exporter[0].Name,
		Type:       resp.Exporter[0].Type,
		Credential: resp.Exporter[0].Credential.Name,
		Config:     configJSON,
		CreatedAt:  resp.Exporter[0].CreatedAt,
		UpdatedAt:  resp.Exporter[0].UpdatedAt,
	})
}

func (e *exporterAPI) deleteExporter(tenant string, w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["name"]
	log.Debugf("Deleting exporter: %s/%s", tenant, name)

	tenantInfo, err := e.tenantAccess.GetByName(tenant)
	if err != nil {
		log.Warnf("Fetching tenant %s failed: %s", tenant, err)
		http.Error(w, fmt.Sprintf("Fetching tenant %s failed: %s", tenant, err), http.StatusInternalServerError)
		return
	}
	if tenantInfo == nil {
		log.Warnf("Tenant %s not found", tenant)
		http.Error(w, fmt.Sprintf("Tenant %s not found", tenant), http.StatusNotFound)
		return
	}
	tenantID := tenantInfo.Tenant[0].ID

	reqid, err := e.exporterAccess.GetIDByName(tenantID, name)
	if err != nil {
		log.Warnf("Get: Credential %s ID failed: %s", name, err)
		http.Error(w, fmt.Sprintf("Getting credential ID failed: %s", err), http.StatusInternalServerError)
		return
	}
	if reqid == nil {
		log.Debugf("Get: Credential %s ID not found", name)
		http.Error(w, fmt.Sprintf("Credential ID not found: %s", name), http.StatusNotFound)
		return
	}

	delid, err := e.exporterAccess.DeleteID(tenantID, *reqid)
	if err != nil {
		log.Warnf("Delete: Exporter %s/%s failed: %s", tenant, name, err)
		http.Error(w, fmt.Sprintf("Deleting exporter failed: %s", err), http.StatusInternalServerError)
		return
	}
	if delid == nil {
		log.Debugf("Delete: Exporter %s/%s not found", tenant, name)
		http.Error(w, fmt.Sprintf("Exporter not found: %s", name), http.StatusNotFound)
		return
	}

	encoder := yaml.NewEncoder(w)
	encoder.Encode(ExporterInfo{Name: *delid})
}

func (e *exporterAPI) listExporterTypesID(tenantID string) (map[string]string, error) {
	// Collect map of existing name->type so that we can decide between insert vs update
	existingTypes := make(map[string]string)
	resp, err := e.exporterAccess.ListID(tenantID)
	if err != nil {
		log.Warnf("Listing exporters failed: %s", err)
		return nil, err
	}
	for _, exporter := range resp.Exporter {
		existingTypes[exporter.Name] = exporter.Type
	}
	return existingTypes, nil
}

// Accepts the tenant name, the name->type mapping of any existing exporters, and the new exporter payload.
// Returns whether the exporter already exists, and any validation error.
func (e *exporterAPI) validateExporterID(
	tenantID string,
	existingTypes map[string]string,
	exporter Exporter,
) (bool, error) {
	// Check that the exporter name is suitable for use in K8s object names
	if err := config.ValidateName(exporter.Name); err != nil {
		return false, err
	}

	if exporter.Credential == "" {
		if err := validateExporterTypes(exporter.Type, nil); err != nil {
			msg := fmt.Sprintf("Invalid exporter input %s: %s", exporter.Name, err)
			log.Debug(msg)
			return false, errors.New(msg)
		}
	} else {
		// Check that the referenced credential exists and has a compatible type for this exporter.
		// If the credential didn't exist, then the graphql insert would fail anyway due to a missing relation,
		// but type mismatches are not validated by graphql.
		reqid, err := e.credentialAccess.GetIDByName(tenantID, exporter.Credential)
		if err != nil {
			return false, fmt.Errorf(
				"failed to read credential %s id referenced in new exporter %s",
				exporter.Credential, exporter.Name,
			)
		}
		if reqid == nil {
			return false, fmt.Errorf(
				"missing credential %s referenced in exporter %s", exporter.Credential, exporter.Name,
			)
		}

		cred, err := e.credentialAccess.GetID(tenantID, *reqid)
		if err != nil {
			return false, fmt.Errorf(
				"failed to read credential %s referenced in new exporter %s",
				exporter.Credential, exporter.Name,
			)
		} else if cred == nil {
			return false, fmt.Errorf(
				"missing credential %s referenced in exporter %s", exporter.Credential, exporter.Name,
			)
		} else if err := validateExporterTypes(exporter.Type, &cred.Credential[0].Type); err != nil {
			return false, fmt.Errorf("invalid exporter input %s: %s", exporter.Name, err)
		}
	}

	var existingType string
	var exists bool
	if existingType, exists = existingTypes[exporter.Name]; exists {
		// Explicitly check and complain if the user tries to change the exporter type
		if exporter.Type != "" && existingType != exporter.Type {
			return false, fmt.Errorf(
				"Exporter '%s' type cannot be updated (current=%s, updated=%s)",
				exporter.Name, existingType, exporter.Type,
			)
		}
	}
	return exists, nil
}
