package store

import (
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

func TestNewMemoryUserStore(t *testing.T) {
	store := NewMemoryUserStore()

	assert.NotNil(t, store)
	assert.NotNil(t, store.users)
	assert.Equal(t, 1, store.nextID)
	assert.Equal(t, 0, len(store.users))
}

func TestMemoryUserStore_Create(t *testing.T) {
	tests := []struct {
		name        string
		user        User
		expectedID  int
		expectError bool
	}{
		{
			name:        "valid user creation",
			user:        User{Name: "John Doe", Email: "john@example.com"},
			expectedID:  1,
			expectError: false,
		},
		{
			name:        "empty name user",
			user:        User{Name: "", Email: "empty@example.com"},
			expectedID:  2,
			expectError: false,
		},
		{
			name:        "empty email user",
			user:        User{Name: "Empty Email", Email: ""},
			expectedID:  3,
			expectError: false,
		},
		{
			name:        "user with existing ID (should be overwritten)",
			user:        User{ID: 999, Name: "Override ID", Email: "override@example.com"},
			expectedID:  4,
			expectError: false,
		},
	}

	store := NewMemoryUserStore()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := store.Create(tt.user)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.expectedID, result.ID)
				assert.Equal(t, tt.user.Name, result.Name)
				assert.Equal(t, tt.user.Email, result.Email)
			}
		})
	}
}

func TestMemoryUserStore_GetByID(t *testing.T) {
	store := NewMemoryUserStore()

	// Create test users
	user1, _ := store.Create(User{Name: "User 1", Email: "user1@example.com"})
	user2, _ := store.Create(User{Name: "User 2", Email: "user2@example.com"})

	tests := []struct {
		name        string
		id          int
		expected    *User
		expectError bool
		errorMsg    string
	}{
		{
			name:        "existing user 1",
			id:          user1.ID,
			expected:    user1,
			expectError: false,
		},
		{
			name:        "existing user 2",
			id:          user2.ID,
			expected:    user2,
			expectError: false,
		},
		{
			name:        "non-existent user",
			id:          999,
			expected:    nil,
			expectError: true,
			errorMsg:    "user not found",
		},
		{
			name:        "zero ID",
			id:          0,
			expected:    nil,
			expectError: true,
			errorMsg:    "user not found",
		},
		{
			name:        "negative ID",
			id:          -1,
			expected:    nil,
			expectError: true,
			errorMsg:    "user not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := store.GetByID(tt.id)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.expected.ID, result.ID)
				assert.Equal(t, tt.expected.Name, result.Name)
				assert.Equal(t, tt.expected.Email, result.Email)
			}
		})
	}
}

func TestMemoryUserStore_GetAll(t *testing.T) {
	tests := []struct {
		name          string
		setupUsers    []User
		expectedCount int
	}{
		{
			name:          "empty store",
			setupUsers:    []User{},
			expectedCount: 0,
		},
		{
			name: "single user",
			setupUsers: []User{
				{Name: "Single User", Email: "single@example.com"},
			},
			expectedCount: 1,
		},
		{
			name: "multiple users",
			setupUsers: []User{
				{Name: "User 1", Email: "user1@example.com"},
				{Name: "User 2", Email: "user2@example.com"},
				{Name: "User 3", Email: "user3@example.com"},
			},
			expectedCount: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := NewMemoryUserStore()

			// Setup test data
			createdUsers := make([]*User, 0, len(tt.setupUsers))
			for _, user := range tt.setupUsers {
				created, _ := store.Create(user)
				createdUsers = append(createdUsers, created)
			}

			result, err := store.GetAll()

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedCount, len(result))

			// Verify all users are present with correct count
			if tt.expectedCount > 0 {
				// Check that we have the right number of users
				assert.Equal(t, tt.expectedCount, len(result))

				// Verify all expected users exist (order may vary due to map iteration)
				expectedNames := make(map[string]bool)
				expectedEmails := make(map[string]bool)
				for _, user := range createdUsers {
					expectedNames[user.Name] = true
					expectedEmails[user.Email] = true
				}

				for _, resultUser := range result {
					assert.True(t, expectedNames[resultUser.Name], "Unexpected user name: %s", resultUser.Name)
					assert.True(t, expectedEmails[resultUser.Email], "Unexpected user email: %s", resultUser.Email)
					assert.NotZero(t, resultUser.ID)
				}
			}
		})
	}
}

