package auth

import (
	"context"
	"testing"
	"time"

	"github.com/ai-agentic-browser/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// AuthServiceTestSuite tests the authentication service
type AuthServiceTestSuite struct {
	suite.Suite
	service *Service
	ctx     context.Context
}

// Helper function to convert string to *string
func stringPtr(s string) *string {
	return &s
}

// SetupSuite initializes the test suite
func (suite *AuthServiceTestSuite) SetupSuite() {
	suite.ctx = context.Background()

	// Initialize auth service with test dependencies
	suite.service = &Service{
		config: config.JWTConfig{
			Secret:             "test-secret",
			Expiry:             time.Hour,
			RefreshTokenExpiry: 24 * time.Hour,
		},
	}
}

// TestPasswordHashing tests password hashing and verification
func (suite *AuthServiceTestSuite) TestPasswordHashing() {
	password := "testpassword123"

	// Test hashing
	hashedPassword, err := suite.service.HashPassword(password)
	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), hashedPassword)
	assert.NotEqual(suite.T(), password, hashedPassword)

	// Test verification
	isValid := suite.service.VerifyPassword(password, hashedPassword)
	assert.True(suite.T(), isValid)

	// Test wrong password
	isValid = suite.service.VerifyPassword("wrongpassword", hashedPassword)
	assert.False(suite.T(), isValid)
}

// TestValidatePasswordComplexity tests password complexity validation
func (suite *AuthServiceTestSuite) TestValidatePasswordComplexity() {
	// Test with nil security config (should pass)
	err := suite.service.validatePasswordComplexity("simplepass")
	assert.NoError(suite.T(), err)
}

// Run the test suite
func TestAuthServiceSuite(t *testing.T) {
	suite.Run(t, new(AuthServiceTestSuite))
}

// Benchmark tests for performance
func BenchmarkHashPassword(b *testing.B) {
	service := &Service{}
	password := "benchmarkpassword123"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := service.HashPassword(password)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkVerifyPassword(b *testing.B) {
	service := &Service{}
	password := "benchmarkpassword123"
	hashedPassword, err := service.HashPassword(password)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		service.VerifyPassword(password, hashedPassword)
	}
}
