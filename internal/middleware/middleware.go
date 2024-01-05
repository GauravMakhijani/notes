package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/GauravMakhijani/notes/internal/jwt"
	gojwt "github.com/golang-jwt/jwt"
	"github.com/sirupsen/logrus"
	"golang.org/x/time/rate"
)

type Response struct {
	Data         interface{} `json:"data,omitempty"`
	ErrorMessage string      `json:"error,omitempty"`
	ErrorCode    int64       `json:"error_code,omitempty"`
}

// WriteServerErrorResponse ...
func WriteServerErrorResponse(ctx context.Context, w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	_, err := w.Write([]byte(fmt.Sprintf("{\"message\":%s}", "internal server error")))
	if err != nil {
		logrus.Error("cannot write server error response body", err)
	}
}
func ErrResponse(ctx context.Context, w http.ResponseWriter, statusCode int, errorCode int64, err error) {

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	w.WriteHeader(statusCode)

	response := Response{
		ErrorMessage: err.Error(),
	}
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		WriteServerErrorResponse(ctx, w)
		return
	}
}

func SetMiddleWareAuthentication(next http.HandlerFunc) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {

		reqKey := r.Header.Get("Authorization")

		if len(reqKey) > 6 && strings.ToUpper(reqKey[0:7]) == "BEARER " {
			reqKey = reqKey[7:]
		}

		userInfo, err := jwt.GetUserInfoFromToken(reqKey)
		if err != nil {
			if e, ok := err.(*gojwt.ValidationError); ok && e.Errors == gojwt.ValidationErrorExpired {
				ErrResponse(r.Context(), rw, http.StatusUnauthorized, 0, errors.New("token expired"))
				return
			}
			ErrResponse(r.Context(), rw, http.StatusUnauthorized, 0, errors.New("invalid token"))
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, "user_id", userInfo.UserID)
		ctx = context.WithValue(ctx, "user_name", userInfo.UserName)
		requestWithValueContext := r.WithContext(ctx)
		next(rw, requestWithValueContext)
	}
}

func RateLimiter(next http.Handler) http.Handler {
	limiter := rate.NewLimiter(10, 20)
	fmt.Println("limiter", limiter.Limit(), limiter.Burst())
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !limiter.Allow() {
			fmt.Println("Request rejected: Too many requests")
			ErrResponse(r.Context(), w, http.StatusTooManyRequests, 0, errors.New("too many requests"))
			return
		} else {
			next.ServeHTTP(w, r)
		}
	})
}
