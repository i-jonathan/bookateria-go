package auth

import (
	"bookateria-api-go/core"
	"fmt"
	"github.com/go-redis/redis/v8"
	"net/http"
)

func AuthorizationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, email := core.GetTokenEmail(w, r)
		redisToken, err := redisClient.Get(ctx, email).Result()
		if err != nil {
			if err == redis.Nil {
				fmt.Println("key does not exists")
				return
			}
			panic(err)
		}
		tokenString, _ := token.SignedString(jwtKey)
		if redisToken != tokenString {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}