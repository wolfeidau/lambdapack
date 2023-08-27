APPNAME = lambdapack-example

export CGO_ENABLED=0
export GOOS=linux
export GOARCH=arm64

.PHONY: default
default: clean build archive package deploy

.PHONY: build
build:
	@mkdir -p dist
	@go build -ldflags "-s -w" -trimpath -o ./dist .

.PHONY: clean
clean:
	@rm -rf dist

.PHONY: archive
archive:
	lambdapack dist/* dist

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