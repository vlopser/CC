all: remove build containers

build:
	@echo "***build"
	docker build -t oauth_proxy:latest     oauth2_proxy/ 
	docker build -t frontend_mocked:latest frontend_mocked/

containers:
	@echo "***containers"
	docker run --privileged -ti -d -p 4180:4180 --name op 	    oauth_proxy
	docker run --privileged -ti -d -p 8080:80   --name frontend frontend_mocked

remove:
	@echo "***remove"
	-docker stop op frontend
	-docker rm -f op frontend