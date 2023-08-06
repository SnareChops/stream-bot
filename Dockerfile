FROM bitnami/golang:1.20 as builder
WORKDIR /home/build
COPY . .
RUN go get . \
    && GOOS=linux GOARCH=amd64 go build -o stream-bot . \
    && chmod +x stream-bot \
    && esbuild --bundle --outfile=ui/index.js ui/src/index.ts

FROM scratch as release
COPY --from=builder /home/build/ui/index.js ui/index.js
COPY --from=builder /home/build/ui/index.css ui/index.css
COPY --from=builder /home/build/ui/index.html ui/index.html
COPY --from=builder /home/build/stream-bot stream-bot
ENTRYPOINT stream-bot