package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/sandeepshukla/server2admin/models"
)

// UserClient is the HTTP client that talks to Server1User
type UserClient struct {
	baseURL    string
	httpClient *http.Client
}

// NewUserClient creates a new client pointing at Server1User's base URL
func NewUserClient(baseURL string) *UserClient {
	return &UserClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second, // fail fast if Server1User is down
		},
	}
}

// GetAllUsers calls GET /users on Server1User
func (c *UserClient) GetAllUsers() ([]models.User, error) {
	resp, err := c.httpClient.Get(c.baseURL + "/users")
	if err != nil {
		return nil, fmt.Errorf("failed to reach Server1User: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Server1User returned status %d", resp.StatusCode)
	}

	var users []models.User
	if err := json.NewDecoder(resp.Body).Decode(&users); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	return users, nil
}

// GetUserByID calls GET /users/{id} on Server1User
func (c *UserClient) GetUserByID(id string) (*models.User, int, error) {
	resp, err := c.httpClient.Get(c.baseURL + "/users/" + id)
	if err != nil {
		return nil, http.StatusBadGateway, fmt.Errorf("failed to reach Server1User: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, http.StatusNotFound, fmt.Errorf("user not found")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, resp.StatusCode, fmt.Errorf("Server1User returned status %d", resp.StatusCode)
	}

	var user models.User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("failed to decode response: %w", err)
	}
	return &user, http.StatusOK, nil
}

// CreateUser calls POST /users on Server1User
func (c *UserClient) CreateUser(req models.CreateUserRequest) (*models.User, int, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("failed to encode request: %w", err)
	}

	resp, err := c.httpClient.Post(
		c.baseURL+"/users",
		"application/json",
		bytes.NewBuffer(body),
	)
	if err != nil {
		return nil, http.StatusBadGateway, fmt.Errorf("failed to reach Server1User: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusConflict {
		return nil, http.StatusConflict, fmt.Errorf("email already exists in Server1User")
	}
	if resp.StatusCode != http.StatusCreated {
		return nil, resp.StatusCode, fmt.Errorf("Server1User returned status %d", resp.StatusCode)
	}

	var user models.User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("failed to decode response: %w", err)
	}
	return &user, http.StatusCreated, nil
}

// UpdateUser calls PUT /users/{id} on Server1User
func (c *UserClient) UpdateUser(id string, req models.UpdateUserRequest) (*models.User, int, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("failed to encode request: %w", err)
	}

	httpReq, err := http.NewRequest(http.MethodPut, c.baseURL+"/users/"+id, bytes.NewBuffer(body))
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("failed to create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, http.StatusBadGateway, fmt.Errorf("failed to reach Server1User: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, http.StatusNotFound, fmt.Errorf("user not found")
	}
	if resp.StatusCode == http.StatusConflict {
		return nil, http.StatusConflict, fmt.Errorf("email already exists in Server1User")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, resp.StatusCode, fmt.Errorf("Server1User returned status %d", resp.StatusCode)
	}

	var user models.User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("failed to decode response: %w", err)
	}
	return &user, http.StatusOK, nil
}

// DeleteUser calls DELETE /users/{id} on Server1User
func (c *UserClient) DeleteUser(id string) (int, error) {
	httpReq, err := http.NewRequest(http.MethodDelete, c.baseURL+"/users/"+id, nil)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return http.StatusBadGateway, fmt.Errorf("failed to reach Server1User: %w", err)
	}
	defer resp.Body.Close()

	// Drain response body
	io.Copy(io.Discard, resp.Body)

	if resp.StatusCode == http.StatusNotFound {
		return http.StatusNotFound, fmt.Errorf("user not found")
	}
	if resp.StatusCode != http.StatusOK {
		return resp.StatusCode, fmt.Errorf("Server1User returned status %d", resp.StatusCode)
	}
	return http.StatusOK, nil
}
