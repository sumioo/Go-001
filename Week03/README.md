# Go 语言实践 - concurrency

## concurrency

在go语言中进行异步调用是非常简单的，只需要在调用函数前加上go关键字，这样函数就以goroutine的方式运行了，
gouroutine所需的资源相比于线程是非常小的，而且创建也非常快，调度也是在用户态进行的，可以说在资源和性能上都比以线程的方式
进行异步调用要好。尽管创建goroutine是如此简单，并且性价比也高，但也不能随意创建，因为goroutine是没有垃圾回收机制的，很容易就
造成goroutine泄露。

在使用goroutine时一定要问自己三个问题

### 1.作为api的设计者不应该假设调用者的调用行为
api的设计者不应该假设调用者的使用方式是异步还是同步，决定权应该由调用者决定。

```go
func ListDirectory(dir string) ([]string, error)

// 程序默认以异步的方式执行
func ListDirectory(dir string) chan string

//对每一个目录或文件都会调用walkFunc,当walkFunc返回费nil error时停止遍历，这样调用者可以确定何时停止调用
//当然，调用者可以选择异步执行
func ListDirectory(dir string, walkFunc func(path string, info os.FileInfo, err error) error) error
```

### 2. goroutine何时退出

never start a goroutine without knowning when it stop

```go
//不会退出，也无法从外部控制退出
func leak() {
    ch := make(chan int)

    go func() {
        val := <-ch
        fmt.Println("receive a value:", val)
    }
}
```

### 如何从外部控制goroutine退出

```go
func serve(addr string, handler http.Handler, stop <-chan struct{}) error {
    s := http.Server {
        Addr: addr,
        Handler: handler
    }

    go func() {//起一个goroutine来监听消息
        <-stop
        s.Shutdown(context.Background())
    }()

    return s.ListenAndServe()
}

func main() {
    done := make(chan error, 2)
    stop := make(chan struct{})
    //只要serveDebug，serveApp任意一个返回错误的话都会往done channel发送消息
    go func(){
        done <- serveDebug(addr, handler, stop)
    }()
     go func(){
        done <- serveApp(addr, handler, stop)
    }()

    var stopped bool
    for i := 0; i < cap(done); i++ {
        if err := <-done; err != nil { //done channel接收到消息后继续执行
            fmt.Println("error: %v", err)
        }

        if !stopped {
            stopped = true
            close(stop) //发送close消息,停止全部
        }
    }
}

//这里的示例有一个缺陷是无法强制退出（控制超时）
```

```go
type Tracker struct{}

func (t *Tracker) Event(data string) {
    time.Sleep(time.Millisecond)
    log.Pringtln(data)
}

type App struct {
    track Tracker
}

func (a *App) Handle(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusCreated)
    
    //Bug: 该goroutine不受管控
    go a.track.Event("this event")
}
```

一般使用context和sync.WaitGroup来管控goroutine

sync.WaitGroup 实现平滑退出
context 实现从外部控制

```go
package main

import (
	"context"
	"fmt"
	"time"
)

type Tracker struct {
	ch   chan string
	stop chan struct{}
}

func (t *Tracker) Event(ctx context.Context, data string) error {
	select {
	case t.ch <- data:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (t *Tracker) Run() {
	for data := range t.ch {
		time.Sleep(1 * time.Second)
		fmt.Println(data)
	}

	t.stop <- struct{}{}
}

func (t *Tracker) Shutdown(ctx context.Context) {
	close(t.ch)
	select {
	case <-t.stop:
	case <-ctx.Done():
	}
}

func NewTracker() *Tracker {
	return &Tracker{ch: make(chan string, 10)}
}

func main() {
	t := NewTracker()
	go t.Run()
	t.Event(context.Background(), "test1")
	t.Event(context.Background(), "test2")
	t.Event(context.Background(), "test3")
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(5*time.Second))
	defer cancel()
	t.Shutdown(ctx)

}

```

## channel 通道

### unbuffered channel

无缓冲信道的本质是保证同步。

发送方发送时无接收者时，发送方阻塞
接收方接收时无发送方，接收方阻塞

*If the channel is unbuffered, the sender blocks until the receiver has received the value*

### buffered channel

缓冲区满时发送方将阻塞，缓冲区空时接收方将阻塞。

buffer == 1 的buffered channel跟unbuffered channel的行为是不同的。


### Context

context可用来协调多个goroutine，控制它们的生命周期，比如级联取消，同步超时。也用来
在api、函数调用间传递元数据，比如request id

使用建议
1. Do not store Contexts inside a struct type; instead, pass a Context explicitly to each function that needs it. The Context should be the first parameter, typically named ctx
2. Do not pass a nil Context, even if a function permits it. Pass context.TODO if you are unsure about which Context to use.
3. Use context Values only for request-scoped data that transits processes and APIs, not for passing optional parameters to functions.
4. The same Context may be passed to functions running in different goroutines; Contexts are safe for simultaneous use by multiple goroutines.