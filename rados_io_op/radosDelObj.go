/*
* @Author: Ins
* @Date:   2018-10-30 16:21:00
* @Last Modified by:   Ins
* @Last Modified time: 2019-01-21 15:34:41
*/
package rados_io_op

import (
    "strconv"
    "github.com/ceph/go-ceph/rados"
)
func deleteObj(ioctx *rados.IOContext, oid string, osize uint64) (error, string) {
    var oid_suffix_list string = "("
    err := ioctx.Delete(oid)
    if err == nil {
        oid_suffix_list += oid
    }
    oid_suffix := osize / MAX_RADOS_BYTES
    for i := uint64(0); i <= oid_suffix; i++ {
        oid_tmp := oid + "_" + strconv.FormatUint(i, 10)
        err_tmp := ioctx.Delete(oid_tmp)
        if err_tmp == nil {
            err = nil
            oid_suffix_list += " " + oid_tmp
        }
    }
    oid_suffix_list += ")"
 
    return err, oid_suffix_list
}

func RadosDelObj(cluster_name string, user_name string, conf_file string, pool_name string, oid string, osize uint64) (bool, []byte) {
    conn, err := NewConn(cluster_name, user_name, conf_file)
    if err != nil {
        return false, []byte(err.Error() + ",error when invoke a new connection:" + ERR_INFO[err.Error()])
    }
    defer conn.Shutdown()

    // open a pool handle
    ioctx, err := conn.OpenIOContext(pool_name)
    if err != nil {
        return false, []byte(err.Error() + ",error when openIOContext:" + ERR_INFO[err.Error()])
    }
    defer ioctx.Destroy()

    // delete a object
    err, oid_suffix_list := deleteObj(ioctx, oid, osize)

    if err != nil {
        return false, []byte(err.Error() + ",error when delete the object:" + ERR_INFO[err.Error()])
    }

    return true, []byte("successfully delete the object:" + oid + oid_suffix_list)
}