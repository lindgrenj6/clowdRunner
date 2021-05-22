all:
	go build -ldflags="-s -w" -o clowdRun

clean:
	rm clowdRun
