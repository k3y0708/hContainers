SHELL := /bin/bash

# Clean Directory -------------------------------------------------------------
clean:
	@echo "Cleaning..."
	@rm -rf dist
	@rm -rf build

# Compile (build without version number) --------------------------------------
compile:
	@echo "Compiling..."
	@cd src && \
	 go build -o ../build/hContainers && \
	 cd ..

compile_windows_64:
	@echo "Compiling for Windows 64..."
	@cd src && \
	 GOOS=windows GOARCH=amd64 go build -o ../build/hContainers-amd64.exe && \
	 cd ..

compile_windows_32:
	@echo "Compiling for Windows 32..."
	@cd src && \
	 GOOS=windows GOARCH=386 go build -o ../build/hContainers-386.exe && \
	 cd ..

compile_linux_64:
	@echo "Compiling for Linux 64..."
	@cd src && \
	 GOOS=linux GOARCH=amd64 go build -o ../build/hContainers-linux-amd64 && \
	 cd ..

compile_linux_32:
	@echo "Compiling for Linux 32..."
	@cd src && \
	 GOOS=linux GOARCH=386 go build -o ../build/hContainers-linux-386 && \
	 cd ..

compile_mac_64:
	@echo "Compiling for Mac 64..."
	@cd src && \
	 GOOS=darwin GOARCH=amd64 go build -o ../build/hContainers-mac-amd64 && \
	 cd ..

compile_all: compile_windows_32 compile_windows_64 compile_linux_32 compile_linux_64 compile_mac_64

# Deploy (build with version number) ------------------------------------------
checkIfDeployable:
	@if [[ "${shell git describe --tags}" =~ ^v[0-9].[0-9].[0-9]$$ ]]; then \
		echo "Deploying..."; \
	 else \
		echo "Not deploying..."; \
		exit 1; \
	 fi

deploy: checkIfDeployable
	@echo "${shell git describe --tags}"
	@cd src && \
	 go build -ldflags \
		"-X github.com/hContainers/hContainers/global.Version=${shell git describe --tags}" \
		 -o ../dist/hContainers && \
	 cd ..

deploy_windows_64: checkIfDeployable
	@echo "Deploying for Windows 64..."
	@cd src && \
	 GOOS=windows GOARCH=amd64 go build -ldflags \
		"-X github.com/hContainers/hContainers/global.Version=${shell git describe --tags}" \
		 -o ../dist/hContainers-amd64.exe && \
	 cd ..

deploy_windows_32: checkIfDeployable
	@echo "Deploying for Windows 32..."
	@cd src && \
	 GOOS=windows GOARCH=386 go build -ldflags \
		"-X github.com/hContainers/hContainers/global.Version=${shell git describe --tags}" \
		 -o ../dist/hContainers-386.exe && \
	 cd ..

deploy_linux_64: checkIfDeployable
	@echo "Deploying for Linux 64..."
	@cd src && \
	 GOOS=linux GOARCH=amd64 go build -ldflags \
		"-X github.com/hContainers/hContainers/global.Version=${shell git describe --tags}" \
		 -o ../dist/hContainers-linux-amd64 && \
	 cd ..

deploy_linux_32: checkIfDeployable
	@echo "Deploying for Linux 32..."
	@cd src && \
	 GOOS=linux GOARCH=386 go build -ldflags \
		"-X github.com/hContainers/hContainers/global.Version=${shell git describe --tags}" \
		 -o ../dist/hContainers-linux-386 && \
	 cd ..

deploy_mac_64: checkIfDeployable
	@echo "Deploying for Mac 64..."
	@cd src && \
	 GOOS=darwin GOARCH=amd64 go build -ldflags \
		"-X github.com/hContainers/hContainers/global.Version=${shell git describe --tags}" \
		 -o ../dist/hContainers-mac-amd64 && \
	 cd ..

deploy_all: deploy_windows_32 deploy_windows_64 deploy_linux_32 deploy_linux_64 deploy_mac_64
	@echo "Deployed for all platforms"
