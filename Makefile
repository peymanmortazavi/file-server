IMAGE_REPO := peymanmo
SERVER_IMAGE_NAME := fs-server
CLIENT_IMAGE_NAME := fsc
QUERY := 

FS_IMAGE := ${IMAGE_REPO}/${SERVER_IMAGE_NAME}
FC_IMAGE := ${IMAGE_REPO}/${CLIENT_IMAGE_NAME}

TARGETS = fs-server fsc

$(TARGETS):
	go build -o bin/$@ cmd/$@/main.go

clean:
	rm -fr bin

build-server-image:
	docker build -t ${FS_IMAGE} -f docker/fs-server/Dockerfile .

run-server-image: build-server-image
	docker run -v $$(pwd)/test:/test:rw -p 6000:6000 -u root ${FS_IMAGE} /opt/fileserver/fs-server --root /test

build-client-image:
	docker build -t ${FC_IMAGE} -f docker/fsc/Dockerfile .

run-client-image: build-client-image
	docker run --network host ${FC_IMAGE} /opt/fileserver/fsc --insecure ${QUERY}

push-server-image:
	docker push ${FS_IMAGE}

push-client-image:
	docker push ${FC_IMAGE}
