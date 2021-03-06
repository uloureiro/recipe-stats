NO_COLOR=\033[0m
OK_COLOR=\033[32;01m
ERROR_COLOR=\033[31;01m

setup-docker:
	@echo "$(OK_COLOR)Setting up the Docker image...$(NO_COLOR)"
	@docker-compose build
	@echo "$(OK_COLOR)Booting container...$(NO_COLOR)"
	@docker-compose up -d
	@echo "$(OK_COLOR)Done!$(NO_COLOR)"

setup-local:
	@echo "$(OK_COLOR)Resolving dependencies...$(NO_COLOR)"
	@go get -d -v ./...
	@go install
	@echo "$(OK_COLOR)Building binary...$(NO_COLOR)"
	@go build
	@tar -xvf sample_data.tar.gz &>/dev/null
	@echo "$(OK_COLOR)Done!$(NO_COLOR)"

docker-up:
	@echo "$(OK_COLOR)Booting container...$(NO_COLOR)"
	@docker-compose up -d
	@echo "$(OK_COLOR)Done!$(NO_COLOR)"

docker-down:
	@-docker-compose down
	@echo "$(ERROR_COLOR)Please come again!$(NO_COLOR)"

run-interactive:
	@docker-compose exec recipe-stats recipe-stats -i

test-local:
	@echo "$(OK_COLOR)Running tests...$(NO_COLOR)"
	@go test tests/*.go -v

test-docker:
	@echo "$(OK_COLOR)Running tests in Docker...$(NO_COLOR)"
	@-docker-compose up -d &>/dev/null &2>/dev/null
	@docker-compose exec recipe-stats go test tests/*.go -v