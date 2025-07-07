package api

import (
	"auth_service/internal/auth"
	"auth_service/internal/storage"
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type Handler struct {
	service *auth.AuthService
	store   *storage.Store
	logger  *logrus.Logger
}

func NewHandler(service *auth.AuthService, store *storage.Store, logger *logrus.Logger) *Handler {
	return &Handler{
		service: service,
		store:   store,
		logger:  logger,
	}
}

// GET /tokens?user_id=
func (h *Handler) GetTokens(c echo.Context) error {
	userIDString := c.QueryParam("user_id")
	if userIDString == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "need user id"})
	}

	userID, err := uuid.Parse(userIDString)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid user id"})
	}

	tokenPair, err := h.service.CreateTokenPair(c.Request().Context(), userID, c.Request().UserAgent(), c.RealIP())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to create tokens"})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"access_token":  tokenPair.AccessToken,
		"refresh_token": tokenPair.RefreshToken,
	})
}

// POST /refresh
func (h *Handler) RefreshTokens(c echo.Context) error {
	var req struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}

	// Биндим тело запроса в нашу структуру.
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	userID, err := h.service.ParseAccessToken(req.AccessToken)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid access token"})
	}

	sessions, err := h.store.FindRefreshSessions(c.Request().Context(), userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
	}
	if len(sessions) == 0 {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "No active sessions found or refresh token is invalid"})
	}

	var currentSession *storage.RefreshSession
	var found bool
	for _, s := range sessions {
		if err := bcrypt.CompareHashAndPassword([]byte(s.TokenHash), []byte(req.RefreshToken)); err == nil {
			currentSession = s
			found = true
			break
		}
	}

	if !found {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid refresh token"})
	}

	if currentSession.UserAgent != c.Request().UserAgent() {
		h.logger.Warnf("ВНИМАНИЕ: User-Agent не подходит для пользователя %s. Revoking all sessions.", userID)
		h.store.RevokeUserSessions(c.Request().Context(), userID)
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "User-Agent mismatch. All sessions have been revoked."})
	}
	if currentSession.IPAddress != c.RealIP() {
		h.logger.Infof("IP адреса не совпадают для пользователя %s. Старый: %s, Новый: %s.", userID, currentSession.IPAddress, c.RealIP())
		go h.sendIPMismatchWebhook(userID, currentSession.IPAddress, c.RealIP())
	}

	if err := h.store.DeleteRefreshSession(c.Request().Context(), currentSession.ID); err != nil {
		h.logger.Errorf("Невозможно удалить прошлую сессию %d: %v", currentSession.ID, err)
	}

	newTokenPair, err := h.service.CreateTokenPair(c.Request().Context(), userID, c.Request().UserAgent(), c.RealIP())
	if err != nil {
		h.logger.Errorf("Ошибка создания новых токенов для пользователя %s: %v", userID, err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create new tokens"})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"access_token":  newTokenPair.AccessToken,
		"refresh_token": newTokenPair.RefreshToken,
	})
}

// GET /me
func (h *Handler) GetMe(c echo.Context) error {
	userID, ok := c.Get("userID").(uuid.UUID)
	if !ok {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "User ID not found in context"})
	}
	return c.JSON(http.StatusOK, map[string]uuid.UUID{"user id": userID})
}

// POST /logout
func (h *Handler) Logout(c echo.Context) error {
	userID, ok := c.Get("userID").(uuid.UUID)
	if !ok {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "User ID not found in context"})
	}

	if err := h.store.RevokeUserSessions(c.Request().Context(), userID); err != nil {
		h.logger.Errorf("Ошибка отзыва сессии для пользователя %s: %v", userID, err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to logout"})
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *Handler) sendIPMismatchWebhook(userID uuid.UUID, oldIP, newIP string) {
	payload := map[string]string{
		"user_id": userID.String(),
		"old_ip":  oldIP,
		"new_ip":  newIP,
	}
	body, _ := json.Marshal(payload)

	resp, err := http.Post("http://localhost:8081/webhook", "application/json", bytes.NewBuffer(body))
	if err != nil {
		h.logger.Errorf("Ошибка отправки webhook для пользователя %s: %v", userID, err)
		return
	}
	defer resp.Body.Close()
	h.logger.Infof("Webhook для пользователя %s отправлен, статус: %s", userID, resp.Status)
}
