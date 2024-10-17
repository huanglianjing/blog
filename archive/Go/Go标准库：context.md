# 1. 介绍

文档地址：https://pkg.go.dev/context

Go 标准库 context 用于在函数调用、API调用、进程调用间传递上下文，context 在多个协程中同时使用是安全的。context 的主要作用是在多个 goroutine 组成的树中同步取消信号，以减少对资源的消耗和占用，其次作用是传递上下文的值。

Go 官方建议不要将 context 定义在其他结构体中，将 context 作为独立的参数传递给需要的函数且应该是首个参数，变量名一般使用 ctx。示例如下：

```go
func DoSomething(ctx context.Context, arg Arg) error {
	// ... use ctx ...
}
```

通过创建一个包含过期时间的 context，传入创建的 goroutine，以控制下层 goroutine 的工作时间，防止下层协程运行时间超过指定时间。如下例子中，handle 协程在900毫秒后就会返回，而在程序运行1秒后 main 函数的 select 才会收到 context 结束的通道信号。如果创建协程时传的时间大于 context 的过期时间，整个程序都会因为上下文过期而终止。

```go
func handle(ctx context.Context, duration time.Duration) {
	select {
	case <-ctx.Done():
		fmt.Println("handle", ctx.Err())
	case <-time.After(duration):
		fmt.Println("process request with", duration)
	}
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	go handle(ctx, 900*time.Millisecond)
	select {
	case <-ctx.Done():
		fmt.Println("main", ctx.Err())
	}
}
```

context 的默认创建函数是 context.Background() 和 context.TODO()，它们都会返回预先初始化的私有变量，context.Background 是上下文默认值，其它上下文应该由它衍生，context.TODO 应该在不确定用哪个上下文时使用。

可以通过 WithValue() 给 context 对象设置 key/value 内容，衍生出新的子上下文，其它地方可以通过 Value() 从 context 对象取出 key 对应的 value。

```go
func handle(ctx context.Context) {
	name := ctx.Value("name")
	fmt.Println("name is", name)
}

func main() {
	ctx := context.Background()
	ctx = context.WithValue(ctx, "name", "moondo")
	handle(ctx)
}
```

# 2. 类型

## 2.1 Context

类型定义：

```go
// Context 上下文
type Context interface {
	Deadline() (deadline time.Time, ok bool) // 被取消的时间
	Done() <-chan struct{} // 返回一个通道，context被取消时关闭
	Err() error // Done未关闭时返回nil，关闭时返回关闭原因
	Value(key any) any // 获取key对应的val
}
```

方法：

```go
// Background 返回一个空的context，不会被取消，没有值，没有deadline，通常用于主函数、初始化、测试、最上层接收请求时
func Background() Context

// TODO 返回一个空的context，不确定会用到context时使用
func TODO() Context

// WithValue 返回父context的复制，并关联key-val
func WithValue(parent Context, key, val any) Context
```

## 2.2 CancelCauseFunc

类型定义：

```go
type CancelCauseFunc func(cause error)
```

## 2.3 CancelFunc

类型定义：

```go
type CancelFunc func()
```

# 3. 函数

```go
// Cause 返回解释为什么context被取消
func Cause(c Context) error

// WithCancel 返回一个父context的复制，父context取消或者cancel被调用都会取消当前context
func WithCancel(parent Context) (ctx Context, cancel CancelFunc)
func WithCancelCause(parent Context) (ctx Context, cancel CancelCauseFunc)

// WithTimeout 带有过期时间
func WithTimeout(parent Context, timeout time.Duration) (Context, CancelFunc)

// WithValue 返回父context的复制，并关联key-val，用于传递过程和API，不要用于传递多个kv值以读取，因为每个kv的储存和查询是顺着父context一层层记录的
func WithValue(parent Context, key, val any) Context
```

