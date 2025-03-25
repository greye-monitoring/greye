package models

import (
	"errors"
	"github.com/stretchr/testify/assert"
	annotations "greye/pkg/annotations/domain/models"
	models2 "greye/pkg/authentication/domain/models"
	modelsHttp "greye/pkg/client/domain/models"
	"greye/pkg/config/domain/models"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
	"time"
)

func TestGetPortUsed(t *testing.T) {
	tests := []struct {
		name         string
		service      *v1.Service
		expectedPort int
	}{
		{
			name: "Valid service with port",
			service: &v1.Service{
				Spec: v1.ServiceSpec{
					Ports: []v1.ServicePort{
						{
							Port: 8080, // The port exposed by the service
						},
					},
				},
			},
			expectedPort: 8080,
		},
		{
			name: "Service with annotation",
			service: &v1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Annotations: map[string]string{
						annotations.Port: "9090",
					},
				},
				Spec: v1.ServiceSpec{
					Ports: []v1.ServicePort{
						{
							Port: 8080, // The port exposed by the service
						},
					},
				},
			},
			expectedPort: 9090,
		},
		{
			name: "Service with multiple ports",
			service: &v1.Service{
				Spec: v1.ServiceSpec{
					Ports: []v1.ServicePort{
						{
							Port: 8080, // The port exposed by the service
						},
						{
							Port: 9080, // The port exposed by the service
						},
					},
				},
			},
			expectedPort: 8080,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetPortUsed(tt.service)
			if result != tt.expectedPort {
				t.Errorf("GetPortUsed() = %v, want %v", result, tt.expectedPort)
			}
		})
	}

}

