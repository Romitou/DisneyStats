package commands

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/romitou/disneystats/database"
	"github.com/romitou/disneystats/database/models"
	"github.com/spf13/cobra"
	"log"
	"os"
	"time"
)

type WaitTime struct {
	AttractionName string            `json:"attractionName"`
	Park           models.DisneyPark `json:"park"`
	Status         string            `json:"status"`
	SingleRider    int16             `json:"singleRider"`
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

			var attractions []models.Attraction
			database.GetDatabase().Model(&models.Attraction{}).Find(&attractions)

			var waitTimes []WaitTime
			for _, attraction := range attractions {
				var attractionWaitTime models.AttractionWaitTime
				database.GetDatabase().Model(&models.AttractionWaitTime{}).Where("attraction_id = ?", attraction.ID).Last(&attractionWaitTime)
				// Si le temps entre la dernière mise à jour et maintenant est supérieur à 1 jour, on ne l'affiche pas
				if attractionWaitTime.Status == "" || (time.Now().Unix()-attractionWaitTime.Time) > 86400 {
					continue
				}
				waitTimes = append(waitTimes, WaitTime{
					AttractionName: attraction.Name,
					Park:           attraction.ParkID,
					Status:         attractionWaitTime.Status,
					SingleRider:    attractionWaitTime.SingleRider,
					WaitTime:       attractionWaitTime.WaitTime,
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
