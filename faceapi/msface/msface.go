package msface

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
)

//FaceAttributes is subtype of JSONResponse
type FaceAttributes struct {
	Emotion map[string]float64 `json:"emotion"`
}

//FaceResponse is returned from Face API if only emotion is requested
type FaceResponse struct {
	FaceID         string         `json:"faceId"`
	FaceRectangle  map[string]int `json:"faceRectangle"`
	FaceAttributes FaceAttributes `json:"faceAttributes"`
}

var uri = "https://westeurope.api.cognitive.microsoft.com/face/v1.0/detect"
var key = ""

func GetEmotion(m io.Reader) (string, error) {
	client := &http.Client{}

	// only request emotion detection
	params := "emotion"
	request, err := http.NewRequest("POST", uri+"?returnFaceAttributes="+params, m)
	request.Header.Add("Content-Type", "application/octet-stream")
	request.Header.Add("Ocp-Apim-Subscription-Key", key)

	// Send the request to the local web service
	resp, err := client.Do(request)
	if err != nil {
		return "", err
	}

	// read response
	respBody, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	var response []FaceResponse
	err = json.Unmarshal(respBody, &response)
	if err != nil {
		return "", err
	}

	// return only first face
	highestID := ""
	highestEmo := 0.0
	for id, emo := range response[0].FaceAttributes.Emotion {
		if emo > highestEmo {
			highestEmo = emo
			highestID = id
		}
	}

	return highestID, nil
}