func TestMemoryUserStore_Update(t *testing.T) {
	store := NewMemoryUserStore()
	existingUser, _ := store.Create(User{Name: "Original User", Email: "original@example.com"})

	tests := []struct {
		name        string
		id          int
		updateUser  User
		expectError bool
		errorMsg    string
	}{
		{
			name:        "successful update",
			id:          existingUser.ID,
			updateUser:  User{Name: "Updated User", Email: "updated@example.com"},
			expectError: false,
		},
		{
			name:        "update non-existent user",
			id:          999,
			updateUser:  User{Name: "Non-existent", Email: "nonexistent@example.com"},
			expectError: true,
			errorMsg:    "user not found",
		},
		{
			name:        "update with empty name",
			id:          existingUser.ID,
			updateUser:  User{Name: "", Email: "empty@example.com"},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := store.Update(tt.id, tt.updateUser)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.id, result.ID)
				assert.Equal(t, tt.updateUser.Name, result.Name)
				assert.Equal(t, tt.updateUser.Email, result.Email)

				// Verify the user was actually updated in store
				retrieved, _ := store.GetByID(tt.id)
				assert.Equal(t, tt.updateUser.Name, retrieved.Name)
				assert.Equal(t, tt.updateUser.Email, retrieved.Email)
			}
		})
	}
}

func TestMemoryUserStore_Delete(t *testing.T) {
	store := NewMemoryUserStore()
	user1, _ := store.Create(User{Name: "User 1", Email: "user1@example.com"})
	user2, _ := store.Create(User{Name: "User 2", Email: "user2@example.com"})

	tests := []struct {
		name        string
		id          int
		expectError bool
		errorMsg    string
	}{
		{
			name:        "delete existing user",
			id:          user1.ID,
			expectError: false,
		},
		{
			name:        "delete non-existent user",
			id:          999,
			expectError: true,
			errorMsg:    "user not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := store.Delete(tt.id)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				assert.NoError(t, err)

				// Verify user was actually deleted
				_, err := store.GetByID(tt.id)
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "user not found")
			}
		})
	}

	// Verify other user still exists
	retrieved, err := store.GetByID(user2.ID)
	assert.NoError(t, err)
	assert.Equal(t, user2.ID, retrieved.ID)
}

func TestMemoryUserStore_ConcurrentAccess(t *testing.T) {
	store := NewMemoryUserStore()
	var wg sync.WaitGroup

	// Test concurrent creates
	numGoroutines := 100
	wg.Add(numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			user := User{Name: fmt.Sprintf("User %d", id), Email: fmt.Sprintf("user%d@example.com", id)}
			_, err := store.Create(user)
			assert.NoError(t, err)
		}(i)
	}

	wg.Wait()

	// Verify all users were created
	users, err := store.GetAll()
	require.NoError(t, err)
	assert.Equal(t, numGoroutines, len(users))

	// Test concurrent reads and writes
	wg.Add(numGoroutines * 2)

	for i := 0; i < numGoroutines; i++ {
		// Concurrent reads
		go func() {
			defer wg.Done()
			_, err := store.GetAll()
			assert.NoError(t, err)
		}()

		// Concurrent updates
		go func(id int) {
			defer wg.Done()
			userID := (id % numGoroutines) + 1 // Use existing user IDs
			user := User{Name: fmt.Sprintf("Updated User %d", id), Email: fmt.Sprintf("updated%d@example.com", id)}
			_, err := store.Update(userID, user)
			assert.NoError(t, err)
		}(i)
	}

	wg.Wait()
}

