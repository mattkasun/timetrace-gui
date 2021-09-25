FROM golang:latest as builder
WORKDIR /go/src/github.com/mattkasun/timetrace-gui

COPY *.go ./
COPY go.* ./
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix temp -ldflags '-extldflags "-static"' .


FROM busybox
WORKDIR /root/
COPY --from=builder /go/src/github.com/mattkasun/timetrace-gui/timetrace-gui .
ADD /images/* images/
ADD /html/* html/
CMD ["./timetrace-gui"]

