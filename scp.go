package main
 
import (
  "fmt"
  "log"
  "os"
  "path"
  "time"
  "net"

  "github.com/pkg/sftp"
 
  "golang.org/x/crypto/ssh"
)

func connect(user, password, host string, port int) (*sftp.Client, error) {
  var (
    auth         []ssh.AuthMethod
    addr         string
    clientConfig *ssh.ClientConfig
    sshClient    *ssh.Client
    sftpClient   *sftp.Client
    err          error
  )
  // get auth method
  auth = make([]ssh.AuthMethod, 0)
  auth = append(auth, ssh.Password(password))
 
  clientConfig = &ssh.ClientConfig{
    User:    user,
    Auth:    auth,
    Timeout: 30 * time.Second,
    HostKeyCallback:func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
  }
 
  // connet to ssh
  addr = fmt.Sprintf("%s:%d", host, port)
 
  if sshClient, err = ssh.Dial("tcp", addr, clientConfig); err != nil {
    return nil, err
  }
 
  // create sftp client
  if sftpClient, err = sftp.NewClient(sshClient); err != nil {
    return nil, err
  }
 
  return sftpClient, nil
}

func remoteCopyFile(path_ string, user_ string, password_ string, host_ string, port_ int) {
	var (
    err        error
    sftpClient *sftp.Client
  )

  // 这里换成实际的 SSH 连接的 用户名，密码，主机名或IP，SSH端口
  sftpClient, err = connect(user_, password_, host_, port_)
  if err != nil {
    log.Fatal(err)
  }
  defer sftpClient.Close()

  // 用来测试的远程文件路径 和 本地文件夹
  var remoteFilePath = path_
  var localDir = "."

  srcFile, err := sftpClient.Open(remoteFilePath)
  if err != nil {
    log.Fatal(err)
  }
  defer srcFile.Close()

  var localFileName = path.Base(remoteFilePath)
  dstFile, err := os.Create(path.Join(localDir, localFileName))
  if err != nil {
    log.Fatal(err)
    }

    defer dstFile.Close()

  if _, err = srcFile.WriteTo(dstFile); err != nil {
    log.Fatal(err)
  }

  fmt.Println("copy file from remote server finished!")
}


func main() {
	// ubuntu@122.51.172.48:/home/ubuntu/go/src/delve
	//              对端路径                               用户名    密码          对端IP        对端端口号
	 remoteCopyFile("/home/lijian/mydir/practice/export1", "xxx", "xxx", "132.232.241.187", 22)

	//remoteCopyFile("/home/ubuntu/test.txt", "ubuntu", "xxx", "127.0.0.1", 22)

}
