clean:
	@echo "Cleaning..."
	@rm -rf build

compile:
	@echo "Compiling..."
	@cd src && \
	 go build -o ../build/otter && \
	 cd ..

compile_windows_64:
	@echo "Compiling for Windows 64..."
	@cd src && \
	 GOOS=windows GOARCH=amd64 go build -o ../build/otter-amd64.exe && \
	 cd ..

compile_windows_32:
	@echo "Compiling for Windows 32..."
	@cd src && \
	 GOOS=windows GOARCH=386 go build -o ../build/otter-386.exe && \
	 cd ..

compile_linux_64:
	@echo "Compiling for Linux 64..."
	@cd src && \
	 GOOS=linux GOARCH=amd64 go build -o ../build/otter-linux-amd64 && \
	 cd ..

compile_linux_32:
	@echo "Compiling for Linux 32..."
	@cd src && \
	 GOOS=linux GOARCH=386 go build -o ../build/otter-linux-386 && \
	 cd ..

compile_mac_64:
	@echo "Compiling for Mac 64..."
	@cd src && \
	 GOOS=darwin GOARCH=amd64 go build -o ../build/otter-mac-amd64 && \
	 cd ..

compile_mac_arm:
	@echo "Compiling for Mac ARM..."
	@cd src && \
	 GOOS=darwin GOARCH=arm64 go build -o ../build/otter-mac-arm64 && \
	 cd ..

compile_all: compile_windows_32 compile_windows_64 compile_linux_32 compile_linux_64 compile_mac_64 compile_mac_arm
