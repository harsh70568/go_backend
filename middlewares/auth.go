package middlewares

import (
	"go_edtech_backend/db"
	"go_edtech_backend/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func AuthCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie("token")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing authentication tokens", "details": err.Error()})
			return
		}

		refreshToken, err := c.Cookie("refresh_token")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing authentication tokens", "details": err.Error()})
			return
		}

		tok, err := utils.VerifyToken(token)
		reftok, referr := utils.VerifyToken(refreshToken)

		/* If both tokens are invalid, ask the user to login again */
		if err != nil && referr != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Tokens expired, Please login again"})
			c.Abort()
			return
		}

		/* If access token is invalid but refresh token is valid, generate new tokens */
		if err != nil {
			claims, ok := reftok.Claims.(jwt.MapClaims)
			if !ok {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse refresh token claims"})
				c.Abort()
				return
			}

			email, ok := claims["email"].(string)
			if !ok {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Inavlid refresh token"})
				c.Abort()
				return
			}

			var storedRefreshToken string
			Qerr := db.DB.Where("email = ?", email).First(&storedRefreshToken).Error
			if Qerr != nil || storedRefreshToken != refreshToken {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Inavlid refresh token"})
				c.Abort()
				return
			}

			newToken, newRefToken, Nerr := utils.GenerateNewTokens(email)
			if Nerr != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating tokens", "details": Nerr.Error()})
				c.Abort()
				return
			}

			Uquery := `Update users SET refresh_token = $1 WHERE email = $2`
			Uerr := db.DB.Exec(Uquery, newRefToken, email)
			if Uerr != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating refresh token"})
				c.Abort()
				return
			}

			c.SetCookie("token", newToken, int(48*time.Hour.Seconds()), "/", "localhost", false, true)
			c.SetCookie("refresh_token", newRefToken, int(240*time.Hour.Seconds()), "/", "localhost", false, true)

			c.Set("email", email)
			c.Next()
		} else {
			/* If access token is valid, extract email from it */
			claims, ok := tok.Claims.(jwt.MapClaims)
			if !ok {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse token claims"})
				c.Abort()
				return
			}

			email, ok := claims["email"].(string)
			if !ok {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Inavlid access token"})
				c.Abort()
				return
			}

			c.Set("email", email)
			c.Next()
		}
	}
}
