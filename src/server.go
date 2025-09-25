package src

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"net"
	"net/http"
	"strconv"
	"strings"
	"yotudo/src/handler"
	"yotudo/src/lib/logger"
	"yotudo/src/settings"
)

// "localhost" | "0.0.0.0"
const listenerIp string = "localhost"

type Server struct {
	httpServer *http.Server
	running    bool
	addr       string
	port       int
}

func NewServer(webAppFS fs.FS) *Server {
	return &Server{
		httpServer: &http.Server{
			Handler: getHandler(webAppFS),
		},
		port: -1,
	}
}

func (s *Server) Start() error {
	if s.running {
		return fmt.Errorf("server is already running")
	}

	s.running = true

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", listenerIp, settings.Global.Server.Port))
	if err != nil {
		return err
	}

	if s.addr, err = s.getLocalIp(); err != nil {
		return err
	}

	if s.port, err = s.getPortFormListener(listener); err != nil {
		return err
	}

	go func() {
		if err := s.httpServer.Serve(listener); err != http.ErrServerClosed {
			logger.Error(err)
		}
	}()

	logger.InfoF("Launching webserver on: http://%s", s.GetAddress())

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	if !s.running {
		return fmt.Errorf("server is not running")
	}

	s.running = false
	s.addr = ""
	s.port = -1

	logger.Info("Closing webserver")

	return s.httpServer.Shutdown(ctx)
}

func (s *Server) IsRunning() bool {
	return s.running
}

func (s *Server) GetAddress() string {
	return fmt.Sprintf("%s:%d", s.addr, s.port)
}

func (s *Server) getLocalIp() (string, error) {
	if listenerIp == "localhost" {
		return listenerIp, nil
	}

	conn, err := net.Dial("udp", "192.168.0.0:80")
	if err != nil {
		return "", fmt.Errorf("server couldn't acquire host IP")
	}

	defer conn.Close()

	localAddress := strings.Split(conn.LocalAddr().String(), ":")
	if len(localAddress) == 0 {
		return "", fmt.Errorf("server couldn't process host IP")
	}

	return localAddress[0], nil
}

func (s *Server) getPortFormListener(listener net.Listener) (int, error) {
	if settings.Global.Server.Port != 0 {
		return settings.Global.Server.Port, nil
	}

	_addr := listener.Addr().String()
	portBeginningIndex := strings.LastIndex(_addr, ":")

	if portBeginningIndex == -1 || len(_addr)-1 == portBeginningIndex {
		return -1, fmt.Errorf("server couldn't acquire host port")
	}

	return strconv.Atoi(_addr[portBeginningIndex+1:])
}

func getHandler(webAppFS fs.FS) http.Handler {
	mainHandler := http.NewServeMux()

	mainHandler.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		switch getPrefix(r.URL.Path) {
		case "/api":
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "*")
			w.Header().Set("Access-Control-Allow-Headers", "*")

			handler.ApiHandlers().ServeHTTP(w, r)
		case "/image":
			handler.NewAssetsHandler().ServeHTTP(w, r)
		default:
			staticFs, err := fs.Sub(webAppFS, "web-app/dist")
			if err != nil {
				logger.Error("Server Error:", err)
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintln(w, "Szerver hiba történt")
				return
			}
			file, err := staticFs.Open(r.URL.Path[1:])
			if err != nil {
				file, err = staticFs.Open("index.html")
				if err != nil {
					logger.Warning("Server Warning:", err)
					w.WriteHeader(http.StatusInternalServerError)
					break
				}

			}

			defer file.Close()

			var contentType string = "text/html"
			if lastIndexOfDot := strings.LastIndex(r.URL.Path, "."); lastIndexOfDot != -1 {
				switch r.URL.Path[lastIndexOfDot:] {
				case ".js":
					contentType = "text/javascript"
				case ".css":
					contentType = "text/css"
				case ".txt":
					contentType = "text/plain"
				}
			}

			w.Header().Set("Content-Type", contentType)
			_, err = io.CopyBuffer(w, file, nil)
			if err != nil {
				logger.Warning("Server Error:", err)
				break
			}
			// http.FileServer(http.FS(staticFs)).ServeHTTP(w, r)
		}
	})

	return mainHandler
}

func getPrefix(urlPath string) string {
	if i := strings.Index(urlPath[1:], "/"); i != -1 {
		return urlPath[:i+1]
	}

	return urlPath
}
