FROM golang:1.14.1-stretch as builder
WORKDIR /go/src/github.com/leakingtapan/sonoff
ADD . .
RUN make clean && make GO_BUILD_EXTR_ENV="GOARM=7 GOARCH=arm"

FROM scratch
COPY --from=builder /go/src/github.com/leakingtapan/sonoff/bin/sonoff /bin/sonoff
COPY --from=builder /go/src/github.com/leakingtapan/sonoff/certs /certs

CMD ["sonoff"]
