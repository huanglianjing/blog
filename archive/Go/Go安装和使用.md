# 1. 安装

## 1.1 CentOS安装Go

```bash
$ yum install golang
```

## 1.2 macOS安装Go

macOS可以通过brew安装Go，需要先安装brew。

```bash
$ brew install go
```



# 2. 工具

Go的工具链通过go命令配合子命令使用。

常用命令：

```
build      编译包和依赖
clean      删除目标文件
doc        显示文档
env        显示go环境变量
fmt        格式化代码
get        下载安装包和依赖
install    编译安装包和依赖
list       列出包
run        编译运行程序
test       测试包
version    显示go版本信息
vet        运行工具vet
```

执行工具

```bash
$ go [commant] // 执行工具
$ go help [commant] // 查看工具文档
```

## run

将一个或多个.go为后缀的源文件进行编译、链接，然后运行生成的可执行文件。

```bash
$ go run helloworld.go
$ go run helloworld.go one two three # 带有运行参数
```

## build

将源文件编译输出成一个可执行的程序，该可执行文件可以直接执行。

```bash
$ go build helloworld.go
$ ./helloworld
```

## env

显示go环境变量。

```bash
$ go env # 显示所有go环境变量
$ go env GOPATH # 显示某个环境变量
```

环境变量GOPATH表示工作空间的根目录，其中有如下子目录：

```
GOPATH/
    src/  源文件
    bin/  可执行程序
    pkg/  编译的包
```

环境变量GOROOT指定Go发行版的根目录，其中提供所有标准库的包。

## get

下载一个包。

```bash
$ go get github.com/golang/lint/golint
$ go get -u github.com/golang/lint/golint # 获取包的最新版本
```

## list

列出可用的包。

```bash
$ go list
$ go list java... # 使用...作通配符匹配子串
```



# 参考

- [《Go程序设计语言》](https://book.douban.com/subject/27044219/)

