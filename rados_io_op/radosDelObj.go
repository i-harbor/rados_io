/*
* @Author: Ins
* @Date:   2018-10-30 16:21:00
* @Last Modified by:   Ins
* @Last Modified time: 2018-10-31 10:26:05
*/
package rados_io_op

import (
    "strconv"
    "github.com/ceph/go-ceph/rados"
)
func deleteObj(ioctx *rados.IOContext, oid string) (error, string) {
    err := ioctx.Delete(oid)
    if err != nil {
        return err, ""
    }
    var flag [MAX_RADOS_SUFFIX]bool
    var oid_suffix_list string = ""
    for i, _ := range flag {
        oid_tmp := oid + "__" + strconv.Itoa(i)
        err = ioctx.Delete(oid_tmp)
        if err != nil {
            flag[i] = true
        } else {
            oid_suffix_list += "," + oid_tmp
        }
        if flag[i] && i >= 5 && flag[i-5] && flag[i-4] && flag[i-3] && flag[i-2] && flag[i-1] {
            break
        }
    }
    return nil, oid_suffix_list
}

func RadosDelObj(cluster_name string, user_name string, conf_file string, pool_name string, oid string) (bool, []byte) {
    conn, err := NewConn(cluster_name, user_name, conf_file)
    if err != nil {
        return false, []byte("error when invoke a new connection:" + err.Error())
    }
    defer conn.Shutdown()

    // open a pool handle
    ioctx, err := conn.OpenIOContext(pool_name)
    if err != nil {
        return false, []byte("error when openIOContext:" + err.Error())
    }
    defer ioctx.Destroy()

    // delete a object 
    err, oid_suffix_list := deleteObj(ioctx, oid)
    if err != nil {
        return false, []byte("error when delete the object:" + err.Error())
    }

    return true, []byte("successfully delete the object:" + oid + oid_suffix_list)
}