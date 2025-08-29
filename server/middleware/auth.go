package middleware

import (
	"context"
	"log"
	"net/http"

	"github.com/dorianneto/bugfy/util"
	"github.com/golang-jwt/jwt/v5"
)

func JWTAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("jwt")
		if err != nil {
			util.WriteError(w, http.StatusUnauthorized, "missing auth token")
			return
		}

		tokenString := cookie.Value
		if tokenString == "" {
			util.WriteError(w, http.StatusUnauthorized, "missing auth token")
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(util.GetEnv("secretKey", "")), nil
		})

		if err != nil || !token.Valid {
			util.WriteError(w, http.StatusUnauthorized, "invalid auth token")
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			util.WriteError(w, http.StatusUnauthorized, "invalid token claims")
			return
		}

		userID, ok := claims["id"].(string)
		if !ok {
			util.WriteError(w, http.StatusUnauthorized, "invalid user ID in token")
			return
		}

		ctx := context.WithValue(r.Context(), "userID", userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func OptionalJWTAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("OptionalJWTAuth: Processing request to %s", r.URL.Path)

		cookie, err := r.Cookie("jwt")
		if err != nil {
			log.Printf("OptionalJWTAuth: No JWT cookie found: %v", err)
			next.ServeHTTP(w, r)
			return
		}

		if cookie.Value == "" {
			log.Printf("OptionalJWTAuth: JWT cookie is empty")
			next.ServeHTTP(w, r)
			return
		}

		log.Printf("OptionalJWTAuth: Found JWT cookie")

		token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				log.Printf("OptionalJWTAuth: Invalid signing method")
				return nil, jwt.ErrSignatureInvalid
			}
			secretKey := util.GetEnv("secretKey", "")
			if secretKey == "" {
				log.Printf("OptionalJWTAuth: Secret key not found in environment")
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(secretKey), nil
		})

		if err != nil {
			log.Printf("OptionalJWTAuth: Error parsing token: %v", err)
			next.ServeHTTP(w, r)
			return
		}

		if !token.Valid {
			log.Printf("OptionalJWTAuth: Token is invalid")
			next.ServeHTTP(w, r)
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			log.Printf("OptionalJWTAuth: Token claims: %+v", claims)
			if userID, ok := claims["id"].(string); ok {
				log.Printf("OptionalJWTAuth: Setting user ID in context: %s", userID)
				ctx := context.WithValue(r.Context(), "userID", userID)
				r = r.WithContext(ctx)
			} else {
				log.Printf("OptionalJWTAuth: No 'id' claim found in token")
			}
		} else {
			log.Printf("OptionalJWTAuth: Failed to parse claims")
		}

		next.ServeHTTP(w, r)
	})
}
