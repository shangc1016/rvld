# 
# 
# 

# 拿到所有的shell脚本
TESTS:=$(wildcard tests/*.sh)

# 编译出rvld可执行文件
build:
	go build

# 测试运行所有的shell脚本
test: build
	$(MAKE) $(TESTS)
	@printf "\e[32mPassed all tests\e[0m\n"

# 运行shell脚本
$(TESTS):
	@echo 'Testing' $@
	@./$@
	@printf '\e[32mOK\e[0m\n'

clean:
	go clean
	rm -rf out/

.PHONY: build clean $(TESTS)