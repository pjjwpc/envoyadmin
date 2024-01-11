FROM golang:1.20 as build
WORKDIR /src
COPY . .
RUN go env -w GOPROXY=https://goproxy.cn,direct
RUN go build

FROM alpine:latest as final
RUN apk --no-cache add ca-certificates libc6-compat

WORKDIR /work
COPY --from=build /src/control-plane .
RUN chmod +x /work/control-plane
ENTRYPOINT [ "/work/control-plane" ]
