
#!/bin/bash

# `test_name`拿到shell文件名   
test_name=$(basename "$0" .sh)
t=out/tests/$test_name

mkdir -p  "$t"
cat << EOF | $CC -o "$t"/a.o -c -xc -
#include <stdio.h>

int main(int argc, char *argv[]) {
    printf("hello, world\n");
    return 0;
}

EOF

# 把C语言的object文件作为参数惨递给链接器rvld
# ./rvld "$t"/a.o

$CC -B. -static "$t"/a.o -o "$t"/out

# -B.这个参数是给编译器的查找路径添加一个路径
# 然后编译器找linker的时候就会有现在当前目录查找名字是ld的linker
# 最后,我们把自己写的rvld软链接过来,取名ld就行