package main

import (
    "github.com/spf13/pflag" 
    "fmt"
    "io"
    "io/ioutil"
    "os"
    "os/exec"
    "bufio"
)

type selgp_args struct{
	start_page int                   //开始页数  -s1
	end_page int                     //结束页数  -e5
	input_file string                //输入文件  [input.txt]
	page_type string                  //指定每页行数-l10 或 换行符-f
	page_len int                     //每页多少行
	destination string               //打印目的（打印机）
}

var args selgp_args

func show_tips() {
	fmt.Println("Correct format: selpg -s[startPageNumber] -e[endPageNumber] [options] [filename]")
	fmt.Println("[options] -l: the number of lines per page (default 72).")
	fmt.Println("[options] -f: paging by page breaks.")
	fmt.Println("[options]     -l and -f are mutually exclusive.")
	fmt.Println("[options] -d: the destination of output.")
	fmt.Println("[filename]  : input file. If it is empty, the input will be from the console.")
}

func get(args *selgp_args) {
	pflag.Usage = show_tips
	//func IntVarP(指向int变量的指针，名称，简写字符串，值，用法字符串)
	pflag.IntVarP(&args.start_page,"start page", "s", 0, "Start page of file")
	pflag.IntVarP(&args.end_page,"end page","e", 0, "End page of file")
	pflag.IntVarP(&args.page_len,"page length", "l", 72, "lines in one page")  //默认72行一页
	//func StringVarP(指向字符串变量的指针，名称，简写字符串，值，用法字符串)
	pflag.StringVarP(&args.page_type,"page type","f", "l", "flag splits page")
	pflag.StringVarP(&args.destination, "destination","d", "", "name of printer")
	
	//解析
	pflag.Parse()

	//获取剩下的输入文件参数（可能没有）
	other_args := pflag.Args()
	if len(other_args) > 0 {
		args.input_file = other_args[0]
	} else {
		args.input_file = ""
	}

/*  //函数参数 a []string
	//参数数量不够
	if len(a) < 3 {
		fmt.Fprintf(os.Stderr, "%s: the num arguments is less 3\n", progname)
		show_tips()
		os.Exit(1)
	}

	//处理第一个参数
	if a[1][0] != '-' || a[1][1] != 's' {
		fmt.Fprintf(os.Stderr, "%s: 1st arg should be -s[start_page]\n", progname)
		show_tips()
		os.Exit(1)
	}
	//提取开始页数
	sp, _ := strconv.Atoi(a[1][2:])
	if sp < 1 {
		fmt.Fprintf(os.Stderr, "%s: start page should not be less than 1 %d\n", progname, sp)
		show_tips()
		os.Exit(1)
	}
	args.start_page = sp
	//处理第二个参数
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

	//其他参数处理
	argindex := 3
	for {
		if argindex > argcount-1 || a[argindex][0] != '-' {
			break
		}
		switch a[argindex][1] {
		case 'l':
			pl, _ := strconv.Atoi(a[argindex][2:])
			if pl < 1 {
				fmt.Fprintf(os.Stderr, "%s: page length should not be less than 1 %d\n", progname, pl)
				show_tips()
				os.Exit(1)
			}
			args.page_len = pl
			argindex++
		case 'f':
			if len(a[argindex]) > 2 {
				fmt.Fprintf(os.Stderr, "%s: option should be \"-f\"\n", progname)
				show_tips()
				os.Exit(1)
			}
			args.page_type = 'f'
			argindex++
		case 'd':
			if len(a[argindex]) <= 2 {
				fmt.Fprintf(os.Stderr, "%s: -d option requires a printer destination\n", progname)
				show_tips()
				os.Exit(1)
			}
			args.destination = a[argindex][2:]
			argindex++
		default:
			fmt.Fprintf(os.Stderr, "%s: unknown option", progname)
			show_tips()
			os.Exit(1)
		}
	}

	if argindex <= argcount-1 {
		args.input_file = a[argindex]
	}
	*/
}

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

func run(args * selgp_args) {
	var cmd * exec.Cmd
	var cmd_in io.WriteCloser
	var cmd_out io.ReadCloser
	if args.destination != "" {
		cmd = exec.Command("bash", "-c", args.destination)
		cmd_in, _ = cmd.StdinPipe()
		cmd_out, _ = cmd.StdoutPipe()
		cmd.Start()
	}

	if args.input_file != "" {       //文件输入
		infile, err := os.Open(args.input_file)
		if err != nil {
			fmt.Println(err)
			os.Exit(5)
		}
		line_num := 1
		page_num := 1
		infile_buf := bufio.NewReader(infile)
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
		}
	} else {   //标准输入
		inconsole := bufio.NewScanner(os.Stdin)
		line_num := 1
		page_num := 1
		out_text := ""

		//读取输入
		for inconsole.Scan(){
			line := inconsole.Text()
			line += "\n"
			if page_num >= args.start_page && page_num <= args.end_page { //输入在给定范围内，放入out_text
				out_text += line
			}
			line_num ++

			//换页
			if args.page_type == "l" && line_num > args.page_len {
				line_num = 1
				page_num++
			} else {
				if string(line) == "\f" {
					page_num ++
				}
			}
		}

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
	}
}

func main() {
    get(&args)          //读取并处理参数
    check(&args)             //检查参数是否合法
    run(&args)               //运行
}