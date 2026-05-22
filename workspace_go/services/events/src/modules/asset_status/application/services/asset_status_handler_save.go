package services

import (
	"events/src/modules/asset_status/domain/entities"

	natsModel "github.com/Mapex-Solutions/mapexGoKit/infrastructure/nats"
)

// parseAssetStatusBatch walks the NATS batch sequentially. For each
// message, sets tenant context, decodes the wire payload, and translates
// it into a ClickHouse row entity. Invalid messages are Rejected (-> DLQ)
// in place. Returns the entities ready for bulk insert and the matching
// msg pointers so the caller can Ack/Nack them as a group.
func (s *AssetStatusService) parseAssetStatusBatch(messages []*natsModel.Message) ([]*entities.AssetStatusEvent, []*natsModel.Message) {
	entitiesBatch := make([]*entities.AssetStatusEvent, 0, len(messages))
	validMessages := make([]*natsModel.Message, 0, len(messages))

	for _, msg := range messages {
		setTenantContext(msg)

		payload, err := parsePersistencePayload(msg.Data)
		if err != nil {
			_ = msg.Reject(err.Error())
			continue
		}

		msg.OrgId = payload.OrgId
		msg.PathKey = payload.PathKey
		msg.EventTrackerId = payload.EventId

		entity, err := payloadToEntity(payload)
		if err != nil {
			_ = msg.Reject(err.Error())
			continue
		}

		entitiesBatch = append(entitiesBatch, entity)
		validMessages = append(validMessages, msg)
	}
	return entitiesBatch, validMessages
}

// nackBatch fans out msg.Nack to every valid message after a bulk-insert
// failure so NATS redelivers them per the retry policy.
func (s *AssetStatusService) nackBatch(messages []*natsModel.Message, err error) {
	for _, msg := range messages {
		_ = msg.Nack(err)
	}
}

// ackBatch fans out msg.Ack to every successfully persisted message.
func (s *AssetStatusService) ackBatch(messages []*natsModel.Message) {
	for _, msg := range messages {
		_ = msg.Ack()
	}
}
