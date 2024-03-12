server:
	cd cmd/server && go build -o ../../keeperServer

client:
	cd cmd/client && go build -o ../../keeperClient

all-clients:
	cd cmd/client && bash ./build.sh

