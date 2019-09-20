# xapp
Demo code repository for golang introduction

### Local DB in Docker

    docker run \
    	-e POSTGRES_USER=xapp \
    	-e POSTGRES_PASSWORD=xapp123 \
    	-e POSTGRES_DB=xapp \
    	-v xappdb:/var/lib/postgresql/data \
    	-p 5432:5432 \
    	postgres:10.6