package about

import "time"

const (
	Application_Title                  string        = "RAFT cluster test application"
	Application_Version                string        = "v0.1.3-beta"
	Application_Configuration_File     string        = "app.config"
	Application_Health_Url             string        = "/actuator/health"
	Application_Sleep_Timeout          time.Duration = 100 * time.Microsecond
	Application_Cancel_Context_Timeout time.Duration = 500 * time.Microsecond
	Database_Directory                 string        = "db"
	Database_Namespace_Root            string        = "Database"
	Configurations_Directory           string        = "conf"
	Configurations_Extension           string        = ".conf"
	SSL_Directory                      string        = "ssl"
	Cluster_State_Directory            string        = "state"
	Protocol_with_out_SSL              string        = "http://"
	Protocol_with_SSL                  string        = "https://"
	API_Context                        string        = "api/service"
)
