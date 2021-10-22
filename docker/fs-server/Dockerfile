FROM golang:alpine3.14 AS builder

WORKDIR /project
COPY . .

RUN go mod download
RUN go build -o fs-server -ldflags='-s -w' cmd/fs-server/main.go

FROM builder AS fs-server

RUN addgroup fs-server -S && \
	adduser fs-server -G fs-server -S -H -D

USER fs-server

COPY --from=builder --chown=fs-server /project/fs-server /opt/fileserver/fs-server

CMD /opt/fileserver/fs-server
