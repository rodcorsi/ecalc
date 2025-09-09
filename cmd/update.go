package main

import (
	"io"
	"net/http"
	"os"
	"os/exec"
	"runtime"

	"github.com/abiosoft/ishell"
)

func update(c *ishell.Context) {
	c.Println("Starting update...")

	var binName string
	switch runtime.GOOS {
	case "linux":
		binName = "ecalc"
	case "windows":
		if runtime.GOARCH == "386" {
			binName = "ecalc32.exe"
		} else {
			binName = "ecalc.exe"
		}
	default:
		c.Printf("Update not supported for this OS: %s", runtime.GOOS)
		return
	}

	url := "https://github.com/rodcorsi/ecalc/releases/download/latest/" + binName
	c.Printf("Downloading from %s...", url)

	resp, err := http.Get(url)
	if err != nil {
		c.Printf("Error downloading update: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		c.Printf("Error downloading update: status %s%s", resp.Status, body)
		return
	}

	exePath, err := os.Executable()
	if err != nil {
		c.Printf("Error getting executable path: %v", err)
		return
	}

	newExePath := exePath + ".new"
	newExe, err := os.Create(newExePath)
	if err != nil {
		c.Printf("Error creating temporary file: %v", err)
		return
	}

	_, err = io.Copy(newExe, resp.Body)
	newExe.Close() // Close after writing
	if err != nil {
		c.Printf("Error writing downloaded content: %v", err)
		os.Remove(newExePath)
		return
	}

	// Make the new file executable (not needed on Windows)
	if runtime.GOOS != "windows" {
		if err := os.Chmod(newExePath, 0755); err != nil {
			c.Printf("Warning: could not make the new executable file executable: %v", err)
		}
	}

	oldExePath := exePath + ".old"
	os.Remove(oldExePath) // remove old backup if it exists

	if err := os.Rename(exePath, oldExePath); err != nil {
		c.Printf("Error moving current executable: %v", err)
		os.Remove(newExePath)
		return
	}

	if err := os.Rename(newExePath, exePath); err != nil {
		c.Printf("Error replacing executable with new version: %v", err)
		// Try to restore
		if err := os.Rename(oldExePath, exePath); err != nil {
			c.Printf("CRITICAL: Failed to restore old executable. Please do it manually from %s", oldExePath)
		}
		return
	}

	c.Println("Update successful! Restarting...")
	c.Printf("The old version is saved as %s. You can remove it if the new version works correctly.\n", oldExePath)

	cmd := exec.Command(exePath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Start(); err != nil {
		c.Printf("Failed to restart application: %v. Please restart manually.\n", err)
	}

	c.Stop()
}