##### log

##### 主要功能如下
* 提供文件滚动，压缩，最多保存多个过期(压缩)文件
* 多级日志分类输出
* 定制多个输出源


##### `file_write.go`
主要提供文件滚动，压缩，最多保存多个过期文件的功能
``` golang
func NewFile(prefix string,  //设置保存文件名的前缀,为空的话使用默认前缀名
            dir string,      //设置默认文件名
            compress int,    //是否压缩过期文件，目前只支持gzip压缩
            maxSize,         //单个日志文件最大限制
            maxArchive int,  //最多保存多少个过期文件
            ) (f *File) {
}


func (f *File) Write(b []byte) (n int, err error) //写入数据

func (f *File) Close()       //关闭输出源
```

##### `log.go`
提供多级日志输出，可设置多个输出源
log.go默认一个输出源都没有, 连stdout或者stderr都不会带，使用时需填充一个
file_write.go 可和log.go组合使用，只要NewFile返回的输出源，填充至NewLog第3个参数就行
``` golang
func NewLog(level string,    //设置日志等级
            procName string, //设置每行日志的tag
            w ...io.Writer,  //设置多个输出源，如果要打印到stdout，这里就写os.Stdout
                             // log.go里面默认一个输出源都没有
) *Log

//输出debug等级日志
func (l *Log) Debugf(format string, a ...interface{}) 

//输出info等级日志
func (l *Log) Infof(format string, a ...interface{}) 

//输出warn等级日志
func (l *Log) Warnf(format string, a ...interface{})

//输出error等日志
func (l *Log) Errorf(format string, a ...interface{}) 
```

##### `tcp_udp_write.go`
ParseSocket可返回tcp 或者 udp 句柄, 可和log.go组合使用，只要返回的句柄填充至log.go第3参数即可
``` golang
func ParseSocket(url string) (io.Writer, error) (io.Writer, error)
```
