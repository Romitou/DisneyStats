package commands

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/romitou/disneystats/database"
	"github.com/romitou/disneystats/database/models"
	"github.com/romitou/disneystats/tasks"
	"github.com/spf13/cobra"
	"log"
	"os"
)

type WaitTime struct {
	AttractionName string            `json:"attractionName"`
	Park           models.DisneyPark `json:"park"`
	Status         string            `json:"status"`
	SingleRider    int16             `json:"singleRider,omitempty"`
	WaitTime       int16             `json:"waitTime"`
}

var DisneyStats = &cobra.Command{
	Use: "disneystats",
	Run: func(cmd *cobra.Command, args []string) {
		err := godotenv.Load()
		if err != nil {
			log.Println("cannot load .env file")
		}

		database.ConnectDatabase()

		r := gin.Default()
		r.GET("/wait-times", func(c *gin.Context) {
			apiWaitTimes := tasks.FetchWaitTimes()

			var waitTimes []WaitTime
			for _, waitTime := range apiWaitTimes {
				var attraction models.Attraction
				database.GetDatabase().Model(&models.Attraction{}).Where("entity_id = ?", waitTime.EntityID).First(&attraction)
				if attraction.ID == 0 {
					continue
				}

				attractionWaitTime, err := waitTime.WaitMinsInt()
				if err != nil {
					log.Println(err)
					continue
				}

				var singleRider int16
				if waitTime.SingleRider.IsAvailable {
					singleRider, err = waitTime.SingleRider.WaitMinsInt()
					if err != nil {
						singleRider = 0
						continue
					}
				}

				waitTimes = append(waitTimes, WaitTime{
					AttractionName: attraction.Name,
					Park:           attraction.ParkID,
					Status:         waitTime.Status,
					SingleRider:    singleRider,
					WaitTime:       attractionWaitTime,
				})
			}

			c.JSON(200, waitTimes)
		})
		r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
	},
}

func Execute() {
	if err := DisneyStats.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
