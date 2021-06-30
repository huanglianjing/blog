# 1. yum

以java为例，先通过命令或者rpm命令查看是否未安装

```bash
$ java

$ rpm -qa | grep java
```

列出yum源可安装的包

```bash
$ yum list java*
```

选择合适的包安装

```bash
$ yum -y install java-1.8.0-openjdk.x86_64
```

检查安装成功

```bash
$ java -version
```

