# Emotion Detection with FER+ in Go

Go program to detect emotion from faces in a video stream. It uses a container built with Azure Machine Learning and the ONNX FER+ model for emotion detection.

Use the following command to start the container:

docker run -d -p 5002:5001 gbaeke/onnxferplus

The container exposes a scoring URI at http://localhost:5002/score.