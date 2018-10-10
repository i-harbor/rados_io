/*
* @Author: Ins
* @Date:   2018-10-10 09:54:12
* @Last Modified by:   Ins
* @Last Modified time: 2018-10-10 10:12:11
*/
package main

import (
    "fmt"
    "os"
    "github.com/ceph/go-ceph/rados"
    "github.com/ceph/go-ceph/rbd"
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

func listImages(ioctx *rados.IOContext, prefix string) {
    imageNames, err := rbd.GetImageNames(ioctx)
    if err != nil {
        fmt.Println("error when getImagesNames", err)
        os.Exit(1)
    }
    fmt.Println(prefix, ":", imageNames)
}

func main() {
    conn, err := newConn()
    if err != nil {
        fmt.Println("error when invoke a new connection:", err)
        return
    }
    defer conn.Shutdown()
    fmt.Println("connect ceph cluster ok!")

    ioctx, err := conn.OpenIOContext("objstore")
    if err != nil {
        fmt.Println("error when openIOContext", err)
        return
    }
    defer ioctx.Destroy()

    listImages(ioctx, "before create new image")

    name := "go-ceph-image"
    img, err := rbd.Create(ioctx, name, 1<<20, 20)
    if err != nil {
        fmt.Println("error when create rbd image", err)
        return
    }
    listImages(ioctx, "after create new image")

    err = img.Remove()
    if err != nil {
        fmt.Println("error when remove image", err)
        return
    }
    listImages(ioctx, "after remove new image")
}
