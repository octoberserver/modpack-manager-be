FROM golang:latest AS builder

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -buildvcs=false -o /build/modpack-manager .

FROM gcr.io/distroless/base-debian12
LABEL authors="october1234"

COPY --from=builder /build/modpack-manager /modpack-manager

EXPOSE 8080

ENTRYPOINT ["/modpack-manager"]
