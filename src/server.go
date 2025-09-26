package src

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"strings"
	"time"
	"yotudo/src/lib/logger"
	"yotudo/src/model"
)

type Server struct {
	cmd    *exec.Cmd
	Config model.ServerConfig
	stdin  io.WriteCloser
}

func (s *Server) Start(ctx context.Context, config *model.ServerConfig) error {
	args := make([]string, 1)
	args[0] = "--host-frontend"

	s.Config.Merge(config)

	if s.Config.IsHosted {
		args = append(args, "--host")
	} else {
		args = append(args, "--localhost")
	}

	if s.Config.Port != 0 {
		args = append(args, "--port", fmt.Sprintf("%d", s.Config.Port))
	}

	serverCmd := exec.CommandContext(context.Background(), "./yotudo-server.exe", args...)

	stdout, err := serverCmd.StdoutPipe()
	if err != nil {
		logger.Error(err)
		return err
	}
	stdin, err := serverCmd.StdinPipe()
	if err != nil {
		logger.Error(err)
		return err
	}
	s.stdin = stdin

	if err := serverCmd.Start(); err != nil {
		logger.Error(err)
		return err
	}

	s.cmd = serverCmd

	errorTrigger := make(chan string)

	go func() {
		buf := make([]byte, 4096)

		for {
			if serverCmd == nil {
				return
			}

			read, err := stdout.Read(buf)
			if err != nil {
				logger.Warning(err)
				return
			}

			if read != 0 {
				var prefix string
				out := string(buf[:read])

				if indexOfColon := strings.Index(out, ":"); indexOfColon != -1 {
					prefix = out[:indexOfColon]
					out = strings.TrimSpace(out[indexOfColon+1:])
				}

				out = strings.ReplaceAll(out, "\n", "")

				switch prefix {
				case "INFO":
					logger.Info(out)
				case "ERR":
					logger.Error(out)
					errorTrigger <- out
					s.Stop()

					return
				case "WARN":
					logger.Warning(out)
				case "EVENT":
					logger.Debug(prefix, "|", out)
					s.processEvent(strings.Split(out, "="))
				default:
					logger.Debug(out)
				}
			}
		}
	}()

	logger.Debug("Waiting for server start and making sure its running by waiting a bit")

	select {
	case err := <-errorTrigger:
		return errors.New(err)
	case <-time.After(time.Second * 2):
	}

	return nil
}

func (s *Server) Stop() error {
	if s.cmd == nil {
		return nil
	}

	s.Config.ListeningOn = ""
	cmd := s.cmd
	s.cmd = nil

	fmt.Fprintln(s.stdin, "exit")
	defer s.stdin.Close()

	if err := cmd.Wait(); err != nil {
		logger.Debug("Error in \"cmd.Wait()\"", err)
		return err
	}

	return nil
}

func (s *Server) processEvent(keyValuePair []string) {
	if len(keyValuePair) != 2 {
		return
	}

	value := keyValuePair[1]

	switch keyValuePair[0] {
	case "listening":
		s.Config.ListeningOn = value
	case "close":
		s.Stop()
	}
}
