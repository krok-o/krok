FROM golang:1.16-alpine as build
WORKDIR /app
COPY . .
RUN go build -o /krok

FROM alpine
RUN apk add -u ca-certificates
COPY --from=build /krok /app/

EXPOSE 9998

RUN mkdir -p /tmp/krok/vault
WORKDIR /app/
ENTRYPOINT [ "/app/krok" ]