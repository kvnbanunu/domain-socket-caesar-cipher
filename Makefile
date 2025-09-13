SERVER = cmd/server.go
CLIENT = cmd/client.go
SERVER_TARGET = bin/server
CLIENT_TARGET = bin/client
BUILD = go build -o
RUN = go run
DEBUG = -d
HELP = -h
DOMAIN_SOCKET = bin/domain.sock
MESSAGE = Hello
SHIFT = 3
CLIENT_ARGS = -p $(DOMAIN_SOCKET) -s $(SHIFT) -i $(MESSAGE)

build-all: clean-all server_build client_build

build-s: clean-s
	@$(BUILD) $(SERVER_TARGET) $(SERVER)

build-c: clean-c
	@$(BUILD) $(CLIENT_TARGET) $(CLIENT)

run-s:
	@$(RUN) $(SERVER) $(DOMAIN_SOCKET)

run-c:
	@$(RUN) $(CLIENT) $(CLIENT_ARGS)

debug-s:
	@$(RUN) $(SERVER) $(DEBUG) $(DOMAIN_SOCKET)

debug-c:
	@$(RUN) $(CLIENT) $(DEBUG) $(CLIENT_ARGS)

help-s:
	@$(RUN) $(SERVER) $(HELP)

help-c:
	@$(RUN) $(CLIENT) $(HELP)

clean-s:
	@rm -f $(SERVER_TARGET)

clean-c:
	@rm -f $(CLIENT_TARGET)

clean-all:
	@rm -rf bin
