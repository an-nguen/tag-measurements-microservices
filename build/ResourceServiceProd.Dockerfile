FROM golang:1.15.2 AS builder

# Create appuser.
ENV USER=appuser
ENV UID=10001

# See https://stackoverflow.com/a/55757473/12429735RUN
RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    "${USER}"


WORKDIR /go/src/
COPY . .

RUN go mod download
RUN go mod verify
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GIN_MODE=release go build -tags=jsoniter -ldflags="-w -s" -o resource_service cmd/resource_service/main.go
RUN readlink -f resource_service

FROM scratch

# Import the user and group files from the builder.
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group
COPY --from=builder /go/src/resource_service /go/src/resource_service
COPY --from=builder /go/src/configs/config_resource.json.prod /go/src/configs/config_resource.json
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

WORKDIR /go/src/

# Use an unprivileged user.
USER appuser:appuser

CMD ["/go/src/resource_service"]
