package main

import (
    "fmt"
    "os"
    "github.com/ceph/go-ceph/rados"
)


func newConn() (*rados.Conn, error) {
    conn, err := rados.NewConnWithClusterAndUser("ceph","client.objstore")
    if err != nil {
        return nil, err
    }

    err = conn.ReadDefaultConfigFile()
    if err != nil {
        return nil, err
    }

    err = conn.Connect()
    if err != nil {
        return nil, err
    }

    return conn, nil
}

func ObjectListFunc(oid string) {
    fmt.Println(oid)
}

func WriteBytesToFile(fname string, bytes []byte) error {
    f, err := os.OpenFile(fname, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
    if err != nil {
        fmt.Println("error when create the file:", err)
    }
    defer f.Close()
    _, err = f.Write(bytes);
    return err
}

func ReadObjectToBytes(ioctx *rados.IOContext, oname string, block_size int, offset uint64) (int, []byte) {
    bytesOut := make([]byte, block_size)
    ret, err := ioctx.Read("nginx", bytesOut, offset)
    if err != nil {
        fmt.Println("error when read the object to bytes:", err)
    }
    bytesOut = bytesOut[:ret]
    return ret, bytesOut
}

func main() {
    conn, err := newConn()
    if err != nil {
        fmt.Println("error when invoke a new connection:", err)
        return
    }
    defer conn.Shutdown()
    // fmt.Println("connect ceph cluster ok!")

    // open a pool handle
    ioctx, err := conn.OpenIOContext("objstore")
    if err != nil {
        fmt.Println("error when openIOContext", err)
        return
    }
    defer ioctx.Destroy()

    // ioctx.ListObjects(ObjectListFunc)

    // read the data and write to the file
    var offset uint64
    block_size := 2048000000
    i := 0
    for offset = 0; ; offset += uint64(block_size) {
        ret, bytesOut := ReadObjectToBytes(ioctx, "nginx", block_size, offset)
        err = WriteBytesToFile("/root/nginx.tar.gz", bytesOut)
        if err != nil {
            fmt.Println("error when write to the file:", err)
        }
        i++
        // fmt.Println(offset)
        if ret != block_size {
            fmt.Println("共切了",i,"个block写入")
            break;
        }
    }
}