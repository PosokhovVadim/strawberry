package handler

import (
	"context"
	"net/http"
	"strconv"

	"github.com/PosokhovVadim/stawberry/internal/domain/entity"
	"github.com/PosokhovVadim/stawberry/internal/handler/dto"
	"github.com/gin-gonic/gin"
)

type AuthService interface {
	Register(ctx context.Context, register entity.Register) (entity.TokenPair, error)
	Login(ctx context.Context, creds entity.Credentials) (entity.TokenPair, error)
	Refresh(ctx context.Context, refresh entity.RefreshSession) (entity.TokenPair, error)
}

type AuthHandler struct {
	authService AuthService
}

func NewAuthHandler(authService AuthService) AuthHandler {
	return AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var dto dto.RegisterReq

	if err := c.ShouldBind(&dto); err != nil {
		handleBindError(c, err)
		return
	}

	tokens, err := h.authService.Register(c.Request.Context(), dto.ConvertToSvc())
	if err != nil {
		handleAuthError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": tokens})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var dto dto.LoginReq

	if err := c.ShouldBind(&dto); err != nil {
		handleBindError(c, err)
		return
	}

	tokens, err := h.authService.Login(c.Request.Context(), dto.ConvertToSvc())
	if err != nil {
		handleAuthError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": tokens})
}

func (h *AuthHandler) RefreshTokens(c *gin.Context) {
	var dto dto.RefreshTokensReq

	if err := c.ShouldBindHeader(&dto); err != nil {
		if err := c.ShouldBindJSON(&dto); err != nil {
			handleBindError(c, err)
			return
		}
	}

	tokens, err := h.authService.Refresh(c.Request.Context(), dto.ConvertToEntity())
	if err != nil {
		handleAuthError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": tokens})
}

func (h *AuthHandler) Me(c *gin.Context) {
	userID, _ := strconv.ParseUint(c.GetString("userID"), 10, 0)

	c.JSON(http.StatusOK, map[string]any{"user_id": userID})
}
