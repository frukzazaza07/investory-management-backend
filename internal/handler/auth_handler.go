package handler

import (
	"investory-management-backend/internal/service"
	"investory-management-backend/pkg/i18n"
	"investory-management-backend/pkg/response"

	"github.com/gofiber/fiber/v3"
)

type authRequest struct {
	Email    string `json:"email" example:"test@test.com"`
	Password string `json:"password" example:"123456"`
}

// Register godoc
// @Summary      Register new user
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body body authRequest true "Register body"
// @Success      201  {object}  response.Response
// @Failure      400  {object}  response.Response
// @Failure      500  {object}  response.Response
// @Router       /auth/register [post]
func Register(c fiber.Ctx) error {
	lang := i18n.Lang(c)
	var req authRequest
	if err := c.Bind().JSON(&req); err != nil {
		return response.BadRequest(c, i18n.T(lang, "err.invalid_request"))
	}
	if err := service.Register(req.Email, req.Password); err != nil {
		return response.InternalError(c, err.Error())
	}
	return response.Created(c, i18n.T(lang, "auth.register_success"), nil)
}

// Login godoc
// @Summary      Login
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body body authRequest true "Login body"
// @Success      200  {object}  response.Response
// @Failure      400  {object}  response.Response
// @Failure      401  {object}  response.Response
// @Router       /auth/login [post]
func Login(c fiber.Ctx) error {
	lang := i18n.Lang(c)
	var req authRequest
	if err := c.Bind().JSON(&req); err != nil {
		return response.BadRequest(c, i18n.T(lang, "err.invalid_request"))
	}
	token, err := service.Login(req.Email, req.Password)
	if err != nil {
		return response.Unauthorized(c, err.Error())
	}
	return response.Success(c, i18n.T(lang, "auth.login_success"), fiber.Map{"token": token})
}
