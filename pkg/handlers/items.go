package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/zardan4/todo-app-gin"
)

type getAllItemsResponse struct { // структура відповіді
	Data []todo.TodoItem
}

func (h *Handler) getItems(c *gin.Context) {
	userid, err := h.getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	listid, err := strconv.Atoi(c.Param(itemIdParam)) // отримуємо айді
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	items, err := h.services.TodoItem.GetAll(userid, listid) // відправляємо запит на створення ліста в сервіс
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, getAllItemsResponse{
		Data: items,
	})
}

func (h *Handler) postItem(c *gin.Context) {
	userid, err := h.getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	listid, err := strconv.Atoi(c.Param(itemIdParam)) // отримуємо айді
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	var input todo.TodoItem
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.services.TodoItem.Create(userid, listid, input) // відправляємо запит на створення айтема в сервіс
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

func (h *Handler) getItem(c *gin.Context) {
	userid, err := h.getUserId(c) // отримуємо айді
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	itemid, err := strconv.Atoi(c.Param(itemIdParam)) // отримуємо айді
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	item, err := h.services.TodoItem.GetById(userid, itemid) // прокидуємо айді далі
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, item)
}

func (h *Handler) updateItem(c *gin.Context) {

	//
	userid, err := h.getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	//
	itemid, err := strconv.Atoi(c.Param(itemIdParam))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	//
	var input todo.UpdateTodoItemInput
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	//
	err = h.services.TodoItem.Update(userid, itemid, input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	//
	c.JSON(http.StatusOK, statusResponse{
		Status: "ok",
	})
}

func (h *Handler) deleteItem(c *gin.Context) {
	userid, err := h.getUserId(c) // отримуємо айді
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	itemid, err := strconv.Atoi(c.Param(itemIdParam)) // отримуємо айді
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	err = h.services.TodoItem.Delete(userid, itemid) // прокидуємо айді далі
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "ok",
	})
}
