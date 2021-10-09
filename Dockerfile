FROM golang:latest as builder
WORKDIR /go/src/github.com/mattkasun/timetrace-gui

COPY *.go ./
COPY go.* ./
RUN go mod tidy
ADD /images/* images/
ADD /html/* html/
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix temp -ldflags '-extldflags "-static"' .


FROM busybox
WORKDIR /root/
COPY --from=builder /go/src/github.com/mattkasun/timetrace-gui/timetrace-gui .
ADD /images/* images/
ADD /html/* html/
RUN mkdir -p .timetrace/projects
RUN mkdir -p .timetrace/records
RUN mkdir -p .timetrace/reports
CMD ["./timetrace-gui"]

