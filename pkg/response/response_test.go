package response_test

import (
	"encoding/json"
	"io"
	"net/http/httptest"
	"testing"

	"investory-management-backend/pkg/response"

	"github.com/gofiber/fiber/v3"
)

func newApp(handler fiber.Handler) *fiber.App {
	app := fiber.New(fiber.Config{})
	app.Get("/", handler)
	return app
}

func decode(t *testing.T, body io.Reader) response.Response {
	t.Helper()
	var r response.Response
	if err := json.NewDecoder(body).Decode(&r); err != nil {
		t.Fatalf("decode: %v", err)
	}
	return r
}

func TestSuccess(t *testing.T) {
	app := newApp(func(c fiber.Ctx) error {
		return response.Success(c, "all good", map[string]string{"k": "v"})
	})
	req := httptest.NewRequest("GET", "/", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("want 200, got %d", resp.StatusCode)
	}
	r := decode(t, resp.Body)
	if r.Status != "success" {
		t.Errorf("want status=success, got %q", r.Status)
	}
	if r.Message != "all good" {
		t.Errorf("want message=all good, got %q", r.Message)
	}
	if r.Data == nil {
		t.Error("want data to be set")
	}
}

func TestCreated(t *testing.T) {
	app := newApp(func(c fiber.Ctx) error {
		return response.Created(c, "item created", nil)
	})
	req := httptest.NewRequest("GET", "/", nil)
	resp, _ := app.Test(req)
	if resp.StatusCode != 201 {
		t.Errorf("want 201, got %d", resp.StatusCode)
	}
	r := decode(t, resp.Body)
	if r.Status != "success" {
		t.Errorf("want status=success, got %q", r.Status)
	}
}

func TestBadRequest(t *testing.T) {
	app := newApp(func(c fiber.Ctx) error {
		return response.BadRequest(c, "bad input")
	})
	req := httptest.NewRequest("GET", "/", nil)
	resp, _ := app.Test(req)
	if resp.StatusCode != 400 {
		t.Errorf("want 400, got %d", resp.StatusCode)
	}
	r := decode(t, resp.Body)
	if r.Status != "error" {
		t.Errorf("want status=error, got %q", r.Status)
	}
	if r.Message != "bad input" {
		t.Errorf("want message=bad input, got %q", r.Message)
	}
}

func TestUnauthorized(t *testing.T) {
	app := newApp(func(c fiber.Ctx) error {
		return response.Unauthorized(c, "forbidden")
	})
	req := httptest.NewRequest("GET", "/", nil)
	resp, _ := app.Test(req)
	if resp.StatusCode != 401 {
		t.Errorf("want 401, got %d", resp.StatusCode)
	}
	r := decode(t, resp.Body)
	if r.Status != "error" {
		t.Errorf("want status=error, got %q", r.Status)
	}
}

func TestNotFound(t *testing.T) {
	app := newApp(func(c fiber.Ctx) error {
		return response.NotFound(c, "not found")
	})
	req := httptest.NewRequest("GET", "/", nil)
	resp, _ := app.Test(req)
	if resp.StatusCode != 404 {
		t.Errorf("want 404, got %d", resp.StatusCode)
	}
	r := decode(t, resp.Body)
	if r.Status != "error" {
		t.Errorf("want status=error, got %q", r.Status)
	}
}

func TestForbidden(t *testing.T) {
	app := newApp(func(c fiber.Ctx) error {
		return response.Forbidden(c, "no access")
	})
	req := httptest.NewRequest("GET", "/", nil)
	resp, _ := app.Test(req)
	if resp.StatusCode != 403 {
		t.Errorf("want 403, got %d", resp.StatusCode)
	}
}

func TestInternalError(t *testing.T) {
	app := newApp(func(c fiber.Ctx) error {
		return response.InternalError(c, "boom")
	})
	req := httptest.NewRequest("GET", "/", nil)
	resp, _ := app.Test(req)
	if resp.StatusCode != 500 {
		t.Errorf("want 500, got %d", resp.StatusCode)
	}
	r := decode(t, resp.Body)
	if r.Status != "error" {
		t.Errorf("want status=error, got %q", r.Status)
	}
}

func TestSuccessNilData(t *testing.T) {
	app := newApp(func(c fiber.Ctx) error {
		return response.Success(c, "empty", nil)
	})
	req := httptest.NewRequest("GET", "/", nil)
	resp, _ := app.Test(req)
	if resp.StatusCode != 200 {
		t.Errorf("want 200, got %d", resp.StatusCode)
	}
}
