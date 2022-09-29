package commands

import (
	"github.com/joho/godotenv"
	"github.com/romitou/disneystats/database"
	"github.com/romitou/disneystats/tasks"
	"github.com/spf13/cobra"
	"log"
)

var fetchCommand = &cobra.Command{
	Use: "fetch",
	Run: func(cmd *cobra.Command, args []string) {
		err := godotenv.Load()
		if err != nil {
			log.Println("cannot load .env file")
		}

		database.ConnectDatabase()

		tasks.UpdateAttractions()
		tasks.UpdateWaitTimes()
	},
}

func init() {
	DisneyStats.AddCommand(fetchCommand)
}
