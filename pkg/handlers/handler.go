package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/zardan4/todo-app-gin/pkg/service"
)

// хендлери звертатимуться до сервісів. в чистій архітектурі БД handler -> service
type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *gin.Engine { // повертає екземпляр сервера
	router := gin.New()

	auth := router.Group("/auth") // роути в групі /auth
	{
		auth.POST("/signup", h.signUp) // auth/signup
		auth.POST("/signin", h.signIn)
	}

	// /api
	api := router.Group("/api", h.userIdentity)
	{
		// /api/lists
		lists := api.Group("/lists")
		{
			lists.GET("/", h.getLists)
			lists.POST("/", h.postList)
			// get lists by id
			lists.GET("/:id", h.getList)
			lists.PUT("/:id", h.updateList)
			lists.DELETE("/:id", h.deleteList)

			// /api/lists/items
			items := lists.Group("/:id/items")
			{
				items.GET("/", h.getItems)
				items.POST("/", h.postItem)
			}
		}

		// в цій групі винесені всі маніпуляції з елеменетами списку, в яких не потрібно мати доступ до айдішніка ліста
		items := api.Group("/items")
		{
			items.GET("/:id", h.getItem)
			items.PUT("/:id", h.updateItem)
			items.DELETE("/:id", h.deleteItem)
		}
	}

	return router
}
