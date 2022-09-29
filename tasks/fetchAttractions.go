package tasks

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/romitou/disneystats/database"
	"github.com/romitou/disneystats/database/models"
	"gorm.io/gorm"
	"log"
	"net/http"
)

const AttractionsQuery = `query activities($market: String!, $types: [String]) {
    activities(market: $market, types: $types) {
      id
      name
      location {
          id
      }
    }
}`

type Attraction struct {
	ID       string             `json:"id"`
	Name     string             `json:"name"`
	Location AttractionLocation `json:"location"`
}

type AttractionLocation struct {
	ID string `json:"id"`
}

func (l AttractionLocation) GetDisneyPark() models.DisneyPark {
	switch l.ID {
	case "P1":
		return models.DisneylandPark
	case "P2":
		return models.WaltDisneyStudios
	}
	return models.Unknown
}

type AttractionResponse struct {
	Data struct {
		Attractions []Attraction `json:"activities"`
	} `json:"data"`
}

func fetchAttractions() []Attraction {
	query, err := GraphQlQuery{
		Query: AttractionsQuery,
		Variables: map[string]interface{}{
			"market": "en-gb",
			"types":  []string{"Attraction"},
		},
	}.ToJSON()
	if err != nil {
		log.Println(err)
		return nil
	}

	response, err := http.Post(DLPGraphQL, "application/json", bytes.NewReader(query))
	if err != nil {
		log.Println(err)
		return nil
	}

	var attractionResponse AttractionResponse
	err = json.NewDecoder(response.Body).Decode(&attractionResponse)
	if err != nil {
		log.Println(err)
		return nil
	}

	return attractionResponse.Data.Attractions
}

func UpdateAttractions() {
	attractions := fetchAttractions()
	if attractions == nil {
		return
	}

	db := database.GetDatabase()

	for i := range attractions {
		err := db.Where(&models.Attraction{
			EntityID: attractions[i].ID,
		}).First(&models.Attraction{}).Error

		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = db.Create(&models.Attraction{
				EntityID: attractions[i].ID,
				Name:     attractions[i].Name,
				ParkID:   attractions[i].Location.GetDisneyPark(),
			}).Error
			if err != nil {
				log.Println("an error occurred while creating attraction: ", err)
			}
		}
	}
}
