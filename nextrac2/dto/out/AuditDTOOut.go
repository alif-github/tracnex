package out

import "time"

type AuditMonitoringResponse struct {
	ID            int64     `json:"id"`
	TableName     string    `json:"table_name"`
	PrimaryKey    int64     `json:"primary_key"`
	Action        int32     `json:"action"`
	Username      string    `json:"username"`
	CreatedName   string    `json:"created_name"`
	CreatedBy     int64     `json:"created_by"`
	CreatedClient string    `json:"created_client"`
	CreatedAt     time.Time `json:"created_at"`
}

type ViewAuditMonitoringResponse struct {
	ID            int64                  `json:"id"`
	TableName     string                 `json:"table_name"`
	UUIDKey       string                 `json:"uuid_key"`
	Data          map[string]interface{} `json:"data"`
	PrimaryKey    int64                  `json:"primary_key"`
	Action        int32                  `json:"action"`
	CreatedBy     int64                  `json:"created_by"`
	CreatedClient string                 `json:"created_client"`
	CreatedAt     time.Time              `json:"created_at"`
}
