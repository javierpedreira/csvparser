build:
	go build -o ./bin/csvparser ./src/csvparser/main.go

run:
	./bin/csvparser ./input ing

restore:
	mv archive/* input
