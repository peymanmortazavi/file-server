FROM golang:alpine3.14 AS builder

WORKDIR /project
COPY . .

RUN go mod download
RUN go build -o fsc -ldflags='-s -w' cmd/fsc/main.go

FROM builder AS fsc

RUN addgroup fsc -S && \
	adduser fsc -G fsc -S -H -D

USER fsc

COPY --from=builder --chown=fsc /project/fsc /opt/fileserver/fsc

CMD /opt/fileserver/fsc
