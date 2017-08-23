all: build

build:
	cd bin/terraform-provider-aviatrix; go install
