package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Toggle to true if sudo is needed to control services
const useSudo = true

// List your systemd service names here (omit .service)
var services = []string{
	"nginx",
	"mysql",        // or "mariadb" depending on your system
	"php8.2-fpm",
}

// Checks if a service is known to systemd by querying list-unit-files
func serviceExists(service string) bool {
	unitName := service
	if !strings.HasSuffix(unitName, ".service") {
		unitName += ".service"
	}

	var args []string
	if useSudo {
		args = []string{"sudo", "systemctl", "list-unit-files", unitName}
	} else {
		args = []string{"systemctl", "list-unit-files", unitName}
	}

	cmd := exec.Command(args[0], args[1:]...)
	output, err := cmd.Output()
	if err != nil {
		return false
	}

	return strings.Contains(string(output), unitName)
}

// Runs the given action (start, stop, etc.) on the service and returns output
func runAction(action, service string) (string, error) {
	var args []string
	if useSudo {
		args = []string{"sudo", "systemctl", action, service}
	} else {
		args = []string{"systemctl", action, service}
	}

	cmd := exec.Command(args[0], args[1:]...)
	output, err := cmd.CombinedOutput()
	return string(output), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: svcctl <start|stop|restart|status>")
		os.Exit(1)
	}

	action := os.Args[1]
	validActions := map[string]bool{
		"start":   true,
		"stop":    true,
		"restart": true,
		"status":  true,
	}

	if !validActions[action] {
		fmt.Printf("❌ Invalid action: %s\n", action)
		fmt.Println("Allowed actions: start, stop, restart, status")
		os.Exit(1)
	}

	fmt.Printf("▶️  Performing '%s' on services...\n\n", action)

	for _, service := range services {
		fmt.Printf("🔍 Checking service: %s\n", service)

		if !serviceExists(service) {
			fmt.Printf("⚠️  Service '%s' not recognized by systemd.\n\n", service)
			continue
		}

		output, err := runAction(action, service)
		if err != nil {
			fmt.Printf("❌ Failed to %s %s: %v\nOutput:\n%s\n\n", action, service, err, output)
		} else {
			fmt.Printf("✅ Successfully %sd %s\n", action, service)
			if action == "status" {
				fmt.Println("--- Status Output ---")
				fmt.Println(strings.TrimSpace(output))
				fmt.Println("---------------------")
			}
			fmt.Println()
		}
	}

	fmt.Println("✅ Done.")
}
