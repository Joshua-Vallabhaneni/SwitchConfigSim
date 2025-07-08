package main

import (
	"encoding/json" // data -> json
	"fmt"
	"net/http" // create web server, http requests
	"os"       // reading files from disk
	"os/exec"  // execute shell scripts from Go code
)

// Global vars to store the current switch state
var currentHostname = "switch1"
var currentStatus = "Operational"
var interfaceStates = map[string]string{
	"eth0": "up",
	"eth1": "down",
	"eth2": "up",
	"ethx": "up",
}

// Represents the current state of switch
type SwitchConfig struct {
	Hostname   string            `json:"hostname"`
	Interfaces map[string]string `json:"interfaces"` // show name and state of each port
	Status     string            `json:"status"`
}

// Version info
type Version struct {
	Name      string `json:"name"`
	Version   string `json:"version"`
	BuiltWith string `json:"builtWith"`
}

// Config change request, can change name, status, or interface state
type ConfigUpdate struct {
	// omitempty to skip if empty
	Hostname       string `json:"hostname,omitempty"`
	Status         string `json:"status,omitempty"`
	Interface      string `json:"interface,omitempty"`
	InterfaceState string `json:"interface_state,omitempty"`
}

func main() {
	fmt.Println("SwitchConfigSim REST API Server")
	fmt.Println("Starting server on http://localhost:8080")
	fmt.Println("API Documentation available at http://localhost:8080/docs")
	fmt.Println()

	// Register API endpoints
	http.HandleFunc("/switch/config", handleSwitchConfig)
	// when visiting /switch/config, call handleSwitchConfig
	http.HandleFunc("/version", handleVersion)
	http.HandleFunc("/", handleRoot)

	// Add Swagger UI endpoints for interactive documentation
	http.HandleFunc("/docs", handleSwaggerUI)
	http.HandleFunc("/openapi.yaml", handleOpenAPISpec)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Printf("Server failed to start: %v\n", err)
	}
}

// function to provide basic API info
func handleRoot(w http.ResponseWriter, r *http.Request) {
	// w http.ResponseWriter -> where to send reponse
	// r *http.Request -> info about the req

	// Indicate browser to expect JSON
	w.Header().Set("Content-Type", "application/json")

	// create interface w/api info to send back to browser
	apiInfo := map[string]interface{}{
		"name":        "SwitchConfigSim REST API",
		"description": "NVIDIA network switch management",
		"endpoints": map[string]string{
			"GET /switch/config": "Display switch configuration",
			"PUT /switch/config": "Update switch configuration",
			"GET /version":       "Show version information",
		},
	}
	// convert api info to json, send to brower
	json.NewEncoder(w).Encode(apiInfo)
}

// method to return version info (like version cmd)
func handleVersion(w http.ResponseWriter, r *http.Request) {
	// Only allow GET reqs
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Indicate browser to expect JSON
	w.Header().Set("Content-Type", "application/json")

	version := Version{
		Name:      "SwitchConfigSim",
		Version:   "v1.0.0",
		BuiltWith: "Go REST API",
	}

	json.NewEncoder(w).Encode(version)
}

