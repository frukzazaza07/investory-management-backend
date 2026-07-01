package handler

import (
	"investory-management-backend/internal/service"
	"investory-management-backend/pkg/i18n"
	"investory-management-backend/pkg/response"

	"github.com/gofiber/fiber/v3"
)

// ListWebhooks godoc
// @Summary      List webhook subscriptions
// @Tags         webhooks
// @Security     BearerAuth
// @Produce      json
// @Success      200 {object} response.Response
// @Router       /api/v1/webhooks [get]
func ListWebhooks(c fiber.Ctx) error {
	lang := i18n.Lang(c)
	subs, err := service.ListWebhooks()
	if err != nil {
		return response.InternalError(c, err.Error())
	}
	return response.Success(c, i18n.T(lang, "webhook.list"), subs)
}

// GetWebhook godoc
// @Summary      Get webhook subscription
// @Tags         webhooks
// @Security     BearerAuth
// @Produce      json
// @Param        id path string true "Webhook ID"
// @Success      200 {object} response.Response
// @Router       /api/v1/webhooks/{id} [get]
func GetWebhook(c fiber.Ctx) error {
	lang := i18n.Lang(c)
	sub, err := service.GetWebhook(c.Params("id"))
	if err != nil {
		return response.NotFound(c, i18n.T(lang, "webhook.not_found"))
	}
	return response.Success(c, i18n.T(lang, "webhook.get"), sub)
}

type webhookRequest struct {
	Name     string   `json:"name"`
	URL      string   `json:"url"`
	Secret   string   `json:"secret"`
	Events   []string `json:"events"`
	IsActive bool     `json:"is_active"`
}

// CreateWebhook godoc
// @Summary      Create webhook subscription
// @Tags         webhooks
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        body body webhookRequest true "Webhook"
// @Success      201 {object} response.Response
// @Router       /api/v1/webhooks [post]
func CreateWebhook(c fiber.Ctx) error {
	lang := i18n.Lang(c)
	var req webhookRequest
	if err := c.Bind().JSON(&req); err != nil {
		return response.BadRequest(c, i18n.T(lang, "err.invalid_body"))
	}
	if req.Name == "" || req.URL == "" || req.Secret == "" {
		return response.BadRequest(c, i18n.T(lang, "err.name_url_secret_required"))
	}
	if len(req.Events) == 0 {
		req.Events = []string{"STOCK_UPDATED", "STOCK_LOW", "STOCK_OUT"}
	}
	sub, err := service.CreateWebhook(req.Name, req.URL, req.Secret, req.Events, req.IsActive)
	if err != nil {
		return response.InternalError(c, err.Error())
	}
	return response.Created(c, i18n.T(lang, "webhook.create"), fiber.Map{
		"id":        sub.ID,
		"name":      sub.Name,
		"url":       sub.URL,
		"events":    sub.Events,
		"is_active": sub.IsActive,
	})
}

// UpdateWebhook godoc
// @Summary      Update webhook subscription
// @Tags         webhooks
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id   path string         true "Webhook ID"
// @Param        body body webhookRequest true "Webhook"
// @Success      200 {object} response.Response
// @Router       /api/v1/webhooks/{id} [put]
func UpdateWebhook(c fiber.Ctx) error {
	lang := i18n.Lang(c)
	var req webhookRequest
	if err := c.Bind().JSON(&req); err != nil {
		return response.BadRequest(c, i18n.T(lang, "err.invalid_body"))
	}
	sub, err := service.UpdateWebhook(c.Params("id"), req.Name, req.URL, req.Secret, req.Events, req.IsActive)
	if err != nil {
		return response.NotFound(c, err.Error())
	}
	return response.Success(c, i18n.T(lang, "webhook.update"), sub)
}

// DeleteWebhook godoc
// @Summary      Delete webhook subscription
// @Tags         webhooks
// @Security     BearerAuth
// @Produce      json
// @Param        id path string true "Webhook ID"
// @Success      200 {object} response.Response
// @Router       /api/v1/webhooks/{id} [delete]
func DeleteWebhook(c fiber.Ctx) error {
	lang := i18n.Lang(c)
	if err := service.DeleteWebhook(c.Params("id")); err != nil {
		return response.NotFound(c, err.Error())
	}
	return response.Success(c, i18n.T(lang, "webhook.delete"), nil)
}

// TestWebhook godoc
// @Summary      Send a test event to webhook
// @Tags         webhooks
// @Security     BearerAuth
// @Produce      json
// @Param        id path string true "Webhook ID"
// @Success      200 {object} response.Response
// @Router       /api/v1/webhooks/{id}/test [post]
func TestWebhook(c fiber.Ctx) error {
	lang := i18n.Lang(c)
	err := service.TestWebhook(c.Params("id"))
	if err != nil {
		// Still return 200 — the test was attempted; caller can check delivery logs
		return response.Success(c, "test delivery attempted (delivery failed: "+err.Error()+")", nil)
	}
	return response.Success(c, i18n.T(lang, "webhook.test_sent"), nil)
}

// GetWebhookLogs godoc
// @Summary      Get delivery logs for a webhook
// @Tags         webhooks
// @Security     BearerAuth
// @Produce      json
// @Param        id path string true "Webhook ID"
// @Success      200 {object} response.Response
// @Router       /api/v1/webhooks/{id}/logs [get]
func GetWebhookLogs(c fiber.Ctx) error {
	lang := i18n.Lang(c)
	logs, err := service.GetWebhookLogs(c.Params("id"))
	if err != nil {
		return response.InternalError(c, err.Error())
	}
	return response.Success(c, i18n.T(lang, "webhook.logs"), logs)
}
