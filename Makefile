build:
	go build .

clean:
	go clean .

ci : \
	fmt \
	clean \
	init \
	build

deploy:
	gcloud app deploy app.yaml \cron.yaml

fmt:
	go fmt .

init:
	go get -u google.golang.org/appengine/...
	go get -u cloud.google.com/go/...

.PHONY: \
	build \
	clean \
	ci \
	deploy \
	fmt \
	init \
	run \
