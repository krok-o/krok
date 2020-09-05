FROM alpine
RUN apk add -u ca-certificates
COPY ./build/linux/amd64/krok /app/

EXPOSE 9998

WORKDIR /app/
ENTRYPOINT [ "/app/krok" ]