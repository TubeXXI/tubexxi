FROM golang:1.25-bookworm

WORKDIR /workspace

ENV DEBIAN_FRONTEND=noninteractive

RUN apt-get update && apt-get install -y --no-install-recommends \
		ca-certificates \
		curl \
		python3 \
		python3-venv \
		python3-pip \
		unzip \
	&& rm -rf /var/lib/apt/lists/*

RUN curl -fsSL -o /tmp/protoc.zip \
		https://github.com/protocolbuffers/protobuf/releases/download/v29.6/protoc-29.6-linux-x86_64.zip \
	&& unzip -o /tmp/protoc.zip -d /usr/local \
	&& rm -f /tmp/protoc.zip

RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.36.11 \
	&& go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.6.1

RUN python3 -m venv /opt/venv \
	&& /opt/venv/bin/pip install --no-cache-dir \
		grpcio-tools==1.78.0 \
		protobuf==6.31.1

ENV PATH="/opt/venv/bin:/go/bin:${PATH}"

COPY generate-proto.sh /usr/local/bin/generate-proto
RUN chmod +x /usr/local/bin/generate-proto

ENTRYPOINT ["generate-proto"]
