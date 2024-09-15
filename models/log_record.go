package models

type LogRecord struct {
	Before *ChangeData `json:"before"`
	After  ChangeData  `json:"after"`
	Op     string      `json:"op"`
	TsMs   int64       `json:"ts_ms"`
}

type ChangeData struct {
	Key   string     `json:"key"`
	Value ObjectData `json:"value"`
}

type ObjectData struct {
	Type   int        `json:"type"`
	Object ObjectInfo `json:"object"`
}

type ObjectInfo struct {
	ID               string                 `json:"id"`
	Host             string                 `json:"host,omitempty"`
	Name             string                 `json:"name"`
	Path             string                 `json:"path,omitempty"`
	Port             int                    `json:"port,omitempty"`
	Tags             []string               `json:"tags,omitempty"`
	Enabled          bool                   `json:"enabled,omitempty"`
	Retries          int                    `json:"retries,omitempty"`
	Protocol         string                 `json:"protocol,omitempty"`
	CreatedAt        int64                  `json:"created_at"`
	UpdatedAt        int64                  `json:"updated_at"`
	ReadTimeout      int                    `json:"read_timeout,omitempty"`
	WriteTimeout     int                    `json:"write_timeout,omitempty"`
	ConnectTimeout   int                    `json:"connect_timeout,omitempty"`
	LastPing         int64                  `json:"last_ping,omitempty"`
	Version          string                 `json:"version,omitempty"`
	Hostname         string                 `json:"hostname,omitempty"`
	ConfigHash       string                 `json:"config_hash,omitempty"`
	ProcessConf      map[string]interface{} `json:"process_conf,omitempty"`
	ConnectionState  *ConnectionState       `json:"connection_state,omitempty"`
	DataPlaneCertID  string                 `json:"data_plane_cert_id,omitempty"`
	Labels           map[string]string      `json:"labels,omitempty"`
	ResourceType     string                 `json:"resource_type,omitempty"`
	Value            string                 `json:"value,omitempty"`
	HealthChecks     *HealthChecks          `json:"healthchecks,omitempty"`
	HashOn           string                 `json:"hash_on,omitempty"`
	HashFallback     string                 `json:"hash_fallback,omitempty"`
	UseSrvName       bool                   `json:"use_srv_name,omitempty"`
	Slots            int                    `json:"slots,omitempty"`
	Algorithm        string                 `json:"algorithm,omitempty"`
	HashOnCookiePath string                 `json:"hash_on_cookie_path,omitempty"`
}

type ConnectionState struct {
	IsConnected bool `json:"is_connected"`
}

type HealthChecks struct {
	Active    ActiveHealthCheck  `json:"active"`
	Passive   PassiveHealthCheck `json:"passive"`
	Threshold int                `json:"threshold"`
}

type ActiveHealthCheck struct {
	Type            string         `json:"type"`
	Healthy         HealthyCheck   `json:"healthy"`
	Unhealthy       UnhealthyCheck `json:"unhealthy"`
	Timeout         int            `json:"timeout"`
	HTTPPath        string         `json:"http_path"`
	Concurrency     int            `json:"concurrency"`
	HTTPSVerifyCert bool           `json:"https_verify_certificate"`
}

type PassiveHealthCheck struct {
	Type      string         `json:"type"`
	Healthy   HealthyCheck   `json:"healthy"`
	Unhealthy UnhealthyCheck `json:"unhealthy"`
}

type HealthyCheck struct {
	Interval     int   `json:"interval"`
	Successes    int   `json:"successes"`
	HTTPStatuses []int `json:"http_statuses"`
}

type UnhealthyCheck struct {
	Interval     int   `json:"interval"`
	Timeouts     int   `json:"timeouts"`
	TCPFailures  int   `json:"tcp_failures"`
	HTTPFailures int   `json:"http_failures"`
	HTTPStatuses []int `json:"http_statuses"`
}
