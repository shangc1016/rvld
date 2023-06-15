
#!/bin/bash


test_name=$(basename "$0" .sh)
t=out/tests/$test_name

mkdir -p  "$t"
cat << EOF | riscv64-unknown-elf-gcc -o "$t"/a.o -c -xc -
#include <stdio.h>

int main(int argc, char *argv[]) {
    printf("hello, world\n");
    return 0;
}

EOF

# 把C语言的object文件作为参数惨递给链接器rvld
./rvld "$t"/a.o