build:
	@go build -o bin/api

run: build
	@./bin/api

seed:
	@go run scripts/seed.go

docker:
	echo "building docker file"
	@docker build -t api .
	echo "running API inside Docker container"
	@docker run -p 3000:3000 api

test:
	@go test -v ./...

testauth:
	@go test -v -run TestAuthenticate ./api

testbooking:
	@go test -v -run Booking ./api
