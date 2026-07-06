package dtos

import (
	"time"

	query "github.com/Mapex-Solutions/MapexOS/contracts/common/query"
)

// OTAPlanQueryDTO is the query for the plans /find endpoint.
type OTAPlanQueryDTO struct {
	query.BaseQueryDTO
	Name   *string `query:"name"`
	Status *string `query:"status"`
}

// OTAExecutionQueryDTO is the query for a plan's executions /find endpoint.
type OTAExecutionQueryDTO struct {
	query.BaseQueryDTO
	State   *string `query:"state"`
	AssetID *string `query:"assetId"`
}

// OTAPlanUpdateRequest is the partial-update body for a plan (only provided
// fields are applied).
type OTAPlanUpdateRequest struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
}

// OTAPlanIdDto is the :planId path parameter.
type OTAPlanIdDto struct {
	PlanId string `params:"planId" validate:"required,mongoid"`
}

// OTAFirmwareIdDto is the :firmwareId path parameter.
type OTAFirmwareIdDto struct {
	FirmwareId string `params:"firmwareId" validate:"required,mongoid"`
}

// FirmwareInitRequest starts a presigned firmware upload. The client declares
// the target template, the firmware version, the file name, the byte size, and
// the SHA-256 it will upload; the platform validates the stored object against
// Size + SHA256 on finalize. Version is the target version the device receives
// in the OTA command (targetVersion) and reports after applying.
//
// SHA256 is the BASE64-encoded SHA-256 digest (the S3 `x-amz-checksum-sha256`
// encoding, NOT hex). The platform signs this exact value into the presigned
// PUT and the client MUST send the same value in the upload header, so any
// other encoding is rejected by the object store.
type FirmwareInitRequest struct {
	TargetTemplateID string `json:"targetTemplateId" validate:"required,min=1"`
	Version          string `json:"version" validate:"required,min=1"`
	Filename         string `json:"filename" validate:"required,min=1"`
	Size             int64  `json:"size" validate:"required,gt=0"`
	SHA256           string `json:"sha256" validate:"required,min=1"`
}

// FirmwareInitResponse returns the firmware artifact id and the short-TTL
// presigned PUT URL the client uses to upload the .bin directly to object
// storage.
type FirmwareInitResponse struct {
	FirmwareID string `json:"firmwareId"`
	UploadURL  string `json:"uploadUrl"`
}

// OTAFirmwareDownloadResponse returns a short-TTL presigned GET URL for the
// firmware artifact referenced by a plan, minted on demand when the operator
// clicks download. The bytes never traverse the Asset MS — the browser fetches
// them directly from object storage. ExpiresAt is when the URL stops working.
type OTAFirmwareDownloadResponse struct {
	URL       string    `json:"url"`
	Filename  string    `json:"filename"`
	ExpiresAt time.Time `json:"expiresAt"`
}

// OTARolloutConfigDTO controls the pacing and auto-abort of a plan rollout.
type OTARolloutConfigDTO struct {
	RatePerMinute     int `json:"ratePerMinute" validate:"required,gt=0"`
	AbortThresholdPct int `json:"abortThresholdPct" validate:"gte=0,lte=100"`
	AbortMinExecuted  int `json:"abortMinExecuted" validate:"gte=0"`
}

// OTAPlanCreateRequest creates an OTA plan migrating the selected assets from
// their source template to the target template carried by the firmware. The
// asset's template switches to the target only when the device reports updated.
// StartAt is OPTIONAL: absent/zero means "run immediately" — the platform
// stamps startAt = now so the start timer fires at once.
type OTAPlanCreateRequest struct {
	FirmwareID       string              `json:"firmwareId" validate:"required,min=1"`
	SourceTemplateID string              `json:"sourceTemplateId" validate:"required,min=1"`
	AssetIDs         []string            `json:"assetIds" validate:"required,min=1,dive,required"`
	StartAt          time.Time           `json:"startAt,omitzero"`
	MaxTime          time.Time           `json:"maxTime" validate:"required"`
	Name             string              `json:"name" validate:"required,min=1"`
	Description      string              `json:"description"`
	RolloutConfig    OTARolloutConfigDTO `json:"rolloutConfig" validate:"required"`
}

// OTAPlanCounters is the per-plan aggregate in responses (mirrors the entity's
// nested counters so the mapper copies it field-for-field).
type OTAPlanCounters struct {
	Total     int `json:"total"`
	Succeeded int `json:"succeeded"`
	Failed    int `json:"failed"`
	TimedOut  int `json:"timedOut"`
}

// OTAFirmwareInfo is the denormalized firmware artifact detail carried on the
// plan DETAIL response so the UI can show the file facts (version, name, size,
// checksum) without a second fetch. Populated by GetPlan from the referenced
// firmware; omitted on the list response.
type OTAFirmwareInfo struct {
	ID                string `json:"id"`
	Version           string `json:"version"`
	Filename          string `json:"filename"`
	Size              int64  `json:"size"`
	Checksum          string `json:"checksum"`
	ChecksumAlgorithm string `json:"checksumAlgorithm"`
	Status            string `json:"status"`
	TargetTemplateID  string `json:"targetTemplateId"`
}

// OTAPlanResponse is the read model for a plan in list/detail responses. Field
// names mirror the entity (ids as strings via ObjectIdToString) so the shared
// mapper populates it without hand-building. Firmware is present only on the
// detail response (GetPlan), nil on the list.
type OTAPlanResponse struct {
	ID               string           `json:"id"`
	OrgID            string           `json:"orgId"`
	Name             string           `json:"name"`
	Description      string           `json:"description"`
	FirmwareID       string           `json:"firmwareId"`
	SourceTemplateID string           `json:"sourceTemplateId"`
	TargetTemplateID string           `json:"targetTemplateId"`
	Status           string           `json:"status"`
	StartAt          time.Time        `json:"startAt"`
	MaxTime          time.Time        `json:"maxTime"`
	Counters         OTAPlanCounters  `json:"counters"`
	Firmware         *OTAFirmwareInfo `json:"firmware,omitempty"`
	Created          time.Time        `json:"created"`
	Updated          time.Time        `json:"updated"`
}

// OTAExecutionResponse is the read model for a single per-device execution.
type OTAExecutionResponse struct {
	ID         string    `json:"id"`
	PlanID     string    `json:"planId"`
	AssetID    string    `json:"assetId"`
	Protocol   string    `json:"protocol"`
	State      string    `json:"state"`
	Percentage int32     `json:"percentage"`
	Error      string    `json:"error,omitempty"`
	Attempts   int       `json:"attempts"`
	QueuedAt   time.Time `json:"queuedAt"`
	Updated    time.Time `json:"updated"`
}

// OTAStatusRequestDTO is the device-facing OTA status report body
// (POST /api/v1/ota/status?ds={dataSourceId} on the HTTP gateway). The gateway
// enriches it with the org resolved from the data-source auth — never from the
// body — and publishes the OTAStatusAdvisory for the Asset MS.
type OTAStatusRequestDTO struct {
	AssetUUID   string `json:"assetUUID" validate:"required,min=1"`
	ExecutionID string `json:"executionId" validate:"required,mongoid"`
	Status      string `json:"status" validate:"required,oneof=downloading downloaded verified updating updated failed"`
	Progress    int32  `json:"progress" validate:"gte=0,lte=100"`
	Error       string `json:"error,omitempty"`
	Message     string `json:"message,omitempty"`
}
