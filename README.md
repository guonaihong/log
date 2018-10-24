##### log

##### 主要功能如下
* 提供文件滚动，压缩，最多保存多个过期(压缩)文件
* 多级日志分类输出
* 定制多个输出源


##### `file_write.go`
主要提供文件滚动，压缩，最多保存多个过期文件的功能
``` golang
func NewFile(prefix string,           //设置保存文件名的前缀,为空的话使用默认前缀名
            dir string,               //设置默认文件名
            compress CompressType,    //是否压缩过期文件，目前只支持gzip压缩(log.Gzip) 不压缩(log.NotCompress)
            maxSize,                  //单个日志文件最大限制
            maxArchive int,           //最多保存多少个过期文件
            ) (f *File) {
}


func (f *File) Write(b []byte) (n int, err error) //写入数据

func (f *File) Close()       //关闭输出源
```

##### `log.go`
提供多级日志输出，可设置多个输出源
log.go默认一个输出源都没有, 连stdout或者stderr都不会带，使用时需填充一个
file_write.go 可和log.go组合使用，只要NewFile返回的输出源，填充至New第3个参数就行
``` golang
func New(level string,    //设置日志等级
            procName string, //设置每行日志的tag
            w ...io.Writer,  //设置多个输出源，如果要打印到stdout，这里就写os.Stdout
                             // log.go里面默认一个输出源都没有
) *Log


//F函数可以指定出错时打印哪一种调用栈,可和ID函数和等级函数组合使用
//特别是调用时包装了Error或者Warn函数时，可用F函数修改下需打印的调用栈
func (l *Log) F(frame int) *Log

//ID函数传递sessionID，可和等级日志函数组合使用
func (l *Log) ID(sessionID string) *Log

//输出debug等级日志
func (l *Log) Debugf(format string, a ...interface{}) 

//不带格式化功能的debug等级日志输出函数
func (l *Log) Debug(a ...interface{});

//输出info等级日志
func (l *Log) Infof(format string, a ...interface{}) 

//不带格式化功能的info等级日志输出函数
func (l *Log) Info(a ...interface{});

//输出warn等级日志
func (l *Log) Warnf(format string, a ...interface{})

//不带格式化功能的warn等级日志输出函数
func (l *Log) Warn(a ...interface{})

//输出error等日志
func (l *Log) Errorf(format string, a ...interface{}) 

//不带格式化功能的error等级日志输出函数
func (l *Log) Error(a ...interface{})
```

##### `tcp_udp_write.go`
ParseSocket可返回tcp 或者 udp 句柄, 可和log.go组合使用，只要返回的句柄填充至log.go第3参数即可
``` golang
func ParseSocket(url string) (io.Writer, error) (io.Writer, error)
```

##### `log.go example`
``` golang

	error2 := func(log *Log, a ...interface{}) {
		log.F(1).Error(a...)
	}

	l := New("debug", "test")

	l.AddWriter(os.Stdout)

	l.Debugf("hello world\n")
	l.Infof("hello world\n")
	l.Info("hello", " world\n")
	l.Warnf("hello world\n")
	l.Warn("hello", " world\n")
	l.Errorf("hello world\n")
	l.Error("hello", " world\n")

	error2(l /* *Log */, "hello2 world2\n")

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		l.ID("1").Debugf("hello world.1\n")
		l.ID("1").Infof("hello world.1\n")
		l.ID("1").Info("hello", " world.1\n")
		l.ID("1").Warnf("hello world.1\n")
		l.ID("1").Warn("hello", " world.1\n")
		l.ID("1").Errorf("hello world.1\n")
		l.ID("1").Error("hello", " world.1\n")
	}()

	l.ID("2").Debugf("hello world.2\n")
	l.ID("2").Infof("hello world.2\n")
	l.ID("2").Info("hello", " world.2\n")
	l.ID("2").Warnf("hello world.2\n")
	l.ID("2").Warn("hello", " world.2\n")
	l.ID("2").Errorf("hello world.2\n")
	l.ID("2").Error("hello", " world.2\n")
	wg.Wait()
```

输出:
```console
[test] [2018-10-21 14:24:43.530399] [debug] hello world
[test] [2018-10-21 14:24:43.530551] [info ] hello world
[test] [2018-10-21 14:24:43.530612] [info ] hello world
[test] [2018-10-21 14:24:43.530629] [warn ] [log_test.go:21] hello world
[test] [2018-10-21 14:24:43.530656] [warn ] [log_test.go:22] hello world
[test] [2018-10-21 14:24:43.530684] [error] [log_test.go:23] hello world
[test] [2018-10-21 14:24:43.530706] [error] [log_test.go:24] hello world
[test] [2018-10-21 14:24:43.530732] [error] [log_test.go:26] hello2 world2
[test] [2018-10-21 14:24:43.530750] [debug] <sid:2> hello world.2
[test] [2018-10-21 14:24:43.530770] [info ] <sid:2> hello world.2
[test] [2018-10-21 14:24:43.530786] [info ] <sid:2> hello world.2
[test] [2018-10-21 14:24:43.530801] [warn ] [log_test.go:44] <sid:2> hello world.2
[test] [2018-10-21 14:24:43.530822] [warn ] [log_test.go:45] <sid:2> hello world.2
[test] [2018-10-21 14:24:43.530844] [error] [log_test.go:46] <sid:2> hello world.2
[test] [2018-10-21 14:24:43.530864] [error] [log_test.go:47] <sid:2> hello world.2
[test] [2018-10-21 14:24:43.530893] [debug] <sid:1> hello world.1
[test] [2018-10-21 14:24:43.530918] [info ] <sid:1> hello world.1
[test] [2018-10-21 14:24:43.530953] [info ] <sid:1> hello world.1
[test] [2018-10-21 14:24:43.531009] [warn ] [log_test.go:35] <sid:1> hello world.1
[test] [2018-10-21 14:24:43.531026] [warn ] [log_test.go:36] <sid:1> hello world.1
[test] [2018-10-21 14:24:43.531041] [error] [log_test.go:37] <sid:1> hello world.1
[test] [2018-10-21 14:24:43.531063] [error] [log_test.go:38] <sid:1> hello world.1
```
