run:
	docker build -t travelite:image .
	docker run -d -p 8080:8080 travelite:image
doc:
	swag init --parseDependency --parseInternal -g cmd/main.go --o docs/