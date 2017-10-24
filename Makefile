install:
	go install -v

gen: validate
	swagger generate server \
		-t ./swagger \
		-f ./swagger/swagger.yml \
		--exclude-main \
		-A hello

validate:
	swagger validate ./swagger/swagger.yml

.PHONY: install gen validate
