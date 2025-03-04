FROM golang:1.23-alpine as build
RUN apk add --no-cache git gcc musl-dev && \
    go install github.com/swaggo/swag/cmd/swag@latest
RUN mkdir /app
ADD . /app/
WORKDIR /app
COPY ./go.mod .
COPY ./go.sum .
ENV GOPROXY https://proxy.golang.org,direct
RUN go mod download


RUN swag init --parseDependency=true

ENV CGO_ENABLED=0
RUN GOOS=linux go build -o main ./main.go

FROM alpine as serve
WORKDIR /app
COPY --from=build /app/main .
COPY --from=build /app/docs ./docs

CMD ["/app/main"]