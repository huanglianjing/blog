# 1. 环境配置

## 1.1 bash配置

在 ~/.bashrc 添加：

```
alias l='ls -l'
alias la='ls -la'
alias vi='vim'
export PS1="\[\e[32m\]\u\[\e[m\]@\[\e[32m\]\h\[\e[m\]:\[\e[31m\]\w\[\e[m\]$ "
```

## 1.2 vim配置

在 ~/.vimrc 添加：

```
syntax on
set number
set nocompatible
set showcmd
set t_Co=256
set tabstop=4
set softtabstop=4
set shiftwidth=4
set cursorline
set showmatch
set hlsearch
set incsearch
set ignorecase
set nowrapscan
set noswapfile
set autochdir
set list
set listchars=tab:>-,trail:-
set shortmess=atI
set backspace=indent,eol,start
if &diff
    colors blue
endif
```

# 2. 软件安装

## 2.1 git

```bash
yum -y install git
```

## 2.2 go

```bash
yum -y install golang
```

## 2.3 nginx

```bash
yum -y install nginx
systemctl start nginx
```

## 2.4 redis

```bash
yum -y install redis
systemctl start redis
```

