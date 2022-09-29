package tasks

import (
	"encoding/json"
	"errors"
	"github.com/romitou/disneystats/database"
	"github.com/romitou/disneystats/database/models"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

const WaitTimesEndpoint = "https://dlp-wt.wdprapps.disney.com/prod/v1/waitTimes"

type SingleRiderEntity struct {
	IsAvailable bool   `json:"isAvailable"`
	WaitMins    string `json:"singleRiderWaitMinutes"`
}

func (e SingleRiderEntity) WaitMinsInt() (int16, error) {
	atoi, err := strconv.Atoi(e.WaitMins)
	if err != nil {
		return 0, err
	}
	return int16(atoi), nil
}

type WaitTimeEntity struct {
	EntityID    string            `json:"entityId"`
	Status      string            `json:"status"`
	WaitMins    string            `json:"postedWaitMinutes"`
	SingleRider SingleRiderEntity `json:"singleRider"`
}

func (e WaitTimeEntity) WaitMinsInt() (int16, error) {
	atoi, err := strconv.Atoi(e.WaitMins)
	if err != nil {
		return 0, err
	}
	return int16(atoi), nil
}

func fetchWaitTimes() []WaitTimeEntity {
	request, err := http.NewRequest("GET", WaitTimesEndpoint, nil)
	if err != nil {
		log.Println(err)
		return nil
	}

	dlpApiKey := os.Getenv("DLP_API_KEY")
	if dlpApiKey == "" {
		return nil
	}

	request.Header.Set("x-api-key", dlpApiKey)

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Println(err)
		return nil
	}

	var waitTimes []WaitTimeEntity
	err = json.NewDecoder(response.Body).Decode(&waitTimes)
	if err != nil {
		log.Println(err)
		return nil
	}

	return waitTimes
}

func UpdateWaitTimes() {
	waitTimes := fetchWaitTimes()
	now := time.Now().Unix()

	for i := range waitTimes {
		var attraction models.Attraction
		err := database.GetDatabase().Where(&models.Attraction{
			EntityID: waitTimes[i].EntityID,
		}).First(&attraction).Error
		if err != nil {
			continue
		}

		var attractionWaitTime models.AttractionWaitTime
		err = database.GetDatabase().Debug().Where(&models.AttractionWaitTime{
			AttractionID: attraction.ID,
		}).Order("time DESC").First(&attractionWaitTime).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			log.Println(err)
			continue
		}

		waitTimeMins, err := waitTimes[i].WaitMinsInt()
		if err != nil {
			log.Println(err)
			continue
		}

		var singleRiderTimeMins int16
		if waitTimes[i].SingleRider.IsAvailable {
			singleRiderTimeMins, err = waitTimes[i].SingleRider.WaitMinsInt()
			if err != nil {
				log.Println(err)
				continue
			}
		}

		update := false
		if attractionWaitTime.Status != waitTimes[i].Status {
			update = true
		} else if attractionWaitTime.WaitTime != waitTimeMins {
			update = true
		} else if attractionWaitTime.SingleRider != singleRiderTimeMins {
			update = true
		}

		if !update {
			continue
		}

		err = database.GetDatabase().Create(&models.AttractionWaitTime{
			Time:        now,
			Attraction:  attraction,
			Status:      waitTimes[i].Status,
			WaitTime:    waitTimeMins,
			SingleRider: singleRiderTimeMins,
		}).Error

		if err != nil {
			log.Println(err)
			continue
		}
	}
}
