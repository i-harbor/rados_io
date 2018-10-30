/*
* @Author: Ins
* @Date:   2018-10-30 16:21:00
* @Last Modified by:   Ins
* @Last Modified time: 2018-10-30 20:34:35
*/
package rados_io_op

import (
    "strconv"
    "github.com/ceph/go-ceph/rados"
)

func writeToObj(ioctx *rados.IOContext, oid string, bytesIn []byte, offset uint64) error {
    var err error = nil
    oid_suffix := offset / MAX_RADOS_BYTES
    offset %= MAX_RADOS_BYTES
    oid_suffix_gap := (offset + uint64(len(bytesIn))) / MAX_RADOS_BYTES
    for oid_suffix_gap > 0 {
        bytesIn_tmp := bytesIn[(oid_suffix_gap * MAX_RADOS_BYTES - offset):]
        err = ioctx.Write(oid + "__" + strconv.FormatUint(oid_suffix + oid_suffix_gap,10), bytesIn_tmp, 0)
        if err != nil {
            return err
        }
        bytesIn = bytesIn[:(oid_suffix_gap * MAX_RADOS_BYTES - offset)]
        oid_suffix_gap -= 1
    }
    if oid_suffix > 0 {
        err = ioctx.Write(oid + "__" + strconv.FormatUint(oid_suffix,10), bytesIn, offset)
    } else {
        err = ioctx.Write(oid, bytesIn, offset)
    }
    
    if err != nil {
        return err
    }
    return nil
}

func RadosToObj(cluster_name string, user_name string, conf_file string, pool_name string, oid string, bytesIn []byte, bytesLen int, mode string, offset uint64) (bool, []byte) {
    conn, err := NewConn(cluster_name, user_name, conf_file)
    if err != nil {
        return false, []byte("error when invoke a new connection:" + err.Error())
    }
    defer conn.Shutdown()

    // open a pool handle
    ioctx, err := conn.OpenIOContext(pool_name)
    if err != nil {
        return false, []byte("error when invoke a new connection:" + err.Error())
    }
    defer ioctx.Destroy()

    // write data to object
    switch mode {
        case "w":
            err = ioctx.Write(oid, bytesIn, offset)
            if err != nil {
                return false, []byte("error when write to object:" + err.Error())
            }
        case "wf":
            err = ioctx.WriteFull(oid, bytesIn)
            if err != nil {
                return false, []byte("error when write full to object:" + err.Error())
            }
        case "wa":
            err = ioctx.Append(oid, bytesIn)
            if err != nil {
                return false, []byte("error when append to object:" + err.Error())
            }
        default:
            return false, []byte("error when write to object: unknown wirte mode : " + mode + ", only ['w' : write; 'wf' :write full; 'wa':write append]")
    }
    return true, []byte("successfully writed(mode : " + mode + ") to object:" + oid)
}