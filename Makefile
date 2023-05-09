clean:
	@echo "Cleaning..."
	@rm -rf build

compile:
	@echo "Compiling..."
	@cd src && \
	 go build -o ../build/otter && \
	 cd ..