func TestNewSchedulerApplicationFromService(t *testing.T) {
	tests := []struct {
		name                      string
		service                   *v1.Service
		defaultValue              *models.Application
		expectedInterval          time.Duration
		expectedProtocol          string
		expectedTimeout           time.Duration
		expectedMaxFails          int
		expectedMonitor           string
		expectedAuth              models2.AuthenticationData
		expectStopMonitoringUntil time.Time
		expectError               bool
	}{
		{
			name: "Valid service with annotations",
			service: &v1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Annotations: map[string]string{
						annotations.IntervalSeconds:        "60",
						annotations.Protocol:               "http",
						annotations.TimeoutSeconds:         "5",
						annotations.MaxFailedRequests:      "3",
						annotations.ForcePodMonitor:        "pod-monitor-instance",
						annotations.Paths:                  "/path1\n/path2",
						annotations.AuthenticationMethod:   "basic",
						annotations.AuthenticationUsername: "user",
						annotations.AuthenticationPassword: "password",
						annotations.StopMonitoringUntil:    "2100-01-01T00:00:00",
					},
				},
				Spec: v1.ServiceSpec{
					Ports: []v1.ServicePort{
						{
							Port: 8080, // The port exposed by the service
						},
					},
				},
			},
			defaultValue: &models.Application{
				IntervalSeconds:   30,
				Protocol:          "https",
				Method:            "GET",
				MaxFailedRequests: 2,
				TimeoutSeconds:    10,
			},
			expectedInterval: 60 * time.Second,
			expectedProtocol: "http",
			expectedTimeout:  5 * time.Second,
			expectedMaxFails: 3,
			expectedMonitor:  "pod-monitor-instance",
			expectedAuth: models2.AuthenticationData{
				Method:   "basic",
				Username: "user",
				Password: "password",
			},
			expectStopMonitoringUntil: time.Date(2100, time.January, 1, 0, 0, 0, 0, time.UTC),
			expectError:               false,
		},
		{
			name: "Service missing some annotations (defaults applied)",
			service: &v1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Annotations: map[string]string{
						annotations.IntervalSeconds: "120",  // Valid annotation
						annotations.Protocol:        "http", // Valid annotation
					},
				},
				Spec: v1.ServiceSpec{
					Ports: []v1.ServicePort{
						{
							Port: 8080, // The port exposed by the service
						},
					},
				},
			},
			defaultValue: &models.Application{
				IntervalSeconds:   30,
				Protocol:          "https",
				Method:            "GET",
				MaxFailedRequests: 2,
				TimeoutSeconds:    10,
			},
			expectedInterval:          120 * time.Second,
			expectedProtocol:          "http",
			expectedTimeout:           10 * time.Second,             // Default value should apply
			expectedMaxFails:          2,                            // Default value should apply
			expectedMonitor:           "",                           // Default value should apply
			expectedAuth:              models2.AuthenticationData{}, // Default value should apply
			expectStopMonitoringUntil: time.Now(),
			expectError:               false,
		},
		{
			name: "Service with invalid max failed requests annotation",
			service: &v1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Annotations: map[string]string{
						annotations.IntervalSeconds:   "60",
						annotations.Protocol:          "http",
						annotations.TimeoutSeconds:    "5",
						annotations.MaxFailedRequests: "invalid", // Invalid value
					},
				},
				Spec: v1.ServiceSpec{
					Ports: []v1.ServicePort{
						{
							Port: 8080, // The port exposed by the service
						},
					},
				},
			},
			defaultValue: &models.Application{
				IntervalSeconds:   30,
				Protocol:          "https",
				Method:            "GET",
				MaxFailedRequests: 2,
				TimeoutSeconds:    10,
			},
			expectedInterval:          60 * time.Second,
			expectedProtocol:          "http",
			expectedTimeout:           5 * time.Second,
			expectedMaxFails:          2,  // Default value should apply due to invalid annotation
			expectedMonitor:           "", // Default value should apply
			expectedAuth:              models2.AuthenticationData{},
			expectStopMonitoringUntil: time.Now(),
			expectError:               false,
		},
		{
			name: "Service with invalid timeout annotation",
			service: &v1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Annotations: map[string]string{
						annotations.IntervalSeconds: "60",
						annotations.Protocol:        "http",
						annotations.TimeoutSeconds:  "invalid", // Invalid value
					},
				},
				Spec: v1.ServiceSpec{
					Ports: []v1.ServicePort{
						{
							Port: 8080, // The port exposed by the service
						},
					},
				},
			},
			defaultValue: &models.Application{
				IntervalSeconds:   30,
				Protocol:          "https",
				Method:            "GET",
				MaxFailedRequests: 2,
				TimeoutSeconds:    10,
			},
			expectedInterval:          60 * time.Second,
			expectedProtocol:          "http",
			expectedTimeout:           10 * time.Second, // Default value should apply
			expectedMaxFails:          2,                // Default value should apply
			expectedMonitor:           "",               // Default value should apply
			expectedAuth:              models2.AuthenticationData{},
			expectStopMonitoringUntil: time.Now(),
			expectError:               false,
		},

		{
			name: "Service without interval annotation",
			service: &v1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Annotations: map[string]string{
						annotations.Protocol: "http",
					},
				},
				Spec: v1.ServiceSpec{
					Ports: []v1.ServicePort{
						{
							Port: 8080, // The port exposed by the service
						},
					},
				},
			},
			defaultValue: &models.Application{
				IntervalSeconds:   60,
				Protocol:          "https",
				Method:            "GET",
				MaxFailedRequests: 2,
				TimeoutSeconds:    10,
			},
			expectedInterval:          60 * time.Second,
			expectedProtocol:          "http",
			expectedTimeout:           10 * time.Second, // Default value should apply
			expectedMaxFails:          2,                // Default value should apply
			expectedMonitor:           "",               // Default value should apply
			expectedAuth:              models2.AuthenticationData{},
			expectStopMonitoringUntil: time.Now(),
			expectError:               false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NewSchedulerApplicationFromService(tt.service, tt.defaultValue)

			// Check if result matches expected values
			assert.Equal(t, tt.expectedInterval, result.Job.Interval)
			assert.Equal(t, tt.expectedProtocol, result.MonitoringHttpRequest.Protocol)
			assert.Equal(t, tt.expectedTimeout, result.MonitoringHttpRequest.Timeout)
			assert.Equal(t, tt.expectedMaxFails, result.MaxFailRequests)
			assert.Equal(t, tt.expectedMonitor, result.ForcePodMonitorInstance)
			assert.Equal(t, tt.expectedAuth, result.MonitoringHttpRequest.Authentication)
			assert.Equal(t, tt.expectStopMonitoringUntil.Day(), result.StopMonitoringUntil.Day())
			assert.Equal(t, tt.expectStopMonitoringUntil.Month(), result.StopMonitoringUntil.Month())
			assert.Equal(t, tt.expectStopMonitoringUntil.Year(), result.StopMonitoringUntil.Year())
		})
	}
}

