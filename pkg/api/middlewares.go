package api

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/aseerkt/go-simple-bank/pkg/token"
	"github.com/gin-gonic/gin"
)

const (
	authorizationHeaderKey = "authorization"
	authorizationTypeKey   = "bearer"
	authPayloadKey         = "auth"
)

func abortUnauthorized(ctx *gin.Context, err error) {
	ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
}

func auth(tm token.Maker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorizationHeader := ctx.GetHeader(authorizationHeaderKey)

		if len(authorizationHeader) < 1 {
			err := errors.New("authorization header not found")
			abortUnauthorized(ctx, err)
			return
		}

		fields := strings.Fields(authorizationHeader)

		if len(fields) < 2 {
			err := errors.New("invalid authorization header format")
			abortUnauthorized(ctx, err)
			return
		}

		authorizationType := strings.ToLower(fields[0])

		if authorizationType != authorizationTypeKey {
			err := fmt.Errorf("unsupported authorization type: %s", authorizationType)
			abortUnauthorized(ctx, err)
			return
		}

		accessToken := fields[1]
		payload, err := tm.VerifyToken(accessToken)

		if err != nil {
			abortUnauthorized(ctx, err)
			return
		}

		ctx.Set(authPayloadKey, payload)
		ctx.Next()
	}
}

func getAuthCtx(c *gin.Context) *token.Payload {
	return c.MustGet(authPayloadKey).(*token.Payload)
}
