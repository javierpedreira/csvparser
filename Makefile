build:
	go build -o csvparser main.go

run:
	cd input && /Applications/LibreOffice.app/Contents/MacOS/soffice --convert-to xlsx *.xls --headless
	./csvparser input
	rm input/*.xls

restore:
	mv archive/* input