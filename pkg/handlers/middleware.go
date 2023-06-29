package handlers

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	authHeader  = "Authorization"
	userCtx     = "userid" // щоб потім не забути
	listIdParam = "id"     // щоб потім не забути
	itemIdParam = "id"     // щоб потім не забути
)

// прослойка аутентифікації. записуємо айді користувача в контекст. таким чином в нас буде доступ до айді користувача всюди :)
func (r *Handler) userIdentity(ctx *gin.Context) { // задамо в якості обробника для /api
	header := ctx.GetHeader(authHeader)
	if header == "" { // перевіряємо чи не пустий токен
		newErrorResponse(ctx, http.StatusUnauthorized, "not authenticated")
		return
	}

	headersParts := strings.Split(header, " ") // частинки токена
	if len(headersParts) != 2 {
		newErrorResponse(ctx, http.StatusUnauthorized, "not correct token")
		return
	}

	// парсимо токен і записуємо значення айді користувача в контекст
	userId, err := r.services.Authorization.ParseToken(headersParts[1]) // отримуємо дані, передаючи частинку хедера, в якій немає "Bearer"
	if err != nil {
		newErrorResponse(ctx, http.StatusUnauthorized, err.Error())
		return
	}

	// додаємо айді користувача в контекст, щоб в наступних обробниках ми могли взаємодіяти з користувачем по цьому айді
	ctx.Set(userCtx, userId)
}

// щоб кожного разу не прописувати довжелезний код
func (r *Handler) getUserId(c *gin.Context) (int, error) {
	userid, ok := c.Get(userCtx)
	if !ok {
		newErrorResponse(c, http.StatusUnauthorized, "user id is not found")
		return 0, errors.New("user id is not found")
	}

	useridInt, ok := userid.(int)
	if !ok {
		newErrorResponse(c, http.StatusUnauthorized, "user id is not correct type")
		return 0, errors.New("user id is not correct type")
	}

	return useridInt, nil
}
