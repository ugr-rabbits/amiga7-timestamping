##### Build stage #####
FROM golang:1.17-alpine3.14 AS build

WORKDIR /root

COPY go.mod go.sum ./
COPY cmd/website/main.go cmd/website/
COPY web/ web/

RUN go build -o website cmd/website/main.go

##### Run stage #####
FROM alpine:3.14

WORKDIR /amiga7-timestamping

COPY --from=build /root/website ./
COPY --from=build /root/web/ web/

EXPOSE 8000

ENTRYPOINT [ "./website" ]