// function to handle both show and set ops
func handleSwitchConfig(w http.ResponseWriter, r *http.Request) {
	// Indicate browser to expect JSON
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		// like show cmd, wants to read info
		handleShowConfig(w)
	case http.MethodPut:
		// like set cmd, wants to update info
		handleSetConfig(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// function to return current switch config
func handleShowConfig(w http.ResponseWriter) {
	// Call shell script 
	scriptPath := "./shell/get_status.sh"
	cmd := exec.Command(scriptPath)

	// exec script and get output
	output, err := cmd.CombinedOutput()
	if err != nil {
		// If script fails, return error to API user
		http.Error(w, fmt.Sprintf("Error getting system status: %v", err), http.StatusInternalServerError)
		return
	}

	// Return current global state
	fmt.Printf("System status from script:\n%s\n", string(output))

	// Return current config in JSON format
	config := SwitchConfig{
		Hostname:   currentHostname,
		Interfaces: interfaceStates,
		Status:     currentStatus,
	}
	json.NewEncoder(w).Encode(config)
}

// function to simulate config changes
func handleSetConfig(w http.ResponseWriter, r *http.Request) {
	// captures user input
	var update ConfigUpdate

	// Read json data and put into update variable
	// check if json is malformed
	err := json.NewDecoder(r.Body).Decode(&update)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// build response to show changes
	changes := make(map[string]string)

	// check if user wants to change hostname
	if update.Hostname != "" {
		// Call hostname script
		scriptPath := "./shell/set_hostname.sh"
		cmd := exec.Command(scriptPath, update.Hostname)

		// Execute hostname change script
		output, err := cmd.CombinedOutput()
		if err != nil {
			// Script failed
			http.Error(w, fmt.Sprintf("Error setting hostname: %v\nScript output: %s", err, string(output)), http.StatusInternalServerError)
			return
		}

		// Script success and update global var and track change
		oldHostname := currentHostname
		currentHostname = update.Hostname
		changes["hostname"] = oldHostname + " -> " + currentHostname

		// Log script exec for debugging
		fmt.Printf("Hostname script output:\n%s\n", string(output))
	}

	// check if user wants to change status
	if update.Status != "" {
		// Update the variable (don't have script for this yet)
		oldStatus := currentStatus
		currentStatus = update.Status
		changes["status"] = oldStatus + " -> " + currentStatus
	}

	// check if user wants to change an interface state
	if update.Interface != "" && update.InterfaceState != "" {
		// Call interface script instead of just updating variable
		scriptPath := "./shell/set_interface.sh"
		cmd := exec.Command(scriptPath, update.Interface, update.InterfaceState)

		// Execute interface change script
		output, err := cmd.CombinedOutput()
		if err != nil {
			// Script failed - return error to API client
			http.Error(w, fmt.Sprintf("Error setting interface: %v\nScript output: %s", err, string(output)), http.StatusInternalServerError)
			return
		}

		// Script succeeded - update our global variable and track change
		oldState := interfaceStates[update.Interface]
		interfaceStates[update.Interface] = update.InterfaceState
		changes["interface_"+update.Interface] = oldState + " -> " + update.InterfaceState

		// Log script execution for debugging
		fmt.Printf("Interface script output:\n%s\n", string(output))
	}

	// send response showing what actually changed
	response := map[string]interface{}{
		"status":  "Config Updated",
		"changes": changes,
		"message": "Configuration changes applied successfully",
	}
	json.NewEncoder(w).Encode(response)
}

// function for OpenAPI spec file
func handleOpenAPISpec(w http.ResponseWriter, r *http.Request) {
	// Only allow GET requests
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Read the OpenAPI file
	yamlContent, err := os.ReadFile("openapi.yaml")
	if err != nil {
		http.Error(w, "OpenAPI specification not found", http.StatusNotFound)
		return
	}

	// Set headers for YAML
	w.Header().Set("Content-Type", "application/yaml")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write(yamlContent)
}

// func for API docs
func handleSwaggerUI(w http.ResponseWriter, r *http.Request) {
	// Only allow GET requests
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Set content type to HTML
	w.Header().Set("Content-Type", "text/html")

	// Embedded Swagger UI HTML
	swaggerHTML := `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    <meta name="description" content="SwitchConfigSim API Documentation" />
    <title>SwitchConfigSim API Documentation</title>
    <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5.9.0/swagger-ui.css" />
    <style>
        body { margin: 0; }
        .topbar { display: none; }
        .swagger-ui .info { margin: 50px 0; }
        .swagger-ui .info .title { 
            color: #76B900; 
            font-size: 36px;
        }
        .swagger-ui .scheme-container {
            background: #76B900;
            box-shadow: 0 1px 2px 0 rgba(0,0,0,.15);
        }
    </style>
</head>
<body>
    <div id="swagger-ui"></div>
    <script src="https://unpkg.com/swagger-ui-dist@5.9.0/swagger-ui-bundle.js" crossorigin></script>
    <script>
        window.onload = () => {
            window.ui = SwaggerUIBundle({
                url: '/openapi.yaml',
                dom_id: '#swagger-ui',
                deepLinking: true,
                presets: [
                    SwaggerUIBundle.presets.apis,
                    SwaggerUIBundle.presets.standalone,
                ],
                plugins: [
                    SwaggerUIBundle.plugins.DownloadUrl
                ],
                tryItOutEnabled: true,
                supportedSubmitMethods: ['get', 'put', 'post', 'delete', 'options', 'head', 'patch'],
                validatorUrl: null,
                docExpansion: 'list',
                operationsSorter: 'alpha',
                tagsSorter: 'alpha'
            });
        };
    </script>
</body>
</html>`

	// Send the HTML to the browser
	w.Write([]byte(swaggerHTML))
}
