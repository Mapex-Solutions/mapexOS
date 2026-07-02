package repositories

import (
	"context"

	"assets/src/modules/assettemplates/domain/entities"

	reqCtx "github.com/Mapex-Solutions/mapexGoKit/microservices/common/context"
)

// FieldVocabularyRepository reads the curated, multi-tenant field
// vocabulary. It is intentionally read-only: management of org-specific
// entries is reserved for a future ticket.
type FieldVocabularyRepository interface {
	// ListFieldVocabulary returns every enabled vocabulary field visible
	// to the caller: the central platform standard (isSystem=true) unioned
	// with the caller's own org-scoped entries, derived from reqContext.
	ListFieldVocabulary(ctx context.Context, reqContext *reqCtx.RequestContext) ([]entities.FieldVocabulary, error)
}
