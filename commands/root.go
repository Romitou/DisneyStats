package commands

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"os"
)

var DisneyStats = &cobra.Command{
	Use: "disneystats",
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("a")
	},
}

func Execute() {
	if err := DisneyStats.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
