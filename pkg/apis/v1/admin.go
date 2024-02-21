package v1

import (
	"fmt"
	"net/url"
	"reflect"
	"strings"
)

// AdminConfig holds the administrative configuration.
type AdminConfig struct {
	// Name of the server. Default is "default-stunnerd".
	Name string `json:"name,omitempty"`
	// LogLevel is the desired log verbosity, e.g.: "stunner:TRACE,all:INFO". Default is
	// "all:INFO".
	LogLevel string `json:"loglevel,omitempty"`
	// MetricsEndpoint is the URI in the form `http://address:port/path` at which HTTP metric
	// requests are served. The scheme (`http://`") is mandatory. Default is to expose no
	// metric endpoints.
	MetricsEndpoint string `json:"metrics_endpoint,omitempty"`
	// HealthCheckEndpoint is the URI of the form `http://address:port` exposed for external
	// HTTP health-checking. A liveness probe responder will be exposed on path `/live` and
	// readiness probe on path `/ready`. The scheme (`http://`) is mandatory, and if no port is
	// specified then the default port is 8086. If ignored, then the default is to enable
	// health-checking at `http://0.0.0.0:8086`. Set to a pointer to an empty string to disable
	// health-checking.
	HealthCheckEndpoint *string `json:"healthcheck_endpoint,omitempty"`
}

// Validate checks a configuration and injects defaults.
func (req *AdminConfig) Validate() error {
	if req.LogLevel == "" {
		req.LogLevel = DefaultLogLevel
	}

	if req.Name == "" {
		req.Name = DefaultStunnerName
	}

	if req.MetricsEndpoint != "" {
		//Metrics endpoint set: validate. The empty string is valid
		if _, err := url.Parse(req.MetricsEndpoint); err != nil {
			return fmt.Errorf("invalid metric server endpoint URL %s: %s",
				req.MetricsEndpoint, err.Error())
		}
	}

	if req.HealthCheckEndpoint == nil {
		// No healtchcheck endpoint given: use default URL
		e := fmt.Sprintf("http://:%d", DefaultHealthCheckPort)
		req.HealthCheckEndpoint = &e
	} else {
		// Healtcheck endpoint set: validate. Empty string is valid
		if _, err := url.Parse(*req.HealthCheckEndpoint); err != nil {
			return fmt.Errorf("invalid health-check server endpoint URL %s: %s",
				*req.HealthCheckEndpoint, err.Error())
		}
	}

	return nil
}

// Name returns the name of the object to be configured.
func (req *AdminConfig) ConfigName() string {
	// Singleton!
	return DefaultAdminName
}

// DeepEqual compares two configurations.
func (req *AdminConfig) DeepEqual(other Config) bool {
	return reflect.DeepEqual(req, other)
}

// DeepCopyInto copies a configuration.
func (req *AdminConfig) DeepCopyInto(dst Config) {
	ret := dst.(*AdminConfig)
	*ret = *req
}

// String stringifies the configuration.
func (req *AdminConfig) String() string {
	status := []string{}
	if req.Name != "" {
		status = append(status, fmt.Sprintf("name=%q", req.Name))
	}
	if req.LogLevel != "" {
		status = append(status, fmt.Sprintf("logLevel=%q", req.LogLevel))
	}
	if req.MetricsEndpoint != "" {
		status = append(status, fmt.Sprintf("metrics=%q", req.MetricsEndpoint))
	}
	if req.HealthCheckEndpoint != nil {
		status = append(status, fmt.Sprintf("health-check=%q", *req.HealthCheckEndpoint))
	}
	return fmt.Sprintf("admin:{%s}", strings.Join(status, ","))
}

// AdminStatus represents the administrative status.
type AdminStatus struct {
	Name                string `json:"name,omitempty"`
	LogLevel            string `json:"loglevel,omitempty"`
	MetricsEndpoint     string `json:"metrics_endpoint,omitempty"`
	HealthCheckEndpoint string `json:"healthcheck_endpoint,omitempty"`
	// licencing status comes here
}

// String returns a string reprsentation of the administrative status.
func (a *AdminStatus) String() string {
	status := []string{}
	if a.LogLevel != "" {
		status = append(status, fmt.Sprintf("logLevel=%q", a.LogLevel))
	}
	if a.MetricsEndpoint != "" {
		status = append(status, fmt.Sprintf("metrics=%q", a.MetricsEndpoint))
	}
	if a.HealthCheckEndpoint != "" {
		status = append(status, fmt.Sprintf("health-check=%q", a.HealthCheckEndpoint))
	}

	// add licencing status here

	return fmt.Sprintf("%s:{%s}", a.Name, strings.Join(status, ","))
}
