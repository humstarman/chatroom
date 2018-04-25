package main

import(
    "fmt"
    "net"
    "log"
    "bufio"
    "flag"
)

type client chan<- string

var (
    entering = make(chan client) 
    leaving = make(chan client) 
    messages = make(chan string)
    portFlag = flag.String("port","8080","the working port")
)

func broadcaster() {
    clients := make(map[client]bool)
    for {
        select {
        case msg := <-messages:
            for cli := range clients {
                cli <- msg
            }

        case cli := <-entering:
            clients[cli] = true

        case cli := <-leaving:
            delete(clients,cli) 
            close(cli)
        }
    }
}

func handleConn(conn net.Conn) {
    ch := make(chan string)
    go clientWriter(conn,ch)
    
    who := conn.RemoteAddr().String()
    ch <- "You are " + who
    messages <- who + " has arrived"
    entering <- ch

    input := bufio.NewScanner(conn)
    for input.Scan() {
        messages <- who + ": " + input.Text()
    }

    leaving <- ch
    messages <- who + " has left"
    conn.Close()
}

func clientWriter(conn net.Conn,ch <-chan string){
    for msg := range ch {
        fmt.Fprintln(conn,msg)
    }
}

func main() {
    flag.Parse()
    sock := fmt.Sprintf("0.0.0.0:%s",*portFlag)
    listener,err := net.Listen("tcp",sock)
    if err != nil {
        log.Fatal(err)
    }
    go broadcaster()
    for {
        conn,err := listener.Accept()
        if err != nil {
            log.Print(err)
            continue
        }
        go handleConn(conn)
    }
}


