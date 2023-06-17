
## 01：搭建环境、解析ELF文件
本节主要是搭建环境：包括shell脚本生成riscv的对象文件作为rvld的输入文件；以及解析ELF文件的section header。  

ELF文件的格式如下面参考，program header主要用于加载到内存中使用的。section header主要是链接的过程使用的。elf header中的eshoff表示section header区域在elf文件中的偏移。e_shnum表示section header的个数。但是这个字段类型是uint16，在比较大的对象文件中，会溢出，这种情况下e_shnum设置为0，此时需要先把第一个section header都进来。第一个sh的size字段才是真正的sh的数量，这个字段的类型是uint64，足够表示了。这段逻辑体现在`pkg/linker/inputfile.go:40`。


最后解析出了所以的sh，放入数组之后，可以通过命令`readelf -S out/tests/hello/a.o` 查看对象文件的section header的数量。我这儿显示的是10个，在代码中assert结果是正确的。

![](https://note-img-1300721153.cos.ap-nanjing.myqcloud.com//md-img202306171632710.png)


#### 第一节参考
- https://en.wikipedia.org/wiki/Executable_and_Linkable_Format?oldformat=true
- https://github.com/corkami/pics/blob/master/binary/ELF101.png


---
