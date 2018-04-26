build:
	go build .
clean:
	@rm -rf fileserver
docker: build
	docker build -t fileserver . 
