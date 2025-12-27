package response

import (
	"github.com/gin-gonic/gin"
)

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

// Success sends a success response
func Success(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// Error sends an error response
func Error(c *gin.Context, statusCode int, message string, err interface{}) {
	c.JSON(statusCode, Response{
		Success: false,
		Message: message,
		Error:   err,
	})
}

// ValidationError sends a validation error response
func ValidationError(c *gin.Context, errors interface{}) {
	c.JSON(400, Response{
		Success: false,
		Message: "Validation failed",
		Error:   errors,
	})
}

// SetAuthCookies sets access token and refresh token cookies
func SetAuthCookies(c *gin.Context, accessToken, refreshToken string, expiresIn int64) {
	// Set access token cookie
	c.SetCookie(
		"access_token", // name
		accessToken,    // value
		int(expiresIn), // maxAge in seconds
		"/",            // path
		"",             // domain (empty = current domain)
		false,          // secure (set true in production with HTTPS)
		true,           // httpOnly
	)

	// Set refresh token cookie - 7 days (604800 seconds)
	c.SetCookie(
		"refresh_token", // name
		refreshToken,    // value
		604800,          // maxAge in seconds (7 days)
		"/",             // path
		"",              // domain
		false,           // secure
		true,            // httpOnly
	)
}

// ClearAuthCookies clears authentication cookies
func ClearAuthCookies(c *gin.Context) {
	c.SetCookie("access_token", "", -1, "/", "", false, true)
	c.SetCookie("refresh_token", "", -1, "/", "", false, true)
}
