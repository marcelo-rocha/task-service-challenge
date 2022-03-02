package server

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/marcelo-rocha/task-service-challenge/domain/entities"
	"github.com/marcelo-rocha/task-service-challenge/domain/user"

	"github.com/dgrijalva/jwt-go"
)

var AuthencationSecretKey []byte

func authenticationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := strings.Split(r.Header.Get("Authorization"), "Bearer ")
		if len(authHeader) != 2 {
			fmt.Println("Malformed token")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Malformed Token"))
		} else {
			jwtToken := authHeader[1]
			token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return AuthencationSecretKey, nil
			})

			if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {

				ctx := r.Context()
				var v interface{}
				var s string
				var ok bool
				var userId int64
				var kind entities.UserKind

				if v, ok = claims["sub"]; !ok {
					goto Next
				}
				if userId, ok = v.(int64); !ok {
					goto Next
				}

				if v, ok = claims["taskschallenge.user_kind"]; !ok {
					goto Next
				}
				if s, ok = v.(string); !ok {
					goto Next
				}
				kind = entities.ToUserKind(s)
				ctx = user.ContextWithUserInfo(ctx, userId, kind)

			Next:
				next.ServeHTTP(w, r.WithContext(ctx))
			} else {
				fmt.Println(err)
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("Unauthorized"))
			}
		}
	})
}