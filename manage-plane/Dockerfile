FROM golang:1.24 as build
WORKDIR /src
COPY . .
RUN go env -w GOPROXY=https://goproxy.cn,direct
RUN go build

FROM alpine:latest as final
RUN apk --no-cache add ca-certificates libc6-compat

WORKDIR /work
COPY --from=build /src/manage-plane .
RUN chmod +x /work/manage-plane
ENTRYPOINT [ "/work/manage-plane" ]
