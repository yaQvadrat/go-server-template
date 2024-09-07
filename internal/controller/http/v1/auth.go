package httpapi

import (
	"app/internal/service"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
)

type authRoutes struct {
	accountService service.Account
}

func newAuthRoutes(accountService service.Account) *authRoutes {
	return &authRoutes{accountService: accountService}
}

type signUpInput struct {
	Username string `json:"username" validate:"required,min=2,max=64"`
	Password string `json:"password" validate:"required,min=8,max=32"`
	Role     string `json:"role" validate:"required"`
}

func (r *authRoutes) signUp(c echo.Context) error {
	var input signUpInput
	if err := c.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	if err := c.Validate(input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid field value")
	}

	id, err := r.accountService.CreateAccount(c.Request().Context(), service.CreateAccountInput{
		Username: input.Username,
		Password: input.Password,
		Role:     input.Role,
	})

	if err != nil {
		if errors.Is(err, service.ErrAccountAlreadyExists) {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		return echo.NewHTTPError(http.StatusInternalServerError, ErrInternalServer.Error())
	}

	return c.JSON(http.StatusCreated, map[string]int{
		"id": id,
	})
}

type signInInput struct {
	Username string `json:"username" validate:"required,min=2,max=64"`
	Password string `json:"password" validate:"required,min=8,max=32"`
}

func (r *authRoutes) signIn(c echo.Context) error {
	var input signInInput
	if err := c.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	if err := c.Validate(input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid field value")
	}

	token, err := r.accountService.GenerateToken(c.Request().Context(), service.GenerateTokenInput{
		Username: input.Username,
		Password: input.Password,
	})

	if err != nil {
		if errors.Is(err, service.ErrAccountNotFound) {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		return echo.NewHTTPError(http.StatusInternalServerError, ErrInternalServer.Error())
	}

	return c.JSON(http.StatusOK, map[string]string{
		"token": token,
	})
}
