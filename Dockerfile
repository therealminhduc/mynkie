FROM golang:1.25-alpine AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go run main.go


FROM nginx:alpine AS runtime

COPY --from=builder /build/public /usr/share/nginx/html
COPY --from=builder /build/static /usr/share/nginx/html/static
COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80
