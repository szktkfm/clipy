package clipy

import (
	"context"
	"io"
	"log"
	"os"
	"path/filepath"
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

func Execute() {
	db := NewRepository("./test.db")
	err := clipboard.Init()
	if err != nil {
		panic(err)
	}
	changed := clipboard.Watch(context.Background(), clipboard.FmtText)

	for b := range changed {
		db.write(b)
	}
}

func KillClipies() {
	dirs, err := os.ReadDir("/proc")
	mypid := os.Getpid()
	if err != nil {
		log.Fatal("error")
	}

	for _, dir := range dirs {
		if !dir.IsDir() || !isNumeric(dir.Name()) || strconv.Itoa(mypid) == dir.Name() {
			continue
		}

		pid := dir.Name()
		comm := filepath.Join("/proc", pid, "comm")
		f, err := os.Open(comm)
		if err != nil {
			log.Fatal("error")
		}
		prs, err := io.ReadAll(f)
		if err != nil {
			log.Fatal("error")
		}

		if strings.Contains(strings.TrimSpace(string(prs)), "clipy") {
			targetPid, err := strconv.Atoi(pid)
			if err != nil {
				log.Fatal("error")
			}

			process, err := os.FindProcess(targetPid)
			if err != nil {
				log.Fatal("error")
			}
			process.Signal(os.Kill)
		}
	}
}
