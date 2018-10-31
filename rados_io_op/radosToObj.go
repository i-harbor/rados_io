/*
* @Author: Ins
* @Date:   2018-10-30 16:21:00
* @Last Modified by:   Ins
* @Last Modified time: 2018-10-31 18:04:54
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
    var oid_suffix_list string
    switch {
        case oid_suffix == 0 && oid_suffix_gap > 0:
            oid_suffix_list = "(" + oid + " to " + oid + "__" + strconv.FormatUint(oid_suffix_gap, 10) + ")"
        case oid_suffix > 0 && oid_suffix_gap == 0:
            oid_suffix_list = "(" + oid + "__" + strconv.FormatUint(oid_suffix, 10)
        case oid_suffix > 0 && oid_suffix_gap >0:
            oid_suffix_list = "(" + oid + "__" + strconv.FormatUint(oid_suffix, 10) + " to " + oid + "__" + strconv.FormatUint(oid_suffix + oid_suffix_gap, 10) + ")"
        default:
            oid_suffix_list = ""
    }

    // write to the rados cyclically if the data size greater MAX_RADOS_BYTES
    for oid_suffix_gap > 0 {
        bytesIn_tmp := bytesIn[(oid_suffix_gap * MAX_RADOS_BYTES - offset):]
        err = ioctx.Write(oid + "__" + strconv.FormatUint(oid_suffix + oid_suffix_gap, 10), bytesIn_tmp, 0)
        if err != nil {
            return err, ""
        }
        bytesIn = bytesIn[:(oid_suffix_gap * MAX_RADOS_BYTES - offset)]
        oid_suffix_gap --
    }
    if oid_suffix > 0 {
        err = ioctx.Write(oid + "__" + strconv.FormatUint(oid_suffix, 10), bytesIn, offset)
    } else {
        err = ioctx.Write(oid, bytesIn, offset)
    }
    
    if err != nil {
        return err, oid_suffix_list
    }

    return nil, oid_suffix_list
}

func writeFulToObj(ioctx *rados.IOContext, oid string, bytesIn []byte) (error, string) {
    // delete the rados fully
    deleteObj(ioctx, oid)

    oid_suffix := uint64(len(bytesIn)) / MAX_RADOS_BYTES

    // get the oids lists writed
    var oid_suffix_list string
    if oid_suffix > 0 {
        oid_suffix_list = "(" + oid + " to " + oid + "__" + strconv.FormatUint(oid_suffix, 10) + ")"
    } else {
        oid_suffix_list = ""
    }

    // write to the rados cyclically if the data size greater MAX_RADOS_BYTES
    for oid_suffix > 0 {
        bytesIn_tmp := bytesIn[(oid_suffix * MAX_RADOS_BYTES):]
        err := ioctx.WriteFull(oid + "__" + strconv.FormatUint(oid_suffix, 10), bytesIn_tmp)
        if err != nil {
            return err, ""
        }
        bytesIn = bytesIn[:(oid_suffix * MAX_RADOS_BYTES)]
        oid_suffix--
    }

    err := ioctx.WriteFull(oid, bytesIn)
    if err != nil {
        return err, ""
    }
    return nil, oid_suffix_list
}

func writeAppendToObj(ioctx *rados.IOContext, oid string, bytesIn []byte) (error, string) {
    data, _ := readObjToBytes(ioctx, oid, int(MAX_RADOS_BYTES), 0)
    offset := len(data)

    // get the real offset of the rados(the end of rados)
    i := 0
    for len(data) == int(MAX_RADOS_BYTES) {
        i++
        data, _ = readObjToBytes(ioctx, oid + "__" + strconv.Itoa(i), int(MAX_RADOS_BYTES), 0)
        offset += len(data)
    }

    // write to rados from the real offset(the end of rados) as append
    err, oid_suffix_list := writeToObj(ioctx, oid, bytesIn, uint64(offset))
    return err, oid_suffix_list
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
    var oid_suffix_list string = ""
    switch mode {
        case "w":
            err, oid_suffix_list = writeToObj(ioctx, oid, bytesIn, offset)
            if err != nil {
                return false, []byte("error when write to object:" + err.Error())
            }
        case "wf":
            err, oid_suffix_list = writeFulToObj(ioctx, oid, bytesIn)
            if err != nil {
                return false, []byte("error when write full to object:" + err.Error())
            }
        case "wa":
            err, oid_suffix_list = writeAppendToObj(ioctx, oid, bytesIn)
            if err != nil {
                return false, []byte("error when append to object:" + err.Error())
            }
        default:
            return false, []byte("error when write to object: unknown wirte mode : " + mode + ", only ['w' : write; 'wf' :write full; 'wa':write append]")
    }
    return true, []byte("successfully writed(mode : " + mode + ") to object:" + oid + oid_suffix_list)
}