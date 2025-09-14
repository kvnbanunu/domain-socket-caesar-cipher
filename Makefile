SERVER = cmd/server/main.go
CLIENT = cmd/client/main.go
SERVER_TARGET = bin/server
CLIENT_TARGET = bin/client
BUILD = go build -o
RUN = go run
DEBUG = -d
HELP = -h
DOMAIN_SOCKET = domain.sock
MESSAGE = Hello
SHIFT = 3
SERVER_ARGS = -p $(DOMAIN_SOCKET)
CLIENT_ARGS = -p $(DOMAIN_SOCKET) -s $(SHIFT) -i $(MESSAGE)
COPY_CONFIG = cp config.json bin/

build-all: clean-all server_build client_build

build-s: clean-s
	@$(BUILD) $(SERVER_TARGET) $(SERVER)
	@$(COPY_CONFIG)

build-c: clean-c
	@$(BUILD) $(CLIENT_TARGET) $(CLIENT)
	@$(COPY_CONFIG)

run-s:
	@$(RUN) $(SERVER) $(SERVER_ARGS)

run-c:
	@$(RUN) $(CLIENT) $(CLIENT_ARGS)

debug-s:
	@$(RUN) $(SERVER) $(DEBUG) $(SERVER_ARGS)

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
