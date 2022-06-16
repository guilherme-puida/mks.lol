FROM golang:1.18-alpine AS builder

WORKDIR /src
COPY . .
RUN ls
RUN CGO_ENABLED=0 go build -o /mks.lol.server ./cmd/mks.lol.server

FROM gcr.io/distroless/static AS final

LABEL maintainer="guilherme-puida <guilhermepmoreira00@gmail.com>"
USER nonroot:nonroot

COPY --from=builder --chown=nonroot:nonroot /mks.lol.server /mks.lol.server

ENV MKS_PORT=8080
ENV MKS_URL=localhost:8080
ENV MKS_SHOULD_RENDER_STATS=true

EXPOSE $MKS_PORT

ENTRYPOINT ["/mks.lol.server"]