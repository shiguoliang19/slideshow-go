#!/bin/sh

docker build -t slideshow-go .

docker run -d -p 8085:8085 slideshow-go

