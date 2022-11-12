package services

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/x-x-x-Ilya/astrologer/internal/models"
)

type NasaClientI interface {
	Picture(date time.Time) (models.Picture, error)
}

type NasaClient struct {
	client ClientServiceI
	apiKey string
	url    string
}

func NewNasaClient(apiKey string, client ClientServiceI) (NasaClientI, error) {
	if client == nil {
		return nil, nilErr("client")
	}

	nasaClient := &NasaClient{
		client,
		apiKey,
		"https://api.nasa.gov",
	}

	return nasaClient, nil
}

type pictureResponse struct {
	URL string `json:"hdurl"`
}

/*
   {
	"copyright":"Ryan Han",
	"date":"2022-11-11",
	"explanation":	"On November 8 the Full Moon turned blood red as it slid through Earth's shadow in a beautiful total lunar eclipse.
					During totality it also passed in front of, or occulted, outer planet Uranus for eclipse viewers located in parts of northern America and Asia.
					For a close-up and wider view these two images were taken just before the occultation began, captured with different telescopes and cameras from
					the same roof top in Shanghai, China. Normally very faint compared to a Full Moon, the tiny, pale, greenish disk of the distant ice giant is just
					to the left of the Moon's edge and about to disappear behind the darkened, red lunar limb. Though only visible from certain locations across planet Earth,
					lunar occultations of planets are fairly common. But for this rare \"lunar eclipse occultation\" to take place, at the time of the total eclipse the outer
					planet had to be both at opposition and very near the ecliptic plane to fall in line with Sun, Earth, and Moon.
					Lunar Eclipse of November 2022: Notable Submissions to APOD  Love Eclipses? (US): Apply to become a NASA Partner Eclipse Ambassador",
	"hdurl":"https://apod.nasa.gov/apod/image/2211/LunarEclipseRyanHan.jpg",
	"media_type":"image",
	"service_version":"v1",
	"title":"Blood Moon, Ice Giant",
	"url":"https://apod.nasa.gov/apod/image/2211/LunarEclipseRyanHan1024.jpg"
}
*/

func (n *NasaClient) Picture(date time.Time) (models.Picture, error) {
	year, month, day := date.Date()
	queryParams := map[string][]string{
		"date":    {fmt.Sprintf("%d-%d-%d", year, month, day)},
		"api_key": {n.apiKey},
	}

	response, err := n.client.Get(n.url+"/planetary/apod/", queryParams)
	defer closeBody(response.Body)
	if err != nil {
		return models.Picture{}, err
	}

	var responseStruct pictureResponse
	err = json.NewDecoder(response.Body).Decode(&responseStruct)
	if err != nil {
		return models.Picture{}, err
	}

	imgResponse, err := n.client.Get(responseStruct.URL, nil)
	if err != nil {
		return models.Picture{}, err
	}
	defer closeBody(imgResponse.Body)

	buffer := make([]byte, imgResponse.ContentLength)
	_, err = io.ReadFull(imgResponse.Body, buffer)
	if err != nil {
		return models.Picture{}, err
	}

	return models.NewPicture(date, buffer), nil
}
