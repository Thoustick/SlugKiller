package mocks

import (
	"github.com/gin-gonic/gin"
)

type MockHandler struct {
	Started        bool
	RoutesRegister bool
}

func (m *MockHandler) Start(_ string) error {
	m.Started = true
	return nil
}

func (m *MockHandler) RegisterRoutes(_ *gin.Engine) {
	m.RoutesRegister = true
}
