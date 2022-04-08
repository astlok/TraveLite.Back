run:
	docker build -t travelite:image .
	docker run -d -p 8080:8080 travelite:image
stop:
	docker rm -f $(shell docker ps -aq)
del:
	docker rmi -f $(shell docker images -a -q)
doc:
	swag init --parseDependency --parseInternal --parseVendor -g cmd/main.go --o docs/