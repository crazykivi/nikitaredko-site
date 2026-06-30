FROM alpine:3.19

RUN apk update && apk add --no-cache ca-certificates tzdata
RUN mkdir -p /app

COPY backend/main /app/
COPY frontend/dist /app/dist

WORKDIR /app
EXPOSE 8080

ENTRYPOINT ["/app/main"]