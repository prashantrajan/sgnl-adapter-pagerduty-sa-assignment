// Copyright 2023 SGNL.ai, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package adapter

import (
	"context"
	"fmt"
	"strconv"

	framework "github.com/sgnl-ai/adapter-framework"
	api_adapter_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
)

const (
	//PagerDuty API Hostname
	APIHost = "https://api.pagerduty.com"

	// MaxPageSize is the maximum page size allowed in a GetPage request.
	//
	// SCAFFOLDING #7-OK - pkg/adapter/validation.go: Update this limit to match the limit of the SoR.
	MaxPageSize = 100

	//PagerDuty's classic pagination REST API permits retrieving a maximum of 10,000 records
	MaxResultSize = 10000
)

// ValidateGetPageRequest validates the fields of the GetPage Request.
func (a *Adapter) ValidateGetPageRequest(ctx context.Context, request *framework.Request[Config]) *framework.Error {
	if err := request.Config.Validate(ctx); err != nil {
		return &framework.Error{
			Message: fmt.Sprintf("Provided config is invalid: %v.", err.Error()),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
		}
	}

	// Validate the requested SoR address
	if request.Address != APIHost {
		return &framework.Error{
			Message: "PagerDuty API URL is invalid.",
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
		}
	}

	// SCAFFOLDING #8-OK - pkg/adapter/validation.go: Modify this validation to match the authn mechanism(s) supported by the SoR.
	if request.Auth == nil || request.Auth.HTTPAuthorization == "" {
		return &framework.Error{
			Message: "PagerDuty auth is missing required token.",
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
		}
	}

	if _, found := ValidEntityExternalIDs[request.Entity.ExternalId]; !found {
		return &framework.Error{
			Message: "Provided entity external ID is invalid.",
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_ENTITY_CONFIG,
		}
	}

	// Validate that at least the unique ID attribute for the requested entity
	// is requested.
	var uniqueIDAttributeFound bool

	for _, attribute := range request.Entity.Attributes {
		if attribute.ExternalId == ValidEntityExternalIDs[request.Entity.ExternalId].uniqueIDAttrExternalID {
			uniqueIDAttributeFound = true

			break
		}
	}

	if !uniqueIDAttributeFound {
		return &framework.Error{
			Message: "Requested entity attributes are missing unique ID attribute.",
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_ENTITY_CONFIG,
		}
	}

	// Validate that no child entities are requested.
	//
	// SCAFFOLDING #9-OK - pkg/adapter/validation.go: Modify this validation if the entity contains child entities.
	if len(request.Entity.ChildEntities) > 0 {
		return &framework.Error{
			Message: "Requested entity does not support child entities.",
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_ENTITY_CONFIG,
		}
	}

	// SCAFFOLDING #10-OK - pkg/adapter/validation.go: Check for Ordered responses.
	// If the datasource doesn't support sorting results by unique ID
	// attribute for the requested entity, check instead that Ordered is set to
	// false.
	if request.Ordered {
		return &framework.Error{
			Message: "Ordered must be set to false.",
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_ENTITY_CONFIG,
		}
	}

	if request.PageSize > MaxPageSize {
		return &framework.Error{
			Message: fmt.Sprintf("Provided page size (%d) exceeds maximum (%d).", request.PageSize, MaxPageSize),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_PAGE_REQUEST_CONFIG,
		}
	}

	if request.Cursor != "" {
		cursorInt, err := strconv.ParseInt(request.Cursor, 10, 64)
		if err != nil {
			return &framework.Error{
				Message: fmt.Sprintf("Invalid cursor value: %v.", err.Error()),
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_PAGE_REQUEST_CONFIG,
			}
		}

		if (request.PageSize + cursorInt) > MaxResultSize {
			return &framework.Error{
				Message: fmt.Sprintf("PagerDuty does not allow requesting more than %d records.", MaxResultSize),
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_PAGE_REQUEST_CONFIG,
			}
		}
	}

	return nil
}
