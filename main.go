package main

import (
	core "betassist.ru/bookmaker-game-parser/internal"
	"betassist.ru/bookmaker-game-parser/internal/database"
	factory "betassist.ru/bookmaker-game-parser/internal/database/repositories"
	"betassist.ru/bookmaker-game-parser/internal/repositories"
	"betassist.ru/bookmaker-game-parser/internal/services"
	"fmt"
	"github.com/spf13/cobra"
	"runtime"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "",
		Short: "Bookmaker game parser service",
		Long:  "Bookmaker game parser service template",
		Run: func(cmd *cobra.Command, args []string) {
			database.Open()
			defer database.Close()
			factory.Register()

			bookmaker, err := repositories.Bookmaker.Select(core.Config.Websocket.Bookmaker)

			if nil != err {
				panic(err)
			}

			core.Config.App.BookmakerID = bookmaker.ID

			services.Websocket.Register()
			go services.Listener.Listen()

			defer services.Listener.Close()

			services.Websocket.Listen()
		},
	}

	var versionCmd = &cobra.Command{
		Use:     "version",
		Aliases: []string{"ver"},
		Short:   "Bookmaker game parser application version information",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Application v %s (Go version: %s)\n", core.Version, runtime.Version())
		},
	}

	rootCmd.AddCommand(versionCmd)

	_ = rootCmd.Execute()
}
