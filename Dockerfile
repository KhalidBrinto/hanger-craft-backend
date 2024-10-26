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

FROM golang:1.19-alpine as dependencies

WORKDIR /app
COPY go.mod go.sum ./


RUN go mod tidy

# FROM dependencies AS build
COPY . ./
RUN CG0_ENABLE=0 go build -o /main -ldflags="-w -s"

# FROM golang:1.19-alpine 
# COPY --from=build /main /main
CMD [ "/main" ]