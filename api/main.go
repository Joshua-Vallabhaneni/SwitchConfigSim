package main

import (
	"encoding/json" // data -> json
	"fmt"
	"net/http" // create web server, http requests
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
	fmt.Println()

	// Register API endpoints
	http.HandleFunc("/switch/config", handleSwitchConfig)
	// when visiting /switch/config, call handleSwitchConfig
	http.HandleFunc("/version", handleVersion)
	http.HandleFunc("/", handleRoot)

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
	// changes SwitchConfig struct to send back to browser
	config := SwitchConfig{
		// fetch from global var
		Hostname:   currentHostname, 
		Interfaces: interfaceStates, 
		Status:     currentStatus,   
	}
	// converts struct to json, layout at top of file has logic for this
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
		oldHostname := currentHostname
		currentHostname = update.Hostname
		changes["hostname"] = oldHostname + " -> " + currentHostname
	}

	// check if user wants to change status
	if update.Status != "" {
		oldStatus := currentStatus
		currentStatus = update.Status
		changes["status"] = oldStatus + " -> " + currentStatus
	}

	// ceck if user wants to change an interface state
	if update.Interface != "" && update.InterfaceState != "" {
		// fetch old state from global var
		oldState := interfaceStates[update.Interface]
		// update global var
		interfaceStates[update.Interface] = update.InterfaceState
		// build response to show changes
		changes["interface_"+update.Interface] = oldState + " -> " + update.InterfaceState
	}

	// send response showing what actually changed
	response := map[string]interface{}{
		"status":  "Config Updated",
		"changes": changes,
		"message": "Configuration changes applied successfully",
	}
	json.NewEncoder(w).Encode(response)
}
