package main

import (
	"os"
	"os/exec"
	"testing"
	"time"
)

func TestForkExec(t *testing.T) {
	// Issue 38824: importing the plugin package causes it hang in forkExec on darwin.

	t.Parallel()
	buildCmd := exec.Command("go", "build", "-o", "forkexec.exe", "./forkexec/main.go")
	buildCmd.Env = append(os.Environ(), "CGO_ENABLED=1")
	out, err := buildCmd.Output()
	if err != nil {
		t.Fatalf("build failed %s, out: %s", err, out)
	}

	var cmd *exec.Cmd
	done := make(chan int, 1)

	go func() {
		for i := 0; i < 10; i++ {
			cmd = exec.Command("./forkexec.exe", "1")
			out, err := cmd.Output()
			if err != nil {
				t.Errorf("running command failed: %v %s", err, string(out))
				break
			}
		}
		done <- 1
	}()
	select {
	case <-done:
	case <-time.After(5 * time.Minute):
		cmd.Process.Kill()
		t.Fatalf("subprocess hang")
	}
}
