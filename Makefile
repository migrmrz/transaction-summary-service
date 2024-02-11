APPNAME := transactions-summary-service
VERSION := 1.0
SUMMARY := Service that reads a CSV file with transactions and creates a summary that could be sent over by email to the user. 

build:
	mkdir -p build
	GOOS=$(GOOS) GOARCH=$(GOARCH) APPNAME=$(APPNAME) ./scripts/build

run:
	./build/$(APPNAME) $(FILE) $(EMAIL)

clean:
	rm -rf build
