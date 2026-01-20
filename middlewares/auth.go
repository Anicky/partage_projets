package middlewares

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func Authentication() gin.HandlerFunc {
	return func(context *gin.Context) {
		authHeader := context.GetHeader("Authorization")

		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized."})

			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrTokenMalformed
			}

			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil || !token.Valid {
			context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token."})

			return
		}

		claim, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			context.JSON(http.StatusUnauthorized, gin.H{"error": "Unable to read token."})

			return
		}

		userID := int(claim["UserID"].(float64))

		context.Set("userID", userID)

		context.Next()
	}
}

func GetUserId(context *gin.Context) *uint {
	userID, ok := context.Get("userID")
	if !ok {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID."})

		return nil
	}

	userIDInt, ok := userID.(int)
	if !ok {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID."})

		return nil
	}

	userIDUint := uint(userIDInt)

	return &userIDUint
}
