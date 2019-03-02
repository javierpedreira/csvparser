build:
	go build -o csvparser main.go

run:
	cd input && /Applications/LibreOffice.app/Contents/MacOS/soffice --convert-to xlsx *.xls --headless
	rm input/*.xls
	./csvparser input

restore:
	mv archive/* input