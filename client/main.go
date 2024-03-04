package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
)

const (
	dashboardURLKey = "PLOT_KEY" // Environment variable key for dashboard URL
)

// stats represents the system statistics.
type stats struct {
	CPUUtilizationPercentage float64 `json:"CPUUtilizationPercentage"`
	CPUError                 string  `json:"CPUError"`
	FreeMemory               uint64  `json:"FreeMemory"`
	TotalMemory              uint64  `json:"TotalMemory"`
	MemoryError              string  `json:"MemoryError"`
	TimeStamp                string  `json:"timeStamp"`
}

// getCPUUtilization retrieves the CPU utilization as a percentage.
func getCPUUtilization() (float64, error) {
	percent, err := cpu.Percent(10*time.Millisecond, false)
	if err != nil {
		return 0, fmt.Errorf("failed to get CPU utilization: %v", err)
	}
	return percent[0], nil
}

// getFreeMemory retrieves the amount of free and total memory.
func getFreeMemory() (uint64, uint64, error) {
	memory, err := mem.VirtualMemory()
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get memory info: %v", err)
	}
	return memory.Free, memory.Total, nil
}

// getSystemInfo collects system information and sends it to the dashboard.
func getSystemInfo() {
	cpuUtilizationPercentage, cpuErr := getCPUUtilization()
	freeMemory, totalMemory, memoryErr := getFreeMemory()
	timeStamp := time.Now().UTC().Local().Format("2006-01-02T15:04:05")

	data := stats{
		CPUUtilizationPercentage: cpuUtilizationPercentage,
		FreeMemory:               freeMemory,
		TotalMemory:              totalMemory,
		TimeStamp:                timeStamp,
	}

	if cpuErr != nil {
		data.CPUError = cpuErr.Error()
		log.Printf("Error getting CPU utilization: %v", cpuErr)
	}

	if memoryErr != nil {
		data.MemoryError = memoryErr.Error()
		log.Printf("Error getting memory info: %v", memoryErr)
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Printf("Error marshalling JSON: %v", err)
		return
	}

	fmt.Println(string(jsonData))

	dashboardURL := os.Getenv(dashboardURLKey)
	if dashboardURL != "" {
		go sendToDashboard(dashboardURL, jsonData) // Sending data asynchronously
	} else {
		log.Printf("Error: %s environment variable not found. Data will not be sent to dashboard", dashboardURLKey)
	}
}

// sendToDashboard sends data to the dashboard URL.
func sendToDashboard(url string, jsonData []byte) {
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Error sending data to dashboard: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Error sending data to dashboard. Status code: %d", resp.StatusCode)
	}
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading .env file: %v", err)
	}

	getSystemInfo()
}
