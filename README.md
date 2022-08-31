# Rusprofile
This project is a simple wrapper over rusprofile.ru site for getting company info. 
It has API via gRPC and REST API (you can see swagger documentation on path /swagger/ and make some test requests).
## How to run in Docker container
GRPC runs on 8081 port, REST gateway proxy on 8080 port by default. You can change it in docker-compose.yml using env variables.   
```bash
docker-compose build
docker-compose up
```

## How to build and run directly
```bash
make
./server
```