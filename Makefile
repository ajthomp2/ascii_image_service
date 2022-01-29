build:
	docker build --tag ascii_art_service .
service:
	docker run -p 8080:8080 ascii_art_service
