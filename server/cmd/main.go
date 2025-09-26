package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	serverModule "yotudo/server"
	"yotudo/server/cmd/utils"
	"yotudo/src/database"
	"yotudo/src/settings"
)

func main() {
	args := utils.ParseArguments(os.Args[1:])
	if args.GetState(utils.Help) {
		fmt.Println("Most még nézz bele a kódba, később majd le lesz írva minden!")
		return
	}

	defer func() {
		fmt.Println("EVENT: close=true")
	}()

	if err := settings.CreateEssentialDirectoriesAndFiles(); err != nil {
		fmt.Println("ERR:", err)
		return
	}

	if _, err := settings.LoadSettings(); err != nil {
		fmt.Println("ERR:", err)
		return
	}

	database.LoadDatabase()

	server := serverModule.NewServer(func(serverOptions *serverModule.ServerOptions) {
		serverOptions.
			Host(args.GetState(utils.Host)).
			HostFrontend(args.GetState(utils.HostFrontend)).
			SetPort(args.Port)
	})

	fmt.Println("INFO: Starting webserver")
	if err := server.Start(); err != nil {
		fmt.Println("ERR:", err)
		return
	}
	defer func() {
		fmt.Println("INFO: Stopping webserver")
		server.Stop(context.Background())
	}()

	fmt.Printf("INFO: Server is listening on http://%s\n", server.GetAddress())
	fmt.Printf("EVENT: listening=%s\n", server.GetAddress())

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	exitSig := utils.WaitForExit()

	for {
		select {
		case <-c:
			return
		case <-exitSig:
			return
		}
	}
}
