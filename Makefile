clean:
	@echo "Cleaning..."
	@rm -rf build

compile:
	@echo "Compiling..."
	@cd src && \
	 go build -o ../build/otter && \
	 cd ..

compile_linux_64:
	@echo "Compiling for Linux 64..."
	@cd src && \
	 GOOS=linux GOARCH=amd64 go build -o ../build/otter-linux-amd64 && \
	 cd ..

compile_mac_64:
	@echo "Compiling for Mac 64..."
	@cd src && \
	 GOOS=darwin GOARCH=amd64 go build -o ../build/otter-mac-amd64 && \
	 cd ..

compile_all: compile_linux_64 compile_mac_64
