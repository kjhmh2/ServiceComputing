package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
)

//参数列表
type selpg_args struct {
	start_page int
	end_page int
	input_file string
	destination string
	page_len int
	page_type int
}

func show_tips() {
	fmt.Println("\nCorrect tips: selpg -s[startPageNumber] -e[endPageNumber] [options] [filename]")
	fmt.Println("[options] -l: the number of lines per page (default 72)")
	fmt.Println("[options] -f: the type and the way to be paged.")
	fmt.Println("[options] -d: the destination of output.")
	fmt.Println("[filename]  : input file.")
}

var sa selpg_args   //获取参数
var progname string //程序名
var argcount int    //参数个数

func process_args(args []string) {
	pflag.Usage = show_tips
	pflag.IntVarP(&sa.start_page,"start", "s", 0, "Start page of file")
	pflag.IntVarP(&sa.end_page,"end","e", 0, "End page of file")
	pflag.IntVarP(&sa.page_len,"linenum", "l", 20, "lines in one page")
	pflag.StringVarP(&sa.page_type,"printdes","f", "l", "flag splits page")
	pflag.StringVarP(&sa.dest, "destination","d", "", "name of printer")
	pflag.Parse()
	othersArg := pflag.Args()
    if len(othersArg) > 0 {
        sa.inFile = othersArg[0]
    } else {
        sa.inFile = ""
    }
	/*
	//参数数量不够
	if len(args) < 3 {
		fmt.Fprintf(os.Stderr, "%s: the num arguments is less 3\n", progname)
		show_tips()
		os.Exit(1)
	}

	//处理第一个参数
	if args[1][0] != '-' || args[1][1] != 's' {
		fmt.Fprintf(os.Stderr, "%s: 1st arg should be -s[start_page]\n", progname)
		show_tips()
		os.Exit(1)
	}
	//提取开始页数
	sp, _ := strconv.Atoi(args[1][2:])
	if sp < 1 {
		fmt.Fprintf(os.Stderr, "%s: start page should not be less than 1 %d\n", progname, sp)
		show_tips()
		os.Exit(1)
	}
	sa.start_page = sp
	//处理第二个参数
	if args[2][0] != '-' || args[2][1] != 'e' {
		fmt.Fprintf(os.Stderr, "%s: 2nd arg should be -e[end_page]\n", progname)
		show_tips()
		os.Exit(1)
	}
	//提取结束页数
	ep, _ := strconv.Atoi(args[2][2:])
	if ep < 1 || ep < sp {
		fmt.Fprintf(os.Stderr, "%s: end page should not be less than 1 %d\n", progname, ep)
		show_tips()
		os.Exit(1)
	}
	sa.end_page = ep

	//其他参数处理
	argindex := 3
	for {
		if argindex > argcount-1 || args[argindex][0] != '-' {
			break
		}
		switch args[argindex][1] {
		case 'l':
			pl, _ := strconv.Atoi(args[argindex][2:])
			if pl < 1 {
				fmt.Fprintf(os.Stderr, "%s: page length should not be less than 1 %d\n", progname, pl)
				show_tips()
				os.Exit(1)
			}
			sa.page_len = pl
			argindex++
		case 'f':
			if len(args[argindex]) > 2 {
				fmt.Fprintf(os.Stderr, "%s: option should be \"-f\"\n", progname)
				show_tips()
				os.Exit(1)
			}
			sa.page_type = 'f'
			argindex++
		case 'd':
			if len(args[argindex]) <= 2 {
				fmt.Fprintf(os.Stderr, "%s: -d option requires a printer destinationination\n", progname)
				show_tips()
				os.Exit(1)
			}
			sa.destination = args[argindex][2:]
			argindex++
		default:
			fmt.Fprintf(os.Stderr, "%s: unknown option", progname)
			show_tips()
			os.Exit(1)
		}
	}

	if argindex <= argcount-1 {
		sa.input_file = args[argindex]
	}
	*/
}

func process_input() {
	var cmd *exec.Cmd
	var cmd_in io.WriteCloser
	var cmd_out io.ReadCloser
	if sa.destination != "" {
		cmd = exec.Command("bash", "-c", sa.destination)
		cmd_in, _ = cmd.StdinPipe()
		cmd_out, _ = cmd.StdoutPipe()
		cmd.Start()
	}
	if sa.input_file != "" {
		inf, err := os.Open(sa.input_file)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		line_count := 1
		page_count := 1
		fin := bufio.NewReader(inf)
		for {
			//读取输入文件中的一行数据
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
					//标准输出
					fmt.Println(string(line))
				} else {
					//文件输出
					fmt.Fprintln(cmd_in, string(line))
				}
			}
			line_count++
			if sa.page_type == 'l' {
				if line_count > sa.page_len {
					line_count = 1
					page_count++
				}
			} else {
				if string(line) == "\f" {
					page_count++
				}
			}
		}
		if sa.destination != "" {
			cmd_in.Close()
			cmdBytes, err := ioutil.ReadAll(cmd_out)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Print(string(cmdBytes))
			cmd.Wait()
		}
	} 
	else {
		//标准输入
		ns := bufio.NewScanner(os.Stdin)
		line_count := 1
		page_count := 1
		out := ""

		for ns.Scan() {
			line := ns.Text()
			line += "\n"
			if page_count >= sa.start_page && page_count <= sa.end_page {
				out += line
			}
			line_count++
			if sa.page_type == 'l' {
				if line_count > sa.page_len {
					line_count = 1
					page_count++
				}
			} 
			else {
				if string(line) == "\f" {
					page_count++
				}
			}
		}
		if sa.destination == "" {
			fmt.Print(out)
		} 
		else {
			fmt.Fprint(cmd_in, out)
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
	args := os.Args
	sa.start_page = 1
	sa.end_page = 1
	sa.input_file = ""
	sa.destination = ""
	sa.page_len = 20
	sa.page_type = 'l'
	argcount = len(args)
	process_args(args)
	process_input()
}