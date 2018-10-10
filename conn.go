/*
* @Author: Ins
* @Date:   2018-10-10 09:54:12
* @Last Modified by:   Ins
* @Last Modified time: 2018-10-10 10:11:46
*/
package main

import (
    "fmt"
    "github.com/ceph/go-ceph/rados" 
)


func main() {
    conn, err := rados.NewConnWithClusterAndUser("ceph","client.objstore")
    if err != nil {
        fmt.Println("error when invoke a new connection:", err)
        return
    }

    err = conn.ReadDefaultConfigFile()//配置文件默认路径：/root/ceph/ceph.conf
    if err != nil {
        fmt.Println("error when read default config file:", err)
        return
    }

    err = conn.Connect()
    if err != nil {
        fmt.Println("error when connect:", err)
        return
    }

    fmt.Println("connect ceph cluster ok!")
    fmt.Println(conn.GetClusterStats())
    conn.Shutdown()
}