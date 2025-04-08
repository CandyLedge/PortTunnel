package main

import (
    "fmt"
    "io"
    "log"
    "net"
    "sync"
)

// forward 函数用于在两个连接之间转发数据
func forward(src, dst net.Conn, wg *sync.WaitGroup) {
    defer wg.Done()
    defer func() {
        src.Close()
        dst.Close()
    }()

    buffer := make([]byte, 4096)
    for {
        n, err := src.Read(buffer)
        if err != nil {
            if err != io.EOF {
                log.Printf("读取数据时出错: %v", err)
            }
            break
        }

        _, err = dst.Write(buffer[:n])
        if err != nil {
            log.Printf("写入数据时出错: %v", err)
            break
        }
    }
}

// handleConnection 函数处理新的客户端连接
func handleConnection(client net.Conn) {
    // 这里填写宿主机局域网目标端口的信息
    targetAddr := "10.15.8.165:10086" // 替换为实际目标地址和端口
    target, err := net.Dial("tcp", targetAddr)
    if err != nil {
        log.Printf("连接目标地址时出错: %v", err)
        client.Close()
        return
    }

    var wg sync.WaitGroup
    wg.Add(2)

    // 启动两个 goroutine 进行双向转发
    go forward(client, target, &wg)
    go forward(target, client, &wg)

    wg.Wait()
}

func main() {
    // 这里填写虚拟机要监听的端口信息
    proxyAddr := "0.0.0.0:9870" // 替换为实际监听地址和端口
    listener, err := net.Listen("tcp", proxyAddr)
    if err != nil {
        log.Fatalf("监听端口时出错: %v", err)
    }
    defer listener.Close()

    fmt.Printf("开始监听 %s\n", proxyAddr)
    for {
        client, err := listener.Accept()
        if err != nil {
            log.Printf("接受连接时出错: %v", err)
            continue
        }
        fmt.Printf("接受来自 %s 的连接\n", client.RemoteAddr())
        go handleConnection(client)
    }
}