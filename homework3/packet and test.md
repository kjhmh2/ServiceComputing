# 在centos中安装go环境

前面的配置过程请见：https://blog.csdn.net/AJANEEE/article/details/100746056



## 包

在之前我们install了一个ServiceComputing程序，它输出"Hello, world"

现在我们想要编写一个库，让ServiceComputing程序使用它

- 创建包目录

  ```
  mkdir $GOPATH/src/github.com/github-user/stringutil
  ```

  

- 在该目录中创建名为`reverse.go` 的文件

  ```csharp
  // stringutil 包含有用于处理字符串的工具函数。
  package stringutil
  
  // Reverse 将其实参字符串以符文为单位左右反转。
  func Reverse(s string) string {
  	r := []rune(s)
  	for i, j := 0, len(r)-1; i < len(r)/2; i, j = i+1, j-1 {
  		r[i], r[j] = r[j], r[i]
  	}
  	return string(r)
  }
  ```

  这个程序的作用是让字符串反转。

- 用go build命令测试该包的编译

  ```
  go build github.com/github-user/stringutil
  ```

- 再使用go install命令将包的对象放到gowork的pkg目录中

  ```
  go install
  ```

  **【结果】**完成这一步之后，在 pkg/linux_amd64/ 文件夹下会生成子文件夹 github.com/github-user（与src下的目录一致），并在其中生成一个文件stringutil.a，即生成的包stringutil

  - **注意：**install生成包的名字是`reverse.go`文件所在的文件夹的名字，或者说导入路径的最后一个元素（即stringutil），而不是reverse

  ![结果](C:\Users\LENOVO\AppData\Roaming\Typora\typora-user-images\1568252389925.png)

- 修改原来的hello.go文件（位于$GOPATH/src/github.com/github-user/ServiceComputing中）

  在这个文件中使用了刚才构造的包

  **注意文件目录和名字要正确**

  ```csharp
  package main
  
  import (
  	"fmt"
  
  	"github.com/gihtub-user/stringutil"
  )
  
  func main() {
  	fmt.Printf(stringutil.Reverse("!oG ,olleH"))
  }
  ```

- 重新安装hello程序

  ```
  go install github.com/github-user/ServiceComputing
  ```

  go工具会安装它所依赖的任何东西，所以stringutil包也会被自动安装

- 再次运行hello，得到reverse之后的信息

  ![结果](C:\Users\LENOVO\AppData\Roaming\Typora\typora-user-images\1568263552395.png)

  注意 `go install` 会将 `stringutil.a` 对象放到 `pkg/linux_amd64` 目录中，它会反映出其源码目录。 这就是在此之后调用 `go` 工具，能找到包对象并避免不必要的重新编译的原因。 `linux_amd64` 这部分能帮助跨平台编译，并反映出你的操作系统和架构。

- 包名

  Go源文件中的第一个语句必须是

  ```
  package 名称
  ```

  这里的 `名称` 即为导入该包时使用的默认名称。 （一个包中的所有文件都必须使用相同的 `名称`。）

  Go的约定是包名为导入路径的最后一个元素：作为 “`crypto/rot13`” 导入的包应命名为 `rot13`。

  可执行命令必须使用 `package main`。

  链接成单个二进制文件的所有包，其包名无需是唯一的，只有导入路径（它们的完整文件名） 才是唯一的。



## 测试

可以通过创建一个名字以 `_test.go` 结尾的，包含名为 `TestXXX` 且签名为 `func (t *testing.T)` 函数的文件来编写测试。 测试框架会运行每一个这样的函数；若该函数调用了像 `t.Error` 或 `t.Fail` 这样表示失败的函数，此测试即表示失败。

- 通过创建 $GOPATH/src/github.com/user/stringutil/reverse_test.go 来为 `stringutil` 添加测试，内容如下

  ```csharp
  package stringutil
  
  import "testing"
  
  func TestReverse(t *testing.T) {
  	cases := []struct {
  		in, want string
  	}{
  		{"Hello, world", "dlrow ,olleH"},
  		{"Hello, 世界", "界世 ,olleH"},
  		{"", ""},
  	}
  	for _, c := range cases {
  		got := Reverse(c.in)
  		if got != c.want {
  			t.Errorf("Reverse(%q) == %q, want %q", c.in, got, c.want)
  		}
  	}
  }
  ```

- 使用 go test运行该测试

  ![结果](C:\Users\LENOVO\AppData\Roaming\Typora\typora-user-images\1568264100117.png)

  如果在包目录下，则可以忽略包路径

  ![结果](C:\Users\LENOVO\AppData\Roaming\Typora\typora-user-images\1568264129236.png)

  此时多了一个PASS信息

