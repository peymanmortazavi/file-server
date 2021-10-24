TARGETS = fs-server fsc

$(TARGETS):
	go build -o bin/$@ cmd/$@/main.go

clean:
	rm -fr bin

build-server-image:
	docker build -t fs-server -f docker/fs-server/Dockerfile .

run-server-image: build-server-image
	docker run -v $$(pwd)/test:/test:rw -p 6000:6000 fs-server /opt/fileserver/fs-server --root /test

build-client-image:
	docker build -t fsc -f docker/fsc/Dockerfile .

run-client-image: build-client-image
	docker run --network host fsc /opt/fileserver/fsc --insecure
