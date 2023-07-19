package main

import (
    "bufio"
    "fmt"
    "net"
    "regexp"
    "strings"
    "time"
)

var addr string

// TCP Server端
// 处理函数
func process(conn net.Conn, ch chan int) {
    defer conn.Close() // 关闭连接

    reader := bufio.NewReader(conn)
    var buf [128]byte
    for {
        n, err := reader.Read(buf[:]) // 读取数据
        if err != nil {
            fmt.Println("read from client failed, err: ", err)
            return
        }
        recvStr := string(buf[:n])
        fmt.Println("Recv msg:", recvStr)
        matched,_ := regexp.MatchString(`[0-9]+\.[0-9]+\.[0-9]+\.[0-9]+:[0-9]+`, recvStr)
        if !matched{
            if strings.Contains(recvStr, "done"){
                //creating mpp done
                fmt.Println("Creating mpp done! Tell remote.")
                ch <- 1
                return
            }else if strings.Contains(recvStr, "clean"){
                // if last test break before connecting to server, but the mpp is created, there will be a succeful result in channel, so clean it
                fmt.Println("clean channel.")
                select {
                case c := <-ch:
                    conn.Write([]byte(fmt.Sprintf("get %d result from channel.", c)))
                case <- time.After(1 * time.Second):
                    conn.Write([]byte("nothing in channel."))
                }
                return
            }else{
                fmt.Println("wrong address!")
                conn.Write([]byte("Wrong address!"))
            }
        }else{
            fmt.Println("recv address:", recvStr)
            // receive condition

            // time.Sleep(time.Second * 1)
            // fmt.Println("Send success.")
            // conn.Write([]byte("success"))
            r := <- ch
            if r == 1 {
                fmt.Println("Send success.")
                conn.Write([]byte("success"))
            }
            addr = recvStr
        }
    }
    return
}

func main() {
    ch := make(chan int)
    defer close(ch)

    listen, err := net.Listen("tcp", ":3333")
    if err != nil {
        fmt.Println("Listen() failed, err: ", err)
        return
    }
    for {
        conn, err := listen.Accept() // 监听客户端的连接请求
        if err != nil {
            fmt.Println("Accept() failed, err: ", err)
            continue
        }
        go process(conn, ch) // 启动一个goroutine来处理客户端的连接请求
    }
}
