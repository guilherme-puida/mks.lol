alias b := build
alias d := dev

build:
    CGO_ENABLED=0 go build

dev: build
    echo "Service starting at http://localhost:8080/"
    ./mks.lol -url localhost:8080 -port 8080

serve: build
    ./mks.lol -https