// Test suite for interface compliance
type UserStoreTestSuite struct {
	suite.Suite
	store UserStore
}

func (suite *UserStoreTestSuite) SetupTest() {
	suite.store = NewMemoryUserStore()
}

func (suite *UserStoreTestSuite) TestCRUDWorkflow() {
	// Create
	user := User{Name: "Test User", Email: "test@example.com"}
	created, err := suite.store.Create(user)
	suite.Require().NoError(err)
	suite.NotZero(created.ID)
	suite.Equal(user.Name, created.Name)
	suite.Equal(user.Email, created.Email)

	// Read
	retrieved, err := suite.store.GetByID(created.ID)
	suite.Require().NoError(err)
	suite.Equal(created.ID, retrieved.ID)
	suite.Equal(created.Name, retrieved.Name)
	suite.Equal(created.Email, retrieved.Email)

	// Update
	updated := User{Name: "Updated User", Email: "updated@example.com"}
	result, err := suite.store.Update(created.ID, updated)
	suite.Require().NoError(err)
	suite.Equal(created.ID, result.ID)
	suite.Equal("Updated User", result.Name)
	suite.Equal("updated@example.com", result.Email)

	// Verify update persisted
	retrieved, err = suite.store.GetByID(created.ID)
	suite.Require().NoError(err)
	suite.Equal("Updated User", retrieved.Name)
	suite.Equal("updated@example.com", retrieved.Email)

	// Delete
	err = suite.store.Delete(created.ID)
	suite.Require().NoError(err)

	// Verify deletion
	_, err = suite.store.GetByID(created.ID)
	suite.Error(err)
	suite.Contains(err.Error(), "user not found")
}

func (suite *UserStoreTestSuite) TestGetAllAfterOperations() {
	// Initially empty
	users, err := suite.store.GetAll()
	suite.Require().NoError(err)
	suite.Equal(0, len(users))

	// Add some users
	user1, _ := suite.store.Create(User{Name: "User 1", Email: "user1@example.com"})
	user2, _ := suite.store.Create(User{Name: "User 2", Email: "user2@example.com"})

	users, err = suite.store.GetAll()
	suite.Require().NoError(err)
	suite.Equal(2, len(users))

	// Delete one user
	suite.store.Delete(user1.ID)

	users, err = suite.store.GetAll()
	suite.Require().NoError(err)
	suite.Equal(1, len(users))
	suite.Equal(user2.ID, users[0].ID)
}

func TestUserStoreCompliance(t *testing.T) {
	suite.Run(t, new(UserStoreTestSuite))
}

// Benchmark tests
func BenchmarkMemoryUserStore_Create(b *testing.B) {
	store := NewMemoryUserStore()
	user := User{Name: "Benchmark User", Email: "bench@example.com"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		store.Create(user)
	}
}

func BenchmarkMemoryUserStore_GetByID(b *testing.B) {
	store := NewMemoryUserStore()
	user, _ := store.Create(User{Name: "Benchmark User", Email: "bench@example.com"})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		store.GetByID(user.ID)
	}
}

func BenchmarkMemoryUserStore_GetAll(b *testing.B) {
	store := NewMemoryUserStore()
	// Create 1000 users for realistic benchmark
	for i := 0; i < 1000; i++ {
		store.Create(User{Name: fmt.Sprintf("User %d", i), Email: fmt.Sprintf("user%d@example.com", i)})
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		store.GetAll()
	}
}

func BenchmarkMemoryUserStore_ConcurrentReads(b *testing.B) {
	store := NewMemoryUserStore()
	// Setup test data
	for i := 0; i < 100; i++ {
		store.Create(User{Name: fmt.Sprintf("User %d", i), Email: fmt.Sprintf("user%d@example.com", i)})
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			store.GetAll()
		}
	})
}
