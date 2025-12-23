package middleware

import (
	"go-starter-app/app/http/utils"
	"go-starter-app/interfaces"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// JWTAuthMiddleware validates JWT tokens
func JWTAuthMiddleware(app interfaces.IAppDependencies) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Retrieve the Authorization header
		authorization := c.GetHeader("Authorization")
		if authorization == "" {
			utils.SendError(c, http.StatusUnauthorized, "Authorization header required", nil)
			c.Abort()
			return
		}

		// Extract the token from the header (e.g., "Bearer <token>")
		tokenString := strings.TrimPrefix(authorization, "Bearer ")
		if tokenString == authorization {
			utils.SendError(c, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized), nil)
			c.Abort()
			return
		}

		// check if token is exists in database
		//session, err := app.GetService().AuthService.GetSessionByToken(tokenString)
		//if err != nil {
		//  utils.SendError(c, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized), "Invalid or expired token")
		//	c.Abort()
		//	return
		//}
		//
		//// Parse the token
		//token, err := jwt.Parse(session.Token, func(token *jwt.Token) (interface{}, error) {
		//	// Validate the signing method
		//	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		//		return nil, errors.New("unexpected signing method")
		//	}
		//
		//	// todo: move jwt secret to config/env
		//	return []byte("lsadkflfjieowfdlkdljskaldskjaodqwkdjlksdladksajd"), nil
		//})

		// Handle parsing errors
		//if err != nil || !token.Valid {
		//  utils.SendError(c, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized), "Invalid or expired token")
		//	c.Abort()
		//	return
		//}

		// Extract claims
		//claims, ok := token.Claims.(jwt.MapClaims)
		//if !ok || !token.Valid {
		//  utils.SendError(c, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized), "Invalid token claims")
		//	c.Abort()
		//	return
		//}

		// Check expiration
		//if exp, ok := claims["exp"].(float64); ok {
		//	if time.Now().Unix() > int64(exp) {
		//  	utils.SendError(c, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized), "Token expired")
		//		c.Abort()
		//		return
		//	}
		//} else {
		//  utils.SendError(c, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized), "Expiration claim missing")
		//	c.Abort()
		//	return
		//}

		//user := map[string]interface{}{
		//	"id":       claims["id"],
		//	"email":    claims["email"],
		//	"username": claims["username"],
		//	"name":     claims["name"],
		//	"exp":      claims["exp"],
		//	"roles":    claims["roles"],
		//}

		// set user claims
		//c.Set("user", user)

		// Proceed to the next middleware or route handler
		c.Next()
	}
}
