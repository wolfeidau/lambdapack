APPNAME = lambdapack-example

.PHONY: default
default: clean build archive package deploy

.PHONY: build
build:
	@mkdir -p dist
	@CGO_ENABLED=0 GOOS=linux GOARCH=arm64 \
		go build -ldflags "-s -w" -trimpath -o ./dist .

.PHONY: clean
clean:
	@rm -rf dist

.PHONY: archive
archive:
	lambdapack --debug dist dist

.PHONY: package
package:
	@echo "Packaging..."
	@aws cloudformation package \
		--template-file template.yaml \
		--output-template-file template.out.yaml \
		--s3-bucket $(DEPLOY_BUCKET) \
		--s3-prefix $(APPNAME)/latest

.PHONY: deploy
deploy:
	@echo "Deploying..."
	@aws cloudformation deploy --template-file template.out.yaml \
		--s3-bucket $(DEPLOY_BUCKET) \
		--s3-prefix $(APPNAME)/latest \
		--tags "service=$(APPNAME)" \
		--stack-name $(APPNAME) --capabilities CAPABILITY_IAM