FROM golang:1-alpine AS builder

RUN apk add --no-cache build-base
WORKDIR /build
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN go build -ldflags='-s -w' -trimpath -o /dist/app
RUN ldd /dist/app | tr -s [:blank:] '\n' | grep ^/ | xargs -I % install -D % /dist/%
RUN ln -s ld-musl-x86_64.so.1 /dist/lib/libc.musl-x86_64.so.1

FROM scratch
ENV FILE_PATH='.'
ENV DB_PATH='gallery.db'
COPY --from=builder /dist /
USER 65534
ENTRYPOINT ["/app"]
