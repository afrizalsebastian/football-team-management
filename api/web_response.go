package api

import (
	"context"

	"github.com/gin-gonic/gin"
)

type ErrorsDetail struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}

type WebResponse[T any] struct {
	MessageId    string         `json:"message_id"`
	HttpCode     int            `json:"http_code"`
	MessageCode  int            `json:"message_code"`
	Message      string         `json:"message,omitempty"`
	Data         T              `json:"data,omitempty"`
	ErrorMessage string         `json:"error_message,omitempty"`
	ErrorsDetail []ErrorsDetail `json:"errors_detail,omitempty"`
}

func SuccessResponse[T any](ctx context.Context, message string, messageCode int, httpCode int, data T) WebResponse[T] {
	msgId, ok := ctx.Value("message_id").(string)
	if !ok || msgId == "" {
		msgId = "Unknown"
	}

	return WebResponse[T]{
		MessageId:   msgId,
		Message:     message,
		HttpCode:    httpCode,
		MessageCode: messageCode,
		Data:        data,
	}
}

func ErrorResponse[T any](ctx context.Context, errorMessage string, messageCode int, httpCode int, errorsDetail []ErrorsDetail) WebResponse[T] {
	msgId, ok := ctx.Value("message_id").(string)
	if !ok || msgId == "" {
		msgId = "Unknown"
	}

	return WebResponse[T]{
		MessageId:    msgId,
		HttpCode:     httpCode,
		MessageCode:  messageCode,
		ErrorMessage: errorMessage,
		ErrorsDetail: errorsDetail,
	}
}

func WriteJSONResponse[T any](c *gin.Context, statusCode int, response WebResponse[T]) {
	c.JSON(statusCode, response)
}
