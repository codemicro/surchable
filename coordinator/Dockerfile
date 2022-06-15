FROM golang:1 as builder

RUN mkdir -p /build/coordinator
ADD . /build/coordinator/
ADD ../go.mod ../go.sum /build/

WORKDIR /build 
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static" -s -w' -o main github.com/codemicro/surchable/coordinator

FROM scratch

COPY --from=builder /build/main /
WORKDIR /run

CMD ["../main"]