package response

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Status string

const (
	StatusSuccess Status = "success"
	StatusError   Status = "error"
)

type Paging struct {
	Page         int  `json:"page"`
	PerPage      int  `json:"per_page"`
	TotalPages   int  `json:"total_pages"`
	TotalRecords int  `json:"total_records"`
	HasNext      bool `json:"has_next"`
	HasPrev      bool `json:"has_prev"`
}

type Response struct {
	Status    `json:"status"`
	Code      int         `json:"code"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`
	Paging    *Paging     `json:"paging,omitempty"`
	Timestamp time.Time  `json:"timestamp"`
}

func NewResponseSuccessPaging(c *gin.Context, result interface{}, paging *Paging) {
	c.JSON(http.StatusOK, Response{
		Status:  StatusSuccess,
		Code:    http.StatusOK,
		Message: "Success",
		Data:    result,
		Paging:  paging,
		Timestamp: time.Now(),
	})
}

func NewResponseSuccess(c *gin.Context, result interface{}) {
	c.JSON(http.StatusOK, Response{
		Status:  StatusSuccess,
		Code:    http.StatusOK,
		Message: "Success",
		Data:    result,
		Timestamp: time.Now(),
	})
}

func NewResponseCreated(c *gin.Context, result interface{}) {
	c.JSON(http.StatusCreated, Response{
		Status:  StatusSuccess,
		Code:    http.StatusCreated,
		Message: "Created",
		Data:    result,
		Timestamp: time.Now(),
	})
}

func NewResponseBadRequest(c *gin.Context, err string) {
	c.JSON(http.StatusBadRequest, Response{
		Status:  StatusError,
		Code:    http.StatusBadRequest,
		Message: err,
		Timestamp: time.Now(),
	})
}

func NewResponseError(c *gin.Context, err string) {
	c.JSON(http.StatusInternalServerError, Response{
		Status:  StatusError,
		Code:    http.StatusInternalServerError,
		Message: err,
		Timestamp: time.Now(),
	})
}

func NewResponseForbidden(c *gin.Context, err string) {
	c.JSON(http.StatusForbidden, Response{
		Status:  StatusError,
		Code:    http.StatusForbidden,
		Message: err,
		Timestamp: time.Now(),
	})
}

func NewResponseUnauthorized(c *gin.Context, err string) {
	c.JSON(http.StatusUnauthorized, Response{
		Status:  StatusError,
		Code:    http.StatusUnauthorized,
		Message: err,
		Timestamp: time.Now(),
	})
}
