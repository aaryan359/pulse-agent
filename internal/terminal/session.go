package terminal

import (
	"os"
	"os/exec"
	"syscall"

	"github.com/creack/pty"
)

type Session struct {
	Cmd *exec.Cmd
	Pty *os.File
}

func StartShell() (*Session, error) {
	shell := "/bin/bash"
	if _, err := os.Stat(shell); err != nil {
		shell = "/bin/sh"
	}

	cmd := exec.Command(shell)
	cmd.Env = os.Environ()
	cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}

	ptmx, err := pty.Start(cmd)
	if err != nil {
		return nil, err
	}

	return &Session{Cmd: cmd, Pty: ptmx}, nil
}

func (s *Session) Write(data []byte) error {
	_, err := s.Pty.Write(data)
	return err
}

func (s *Session) ReadLoop(fn func([]byte)) {
	buf := make([]byte, 4096)
	for {
		n, err := s.Pty.Read(buf)
		if err != nil {
			return
		}
		fn(buf[:n])
	}
}

func (s *Session) Close() {
	if s.Pty != nil {
		_ = s.Pty.Close()
	}
	if s.Cmd != nil && s.Cmd.Process != nil {
		_ = s.Cmd.Process.Kill()
	}
}
