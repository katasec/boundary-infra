.ONESHELL:
.DEFAULT_GOAL := help

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

infraup:  		## Does a Pulumi Up for the infrasutructure
	pulumi up -y -s infra
infrad: 	## Pulumi Destroy for the Pulumi Infrastructure
	pulumi destroy -s infra -f -y
	# pulumi stack rm infra -y
appup: 		## Configure Boundary App
	pulumi up -y -s app
appd: 	## Pulumi Destroy Boundary App config
	pulumi destroy -s app -f -y
	# pulumi stack rm app -y
build:  ## test build
	go build -o /dev/null