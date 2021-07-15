# 1. 基础概念

每个项目整个文件夹称为一个仓库repository，仓库中包含零到多个文件，仓库的改动由链状的提交commit所表示，这些提交又分布在不同的分支branch上。每个仓库都有一个起始默认的仓库master/main，然后新建和删除不同的分支，分支之间可以相互合并，以达到多人共同开发同一个项目。

HEAD表示当前指向的提交。

### 远程仓库

origin表示远程仓库，每个仓库包含本地仓库和对应的远程仓库，例如对应为dev和origin/dev，可以从远程仓库分支拉取更新到本地，以及将本地提交的修改推送到远程分支。

### 文件状态



# 2. 配置

## 2.1 GitHub配置密钥

## 2.2 配置文件

### .gitconfig

### .gitignore

### .git-credentials



# 3. 命令

Git的命令有很多，熟练使用常用的即可，并且通过桌面软件如GitHub Desktop，也可以方便地完成很多操作。

通过`git <cmd> -h`可以看到对应命令的参考文档。



## 3.1 仓库

### init

在当前文件夹创建仓库，并生成.git/管理版本库。

```bash
$ git init
```

### clone

从远端克隆仓库到本地。

```bash
$ git clone <url> #克隆仓库
$ git clone -b <branch> <url> #克隆某个分支
```

### remote

```bash
$ git remote #远程库名称，一般为origin
$ git remote -v #查看远程库地址
$ git remote add <name> <url> #添加新的远程仓库
```



## 3.2 拉取上传

### pull

拉取远程仓库代码更新，合并到本地。

```bash
$ git pull
$ git pull --rebase #变基后再合并
```

### fetch

拉取远程仓库代码更新到对应的远程分支，但不合并到本地分支。

```bash
$ git fetch
$ git fetch <remote> #从远程库抓取
```

### push



## 3.3 分支

### branch

### merge

将某个分支合并到本分支。

```bash
$ git merge <branch> 
```

### rebase

### cherry-pick

### cherry



## 3.4 状态

### status

### log

### reflog

### diff

### difftool



## 3.5 编辑提交

### add

### rm

### mv

### commit

### checkout

### reset

### clean

### stash

### blame

### whatchanged



## 3.6 标签

### tag

### show



# 4. 参考

- [Git - Reference](https://git-scm.com/docs)
- [Git - Book](https://git-scm.com/book/zh/v2)

