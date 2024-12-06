package clipy

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"golang.design/x/clipboard"
)

func isNumeric(s string) bool {
	for _, v := range s {
		if v < '0' || '9' < v {
			return false
		}
	}
	return true
}

func Execute(dbpath string) error {
	db, err := NewRepository(dbpath)
	if err != nil {
		return err
	}
	err = clipboard.Init()
	if err != nil {
		return err
	}
	changed := clipboard.Watch(context.Background(), clipboard.FmtText)

	for b := range changed {
		db.write(b)
	}

	return nil
}

func StopAllInstances() error {
	mypid := os.Getpid()

	switch runtime.GOOS {
	case "linux":
		if err := stopAllInstancesLinux(mypid); err != nil {
			return err
		}
	case "darwin":
		if err := stopAllInstancesDarwin(mypid); err != nil {
			return err
		}

	default:
		return fmt.Errorf("unsupported operating system: %v", runtime.GOOS)
	}
	return nil
}

func stopAllInstancesDarwin(mypid int) error {
	cmd := exec.Command("ps", "-A")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("error executing ps command: %v", err)
	}

	lines := strings.Split(out.String(), "\n")
	for _, line := range lines {
		if strings.Contains(line, "clipy") {
			fields := strings.Fields(line)
			if len(fields) > 0 {
				pid, _ := strconv.Atoi(fields[0])

				if mypid == pid {
					continue
				}

				err = killProcess(pid)
				if err != nil {
					return fmt.Errorf("failed to kill process %d: %v", pid, err)
				}
			}
		}
	}
	return nil
}

func stopAllInstancesLinux(mypid int) error {
	dirs, err := os.ReadDir("/proc")
	if err != nil {
		return err
	}
	for _, dir := range dirs {
		if !dir.IsDir() || !isNumeric(dir.Name()) || strconv.Itoa(mypid) == dir.Name() {
			continue
		}

		pid := dir.Name()
		comm := filepath.Join("/proc", pid, "comm")
		f, err := os.Open(comm)
		if err != nil {
			return err
		}
		prs, err := io.ReadAll(f)
		if err != nil {
			return err
		}

		if strings.Contains(strings.TrimSpace(string(prs)), "clipy") {
			targetPid, err := strconv.Atoi(pid)
			if err != nil {
				return err
			}
			err = killProcess(targetPid)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func killProcess(pid int) error {
	process, err := os.FindProcess(pid)
	if err != nil {
		return err
	}
	err = process.Signal(os.Kill)
	if err != nil {
		return err
	}
	return nil
}
