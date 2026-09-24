FROM golang:1.25-alpine AS build
WORKDIR /src

RUN apk add --no-cache ca-certificates git

COPY go.mod ./
COPY go.sum* ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" \
  -o /out/zoho-cliq-release-notifier \
  ./cmd/zoho-cliq-release-notifier

FROM scratch
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /out/zoho-cliq-release-notifier /zoho-cliq-release-notifier
ENTRYPOINT ["/zoho-cliq-release-notifier"]
