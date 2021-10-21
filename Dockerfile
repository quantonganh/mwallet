FROM alpine:3.14
RUN apk add --no-cache ca-certificates
COPY mwallet .
EXPOSE 8080
ENTRYPOINT [ "./mwallet" ]