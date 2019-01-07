# Emotion Detection with FER+ in Go

Go program to detect emotion from faces in a video stream. It uses a container built with Azure Machine Learning and the ONNX FER+ model for emotion detection.

Use the following command to start the container:

docker run -d -p 5002:5001 gbaeke/onnxferplus

The container exposes a scoring URI at http://localhost:5002/score. Scoring URI can be set with environment variable SCOREURI

By default, the webcam capture is shown in a window. With the environment variable VIDEO=0, you can turn this off.

Code requires:

- github.com/disintegration/imaging
- gocv.io/x/gocv

Also install Open CV. See [GoCV](https://gocv.io/) for more info

See [blog post](https://blog.baeke.info/2019/01/06/detecting-emotions-with-fer/) for more information.

Run in a container as follows:

docker run -it --rm --device=/dev/video0 --env SCOREURI="YOUR-SCORE-URI" --env VIDEO=0 gbaeke/emo