build:
	go build -o ./bin/csvparser ./src/csvparser/main.go

preprcess:
	cd input && /Applications/LibreOffice.app/Contents/MacOS/soffice --convert-to xlsx *.xls --headless
	rm input/*.xls

run:
	./bin/csvparser ./input ing

restore:
	mv archive/* input
