FROM golang:1.25-alpine AS builder

RUN apk add --no-cache gcc musl-dev sqlite-dev

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=1 GOOS=linux go build -o forum .

FROM alpine:latest
RUN apk add --no-cache sqlite-libs

WORKDIR /app
COPY --from=builder /app/forum .
COPY schema.sql .
COPY templates/ templates/
COPY static/ static/

RUN mkdir -p uploads

EXPOSE 8080
CMD ["./forum"]