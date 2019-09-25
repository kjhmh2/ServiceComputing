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

设计一个结构体记录存放的参数。包括开始页、结束页、输入文件的名字、输出、一页的长度、分页方式

```go
type selpg_args struct {
	start_page  int
	end_page    int
	in_filename string
	dest        string
	page_len    int
	page_type   int
}
```
设置全局变量记录读取到的参数、程序名和参数个数

```go
var sa selpg_args   
var progname string 
var argcount int
```



### 读取并处理参数

使用**os.Args**[（参考资料）](https://blog.csdn.net/guanchunsheng/article/details/79612153)获取用户输入的参数，得到一个包含参数的数组，数组中的每个元素是string类型的。

参数处理可以使用**pflag**包来解析命令的参数。[参考资料](https://studygolang.com/articles/5608)

要import pflag包，要先在本地安装pflag
```bash
go get github.com/spf13/pflag
```
然后可以通过下面的代码进行参数值的绑定,通过 **pflag.Parse()**方法让pflag 对标识和参数进行解析。之后就可以直接使用绑定的值。
```go
pflag.IntVarP(&sa.start_page,"start", "s", 0, "Start page of file")
pflag.IntVarP(&sa.end_page,"end","e", 0, "End page of file")
pflag.IntVarP(&sa.page_len,"linenum", "l", 20, "lines in one page")
pflag.StringVarP(&sa.page_type,"printdes","f", "l", "flag splits page")
pflag.StringVarP(&sa.dest, "destination","d", "", "name of printer")
pflag.Parse()
```
提示信息，如果给出的参数不正确或者需要查看帮助 -help，那么会给出这里指定的字符串
```go
pflag.Usage = show_tips
```
通过pflag.NArg()可以知道是否有要进行操作的文件。如果是pflag解析不了的类型参数。我们称这种参数为non-flag参数，flag解析遇到non-flag参数就停止了。pflag提供了Arg(i),Args()来获取non-flag参数，NArg()来获取non-flag的个数。所以可以使用`pflag.Arg(0)`来获取输入的文件路径

##### **另外**

如果手动判断每个参数，则需要判断参数个数、参数格式是否符合要求。不符合要求则输出错误信息，终止程序；符合要求则赋给参数结构体的每个变量中。

##### 参数少于3个

```go
if len(args) < 3 {
	fmt.Fprintf(os.Stderr, "%s: the num arguments is less 3\n", progname)
	show_tips()
	os.Exit(1)
}
```

##### 第一个参数：开始页数

先判断第一个参数是不是"-s"形式

```go
if args[1][0] != '-' || args[1][1] != 's' {
	fmt.Fprintf(os.Stderr, "%s: 1st arg should be -sstart_page\n", progname)
	show_tips()
	os.Exit(1)
}
```
如果是，则从数组中提取出参数
```go
	sp, _ := strconv.Atoi(args[1][2:])
	if sp < 1 {
		fmt.Fprintf(os.Stderr, "%s: start page should not be less than 1 %d\n", progname, sp)
		show_tips()
		os.Exit(1)
	}
	sa.start_page = sp
```

##### 第二个参数：结束页数

和start_page很类似，先判断是否符合"-e"格式再提取出参数

```go
if args[2][0] != '-' || args[2][1] != 'e' {
	fmt.Fprintf(os.Stderr, "%s: 2nd arg should be -e[end_page]\n", progname)
	show_tips()
	os.Exit(1)
}

ep, _ := strconv.Atoi(args[2][2:])
if ep < 1 || ep < sp {
	fmt.Fprintf(os.Stderr, "%s: end page should not be less than 1 %d\n", progname, ep)
	show_tips()
	os.Exit(1)
}
sa.end_page = ep
```

##### option参数

如果参数个数大于三，需要获取其他参数信息，可能是"-l", "-f", "-d"。需要确认是以"-"开头才能进行参数判断

- 循环读取

```go
	for {
		if argindex > argcount-1 || args[argindex][0] != '-' {
			break
		}
		switch args[argindex][1] {
```
- -l
  判断-l后面跟着的数字是否符合格式

```go
case 'l':
	pl, _ := strconv.Atoi(args[argindex][2:])
	if pl < 1 {
		fmt.Fprintf(os.Stderr, "%s: page length should not be less than 1 %d\n", progname, pl)
		show_tips()
		os.Exit(1)
	}
	sa.page_len = pl
	argindex++
```
- -f
  -f参数后面不跟数字，所以判断-f参数的长度来判断是否合法

```go
		case 'f':
			if len(args[argindex]) > 2 {
				fmt.Fprintf(os.Stderr, "%s: option should be \"-f\"\n", progname)
				show_tips()
				os.Exit(1)
			}
			sa.page_type = 'f'
			argindex++
```
- -d
  根据说明，selpg不检查destination目的地，但是要确保-d后面跟着目的地

```go
			if len(args[argindex]) == 2 {
				fmt.Fprintf(os.Stderr, "%s: -d option requires a printer destination\n", progname)
				Usage()
				os.Exit(1)
			}
			sa.destination = args[argindex][2:]
			argindex++
```
- 不合法参数
  如果不是以上的任何一种情况，需要报错处理

```go
	default:
		fmt.Fprintf(os.Stderr, "%s: unknown option", progname)
		Usage()
		os.Exit(1)
```

- 输入文件
在上面的循环中，没有以“-”开头的参数，跳出。
此时可能还有文件作为输入，或者没有（此时为标准输入），需要判断
```go
if argindex <= argcount-1 {
	sa.input_file = args[argindex]
}
```



### 从标准输入或文件中获取输入然后输出到标准输出或文件中

处理的过程即从某处读入内容，然后按照一定的格式输出到某处

##### 实现-d

如果有-d参数，则需要设置pipe[（参考资料）](https://godoc.org/os/exec#example-Command)

```go
var cmd *exec.Cmd
var cmd_in io.WriteCloser
var cmd_out io.ReadCloser
if sa.destination != "" {
	cmd = exec.Command("bash", "-c", sa.destination)
	cmd_in, _ = cmd.StdinPipe()
	cmd_out, _ = cmd.StdoutPipe()
	cmd.Start()
}
```
使用os/exec包，可以执行外部命令，将输出的数据作为外部命令的输入。使用exec.Command设定要执行的外部命令，cmd.StdinPipe()返回连接到command标准输入的管道pipe，cmd.Start()使某个命令开始执行，但是并不等到他执行结束。
##### 实现-l

使用页数计数器，在满足一页的条件后页数计数器增加，判断页数是否在范围内，不是则继续读入下一行数据，否则结束读取数据。

```go
			line, _, err := fin.ReadLine()
			if err != io.EOF && err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			if err == io.EOF {
				break
			}
			if page_count >= sa.start_page && page_count <= sa.end_page {
				if sa.destination == "" {
					fmt.Println(string(line))
				} 
				else {
					fmt.Fprintln(cmd_in, string(line))
				}
			}
			line_count++
```
从输入中每次读取一行，然后对每一行进行计数，当行数到达-l后的数字，页数增加，判断页数是否在范围内然后输出。

```go
				if line_count > sa.page_len {
					line_count = 1
					page_count++
				}
```

##### 实现-f

当有-f参数时，将sa.page_type赋值为’f’，从输入中每次读取一行，如果一行的字符为’\f’则页数计数增加，判断页数是否在范围内然后输出。

```go
if string(line) == "\f" {
	page_count++
}
```



### 测试

##### 首先install selpg

这里有一个坑，go install会报错

![1](https://github.com/JanelYul/ServiceComputing/raw/master/pics/1.PNG)

后来查了资料发现是对selpg所在的上一级文件夹进行Install（因为install的是一个包），而不是对go源文件进行install。https://stackoverflow.com/questions/25216765/gobin-not-set-cannot-run-go-install

![2](https://github.com/JanelYul/ServiceComputing/raw/master/pics/2.PNG)

##### 按照文档中“selpg使用”章节进行测试

##### selpg -s1 -e1 test.txt

此时默认打印72行

![t1](https://github.com/JanelYul/ServiceComputing/raw/master/pics/t1.PNG)

##### selpg -s1 -e1 -l5 test.txt

此时显示5行

![t2](https://github.com/JanelYul/ServiceComputing/raw/master/pics/t2.PNG)

**selpg -s1 -e1 < test.txt**

![t3](https://github.com/JanelYul/ServiceComputing/raw/master/pics/t3.PNG)

**selpg -s1 -e1 -l12 test.txt >out.txt**

此时文件中的内容定向输出到out.txt文件中

![t4](https://github.com/JanelYul/ServiceComputing/raw/master/pics/t4.PNG)

**selpg -s1 -e0 -l12 test.txt >/dev/null**

报错

![t5](https://github.com/JanelYul/ServiceComputing/raw/master/pics/t5.PNG)

##### selpg -s1 -e1 -l5 test.txt | cat -n

管道重定向

![t6](https://github.com/JanelYul/ServiceComputing/raw/master/pics/t6.PNG)