FROM golang:1.22

WORKDIR /app

COPY . .
RUN go mod download

RUN ls -al
RUN go build main.go

EXPOSE 8080
CMD [ "./main" ]