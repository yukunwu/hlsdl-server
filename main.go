package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
)

var cmd = &cobra.Command{
	Use:          "hlsdl_server",
	RunE:         cmdF,
	SilenceUsage: true,
}

func main() {
	cmd.Flags().IntP("port", "p", 8080, "Service listen ports")
	cmd.Flags().StringP("dir", "d", "./download", "The directory where the file will be stored")
	cmd.Flags().IntP("workers", "w", 3, "Number of workers to execute concurrent operations")
	cmd.SetArgs(os.Args[1:])
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func cmdF(command *cobra.Command, args []string) error {
	port, err := command.Flags().GetInt("port")
	if err != nil {
		return err
	}
	dir, err := command.Flags().GetString("dir")
	if err != nil {
		return err
	}
	workers, err := command.Flags().GetInt("workers")
	if err != nil {
		return err
	}
	r := gin.Default()
	r.POST("/download", download)
	r.GET("/info", query)
	go func() {
		ins := GetInstance()
		ins.Run(workers, dir)
	}()
	return r.Run(fmt.Sprintf(":%d", port)) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
