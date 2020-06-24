package jwt

import (
	"context"
	"errors"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/go-kit/kit/endpoint"
	"google.golang.org/grpc/metadata"
)

var srecretKey = []byte("bangumi-go-nan^8193")

type BangumiCustomClaims struct {
	UserId int32  `json:"userId"`
	Name   string `json:"name"`

	jwt.StandardClaims
}

func JwtKeyFunc(token *jwt.Token) (interface{}, error) {
	return srecretKey, nil
}

func Sign(name string, uid int32) (string, error) {
	expAt := time.Now().Add(time.Duration(24) * time.Hour).Unix()

	claims := BangumiCustomClaims{
		UserId: uid,
		Name:   name,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expAt,
			Issuer:    "system",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(srecretKey)
}

func AuthMiddleware() endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			md, ok := metadata.FromIncomingContext(ctx)
			var tokenString string
			if ok && len(md["authorization"][0]) != 0 {
				temp := strings.Split(md["authorization"][0], " ")
				if len(temp) > 1 {
					tokenString = temp[1]
				}
			} else {
				return nil, errors.New("token check error")
			}

			token, err := jwt.ParseWithClaims(tokenString, &BangumiCustomClaims{}, JwtKeyFunc)
			if token.Valid {
				return next(ctx, request)
			}
			return nil, errors.New("token invalid")
		}
	}
}
