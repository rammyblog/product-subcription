package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/go-chi/render"
	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/rammyblog/go-product-subscriptions/models"
	"github.com/rammyblog/go-product-subscriptions/response"
)

func CreateJwtToken(user models.User) (string, error) {
	claims := &jwt.MapClaims{
		"expiresAt": 15000,
		"id":        user.ID,
	}
	secret := os.Getenv("JWT_SECRET")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func validateJWT(tokenString string) (*jwt.Token, error) {
	secret := os.Getenv("JWT_SECRET")
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don"t forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(secret), nil
	})
}

func GetUserIdFromToken(bearerToken string) (uint, error) {
	tokenString := strings.Split(bearerToken, " ")[1]
	token, err := validateJWT(tokenString)
	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if ok && token.Valid {
		uid, err := strconv.ParseUint(fmt.Sprintf("%.0f", claims["id"]), 10, 32)
		if err != nil {
			return 0, err
		}
		return uint(uid), nil
	}
	return 0, nil

}

func JwtAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bearerToken := r.Header.Get("Authorization")
		if bearerToken == "" {
			render.Render(w, r, response.ErrInvalidRequest(fmt.Errorf("authorization header is required")))
			return
		}
		_, err := GetUserIdFromToken(bearerToken)
		if err != nil {
			render.Render(w, r, response.ErrInvalidRequest(fmt.Errorf("unauthorized")))
			return
		}
		next.ServeHTTP(w, r)
	})
}
