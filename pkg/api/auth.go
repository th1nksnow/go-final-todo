package api

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	jwtSecretKey    = "todo_jwt_secret" // vulnerability
	tokenCookieName = "token"
	tokenDuration   = 8 * time.Hour
)

type signinRequest struct {
	Password string `json:"password"`
}

type signinResponse struct {
	Token string `json:"token,omitempty"`
	Error string `json:"error,omitempty"`
}

type Claims struct {
	PasswordHash string `json:"password_hash"`
	jwt.RegisteredClaims
}

func auth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !cfgApi.authEnabled {
			next(w, r)
			return
		}

		password := cfgApi.password

		var tokenString string
		cookie, err := r.Cookie(tokenCookieName)
		if err == nil {
			tokenString = cookie.Value
		} else {
			// Пробуем получить токен из заголовка Authorization
			authHeader := r.Header.Get("Authorization")
			if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
				tokenString = strings.TrimPrefix(authHeader, "Bearer ")
			}
		}

		valid, err := validateToken(tokenString, password)
		if !valid || err != nil {
			http.Error(w, "authentication required", http.StatusUnauthorized)
			return
		}

		next(w, r)
	})
}

func generateToken(password string) (string, error) {
	hash := sha256.Sum256([]byte(password))
	passwordHash := hex.EncodeToString(hash[:])

	claims := &Claims{
		PasswordHash: passwordHash,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Подписываем токен
	signedToken, err := token.SignedString([]byte(jwtSecretKey))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %v", err)
	}

	return signedToken, nil
}

func validateToken(tokenString, password string) (bool, error) {
	if tokenString == "" {
		return false, errors.New("token is empty")
	}

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (any, error) {
		// Проверяем метод подписи
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtSecretKey), nil
	})

	if err != nil {
		return false, fmt.Errorf("failed to parse token: %v", err)
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		// Сравниваем хэш пароля из токена с текущим
		currentHash := sha256.Sum256([]byte(password))
		currentHashStr := hex.EncodeToString(currentHash[:])

		if claims.PasswordHash == currentHashStr {
			return true, nil
		}
		return false, errors.New("password hash mismatch")
	}

	return false, errors.New("invalid token claims")
}

func signinHandler(w http.ResponseWriter, r *http.Request) {
	var response signinResponse

	if r.Method != http.MethodPost {
		response.Error = "method not allowed"
		writeJSON(w, response, http.StatusMethodNotAllowed)
		return
	}

	var req signinRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error = fmt.Sprintf("failed to decode JSON: %v", err)
		writeJSON(w, response, http.StatusBadRequest)
		return
	}

	if !cfgApi.authEnabled {
		response.Error = "authentication is not configured"
		writeJSON(w, response, http.StatusBadRequest)
		return
	}

	expectedPassword := cfgApi.password
	if req.Password != expectedPassword {
		response.Error = "invalid password"
		writeJSON(w, response, http.StatusUnauthorized)
		return
	}

	token, err := generateToken(expectedPassword)
	if err != nil {
		response.Error = fmt.Sprintf("failed to generate token: %v", err)
		writeJSON(w, response, http.StatusInternalServerError)
		return
	}

	response.Token = token
	writeJSON(w, response, http.StatusOK)
}
