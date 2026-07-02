package services

import (
	"assets/src/modules/assettemplates/application/constants"
	"assets/src/modules/assettemplates/application/dtos"
	domainconstants "assets/src/modules/assettemplates/domain/constants"
	"assets/src/modules/assettemplates/domain/entities"
)

// groupVocabularyByCategory buckets vocabulary entries by their category
// key, preserving the repository order within each bucket.
func (s *AssetTemplateService) groupVocabularyByCategory(fields []entities.FieldVocabulary) map[string][]entities.FieldVocabulary {
	byCategory := make(map[string][]entities.FieldVocabulary)
	for i := range fields {
		f := fields[i]
		byCategory[f.Category] = append(byCategory[f.Category], f)
	}
	return byCategory
}

// buildVocabularyGroups assembles the response groups in the fixed category
// order, omitting categories with no fields, and resolves each field hint
// plus the group label to the requested language (en-US fallback).
func (s *AssetTemplateService) buildVocabularyGroups(byCategory map[string][]entities.FieldVocabulary, lang string) []dtos.FieldVocabularyGroup {
	groups := make([]dtos.FieldVocabularyGroup, 0, len(domainconstants.FieldVocabularyCategoryOrder))
	for _, category := range domainconstants.FieldVocabularyCategoryOrder {
		items := byCategory[category]
		if len(items) == 0 {
			continue
		}
		groups = append(groups, dtos.FieldVocabularyGroup{
			Category: category,
			Label:    s.resolveCategoryLabel(category, lang),
			Fields:   s.mapVocabularyFields(items, lang),
		})
	}
	return groups
}

// mapVocabularyFields converts entities to response fields, resolving each
// hint to the requested language.
func (s *AssetTemplateService) mapVocabularyFields(items []entities.FieldVocabulary, lang string) []dtos.FieldVocabularyField {
	out := make([]dtos.FieldVocabularyField, len(items))
	for i := range items {
		f := items[i]
		out[i] = dtos.FieldVocabularyField{
			Value: f.Value,
			Hint:  s.resolveLocalized(f.Hint, lang),
			Type:  f.Type,
			Unit:  f.Unit,
		}
	}
	return out
}

// resolveCategoryLabel returns the localized label for a category, falling
// back to en-US, then to the raw category key when no label is configured.
func (s *AssetTemplateService) resolveCategoryLabel(category string, lang string) string {
	labels, ok := constants.FieldVocabularyCategoryLabels[category]
	if !ok {
		return category
	}
	if label := s.resolveLocalized(labels, lang); label != "" {
		return label
	}
	return category
}

// resolveLocalized picks the value for the requested locale, falling back to
// en-US when the locale is empty or missing. Returns "" when neither exists.
func (s *AssetTemplateService) resolveLocalized(values map[string]string, lang string) string {
	if values == nil {
		return ""
	}
	if lang != "" {
		if v, ok := values[lang]; ok {
			return v
		}
	}
	return values[constants.DefaultLocale]
}
