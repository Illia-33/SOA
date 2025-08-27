FROM golang:1.24

ARG PROTOC_VERSION=25.1

RUN apt-get update
RUN apt-get install -y unzip curl

# Install protoc
RUN curl -LO https://github.com/protocolbuffers/protobuf/releases/download/v${PROTOC_VERSION}/protoc-${PROTOC_VERSION}-linux-x86_64.zip \
  && unzip protoc-${PROTOC_VERSION}-linux-x86_64.zip -d /usr/local \
  && rm -f protoc-${PROTOC_VERSION}-linux-x86_64.zip
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.31 \
  && go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3


WORKDIR /source

COPY --from=go-mods . .

RUN go mod download
RUN make autogen

ENTRYPOINT [ "/bin/bash" ]
