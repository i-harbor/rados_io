/*
* @Author: Ins
* @Date:   2018-10-30 16:21:00
* @Last Modified by:   Ins
* @Last Modified time: 2019-01-21 15:27:33
*/
package rados_io_op

import (
    "strconv"
    "github.com/ceph/go-ceph/rados"
)

func writeToObj(ioctx *rados.IOContext, oid string, bytesIn []byte, offset uint64) (error, string) {
    var err error = nil
    oid_suffix := offset / MAX_RADOS_BYTES
    offset %= MAX_RADOS_BYTES
    oid_suffix_gap := (offset + uint64(len(bytesIn))) / MAX_RADOS_BYTES

    // get the oids lists writed
    var oid_suffix_list string = ""
    switch {
        case oid_suffix == 0 && oid_suffix_gap > 0:
            oid_suffix_list = "(" + oid + " to " + oid + "_" + strconv.FormatUint(oid_suffix_gap, 10) + ")"
        case oid_suffix > 0 && oid_suffix_gap == 0:
            oid_suffix_list = "(" + oid + "_" + strconv.FormatUint(oid_suffix, 10)
        case oid_suffix > 0 && oid_suffix_gap >0:
            oid_suffix_list = "(" + oid + "_" + strconv.FormatUint(oid_suffix, 10) + " to " + oid + "_" + strconv.FormatUint(oid_suffix + oid_suffix_gap, 10) + ")"
        default:
            oid_suffix_list = ""
    }

    // write to the rados cyclically if the data size greater MAX_RADOS_BYTES
    for oid_suffix_gap > 0 {
        bytesIn_tmp := bytesIn[(oid_suffix_gap * MAX_RADOS_BYTES - offset):]
        err = ioctx.Write(oid + "_" + strconv.FormatUint(oid_suffix + oid_suffix_gap, 10), bytesIn_tmp, 0)
        if err != nil {
            return err, ""
        }
        bytesIn = bytesIn[:(oid_suffix_gap * MAX_RADOS_BYTES - offset)]
        oid_suffix_gap --
    }
    if oid_suffix > 0 {
        err = ioctx.Write(oid + "_" + strconv.FormatUint(oid_suffix, 10), bytesIn, offset)
    } else {
        err = ioctx.Write(oid, bytesIn, offset)
    }
    
    if err != nil {
        return err, oid_suffix_list
    }

    return nil, oid_suffix_list
}



func RadosToObj(cluster_name string, user_name string, conf_file string, pool_name string, oid string, bytesIn []byte, offset uint64) (bool, []byte) {
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

    // write data to object
    var oid_suffix_list string = ""

    err, oid_suffix_list = writeToObj(ioctx, oid, bytesIn, offset)
    if err != nil {
        return false, []byte(err.Error() + ",error when write to object:" + ERR_INFO[err.Error()])
    }

    return true, []byte("successfully write to object:" + oid + oid_suffix_list)
}