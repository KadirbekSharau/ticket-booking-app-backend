package helpers

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type Response struct {
	Message string `json:"message"`
}

func NewResponse(message string) Response {
	return Response{message}
}

type responseError struct {
	Detail string `json:"detail"`
}

func NewErrorResponse(c *gin.Context, statusCode int, message string) {
	tc := cases.Title(language.English)
	words := strings.Split(message, " ")
	words[0] = tc.String(words[0])
	message = strings.Join(words, " ")
	logrus.Error(message)
	c.AbortWithStatusJSON(statusCode, Response{message})
}
