package middleware

import (
	"strings"

	"investory-management-backend/pkg/i18n"

	"github.com/gofiber/fiber/v3"
)

// Language detects the preferred language from the `lang` query param or
// the `Accept-Language` header and stores it in c.Locals("lang").
// Supported: "th", "en" (default).
func Language(c fiber.Ctx) error {
	lang := i18n.EN

	if q := c.Query("lang"); q != "" {
		lang = normalize(q)
	} else if h := c.Get("Accept-Language"); h != "" {
		lang = normalize(h)
	}

	c.Locals("lang", lang)
	return c.Next()
}

func normalize(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	// Accept-Language may be "th-TH,th;q=0.9,en;q=0.8" — take the first tag
	if idx := strings.IndexAny(s, ",;"); idx != -1 {
		s = s[:idx]
	}
	// Trim region subtag: "th-TH" → "th"
	if idx := strings.Index(s, "-"); idx != -1 {
		s = s[:idx]
	}
	if s == i18n.TH {
		return i18n.TH
	}
	return i18n.EN
}
