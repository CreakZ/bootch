FROM golang:1.23.4 as builder

WORKDIR /bootch

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./main ./main.go


FROM scratch
COPY --from=builder /bootch/main /bootch/main
COPY --from=builder /bootch/cfg.yaml /bootch/cfg.yaml 

WORKDIR /bootch
CMD [ "./main" ]
