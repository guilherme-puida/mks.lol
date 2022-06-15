set dotenv-load

bin := "./bin"
server := "mks.lol.server"

default:
    just --list

# Builds the server binary
build-server:
    mkdir -p {{bin}}
    cd {{bin}}
    go build -v ./cmd/{{server}}
    mv {{server}} {{bin}}
    @echo "Server executable built to {{bin}}/{{server}}"

# Starts the server.
run-server: build-server
    @echo "Starting server..."
    {{bin}}/{{server}}

clean:
    rm -rf {{bin}}