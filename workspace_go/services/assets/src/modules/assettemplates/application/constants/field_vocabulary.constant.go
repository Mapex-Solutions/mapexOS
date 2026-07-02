package constants

// DefaultLocale is the fallback locale used when the requested language is
// absent from a hint or category label map.
const DefaultLocale = "en-US"

// FieldVocabularyCategoryLabels holds the bilingual display label for each
// category, keyed by locale. Used to resolve the group label to the
// requested language (en-US fallback).
var FieldVocabularyCategoryLabels = map[string]map[string]string{
	"climate":             {"en-US": "Climate", "pt-BR": "Clima"},
	"air-quality":         {"en-US": "Air Quality", "pt-BR": "Qualidade do Ar"},
	"presence-counting":   {"en-US": "Presence & Counting", "pt-BR": "Presença & Contagem"},
	"agriculture-liquids": {"en-US": "Agriculture & Liquids", "pt-BR": "Agricultura & Líquidos"},
	"electrical":          {"en-US": "Electrical", "pt-BR": "Elétrico"},
	"mechanical-position": {"en-US": "Mechanical & Position", "pt-BR": "Mecânico & Posição"},
	"weather":             {"en-US": "Weather", "pt-BR": "Clima Externo"},
	"device-health":       {"en-US": "Device Health", "pt-BR": "Saúde do Dispositivo"},
}
