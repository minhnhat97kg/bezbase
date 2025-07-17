package i18n

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

//go:embed locales/*.json
var localeFS embed.FS

// Bundle holds the i18n bundle
var Bundle *i18n.Bundle

// LocalizerKey is the context key for the localizer
type LocalizerKey string

const (
	LocalizerContextKey LocalizerKey = "localizer"
	DefaultLanguage     string       = "en"
)

// Initialize initializes the i18n bundle
func Initialize() error {
	Bundle = i18n.NewBundle(language.English)
	Bundle.RegisterUnmarshalFunc("json", json.Unmarshal)

	// Load English translations
	if err := loadLanguageFile("en"); err != nil {
		return fmt.Errorf("failed to load English translations: %w", err)
	}

	// Load Vietnamese translations
	if err := loadLanguageFile("vi"); err != nil {
		log.Printf("Warning: failed to load Vietnamese translations: %v", err)
	}

	return nil
}

// loadLanguageFile loads a language file from the embedded filesystem
func loadLanguageFile(lang string) error {
	filename := fmt.Sprintf("locales/%s.json", lang)
	data, err := localeFS.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read %s: %w", filename, err)
	}

	if _, err := Bundle.ParseMessageFileBytes(data, filename); err != nil {
		return fmt.Errorf("failed to parse %s: %w", filename, err)
	}

	return nil
}

// GetLocalizer returns a localizer for the given languages
func GetLocalizer(languages ...string) *i18n.Localizer {
	if len(languages) == 0 {
		languages = []string{DefaultLanguage}
	}
	return i18n.NewLocalizer(Bundle, languages...)
}

// GetLocalizerFromContext returns a localizer from the context
func GetLocalizerFromContext(ctx context.Context) *i18n.Localizer {
	if localizer, ok := ctx.Value(LocalizerContextKey).(*i18n.Localizer); ok {
		return localizer
	}
	return GetLocalizer(DefaultLanguage)
}

// T translates a message using the localizer from context
func T(ctx context.Context, messageID string, templateData ...interface{}) string {
	localizer := GetLocalizerFromContext(ctx)
	
	var data interface{}
	if len(templateData) > 0 {
		data = templateData[0]
	}

	message, err := localizer.Localize(&i18n.LocalizeConfig{
		MessageID:    messageID,
		TemplateData: data,
	})
	if err != nil {
		log.Printf("i18n translation error for messageID '%s': %v", messageID, err)
		return messageID // Return the message ID if translation fails
	}
	return message
}

// TWithDefault translates a message with a default fallback
func TWithDefault(ctx context.Context, messageID string, defaultMessage string, templateData ...interface{}) string {
	localizer := GetLocalizerFromContext(ctx)
	
	var data interface{}
	if len(templateData) > 0 {
		data = templateData[0]
	}

	message, err := localizer.Localize(&i18n.LocalizeConfig{
		MessageID:      messageID,
		DefaultMessage: &i18n.Message{ID: messageID, Other: defaultMessage},
		TemplateData:   data,
	})
	if err != nil {
		log.Printf("i18n translation error for messageID '%s': %v", messageID, err)
		return defaultMessage
	}
	return message
}

// ParseAcceptLanguage parses the Accept-Language header and returns preferred languages
func ParseAcceptLanguage(acceptLanguage string) []string {
	if acceptLanguage == "" {
		return []string{DefaultLanguage}
	}

	// Simple parsing of Accept-Language header
	// For production, consider using a more robust parser
	languages := strings.Split(acceptLanguage, ",")
	result := make([]string, 0, len(languages))
	
	for _, lang := range languages {
		// Remove quality values (e.g., "en-US;q=0.9" -> "en-US")
		lang = strings.TrimSpace(strings.Split(lang, ";")[0])
		if lang != "" {
			// Convert to lowercase and handle language codes
			lang = strings.ToLower(lang)
			result = append(result, lang)
		}
	}
	
	if len(result) == 0 {
		return []string{DefaultLanguage}
	}
	
	return result
}

// GetSupportedLanguages returns the list of supported languages
func GetSupportedLanguages() []string {
	return []string{"en", "vi"}
}

// IsLanguageSupported checks if a language is supported
func IsLanguageSupported(lang string) bool {
	supported := GetSupportedLanguages()
	for _, supportedLang := range supported {
		if lang == supportedLang {
			return true
		}
	}
	return false
}