FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o restgeld-backend .

FROM gcr.io/distroless/static:nonroot
WORKDIR /app
COPY --from=builder /app/restgeld-backend .
EXPOSE 8080
ENTRYPOINT ["/app/restgeld-backend"]
