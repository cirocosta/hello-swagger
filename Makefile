install:
	go install -v


fmt:
	go fmt


gen: validate
	swagger generate server \
		-t ./swagger \
		-f ./swagger/swagger.yml \
		--exclude-main \
		-A hello


validate:
	swagger validate ./swagger/swagger.yml


.PHONY: install fmt gen validate
