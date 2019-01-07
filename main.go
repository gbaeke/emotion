package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/disintegration/imaging"
	//"github.com/gbaeke/emotion/faceapi/msface"
	"gocv.io/x/gocv"
)

//InputData is sent to FER+ model
type InputData struct {
	Data [1][1][64][64]uint8 `json:"data"`
}

//OutputData is received from FER+ model
type OutputData struct {
	Result []float64 `json:"result"`
	Time   float64   `json:"time"`
}

func main() {
	scoreuri, ok := os.LookupEnv("SCOREURI")
	fmt.Println(scoreuri)
	if !ok || scoreuri == "" {
		scoreuri = "http://localhost:5002/score"
	}

	deviceID := 0
	xmlFile := "haarcascade_frontalface_default.xml"

	// open webcam
	webcam, err := gocv.VideoCaptureDevice(deviceID)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer webcam.Close()

	// open display window
	//window := gocv.NewWindow("Face Detection with FER+")
	//defer window.Close()

	// captured image ends up in below image matrix
	img := gocv.NewMat()
	defer img.Close()

	// color for the rect when faces detected
	green := color.RGBA{0, 255, 0, 0}

	// load classifier to recognize faces
	classifier := gocv.NewCascadeClassifier()
	defer classifier.Close()

	if !classifier.Load(xmlFile) {
		fmt.Printf("Error reading cascade file: %v\n", xmlFile)
		return
	}

	fmt.Printf("start reading camera device: %v\n", deviceID)
	frameCount := 0 //used to detect emotion every 2nd frame
	emotion := ""   //emotion displayed on screen
	for {
		if ok := webcam.Read(&img); !ok {
			fmt.Printf("cannot read device %d\n", deviceID)
			return
		}
		if img.Empty() {
			continue
		}

		// detect faces
		rects := classifier.DetectMultiScaleWithParams(img, 1.1, 5, 0, image.Point{100, 100},
			image.Point{300, 300})
		frameCount++

		// only look at first face found
		if len(rects) > 0 {
			r := rects[0]

			// draw green rectangle around the face
			gocv.Rectangle(&img, r, green, 3)

			// get mat of face region; copy to a new mat
			faceRegion := img.Region(r)
			face := gocv.NewMat()
			faceRegion.CopyTo(&face)

			// convert new mat with just the face to image
			emoImg, err := face.ToImage()
			emoImg = resizeImage(emoImg, 64, 64)

			// get emotion
			if err == nil && frameCount%2 == 0 {

				//use FER+
				emotion = getEmotion(emoImg, scoreuri)

				//use Microsoft Face API; encode mat to JPG and convert to io.Reader
				//encodedImage, _ := gocv.IMEncode(gocv.JPEGFileExt, face)
				//emotion, err = msface.GetEmotion(bytes.NewReader(encodedImage))
				//if err != nil {
				//	log.Println(err)
				//}
			}

			// add text to webcam image
			size := gocv.GetTextSize(emotion, gocv.FontHersheyPlain, 1.5, 3)
			pt := image.Pt(r.Min.X+(r.Min.X/2)-(size.X/2), r.Min.Y-2)
			gocv.PutText(&img, emotion, pt, gocv.FontHersheyPlain, 1.2, green, 2)

		}

		// show the image in the window, and wait 1 millisecond
		//window.IMShow(img)
		//if window.WaitKey(1) >= 0 {
		//break
		//}
	}

}

func getEmotion(m image.Image, scoreuri string) string {

	// multidim array as input tensor
	var BCHW [1][1][64][64]uint8

	for x := 0; x < 64; x++ {
		for y := 0; y < 64; y++ {
			// get RGB values
			r, g, b, _ := m.At(x, y).RGBA()
			rs := uint8(r >> 8)
			rg := uint8(g >> 8)
			rb := uint8(b >> 8)

			// set grayscale value at yw
			BCHW[0][0][y][x] = rs>>2 + rg>>1 + rb>>2

		}
	}

	// input is struct with 4D array
	input := InputData{
		Data: BCHW,
	}

	// Create JSON from input struct - inputJSON will be sent to model
	inputJSON, _ := json.Marshal(input)
	body := bytes.NewBuffer(inputJSON)

	// Create the HTTP request - no need for auth with local FER+ container
	client := &http.Client{}
	request, err := http.NewRequest("POST", scoreuri, body)
	request.Header.Add("Content-Type", "application/json")

	// Send the request to the web service
	resp, err := client.Do(request)
	if err != nil {
		log.Fatal("Error calling scoring URI: ", err)
	}

	// read response
	respBody, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	//Unmarshal returned JSON data
	var modelResult OutputData
	err = json.Unmarshal(respBody, &modelResult)
	if err != nil {
		log.Fatal("Error unmarshalling JSON response ", err)
	}

	// highest result
	maxProb := 0.0
	maxIndex := 0
	for index, prob := range modelResult.Result {
		if prob > maxProb {
			maxProb = prob
			maxIndex = index
		}
	}

	categories := map[int]string{0: "neutral", 1: "happy", 2: "surprise", 3: "sadness",
		4: "anger", 5: "disgust", 6: "fear", 7: "contempt"}

	fmt.Println("Highest prob is", maxProb, "at", maxIndex, "(inference time:", modelResult.Time, ")")
	return categories[maxIndex]
}

func resizeImage(m image.Image, width, height int) image.Image {
	// resize image
	m = imaging.Resize(m, width, height, imaging.Linear)

	return m
}
