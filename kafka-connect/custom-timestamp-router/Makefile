.PHONY: install-maven build test

install-maven:
	brew install maven

build:
	mvn clean package

test: build
	mvn test