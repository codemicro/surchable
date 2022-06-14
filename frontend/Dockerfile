FROM golang:1 as builder

RUN mkdir /build 
ADD . /build/

WORKDIR /build 
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static" -s -w' -o main github.com/codemicro/surchable/crawler

FROM scratch

COPY --from=builder /build/main /
WORKDIR /run

CMD ["../main"]