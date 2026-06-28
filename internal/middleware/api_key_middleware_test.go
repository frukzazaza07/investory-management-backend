package middleware_test

import (
	"net/http/httptest"
	"os"
	"testing"

	"investory-management-backend/internal/middleware"

	"github.com/gofiber/fiber/v3"
)

func newKeyApp() *fiber.App {
	app := fiber.New(fiber.Config{})
	app.Get("/", middleware.POSAPIKey, func(c fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})
	return app
}

func TestPOSAPIKey_ValidHeader(t *testing.T) {
	os.Setenv("POS_API_KEY", "secret-key")
	defer os.Unsetenv("POS_API_KEY")

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("X-API-Key", "secret-key")
	resp, err := newKeyApp().Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("want 200 with correct header key, got %d", resp.StatusCode)
	}
}

func TestPOSAPIKey_ValidQuery(t *testing.T) {
	os.Setenv("POS_API_KEY", "secret-key")
	defer os.Unsetenv("POS_API_KEY")

	req := httptest.NewRequest("GET", "/?api_key=secret-key", nil)
	resp, _ := newKeyApp().Test(req)
	if resp.StatusCode != 200 {
		t.Errorf("want 200 with correct query key, got %d", resp.StatusCode)
	}
}

func TestPOSAPIKey_WrongKey(t *testing.T) {
	os.Setenv("POS_API_KEY", "secret-key")
	defer os.Unsetenv("POS_API_KEY")

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("X-API-Key", "wrong-key")
	resp, _ := newKeyApp().Test(req)
	if resp.StatusCode != 401 {
		t.Errorf("want 401 with wrong key, got %d", resp.StatusCode)
	}
}

func TestPOSAPIKey_MissingKey(t *testing.T) {
	os.Setenv("POS_API_KEY", "secret-key")
	defer os.Unsetenv("POS_API_KEY")

	req := httptest.NewRequest("GET", "/", nil)
	resp, _ := newKeyApp().Test(req)
	if resp.StatusCode != 401 {
		t.Errorf("want 401 when no key provided, got %d", resp.StatusCode)
	}
}

func TestPOSAPIKey_EmptyEnvAlwaysRejects(t *testing.T) {
	os.Setenv("POS_API_KEY", "")

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("X-API-Key", "any-key")
	resp, _ := newKeyApp().Test(req)
	if resp.StatusCode != 401 {
		t.Errorf("want 401 when POS_API_KEY is unset, got %d", resp.StatusCode)
	}
}

func TestPOSAPIKey_HeaderTakesPrecedenceOverQuery(t *testing.T) {
	os.Setenv("POS_API_KEY", "correct")
	defer os.Unsetenv("POS_API_KEY")

	// header is wrong, query is correct — header is read first
	req := httptest.NewRequest("GET", "/?api_key=correct", nil)
	req.Header.Set("X-API-Key", "wrong")
	resp, _ := newKeyApp().Test(req)
	// middleware reads header first; header is wrong so it should reject
	if resp.StatusCode != 401 {
		t.Errorf("want 401 when header key is wrong even if query is correct, got %d", resp.StatusCode)
	}
}
