FROM golang:1.24-alpine AS builder

# copy app files
WORKDIR /go/src
COPY . .

# download GO modules
RUN go mod download
RUN go mod verify

# compile app
RUN CGO_ENABLED=0 GOOS=linux go build -o /go/bin/backend

# 2nd stage - run
FROM alpine:3.21

# copy app files
# COPY --from=builder /go/src/.env .
COPY --from=builder /go/bin/backend /go/bin/backend

EXPOSE 8080

# run app
CMD ["/go/bin/backend"]