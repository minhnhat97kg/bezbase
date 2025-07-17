package i18n

import (
	"context"
	"fmt"
	"log"
)

// TestI18nFunctionality tests the i18n functionality
func TestI18nFunctionality() {
	fmt.Println("Testing i18n functionality...")
	
	// Initialize i18n
	if err := Initialize(); err != nil {
		log.Fatal("Failed to initialize i18n:", err)
	}
	
	// Test English translations
	fmt.Println("\n--- English Translations ---")
	enLocalizer := GetLocalizer("en")
	ctx := context.WithValue(context.Background(), LocalizerContextKey, enLocalizer)
	
	t := NewTranslator(ctx)
	fmt.Printf("Invalid request body: %s\n", t.InvalidRequestBody())
	fmt.Printf("Invalid credentials: %s\n", t.InvalidCredentials())
	fmt.Printf("User not found: %s\n", t.UserNotFound())
	fmt.Printf("Username already taken: %s\n", t.UsernameAlreadyTaken())
	fmt.Printf("Password changed successfully: %s\n", t.PasswordChangedSuccessfully())
	fmt.Printf("Server status: %s\n", t.Status("server_running"))
	
	// Test Vietnamese translations
	fmt.Println("\n--- Vietnamese Translations ---")
	viLocalizer := GetLocalizer("vi")
	ctx = context.WithValue(context.Background(), LocalizerContextKey, viLocalizer)
	
	t = NewTranslator(ctx)
	fmt.Printf("Invalid request body: %s\n", t.InvalidRequestBody())
	fmt.Printf("Invalid credentials: %s\n", t.InvalidCredentials())
	fmt.Printf("User not found: %s\n", t.UserNotFound())
	fmt.Printf("Username already taken: %s\n", t.UsernameAlreadyTaken())
	fmt.Printf("Password changed successfully: %s\n", t.PasswordChangedSuccessfully())
	fmt.Printf("Server status: %s\n", t.Status("server_running"))
	
	// Test error with template data
	fmt.Println("\n--- Template Data Test ---")
	enCtx := context.WithValue(context.Background(), LocalizerContextKey, GetLocalizer("en"))
	enT := NewTranslator(enCtx)
	fmt.Printf("Insufficient permissions (EN): %s\n", enT.InsufficientPermissions("admin"))
	
	viCtx := context.WithValue(context.Background(), LocalizerContextKey, GetLocalizer("vi"))
	viT := NewTranslator(viCtx)
	fmt.Printf("Insufficient permissions (VI): %s\n", viT.InsufficientPermissions("admin"))
	
	// Test Accept-Language parsing
	fmt.Println("\n--- Accept-Language Parsing Test ---")
	languages := ParseAcceptLanguage("en-US,en;q=0.9,vi;q=0.8,fr;q=0.7")
	fmt.Printf("Parsed languages: %v\n", languages)
	
	// Test supported languages
	fmt.Println("\n--- Supported Languages ---")
	fmt.Printf("Supported languages: %v\n", GetSupportedLanguages())
	fmt.Printf("Is 'en' supported: %v\n", IsLanguageSupported("en"))
	fmt.Printf("Is 'vi' supported: %v\n", IsLanguageSupported("vi"))
	fmt.Printf("Is 'fr' supported: %v\n", IsLanguageSupported("fr"))
	
	fmt.Println("\ni18n functionality test completed successfully!")
}