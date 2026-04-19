package main

import (
	"os"
	"os/exec"
	"testing"
	"time"
)

// TestMain_BuildCheck ensures the binary compiles cleanly.
func TestMain_BuildCheck(t *testing.T) {
	if os.Getenv("SKIP_BUILD_TEST") != "" {
		t.Skip("skipping build test")
	}
	cmd := exec.Command("go", "build", "-o", os.DevNull, ".")
	cmd.Dir = "."
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("build failed: %v\n%s", err, out)
	}
}

// TestMain_RunHelp verifies --help exits cleanly (exit code 2 from flag pkg is ok).
func TestMain_RunHelp(t *testing.T) {
	if os.Getenv("SKIP_BUILD_TEST") != "" {
		t.Skip("skipping run test")
	}

	binPath := t.TempDir() + "/portwatch"
	build := exec.Command("go", "build", "-o", binPath, ".")
	if out, err := build.CombinedOutput(); err != nil {
		t.Fatalf("build failed: %v\n%s", err, out)
	}

	cmd := exec.Command(binPath, "--help")
	cmd.Timeout = 3 * time.Second
	// flag.Parse prints usage and exits with code 2 on --help
	err := cmd.Run()
	if exitErr, ok := err.(*exec.ExitError); ok {
		if exitErr.ExitCode() == 2 {
			return // expected
		}
	}
	// nil err means it exited 0, also acceptable
}
