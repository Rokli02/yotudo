package server

import (
	"context"
	"embed"
	"fmt"
	"io"
	"io/fs"
	"net"
	"net/http"
	"strconv"
	"strings"
	"yotudo/server/handler"
	"yotudo/src/lib/logger"
)

//go:embed all:web-app/dist
var webApp embed.FS

type Server struct {
	httpServer *http.Server
	options    ServerOptions
	running    bool
}

func NewServer(optionFuncs ...ServerOptionsFunc) *Server {
	options := ServerOptions{
		listenerIp: "localhost",
		port:       0,
	}

	for _, optionFunc := range optionFuncs {
		optionFunc(&options)
	}

	return &Server{
		httpServer: &http.Server{
			Handler: getHandler(webApp, options),
		},
		options: options,
	}
}

func (s *Server) Start() error {
	if s.running {
		return fmt.Errorf("server is already running")
	}

	s.running = true

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.options.listenerIp, s.options.port))
	if err != nil {
		return err
	}

	if s.options.addr, err = s.getLocalIp(); err != nil {
		return err
	}

	if s.options.port, err = s.getPortFormListener(listener); err != nil {
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
	s.options.addr = ""

	logger.Info("Closing webserver")

	return s.httpServer.Shutdown(ctx)
}

func (s *Server) IsRunning() bool {
	return s.running
}

func (s *Server) GetAddress() string {
	return fmt.Sprintf("%s:%d", s.options.addr, s.options.port)
}

func (s *Server) getLocalIp() (string, error) {
	if s.options.listenerIp == "localhost" {
		return s.options.listenerIp, nil
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
	if s.options.port != 0 {
		return s.options.port, nil
	}

	_addr := listener.Addr().String()
	portBeginningIndex := strings.LastIndex(_addr, ":")

	if portBeginningIndex == -1 || len(_addr)-1 == portBeginningIndex {
		return -1, fmt.Errorf("server couldn't acquire host port")
	}

	return strconv.Atoi(_addr[portBeginningIndex+1:])
}

func getHandler(webAppFS fs.FS, options ServerOptions) http.Handler {
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
			if !options.hostFrontend {
				return
			}

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
