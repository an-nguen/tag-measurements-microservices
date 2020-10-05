FROM golang:alpine AS builder

RUN apk update && apk add --no-cache git
RUN apk --no-cache add tzdata

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
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o clean_service cmd/clean_service/main.go
RUN readlink -f clean_service

FROM scratch

# Import the user and group files from the builder.
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group

COPY --from=builder /go/src/clean_service /go/src/clean_service
COPY --from=builder /go/src/configs/config_clean.json.prod /go/src/configs/config_clean.json
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
ENV TZ=Europe/Moscow
WORKDIR /go/src/

# Use an unprivileged user.
USER appuser:appuser

CMD ["/go/src/clean_service"]
