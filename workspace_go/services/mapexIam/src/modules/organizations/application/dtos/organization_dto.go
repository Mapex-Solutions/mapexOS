package dtos

import (
	v1 "github.com/Mapex-Solutions/MapexOS/contracts/services/mapexIam/organizations"
)

type (
	OrganizationIdDto                = v1.OrganizationId
	AddressDTO                       = v1.Address
	AuthConfigDTO                    = v1.AuthConfig
	AccessPolicyDTO                  = v1.AccessPolicy
	CreateOrganizationDto            = v1.OrganizationCreate
	UpdateOrganizationDto            = v1.OrganizationUpdate
	OrganizationQueryDto             = v1.OrganizationQuery
	OrganizationResponse             = v1.OrganizationResponse
	TreeQueryDto                     = v1.TreeQuery
	TreeItemDto                      = v1.TreeItem
	TreeResponseDto                  = v1.TreeResponse
)