func TestGenerateJobSchedulerApplication(t *testing.T) {
	tests := []struct {
		name     string
		input    SchedulerApplication
		expected time.Duration
	}{
		{
			name: "Valid Host",
			input: SchedulerApplication{
				MonitoringHttpRequest: modelsHttp.MonitoringHttpRequest{
					Interval: 60 * time.Second,
				},
			},
			expected: 60 * time.Second,
		},
		{
			name: "Empty Host",
			input: SchedulerApplication{
				MonitoringHttpRequest: modelsHttp.MonitoringHttpRequest{
					Interval: 30 * time.Second,
				},
			},
			expected: 30 * time.Second,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GenerateJobSchedulerApplication(tt.input)
			if got.Job.Ticker == nil || got.Job.Quit == nil || got.Job.Interval != tt.expected {
				t.Errorf("GenerateJobSchedulerApplication() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name    string
		input   SchedulerApplication
		wantErr error
	}{
		{
			name: "Valid input",
			input: SchedulerApplication{
				MaxFailRequests: 1,
				MonitoringHttpRequest: modelsHttp.MonitoringHttpRequest{
					Protocol: "http",
					Interval: 60,
					Timeout:  5,
				},
			},
			wantErr: nil,
		},
		{
			name: "Invalid MaxFailRequests",
			input: SchedulerApplication{
				MaxFailRequests: 0,
				MonitoringHttpRequest: modelsHttp.MonitoringHttpRequest{
					Protocol: "http",
					Timeout:  5,
				},
			},
			wantErr: errors.New("MaxFailedRequests must be greater than 0"),
		},
		{
			name: "Invalid Protocol",
			input: SchedulerApplication{
				MaxFailRequests: 1,
				MonitoringHttpRequest: modelsHttp.MonitoringHttpRequest{
					Protocol: "",
					Timeout:  5,
				},
			},
			wantErr: errors.New("Protocol must be provided"),
		},
		{
			name: "Invalid Timeout",
			input: SchedulerApplication{
				MaxFailRequests: 1,
				MonitoringHttpRequest: modelsHttp.MonitoringHttpRequest{
					Protocol: "http",
					Timeout:  0,
				},
			},
			wantErr: errors.New("TimeoutSeconds must be greater or equal to 1 second"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.input.Validate()
			if err != nil && err.Error() != tt.wantErr.Error() {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetSvcHostname(t *testing.T) {
	tests := []struct {
		name     string
		input    SchedulerApplication
		expected string
	}{
		{
			name: "Valid Host",
			input: SchedulerApplication{
				MonitoringHttpRequest: modelsHttp.MonitoringHttpRequest{
					Host: "localhost",
				},
			},
			expected: "localhost",
		},
		{
			name: "Empty Host",
			input: SchedulerApplication{
				MonitoringHttpRequest: modelsHttp.MonitoringHttpRequest{
					Host: "",
				},
			},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.input.GetSvcHostname()
			if got != tt.expected {
				t.Errorf("GetSvcHostname() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestAddPortToForcePodMonitorInstanceIfMissing(t *testing.T) {
	tests := []struct {
		name     string
		input    SchedulerApplication
		expected string
	}{
		{
			name: "No ForcePodMonitorInstance",
			input: SchedulerApplication{
				ForcePodMonitorInstance: "",
				MonitoringHttpRequest: modelsHttp.MonitoringHttpRequest{
					Port: 8080,
				}},
			expected: "",
		},
		{
			name: "ForcePodMonitorInstance without port",
			input: SchedulerApplication{
				ForcePodMonitorInstance: "example.com",
				MonitoringHttpRequest: modelsHttp.MonitoringHttpRequest{
					Port: 8080,
				}},
			expected: "example.com:8080",
		},
		{
			name: "ForcePodMonitorInstance with port",
			input: SchedulerApplication{
				ForcePodMonitorInstance: "example.com:9090",
				MonitoringHttpRequest: modelsHttp.MonitoringHttpRequest{
					Port: 8080,
				},
			},
			expected: "example.com:9090",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.input.AddPortToForcePodMonitorInstanceIfMissing()
			if got != tt.expected {
				t.Errorf("AddPortToForcePodMonitorInstanceIfMissing() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestSchedulerApplication_AddPortToForcePodMonitorInstanceIfMissing(t *testing.T) {
	tests := []struct {
		name                    string
		forcePodMonitorInstance string
		port                    int
		want                    string
	}{
		{
			name:                    "Empty instance name",
			forcePodMonitorInstance: "",
			port:                    8080,
			want:                    "",
		},
		{
			name:                    "Instance name without port",
			forcePodMonitorInstance: "my-pod-123",
			port:                    8080,
			want:                    "my-pod-123:8080",
		},
		{
			name:                    "Instance name with port",
			forcePodMonitorInstance: "my-pod-123:9090",
			port:                    8080,
			want:                    "my-pod-123:9090",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := SchedulerApplication{
				ForcePodMonitorInstance: tt.forcePodMonitorInstance,
				MonitoringHttpRequest: modelsHttp.MonitoringHttpRequest{
					Port: tt.port,
				},
			}

			got := s.AddPortToForcePodMonitorInstanceIfMissing()
			if got != tt.want {
				t.Errorf("AddPortToForcePodMonitorInstanceIfMissing() = %v, want %v", got, tt.want)
			}
		})
	}
}
