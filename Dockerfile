FROM bitnami/golang:1.20 as builder
WORKDIR /home/build
COPY . .
RUN curl -fsSL https://esbuild.github.io/dl/latest | sh \ 
    && mv ./esbuild /bin/esbuild \
    && chmod +x /bin/esbuild \
    && go get . \
    && GOOS=linux GOARCH=amd64 go build -o stream-bot . \
    && chmod +x stream-bot \
    && esbuild --bundle --outfile=ui/index.js ui/src/index.ts

# FROM scratch as release
FROM bitnami/minideb as release
RUN apt update -y && apt install -y ca-certificates
WORKDIR /home/bot
COPY --from=builder /home/build/ui/index.js ui/index.js
COPY --from=builder /home/build/ui/index.css ui/index.css
COPY --from=builder /home/build/ui/index.html ui/index.html
COPY --from=builder /home/build/stream-bot stream-bot
ENTRYPOINT ["./stream-bot"]