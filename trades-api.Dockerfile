# First we build the ent schemas
FROM golang:1.21.4 as builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY cmd ./cmd
COPY cmd/cmd_api/main.go ./cmd/cmd_api/main.go
COPY internal ./internal

RUN go generate ./internal/dbschema
RUN CGO_ENABLED=0 go build -o app ./cmd/cmd_api/main.go

# Use distroless as minimal base image to package the manager binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM gcr.io/distroless/static:nonroot
WORKDIR /

COPY --from=builder /src/app /app

USER nonroot:nonroot

ENTRYPOINT ["/app"]
