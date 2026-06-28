package middleware_test

import (
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"investory-management-backend/internal/middleware"
	"investory-management-backend/internal/service"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gofiber/fiber/v3"
)

const testSecret = "test-jwt-secret"

func makeToken(t *testing.T, secret string, expiry time.Duration) string {
	t.Helper()
	claims := service.Claims{
		UserID: 1,
		Email:  "test@example.com",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(secret))
	if err != nil {
		t.Fatalf("sign token: %v", err)
	}
	return signed
}

func newJWTApp() *fiber.App {
	app := fiber.New(fiber.Config{})
	app.Get("/", middleware.Protected, func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"user_id": c.Locals("userID"),
			"email":   c.Locals("email"),
		})
	})
	return app
}

func TestProtected_ValidToken(t *testing.T) {
	os.Setenv("JWT_SECRET", testSecret)
	defer os.Unsetenv("JWT_SECRET")

	token := makeToken(t, testSecret, time.Hour)
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := newJWTApp().Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("want 200 with valid token, got %d", resp.StatusCode)
	}
}

func TestProtected_MissingToken(t *testing.T) {
	os.Setenv("JWT_SECRET", testSecret)
	defer os.Unsetenv("JWT_SECRET")

	req := httptest.NewRequest("GET", "/", nil)
	resp, _ := newJWTApp().Test(req)
	if resp.StatusCode != 401 {
		t.Errorf("want 401 when no token, got %d", resp.StatusCode)
	}
}

func TestProtected_WrongSecret(t *testing.T) {
	os.Setenv("JWT_SECRET", testSecret)
	defer os.Unsetenv("JWT_SECRET")

	token := makeToken(t, "different-secret", time.Hour)
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, _ := newJWTApp().Test(req)
	if resp.StatusCode != 401 {
		t.Errorf("want 401 with wrong secret, got %d", resp.StatusCode)
	}
}

func TestProtected_ExpiredToken(t *testing.T) {
	os.Setenv("JWT_SECRET", testSecret)
	defer os.Unsetenv("JWT_SECRET")

	token := makeToken(t, testSecret, -time.Minute)
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, _ := newJWTApp().Test(req)
	if resp.StatusCode != 401 {
		t.Errorf("want 401 with expired token, got %d", resp.StatusCode)
	}
}

func TestProtected_MalformedToken(t *testing.T) {
	os.Setenv("JWT_SECRET", testSecret)
	defer os.Unsetenv("JWT_SECRET")

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer not.a.real.token")

	resp, _ := newJWTApp().Test(req)
	if resp.StatusCode != 401 {
		t.Errorf("want 401 with malformed token, got %d", resp.StatusCode)
	}
}

func TestProtected_NoBearerPrefix(t *testing.T) {
	os.Setenv("JWT_SECRET", testSecret)
	defer os.Unsetenv("JWT_SECRET")

	token := makeToken(t, testSecret, time.Hour)
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", token) // missing "Bearer " prefix

	resp, _ := newJWTApp().Test(req)
	if resp.StatusCode != 401 {
		t.Errorf("want 401 without Bearer prefix, got %d", resp.StatusCode)
	}
}
