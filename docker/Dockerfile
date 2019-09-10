FROM golang:1.12.0-stretch as build_base

RUN apt-get update \
	&& apt-get install -y \
		git \
		openssh-client \
	&& rm -rf /var/lib/apt/lists/*

ENV GO111MODULE=on \
	CGO_ENABLED=1 \
	GOOS=linux \
	GOARCH=amd64

WORKDIR /go/src/github.com/callistaenterprise/xapp

# docker cache go modules
COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN	go build -a \
	-o bin/xapp $PWD/cmd/xapp

FROM alpine:latest

RUN apk --update add \
    curl \
	tzdata \
	ca-certificates \
	&& rm -rf /var/cache/apk/*

ENV TZ=Europe/Stockholm

WORKDIR /app

COPY --from=build_base /go/src/github.com/callistaenterprise/xapp/bin/* /usr/local/bin/xapp

# Support static build docker binary
RUN mkdir /lib64 \
&& ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2

RUN mkdir /tmp/xapp-images

HEALTHCHECK --interval=3s --timeout=3s --start-period=3s --retries=5 CMD curl -sSf http://localhost:9090/health

CMD ["xapp"]
