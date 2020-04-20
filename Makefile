# Makefile variables
PROJECT_NAME=go_interview

# go build
compile:
	mkdir -p dist
	go build -o dist/${PROJECT_NAME}

run:
	./dist/${PROJECT_NAME}