serve:
	dev_appserver.py gae

deploy: copy-deps really-deploy rm-deps

really-deploy:
	GOPATH=$(CURDIR) gcloud app deploy src/app

ROOTS = $(wildcard $(CURDIR)/vendor/*)
copy-deps:
	for root in $(ROOTS); do\
		cp -r $$root src; \
	done

rm-deps:
	for root in $(ROOTS); do\
		echo $$root | xargs basename | sed "s/\(.*\)/src\/\1/" | xargs rm -rf; \
	done

update-deps:
	dep ensure
	dep prune