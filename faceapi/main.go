package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/gbaeke/emotion/faceapi/msface"
)

func main() {
	imagePath := flag.String("image", "", "path to image")
	flag.Parse()

	if *imagePath == "" {
		log.Fatal("Please specify image path. Use --help for help.")
	}

	//read image and return reader
	// read the image file
	m, err := os.Open(*imagePath)
	if err != nil {
		log.Fatalln("Error opening image", err)
	}
	defer m.Close()

	emotion, err := msface.GetEmotion(m)

	fmt.Println("probably", emotion)
}
