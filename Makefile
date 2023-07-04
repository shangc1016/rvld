# 
# 
# 
VERSION=0.1.0
# 拿到上一个commit的哈希
COMMIT_ID=$(shell git rev-list -1 HEAD)

# 拿到所有的shell脚本
TESTS:=$(wildcard tests/*.sh)

# 编译出rvld可执行文件
build:
# -ldflags 注入参数，也就是main包中的version变量
	@go build -ldflags "-X main.version=${VERSION}-${COMMIT_ID}"
	@if [ -f ld ]; then rm ld; fi;
	@ln -s rvld ld

# 测试运行所有的shell脚本
test: build
	@CC="riscv64-unknown-linux-gnu-gcc" \
	$(MAKE) $(TESTS)
	@printf "\e[32mPassed all tests\e[0m\n"

# 运行shell脚本
$(TESTS):
	@echo 'Testing' $@
	@./$@
	@printf '\e[32mOK\e[0m\n'

clean:
	go clean
	rm -rf out/ ld

.PHONY: build clean $(TESTS)