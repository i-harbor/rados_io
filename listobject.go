/*
* @Author: Ins
* @Date:   2018-10-10 09:54:12
* @Last Modified by:   Ins
* @Last Modified time: 2018-10-16 10:12:20
*/
package main

import (
    "fmt"
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

    //列出pool中所有objects
    ioctx.ListObjects(ObjectListFunc)

    bytesOut := make([]byte, 100)
    ret, err := ioctx.Read("abc", bytesOut, 0)
    bytesOut = bytesOut[:ret]
    fmt.Println(string(bytesOut))
}