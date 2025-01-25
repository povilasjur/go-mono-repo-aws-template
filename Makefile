# Define variables
SERVICES=$(shell ls -d services/*/)
BUILD_DIR=target

# Define the binary output for each service
BINARIES=$(patsubst services/%/,$(BUILD_DIR)/%/bootstrap,$(SERVICES))

# Define the zip output for each service
ZIPS=$(patsubst $(BUILD_DIR)/%/bootstrap,$(BUILD_DIR)/%-service.zip,$(BINARIES))

# Default target runs everything
all: clean build package deploy

# Clean up build directories and zips
clean:
	@echo "Cleaning up..."
	@rm -rf $(BUILD_DIR)
	@find services -name '*.zip' -exec rm -f {} \;

# Build each service binary only if the binary doesn't exist
$(BUILD_DIR)/%/bootstrap:
	@echo "Testing $*..."
	@cd services/$* && go test ./... || exit 1
	@echo "Building $*..."
	@mkdir -p $(BUILD_DIR)/$*
	@cd services/$* && GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ../../$(BUILD_DIR)/$*/bootstrap

# Build target now depends on the binaries
build: $(BINARIES)

# Package each service as a zip file (only if the binary exists)
package: $(ZIPS)

$(BUILD_DIR)/%-service.zip: $(BUILD_DIR)/%/bootstrap
	@echo "Packaging $*..."
	@cd $(BUILD_DIR)/$* && zip -j ../$*-service.zip bootstrap

# Deploy each service to S3 (ensure binaries and zip files are present)
deploy: $(ZIPS)
	@echo "All binaries and zip files are ready. Proceeding with deployment..."
	@echo "Terraform scripts are not provided in this sample"
	# cd infra/terraform && terraform init && terraform apply -auto-approve
