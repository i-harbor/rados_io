/*
* @Author: Ins
* @Date:   2018-10-10 09:54:12
* @Last Modified by:   Ins
* @Last Modified time: 2018-10-10 10:12:05
*/
package main

import (
    "fmt"
    "os"
    "io/ioutil"
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

func ReadFileToBytes(filePth string) ([]byte, error) {
    f, err := os.Open(filePth)
    if err != nil {
        return nil, err
    }
    defer f.Close()
    return ioutil.ReadAll(f)
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
    fmt.Println("before add a object. object list:")
    ioctx.ListObjects(ObjectListFunc)

    //open a file
    bytesIn, err := ReadFileToBytes("/root/nginx-1.13.12.tar.gz")
    // fmt.Println(string(bytesIn))

    // write data to object
    err = ioctx.Write("nginx", bytesIn, 0)

    fmt.Println("after add a object. object list:")
    ioctx.ListObjects(ObjectListFunc)

}