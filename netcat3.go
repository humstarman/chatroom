package main

import(
    "fmt"
    "io"
    "log"
    "net"
    "os"
    "flag"
)

var (
    ipFlag = flag.String("ip","localhost","IP address to connect")
    portFlag = flag.String("port","8000","Port to use")
)

func main() {
    flag.Parse()
    dst := fmt.Sprintf("%s:%s",*ipFlag,*portFlag)
    conn, err := net.Dial("tcp",dst)
    if err != nil {
        log.Fatal(err)
    }
    done := make(chan struct{})
    //go mustCopy(os.Stdout,conn)
    go func() {
        io.Copy(os.Stdout,conn)
        log.Println("done")
        done <- struct{}{}
    }()
    mustCopy(conn,os.Stdin)
    conn.Close()
    <- done
}

func mustCopy(dst io.Writer,src io.Reader) {
    if _,err := io.Copy(dst,src); err != nil {
        log.Fatal(err)
    }
}
