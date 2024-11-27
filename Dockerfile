##
## STEP 1 - BUILD
##

FROM golang:1.23.3-alpine3.20 as build

WORKDIR /app

COPY go.mod ./

RUN go mod download

COPY . ./

RUN CG0_ENABLED=0 go build -o /main -ldflags="-w -s"

##
## STEP 2 - DEPLOY
##

FROM alpine:latest

WORKDIR /
COPY --from=build /main /main

EXPOSE 3000

ENTRYPOINT [ "/main" ]