package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "feature-flags",
	Short: "Feature Flags Service CLI",
	Long:  `CLI для управления сервисом feature-flags (запуск сервера, миграции базы данных и др).`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
