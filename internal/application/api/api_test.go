package api_test

import (
	"errors"
	"testing"

	applicationAPI "github.com/felipewom/go-hexagonal/internal/application/api"
	"github.com/felipewom/go-hexagonal/internal/ports" // Required for ports.DbPort
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// --- Mocks ---

// MockArithmetic is a mock for the api.Arithmetic interface
type MockArithmetic struct {
	mock.Mock
}

func (m *MockArithmetic) Addition(a, b int32) (int32, error) {
	args := m.Called(a, b)
	// Ensure correct type assertion for return value if it's not nil
	if args.Get(0) == nil {
		return 0, args.Error(1)
	}
	return args.Get(0).(int32), args.Error(1)
}

func (m *MockArithmetic) Subtraction(a, b int32) (int32, error) {
	args := m.Called(a, b)
	if args.Get(0) == nil {
		return 0, args.Error(1)
	}
	return args.Get(0).(int32), args.Error(1)
}

func (m *MockArithmetic) Multiplication(a, b int32) (int32, error) {
	args := m.Called(a, b)
	if args.Get(0) == nil {
		return 0, args.Error(1)
	}
	return args.Get(0).(int32), args.Error(1)
}

func (m *MockArithmetic) Division(a, b int32) (int32, error) {
	args := m.Called(a, b)
	if args.Get(0) == nil {
		return 0, args.Error(1)
	}
	return args.Get(0).(int32), args.Error(1)
}

// MockDbPort is a mock for the ports.DbPort interface
type MockDbPort struct {
	mock.Mock
	// We don't need to mock CloseDbConnection for these tests,
	// but AddToHistory is crucial.
}

func (m *MockDbPort) CloseDbConnection() {
	// This method is not critical for the logic being tested in api operations,
	// but it's part of the interface.
	m.Called()
}

func (m *MockDbPort) AddToHistory(answer int32, operation string) error {
	args := m.Called(answer, operation)
	return args.Error(0)
}

// --- Tests ---

func TestApplication_GetAddition(t *testing.T) {
	mockArith := new(MockArithmetic)
	mockDb := new(MockDbPort)
	app := applicationAPI.NewApplication(mockDb, mockArith)

	mockArith.On("Addition", int32(5), int32(3)).Return(int32(8), nil)
	mockDb.On("AddToHistory", int32(8), "addition").Return(nil)

	result, err := app.GetAddition(5, 3)

	assert.NoError(t, err)
	assert.Equal(t, int32(8), result)
	mockArith.AssertExpectations(t)
	mockDb.AssertExpectations(t)
}

func TestApplication_GetAddition_ArithError(t *testing.T) {
	mockArith := new(MockArithmetic)
	mockDb := new(MockDbPort) // DB is not used if arith fails
	app := applicationAPI.NewApplication(mockDb, mockArith)

	expectedError := errors.New("arithmetic error")
	mockArith.On("Addition", int32(5), int32(3)).Return(int32(0), expectedError)
	// mockDb.AddToHistory should not be called

	result, err := app.GetAddition(5, 3)

	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Equal(t, int32(0), result)
	mockArith.AssertExpectations(t)
	mockDb.AssertNotCalled(t, "AddToHistory", mock.Anything, mock.Anything)
}

func TestApplication_GetAddition_DbError(t *testing.T) {
	mockArith := new(MockArithmetic)
	mockDb := new(MockDbPort)
	app := applicationAPI.NewApplication(mockDb, mockArith)

	expectedError := errors.New("db error")
	mockArith.On("Addition", int32(5), int32(3)).Return(int32(8), nil)
	mockDb.On("AddToHistory", int32(8), "addition").Return(expectedError)

	result, err := app.GetAddition(5, 3)

	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Equal(t, int32(0), result) // Expect 0 on error
	mockArith.AssertExpectations(t)
	mockDb.AssertExpectations(t)
}

// Similar tests can be written for GetSubtraction, GetMultiplication, GetDivision

func TestApplication_GetDivision_Success(t *testing.T) {
	mockArith := new(MockArithmetic)
	mockDb := new(MockDbPort)
	app := applicationAPI.NewApplication(mockDb, mockArith)

	mockArith.On("Division", int32(6), int32(3)).Return(int32(2), nil)
	mockDb.On("AddToHistory", int32(2), "division").Return(nil)

	result, err := app.GetDivision(6, 3)

	assert.NoError(t, err)
	assert.Equal(t, int32(2), result)
	mockArith.AssertExpectations(t)
	mockDb.AssertExpectations(t)
}

// Example for Division by zero, assuming core.ErrDivisionByZero exists and is propagated
// For this test, we'll use a generic error from the mock arithmetic
var ErrMockDivisionByZero = errors.New("mock division by zero")

func TestApplication_GetDivision_CoreError(t *testing.T) {
	mockArith := new(MockArithmetic)
	mockDb := new(MockDbPort)
	app := applicationAPI.NewApplication(mockDb, mockArith)

	mockArith.On("Division", int32(6), int32(0)).Return(int32(0), ErrMockDivisionByZero)
	// AddToHistory should not be called

	result, err := app.GetDivision(6, 0)

	assert.Error(t, err)
	assert.Equal(t, ErrMockDivisionByZero, err)
	assert.Equal(t, int32(0), result)
	mockArith.AssertExpectations(t)
	mockDb.AssertNotCalled(t, "AddToHistory", mock.Anything, mock.Anything)
}
