package todo

type User struct {
	Id       int    `json:"-"`
	Name     string `json:"name" binding:"required"`     // binding required валідує наявність поля в тілі запита. реалізується за допомогою GIN
	Username string `json:"username" binding:"required"` // binding required валідує наявність поля в тілі запита. реалізується за допомогою GIN
	Password string `json:"password" binding:"required"` // binding required валідує наявність поля в тілі запита. реалізується за допомогою GIN
}
