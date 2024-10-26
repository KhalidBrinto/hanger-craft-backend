# FROM golang:1.19-alpine3.17 as build

# WORKDIR app/

# COPY src/ .
# RUN ["go", "mod", "tidy"]
# RUN ["go", "build", "main/runner.go"]



# FROM alpine:latest
# WORKDIR /app
# COPY --from=build /go/app/runner .
# CMD ["./runner"]
# EXPOSE 3000

FROM golang:1.23.2-alpine3.20 as dependencies

WORKDIR /app
COPY go.mod go.sum ./


RUN go mod tidy

# FROM dependencies AS build
COPY . ./
RUN CG0_ENABLE=0 go build -o /main -ldflags="-w -s"

# FROM golang:1.19-alpine 
# COPY --from=build /main /main
EXPOSE 3000
CMD [ "/main" ]
