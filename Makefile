BIN_DIR=bin

build:
	@echo "Building client..."
	go build -o pomodoro ./main.go

	@echo "Building daemon..."
	go build -o pomodorod ./daemon

daemon:
	@echo "Running daemon..."
	./pomodorod
client:
	@echo "Running client..."
	./pomodoro start -w