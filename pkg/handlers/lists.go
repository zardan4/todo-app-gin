package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/zardan4/todo-app-gin"
)

type getAllListsResponse struct { // структура відповіді
	Data []todo.TodoList
}

func (h *Handler) getLists(c *gin.Context) {
	userid, err := h.getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	lists, err := h.services.TodoList.GetAll(userid) // відправляємо запит на створення ліста в сервіс
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, getAllListsResponse{
		Data: lists,
	})
}

func (h *Handler) postList(c *gin.Context) {
	userid, err := h.getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	var input todo.TodoList
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.services.TodoList.Create(userid, input) // відправляємо запит на створення ліста в сервіс
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

func (h *Handler) getList(c *gin.Context) {
	userid, err := h.getUserId(c) // отримуємо айді
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	listid, err := strconv.Atoi(c.Param(listIdParam))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	list, err := h.services.TodoList.GetById(userid, listid) // прокидуємо айді далі
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, list)
}

func (h *Handler) updateList(c *gin.Context) {
	userid, err := h.getUserId(c) // отримуємо айді
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	listid, err := strconv.Atoi(c.Param(listIdParam)) // отримуємо айді
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	var input todo.UpdateTodoListInput // ліст
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	err = h.services.TodoList.Update(userid, listid, input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "ok",
	})
}

func (h *Handler) deleteList(c *gin.Context) {
	userid, err := h.getUserId(c) // отримуємо айді
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	listid, err := strconv.Atoi(c.Param(listIdParam)) // отримуємо айді
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	err = h.services.TodoList.Delete(userid, listid)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "ok",
	})
}
