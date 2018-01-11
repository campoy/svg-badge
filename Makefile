serve:
	GOPATH=$(CURDIR) dev_appserver.py src/app
deploy:
	GOPATH=$(CURDIR) gcloud app deploy src/app