#!/bin/sh

docker build -t slideshow-go .

docker run -d -p 8085:8085 -v /opt/assets/slideshow-go/src/assets/images:/data/assets/images --restart=always slideshow-go
