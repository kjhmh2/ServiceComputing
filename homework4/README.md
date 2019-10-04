# homework4 用 golang 开发selpg

### 参考资料

[开发Linux命令行实用程序](https://www.ibm.com/developerworks/cn/linux/shell/clutil/index.html)

[CLI 命令行实用程序开发基础](https://blog.csdn.net/C486C/article/details/82990187)

### 程序要求
selpg从标准输入或从作为命令行参数给出的文件名读取文本输入。它允许用户指定来自该输入并随后将被输出的页面范围。
#### 输入

在命令行上指定的文件名。例如：

```bash
$ command input_file
```

标准输入：

```bash
$ command
```

使用 shell 操作符“<”（重定向标准输入），也可将标准输入重定向为来自文件：

```bash
$ command < input_file
```

使用 shell 操作符“|”（pipe）也可以使标准输入来自另一个程序的标准输出：

```bash
$ other_command | command
```


#### 输出

输出应该被写至标准输出，缺省情况下标准输出同样也是终端（也就是用户的屏幕）：

```bash
$ command
```

使用 shell 操作符“>”（重定向标准输出）可以将标准输出重定向至文件：

```bash
$ command > output_file
```

使用“|”操作符，command 的输出可以成为另一个程序的标准输入：

```bash
$ command | other_command
```


#### 命令行参数

##### “-sNumber”和“-eNumber”强制选项：

selpg 要求用户用两个命令行参数“-sNumber”（例如，“-s10”表示从第 10 页开始）和“-eNumber”（例如，“-e20”表示在第 20 页结束）指定要抽取的页面范围的起始页和结束页。selpg 对所给的页号进行合理性检查；换句话说，它会检查两个数字是否为有效的正整数以及结束页是否不小于起始页。这两个选项，“-sNumber”和“-eNumber”是强制性的，而且必须是命令行上在命令名 selpg 之后的头两个参数：

##### “-lNumber”和“-f”可选选项：

selpg 可以处理两种输入文本：

**类型 1**：该类文本的页行数固定（每页 72 行）。这是缺省类型，因此不必给出选项进行说明。该缺省值可以用“-lNumber”选项覆盖，如下所示：

```bash
$ selpg -s10 -e20 -l66 
```
这表明页有固定长度，每页为 66 行。

**类型 2**：该类型文本的页由 ASCII 换页字符（十进制数值为 12，在 C 中用“\f”表示）定界。类型 2 格式由“-f”选项表示，如下所示：

```bash
$ selpg -s10 -e20 -f ...
```
该命令告诉 selpg 在输入中寻找换页符，并将其作为页定界符处理。
注：“-lNumber”和“-f”选项是互斥的。

##### “-dDestination”可选选项：

selpg 还允许用户使用“-dDestination”选项将选定的页直接发送至打印机。这里，“Destination”应该是 lp 命令“-d”选项（请参阅“man lp”）可接受的打印目的地名称。该目的地应该存在——selpg 不检查这一点。

```bash
$ lp -dDestination
```

的管道以便输出，并写至该管道而不是标准输出：

```bash
selpg -s10 -e20 -dlp1
```

该命令将选定的页作为打印作业发送至 lp1 打印目的地。



# 程序设计

### 参数结构体

设计一个结构体记录存放的参数。包括开始页、结束页、输入文件的名字、一页的行数、分页方式

```go
type selgp_args struct{
	start_page int                   //开始页数  -s1
	end_page int                     //结束页数  -e5
	input_file string                //输入文件  [input.txt]
	page_type string                  //指定每页行数-l10 或 换行符-f
	page_len int                     //每页多少行
	destination string               //打印目的（打印机）
}
```
设置全局变量记录读取到的参数

```go
var args selgp_args
```



### 主函数

包括三个部分

```go
func main() {
    get(&args)          //读取并处理参数
    check(&args)             //检查参数是否合法
    run(&args)               //运行
}
```



### 读取并处理参数

#### Golang 的支持

使用os，flag包，最简单处理参数的代码

```go
package main

import (
    "fmt"
    "os"
)

func main() {
    for i, a := range os.Args[1:] {
        fmt.Printf("Argument %d is %s\n", i+1, a)
    }

}
```

可以获取命令行结果如下

![0](https://github.com/JaneYul/ServiceComputing/raw/master/homework4/pics/0.PNG)

使用flag包的代码：

```go
package main

import (
    "flag" 
    "fmt"
)

func main() {
    var port int
    flag.IntVar(&port, "p", 8000, "specify port to use.  defaults to 8000.")
    flag.Parse()

    fmt.Printf("port = %d\n", port)
    fmt.Printf("other args: %+v\n", flag.Args())
}
```

![0.1](https://github.com/JaneYul/ServiceComputing/raw/master/homework4/pics/0.1.PNG)

[参考资料]([http://blog.studygolang.com/2013/02/%E6%A0%87%E5%87%86%E5%BA%93-%E5%91%BD%E4%BB%A4%E8%A1%8C%E5%8F%82%E6%95%B0%E8%A7%A3%E6%9E%90flag/](http://blog.studygolang.com/2013/02/标准库-命令行参数解析flag/))

#### pflag包

使用[**os.Args**](https://blog.csdn.net/guanchunsheng/article/details/79612153)获取用户输入的参数，得到一个包含参数的数组，数组中的每个元素是string类型的。

参数处理可以使用[**pflag**](https://godoc.org/github.com/spf13/pflag)包来解析命令的参数。

##### 安装

要import pflag包，要先在本地安装pflag

```bash
go get github.com/spf13/pflag
```
注意安装后使用import的是github.com/spf13/pflag

##### 提示信息

如果给出的参数不正确或者需要查看帮助 -help，那么会给出这里指定的字符串

```
pflag.Usage = show_tips
```

##### 绑定参数

```go
	//func IntVarP(指向int变量的指针，名称，简写字符串，值，用法字符串)
	pflag.IntVarP(&sa.start_page,"start_page", "s", 0, "Start page of file")
	pflag.IntVarP(&sa.end_page,"end_page","e", 0, "End page of file")
	pflag.IntVarP(&sa.page_len,"page_len", "l", 72, "lines in one page")  //默认72行一页
	//func StringVarP(指向字符串变量的指针，名称，简写字符串，值，用法字符串)
	pflag.StringVarP(&sa.page_type,"page_type","f", "l", "flag splits page")
	pflag.StringVarP(&sa.destination, "destination","d", "", "name of printer")
```

##### 获取剩下的参数

除了上述指定的参数，还有一个文件名参数是没有指定的，需要另外获取

```go
//获取剩下的输入文件参数（可能没有）
other_args := pflag.Args()
if len(other_args) > 0 {
	args.input_file = other_args[0]
} else {
	args.input_file = ""
}
```

##### 解析

使用**pflag.Parse()**让pflag 将命令行解析为定义的标志。之后就可以直接使用绑定的值。

```go
pflag.Parse()
```
解析完以上参数后，可以通过调用pflag.Args()来获得未定义但输入了的参数，这里是文件名。

#### **另外**

如果手动判断每个参数，则需要判断参数个数、参数格式是否符合要求。不符合要求则输出错误信息，终止程序；符合要求则赋给参数结构体的每个变量中。

此时需要通过 **a := os.Args** 先获取命令行中的参数并作为参数传入 get 函数中。需要程序名变量progname和记录参数个数argcount int

##### 参数少于3个

```go
if len(a) < 3 {
	fmt.Fprintf(os.Stderr, "%s: the num arguments is less 3\n", progname)
	show_tips()
	os.Exit(1)
}
```

##### 第一个参数：开始页数

先判断第一个参数是不是"-s"形式

```go
if a[1][0] != '-' || a[1][1] != 's' {
	fmt.Fprintf(os.Stderr, "%s: 1st arg should be -s[start_page]\n", progname)
	show_tips()
	os.Exit(1)
}
```
如果是，则从数组中提取出参数
```go
sp, _ := strconv.Atoi(a[1][2:])
if sp < 1 {
	fmt.Fprintf(os.Stderr, "%s: start page should not be less than 1 %d\n", progname, sp)
	show_tips()
	os.Exit(1)
}
args.start_page = sp
```

##### 第二个参数：结束页数

和start_page很类似，先判断是否符合"-e"格式再提取出参数

```go
if a[2][0] != '-' || a[2][1] != 'e' {
	fmt.Fprintf(os.Stderr, "%s: 2nd arg should be -e[end_page]\n", progname)
	show_tips()
	os.Exit(1)
}
//提取结束页数
ep, _ := strconv.Atoi(a[2][2:])
if ep < 1 || ep < sp {
	fmt.Fprintf(os.Stderr, "%s: end page should not be less than 1 %d\n", progname, ep)
	show_tips()
	os.Exit(1)
}
args.end_page = ep
```

##### option参数

如果参数个数大于三，需要获取其他参数信息，可能是"-l", "-f", "-d"。需要确认是以"-"开头才能进行参数判断

- 循环读取

```go
for {
	if argindex > argcount-1 || a[argindex][0] != '-' {
		break
	}
	switch a[argindex][1] {
```
- -l
  判断-l后面跟着的数字是否符合格式

```go
case 'l':
	pl, _ := strconv.Atoi(a[argindex][2:])
	if pl < 1 {
		fmt.Fprintf(os.Stderr, "%s: page length should not be less than 1 %d\n", progname, pl)
		show_tips()
		os.Exit(1)
	}
	args.page_len = pl
	argindex++
```
- -f
  -f参数后面不跟数字，所以判断-f参数的长度来判断是否合法

```go
case 'f':
	if len(a[argindex]) > 2 {
		fmt.Fprintf(os.Stderr, "%s: option should be \"-f\"\n", progname)
		show_tips()
		os.Exit(1)
	}
	args.page_type = 'f'
	argindex++
```
- -d
  根据说明，selpg不检查destination目的地，但是要确保-d后面跟着目的地

```go
case 'd':
	if len(a[argindex]) <= 2 {
		fmt.Fprintf(os.Stderr, "%s: -d option requires a printer destination\n", progname)
		show_tips()
		os.Exit(1)
	}
	args.destination = a[argindex][2:]
	argindex++
```
- 不合法参数
  如果不是以上的任何一种情况，需要报错处理

```go
default:
	fmt.Fprintf(os.Stderr, "%s: unknown option", progname)
	show_tips()
	os.Exit(1)
```

- 输入文件
在上面的循环中，没有以“-”开头的参数，跳出。
此时可能还有文件作为输入，或者没有（此时为标准输入），需要判断
```go
if argindex <= argcount-1 {
	args.input_file = a[argindex]
}
```



### 检查参数

为了程序运行能够正确运行，需要检查参数是否合法，包括开始页数、结束页数和长度大于0，开始页数小于结束页数，-f 和 -l 参数不能同时使用。

如果参数不合法，输出错误信息并终止程序。

```go
func check(args * selgp_args){
	if args.start_page <= 0 || args.end_page <= 0 {
		fmt.Fprintf(os.Stderr, "[Error] start page and end page must be larger than 0")
		os.Exit(1)
	} else if args.start_page > args.end_page {
		fmt.Fprintf(os.Stderr, "[Error] end page must be not less than start page\n")
		os.Exit(2)
	} else if args.page_type == "f" && args.page_len != 72 {
		fmt.Fprintf(os.Stderr, "[Error] -l and -f must be mutually exclusive\n")
		os.Exit(3)
	} else if args.page_len <= 0 {
		fmt.Fprintf(os.Stderr, "[Error] page length must be larger than 0\n")
		os.Exit(4)
	} else {
		fmt.Printf("start_page: %d\n end_page: %d\n input_file: %s\n page_type: %s\npage_len: %d\n destination: %s\n", args.start_page, args.end_page, args.input_file, args.page_type, args.page_len, args.destination)
	}
}
```



### selpg执行逻辑

处理的过程即从某处读入内容，然后按照一定的格式输出到某处

#### -d

如果有-d参数，则需要设置pipe[（参考资料）](https://godoc.org/os/exec#example-Command)

```go
var cmd * exec.Cmd
var cmd_in io.WriteCloser
var cmd_out io.ReadCloser
if args.destination != "" {
	cmd = exec.Command("bash", "-c", args.destination)
	cmd_in, _ = cmd.StdinPipe()
	cmd_out, _ = cmd.StdoutPipe()
	cmd.Start()
}
```
使用os/exec包，可以执行外部命令，将输出的数据作为外部命令的输入。使用exec.Command设定要执行的外部命令，cmd.StdinPipe()返回连接到command标准输入的管道pipe，cmd.Start()使某个命令开始执行，但是并不等到他执行结束。
#### 输入

通过input_file是否为空来判断文件输入还是标准输入。另外，需要一行一行处理数据

##### 文件输入

文件输入需要打开文件，判断打开文件时是否发生错误，并将文件中读取到的内容放入缓冲中

```go
if args.input_file != "" {       //文件输入
	infile, err := os.Open(args.input_file)
	if err != nil {
    	fmt.Println(err)
		os.Exit(5)
	}
	infile_buf := bufio.NewReader(infile)
```

循环读取文件中的每一行

```go
for {
	//读取输入文件中的一行
	line, _, err := infile_buf.ReadLine()
	if err != io.EOF && err != nil {
	fmt.Println(err)
		os.Exit(6)
	} 
	if err == io.EOF {
		break
	}
```

##### 标准输入

标准输入则从控制台中读取到缓冲中

```go
else {   //标准输入
	inconsole := bufio.NewScanner(os.Stdin)
```

循环读每一行

```go
for inconsole.Scan(){
	line := inconsole.Text()
	line += "\n"
	if page_num >= args.start_page && page_num <= args.end_page { //输入在给定范围内，放入out_text
		out_text += line
	}
```

无论是哪种输入，对于-l和-f的实现是一样的

#### -l 和 -f

对行数和页数进行计数，如果读取的行数在start_page和end_page之间则将这一行输出

对于-l：当行数到达-l后的数字，页数加1，判断页数是否在范围内，不是则继续读入下一行数据，否则结束。

对于-f：读取一行后需要判断是否为换页符，如果是则换页

```go
//如果在打印范围内则输出到相应地方
if page_num >= args.start_page && page_num <= args.end_page {
	//输出
	}
	line_num++
	//换页
	if args.page_type == "l" && line_num > args.page_len {  //-l换页
		line_num = 1
		page_num ++
	} else {
		if string(line) == "\f" {    //换页符换页
			page_num ++
		}
	}
```
#### 输出

对于上面个代码段的注释//输出，文件输入和标准输入的输出处理方法是不一样的

##### 文件输入

如果是输出到文件中则先放在缓冲中，如果不是则直接标准输出

```go
//如果在打印范围内则输出到相应地方
if page_num >= args.start_page && page_num <= args.end_page {
	if args.destination == "" {
		//标准输出
		fmt.Println(string(line))
	} else {
		//文件输出
		fmt.Fprintln(cmd_in, string(line))
	}
}
```

##### 标准输入

如果是标准输出，则最后输出读取的最后一行即可；如果是文件输出，则需要输出到文件中，同时要结束标准输入。

```go
//输出
if args.destination == "" {
	fmt.Print(out_text)
} else {
	fmt.Fprint(cmd_in, out_text)
	cmd_in.Close()
	cmdBytes, err := ioutil.ReadAll(cmd_out)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Print(string(cmdBytes))
	cmd.Wait()
}
```



### 测试

##### 首先install selpg

这里有一个坑，go install会报错

![1](https://github.com/JaneYul/ServiceComputing/raw/master/homework4/pics/1.PNG)

后来查了资料发现是对selpg所在的上一级文件夹进行Install（因为install的是一个包），而不是对go源文件进行install。https://stackoverflow.com/questions/25216765/gobin-not-set-cannot-run-go-install

![2](https://github.com/JaneYul/ServiceComputing/raw/master/homework4/pics/2.PNG)

##### 按照文档中“selpg使用”章节进行测试

##### selpg -s1 -e1 test.txt

此时默认打印72行

![t1](https://github.com/JaneYul/ServiceComputing/raw/master/homework4/pics/t1.PNG)

##### selpg -s1 -e1 -l5 test.txt

此时显示5行

![t2](https://github.com/JaneYul/ServiceComputing/raw/master/homework4/pics/t2.PNG)

**selpg -s1 -e1 < test.txt**

![t3](https://github.com/JaneYul/ServiceComputing/raw/master/homework4/pics/t3.PNG)

**selpg -s1 -e1 -l12 test.txt >out.txt**

此时文件中的内容定向输出到out.txt文件中

![t4](https://github.com/JaneYul/ServiceComputing/raw/master/homework4/pics/t4.PNG)

**selpg -s2 -e1 test.txt**

报错

![t5](https://github.com/JaneYul/ServiceComputing/raw/master/homework4/pics/t5.PNG)

##### selpg -s1 -e1 -l5 test.txt | cat -n

管道重定向

![t6](https://github.com/JaneYul/ServiceComputing/raw/master/homework4/pics/t6.PNG)