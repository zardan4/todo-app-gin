package service

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/zardan4/todo-app-gin"
	"github.com/zardan4/todo-app-gin/pkg/repository"
)

const (
	_salt       = "aasdasdasdasdasdasd"
	_signingKey = "NIv7nBV8svhX8OVSNVUYSZBSV87nxvyBSVSO7v"
	_tokenTTL   = 12 * time.Hour
)

type tokenClaims struct {
	jwt.StandardClaims     // дефолтні клеймси
	UserId             int `json:"user_id"`
}

type AuthService struct {
	repo repository.Authorization
}

func NewAuthService(repo repository.Authorization) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) CreateUser(user todo.User) (int, error) {
	user.Password = hashPassword(user.Password) // фігачимо хеш

	// на цьому моменті ми реалізували всю бізнес-логіку(хешування), тепер передаємо ще на шар нижче, в repository, для роботи з БД
	return s.repo.CreateUser(user)
}

// логіка створення JWT-токена
func (s *AuthService) GenerateToken(username string, password string) (string, error) {
	user, err := s.repo.GetUser(username, hashPassword(password))
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{ // claims - об'єкт з набором різних полів
			ExpiresAt: time.Now().Add(_tokenTTL).Unix(), // встановимо час дії токену на 12 годин
			IssuedAt:  time.Now().Unix(),                // час, коли токен був згенерений
		},
		user.Id, // додатковий клейм з айді користувача
	}) // token generation

	return token.SignedString([]byte(_signingKey)) // передаємо підпис, за допомогою якого буде розшифровуватися отримання токена
}

// логіка парсингу JWT-токена. повертає ID користувача при корректному токені
func (s *AuthService) ParseToken(token string) (int, error) {
	// парсимо токен
	tokenResp, err := jwt.ParseWithClaims(token, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok { // перевіряємо на метод підпису
			return nil, errors.New("invalid signing method")
		}

		return []byte(_signingKey), nil // повертаємо ключ, яким підписується токен
	})
	if err != nil {
		return 0, err
	}

	claims, ok := tokenResp.Claims.(*tokenClaims) // отримуємо claims(інфа в токені)
	if !ok {
		return 0, errors.New("error while getting claims")
	}

	return claims.UserId, nil
}

func hashPassword(password string) string {
	// хешуємо пароль, використовуючи алгоритм sha1
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(_salt))) // add salt to hash
}
