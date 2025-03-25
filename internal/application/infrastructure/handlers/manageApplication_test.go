package handlers

import (
	"encoding/json"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"greye/internal/application/domain/models"
	models2 "greye/pkg/client/domain/models"
	valPort "greye/pkg/validator/domain/ports"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

// Mocking dependencies
type MockSchedulerData struct {
	mock.Mock
}

func (m *MockSchedulerData) GetApplication(url string) (map[string]models.SchedulerApplication, error) {
	args := m.Called(url)
	return args.Get(0).(map[string]models.SchedulerApplication), args.Error(1)
}

func (m *MockSchedulerData) MonitorApplication(app *models.SchedulerApplication, startupPhase bool) error {
	args := m.Called(app, startupPhase)
	return args.Error(0)
}

func (m *MockSchedulerData) DeleteApplicationFromUrl(service string) error {
	args := m.Called(service)
	return args.Error(0)
}

type MockLogger struct {
	mock.Mock
}

func (m *MockLogger) Error(msg string, args ...interface{}) {
	m.Called(msg, args)
}

func (m *MockLogger) Warn(msg string, args ...interface{}) {
	m.Called(msg, args)
}

func (m *MockLogger) Info(msg string, args ...interface{}) {
	m.Called(msg, args)
}

func (m *MockLogger) Debug(msg string, args ...interface{}) {
	m.Called(msg, args)
}

func (m *MockLogger) Trace(msg string, args ...interface{}) {
	m.Called(msg, args)
}

func TestUnscheduleApplication(t *testing.T) {
	// Create mock dependencies
	mockSchedulerData := new(MockSchedulerData)
	mockLogger := new(MockLogger)

	mockLogger.On("Error", "%s", mock.Anything).Return()

	// Initialize handler with mocked dependencies
	handler := &ApplicationHdl{
		schedulerData: mockSchedulerData,
		logger:        mockLogger,
	}

	// Create test cases
	tests := []struct {
		service         string
		mockDeleteError error
		expectedStatus  int
		expectedMessage string
	}{
		{
			service:         "my-service",
			mockDeleteError: nil, // No error when deleting
			expectedStatus:  fiber.StatusOK,
			expectedMessage: "Service unscheduled",
		},
		{
			service:         "my-service",
			mockDeleteError: errors.New("delete error"), // Simulate error
			expectedStatus:  fiber.StatusInternalServerError,
			expectedMessage: "Error unscheduling application",
		},
	}

	for _, tt := range tests {
		t.Run(tt.service, func(t *testing.T) {
			// Decode service in test to simulate URL decoding
			decodedService, err := url.QueryUnescape(tt.service)

			assert.NoError(t, err)

			// Mock DeleteApplicationFromUrl behavior
			mockSchedulerData.On("DeleteApplicationFromUrl", decodedService).Return(tt.mockDeleteError).Once()

			// Create a Fiber app for the test
			app := fiber.New()

			// Register the handler for the DELETE route
			app.Delete("/api/v1/application/monitor/:service", handler.UnscheduleApplication)

			// Send the request (encoded service)
			req := httptest.NewRequest(http.MethodDelete, "/api/v1/application/monitor/"+tt.service, nil)
			resp, err := app.Test(req)

			// Ensure no error occurred
			assert.NoError(t, err)

			// Check the response status and message
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
			var responseBody models.EasyResponse
			err = json.NewDecoder(resp.Body).Decode(&responseBody)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedMessage, responseBody.Message) // Uncomment this line
		})
	}
}

func TestGetApplicationMonitored(t *testing.T) {
	// Create mock dependencies
	mockSchedulerData := new(MockSchedulerData)
	mockLogger := new(MockLogger)

	// Configure logger mock
	mockLogger.On("Error", mock.AnythingOfType("string"), mock.Anything).Return()

	// Initialize handler with mocked dependencies
	handler := &ApplicationHdl{
		schedulerData: mockSchedulerData,
		logger:        mockLogger,
	}

	// Create test cases
	tests := []struct {
		name           string
		url            string
		mockReturnData map[string]models.SchedulerApplication
		mockReturnErr  error
		expectedStatus int
	}{
		{
			name: "Success case",
			url:  "http://example.com",
			mockReturnData: map[string]models.SchedulerApplication{
				"app1": {
					MonitoringHttpRequest: models2.MonitoringHttpRequest{
						Name: "app1",
						Host: "http://example.com",
						Path: []string{"/app1"},
					},
				},
			},
			mockReturnErr:  nil,
			expectedStatus: fiber.StatusOK,
		},
		{
			name:           "Error case",
			url:            "http://example.com",
			mockReturnData: map[string]models.SchedulerApplication{},
			mockReturnErr:  errors.New("failed to get application"),
			expectedStatus: fiber.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Mock GetApplication behavior
			mockSchedulerData.On("GetApplication", tt.url).Return(tt.mockReturnData, tt.mockReturnErr).Once()

			// Create a Fiber app for the test
			app := fiber.New()

			// Register the handler for the GET route
			app.Get("/api/v1/application/monitor", handler.GetApplicationMonitored)

			// Create URL with query parameter
			reqURL := "/api/v1/application/monitor"
			if tt.url != "" {
				reqURL = reqURL + "?url=" + tt.url
			}

			// Send the request
			req := httptest.NewRequest(http.MethodGet, reqURL, nil)
			resp, err := app.Test(req)

			// Ensure no error occurred
			assert.NoError(t, err)

			// Check the response status
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			if tt.mockReturnErr == nil {
				// For success case, check returned data
				var responseBody map[string]models.SchedulerApplication
				err = json.NewDecoder(resp.Body).Decode(&responseBody)
				assert.NoError(t, err)
				assert.Equal(t, tt.mockReturnData, responseBody)
			} else {
				// For error case, check error message
				var responseBody models.EasyResponse
				err = json.NewDecoder(resp.Body).Decode(&responseBody)
				assert.NoError(t, err)
				assert.Equal(t, "Failed to retrieve application data", responseBody.Message)
				assert.Equal(t, tt.mockReturnErr.Error(), responseBody.Error)
			}

			// Verify mocks were called as expected
			mockSchedulerData.AssertExpectations(t)
		})
	}
}

