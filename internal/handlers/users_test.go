package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/dazraf/go-api-example/internal/store"
)

// MockUserStore for testing
type MockUserStore struct {
	mock.Mock
}

func (m *MockUserStore) GetAll() ([]store.User, error) {
	args := m.Called()
	return args.Get(0).([]store.User), args.Error(1)
}

func (m *MockUserStore) GetByID(id int) (*store.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*store.User), args.Error(1)
}

func (m *MockUserStore) Create(user store.User) (*store.User, error) {
	args := m.Called(user)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*store.User), args.Error(1)
}

func (m *MockUserStore) Update(id int, user store.User) (*store.User, error) {
	args := m.Called(id, user)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*store.User), args.Error(1)
}

func (m *MockUserStore) Delete(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func setupTestRouter(userStore store.UserStore) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	handler := NewUserHandler(userStore)

	v1 := router.Group("/api/v1")
	{
		v1.GET("/users", handler.GetUsers)
		v1.GET("/users/:id", handler.GetUser)
		v1.POST("/users", handler.CreateUser)
		v1.PUT("/users/:id", handler.UpdateUser)
		v1.DELETE("/users/:id", handler.DeleteUser)
	}

	return router
}

func TestUserHandler_GetUsers(t *testing.T) {
	tests := []struct {
		name           string
		setupMock      func(*MockUserStore)
		expectedStatus int
		expectedBody   func(t *testing.T, body string)
	}{
		{
			name: "successful get all users",
			setupMock: func(m *MockUserStore) {
				users := []store.User{
					{ID: 1, Name: "John Doe", Email: "john@example.com"},
					{ID: 2, Name: "Jane Smith", Email: "jane@example.com"},
				}
				m.On("GetAll").Return(users, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: func(t *testing.T, body string) {
				var users []store.User
				err := json.Unmarshal([]byte(body), &users)
				require.NoError(t, err)
				assert.Equal(t, 2, len(users))
				assert.Equal(t, "John Doe", users[0].Name)
				assert.Equal(t, "Jane Smith", users[1].Name)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStore := new(MockUserStore)
			tt.setupMock(mockStore)

			router := setupTestRouter(mockStore)

			req, err := http.NewRequest("GET", "/api/v1/users", nil)
			require.NoError(t, err)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			tt.expectedBody(t, w.Body.String())

			mockStore.AssertExpectations(t)
		})
	}
}

func TestUserHandler_CreateUser(t *testing.T) {
	tests := []struct {
		name           string
		payload        interface{}
		setupMock      func(*MockUserStore)
		expectedStatus int
		expectedBody   func(t *testing.T, body string)
	}{
		{
			name:    "successful user creation",
			payload: store.User{Name: "John Doe", Email: "john@example.com"},
			setupMock: func(m *MockUserStore) {
				inputUser := store.User{Name: "John Doe", Email: "john@example.com"}
				createdUser := &store.User{ID: 1, Name: "John Doe", Email: "john@example.com"}
				m.On("Create", inputUser).Return(createdUser, nil)
			},
			expectedStatus: http.StatusCreated,
			expectedBody: func(t *testing.T, body string) {
				var user store.User
				err := json.Unmarshal([]byte(body), &user)
				require.NoError(t, err)
				assert.Equal(t, 1, user.ID)
				assert.Equal(t, "John Doe", user.Name)
				assert.Equal(t, "john@example.com", user.Email)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStore := new(MockUserStore)
			tt.setupMock(mockStore)

			router := setupTestRouter(mockStore)

			var payload []byte
			var err error

			if str, ok := tt.payload.(string); ok {
				payload = []byte(str)
			} else {
				payload, err = json.Marshal(tt.payload)
				require.NoError(t, err)
			}

			req, err := http.NewRequest("POST", "/api/v1/users", bytes.NewBuffer(payload))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			tt.expectedBody(t, w.Body.String())

			mockStore.AssertExpectations(t)
		})
	}
}

// Integration test with real store
func TestUserHandler_Integration_FullCRUDWorkflow(t *testing.T) {
	realStore := store.NewMemoryUserStore()
	router := setupTestRouter(realStore)

	// Create user
	user := store.User{Name: "Integration Test User", Email: "integration@example.com"}
	payload, _ := json.Marshal(user)

	req, _ := http.NewRequest("POST", "/api/v1/users", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var createdUser store.User
	_ = json.Unmarshal(w.Body.Bytes(), &createdUser)
	assert.NotZero(t, createdUser.ID)

	// Get user
	req, _ = http.NewRequest("GET", fmt.Sprintf("/api/v1/users/%d", createdUser.ID), nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// Get all users
	req, _ = http.NewRequest("GET", "/api/v1/users", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	var users []store.User
	_ = json.Unmarshal(w.Body.Bytes(), &users)
	assert.Equal(t, 1, len(users))

	// Update user
	updatedUser := store.User{Name: "Updated Integration User", Email: "updated@example.com"}
	payload, _ = json.Marshal(updatedUser)
	req, _ = http.NewRequest("PUT", fmt.Sprintf("/api/v1/users/%d", createdUser.ID), bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// Delete user
	req, _ = http.NewRequest("DELETE", fmt.Sprintf("/api/v1/users/%d", createdUser.ID), nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNoContent, w.Code)

	// Verify deletion
	req, _ = http.NewRequest("GET", fmt.Sprintf("/api/v1/users/%d", createdUser.ID), nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}
