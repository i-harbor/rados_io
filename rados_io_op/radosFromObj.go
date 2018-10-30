/*
* @Author: Ins
* @Date:   2018-10-30 16:21:00
* @Last Modified by:   Ins
* @Last Modified time: 2018-10-30 20:34:19
*/
package rados_io_op

import (
    "strconv"
    "github.com/ceph/go-ceph/rados"
)

func readObjectToBytes(ioctx *rados.IOContext, oid string, block_size int, offset uint64) (int, []byte, error) {
    bytesOut := make([]byte, block_size)
    ret, err := ioctx.Read(oid, bytesOut, offset)
    if err != nil {
        return -1, bytesOut, err
    }
    bytesOut = bytesOut[:ret]
    return ret, bytesOut, err
}

func RadosFromObj(cluster_name string, user_name string, conf_file string, pool_name string, block_size int, oid string, offset uint64) (bool, []byte) {
    conn, err := NewConn(cluster_name, user_name, conf_file)
    if err != nil {
        return false, []byte("error when invoke a new connection:" + err.Error())
    }
    defer conn.Shutdown()

    // open a pool handle
    ioctx, err := conn.OpenIOContext(pool_name)
    if err != nil {
        return false, []byte("error when openIOContext" + err.Error())
    }
    defer ioctx.Destroy()

    // read the data and write to the file
    if offset / MAX_RADOS_BYTES > 0 {
        oid += "__" + strconv.FormatUint(offset / MAX_RADOS_BYTES,10)
        offset %= MAX_RADOS_BYTES
    }
    ret, bytesOut, err := readObjectToBytes(ioctx, oid, block_size, offset)
    if ret == -1 {
        return false, []byte("error when read the object(oid:" + oid + ") to bytes:" + err.Error())
    }
    return true, bytesOut
}