type MockValidator struct {
	mock.Mock
}

func (m *MockValidator) Struct(s valPort.Evaluable) error {
	args := m.Called(s)
	return args.Error(0)
}

func TestMonitoringApplicationPut(t *testing.T) {
	// Create mock dependencies
	mockSchedulerData := new(MockSchedulerData)
	mockLogger := new(MockLogger)
	mockValidator := new(MockValidator)

	// Configure logger mock
	mockLogger.On("Error", mock.AnythingOfType("string"), mock.Anything).Return()

	// Initialize handler with mocked dependencies
	handler := &ApplicationHdl{
		schedulerData: mockSchedulerData,
		logger:        mockLogger,
		validator:     mockValidator,
	}

	// Create test cases
	tests := []struct {
		name               string
		requestBody        string
		mockValidateError  error
		mockMonitorError   error
		expectedStatusCode int
		expectedMessage    string
	}{
		{
			name: "Success case",
			requestBody: `[{
				"name": "test-app",
				"host": "http://example.com",
				"path": ["/api"],
				"method": "GET",
				"headers": {"Content-Type": "application/json"},
				"expectedStatus": 200
			}]`,
			mockValidateError:  nil,
			mockMonitorError:   nil,
			expectedStatusCode: fiber.StatusOK,
			expectedMessage:    "Application added to monitoring successfully",
		},
		{
			name:               "Invalid request body",
			requestBody:        `invalid json`,
			mockValidateError:  nil,
			mockMonitorError:   nil,
			expectedStatusCode: fiber.StatusBadRequest,
			expectedMessage:    "Invalid request body",
		},
		{
			name: "Validation error",
			requestBody: `[{
				"name": "test-app",
				"host": "http://example.com"
			}]`,
			mockValidateError:  errors.New("validation error"),
			mockMonitorError:   nil,
			expectedStatusCode: fiber.StatusInternalServerError,
			expectedMessage:    "Error validating application",
		},
		{
			name: "Monitor application error",
			requestBody: `[{
				"name": "test-app",
				"host": "http://example.com",
				"path": ["/api"],
				"method": "GET",
				"headers": {"Content-Type": "application/json"},
				"expectedStatus": 200
			}]`,
			mockValidateError:  nil,
			mockMonitorError:   errors.New("monitoring error"),
			expectedStatusCode: fiber.StatusInternalServerError,
			expectedMessage:    "Error while adding application to monitor",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup app and route
			app := fiber.New()
			app.Put("/api/v1/application/monitor", handler.MonitoringApplication)

			if tt.requestBody != "invalid json" {
				// Setup validator mock for valid JSON
				mockValidator.On("Struct", mock.AnythingOfType("*models.SchedulerApplication")).Return(tt.mockValidateError).Once()

				// Setup monitor application mock for cases where validation passes
				if tt.mockValidateError == nil {
					mockSchedulerData.On("MonitorApplication", mock.AnythingOfType("*models.SchedulerApplication"), false).Return(tt.mockMonitorError).Once()
				}
			}

			// Create request
			req := httptest.NewRequest(http.MethodPut, "/api/v1/application/monitor", strings.NewReader(tt.requestBody))
			req.Header.Set("Content-Type", "application/json")

			// Test request
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatusCode, resp.StatusCode)

			// Parse response
			var responseBody models.EasyResponse
			err = json.NewDecoder(resp.Body).Decode(&responseBody)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedMessage, responseBody.Message)

			// Verify mocks were called as expected
			mockValidator.AssertExpectations(t)
			mockSchedulerData.AssertExpectations(t)
		})
	}
}
