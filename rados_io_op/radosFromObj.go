/*
* @Author: Ins
* @Date:   2018-10-30 16:21:00
* @Last Modified by:   Ins
* @Last Modified time: 2018-10-31 16:20:09
*/
package rados_io_op

import (
    "strconv"
    "bytes"
    "github.com/ceph/go-ceph/rados"
)

func readObjectToBytes(ioctx *rados.IOContext, oid string, block_size int, offset uint64) ([]byte, error) {
    var err error = nil
    oid_suffix := offset / MAX_RADOS_BYTES
    oid_tmp := oid
    if oid_suffix > 0 {
        oid += "__" + strconv.FormatUint(oid_suffix,10)
    }
    offset %= MAX_RADOS_BYTES
    
    bytesOut := make([]byte, block_size)
    ret, err := ioctx.Read(oid, bytesOut, offset)
    if err != nil {
        return bytesOut, err
    }
    bytesOut = bytesOut[:ret]

    if(offset + uint64(block_size) > MAX_RADOS_BYTES){
        bytesOut_tmp := make([]byte, block_size)
        ret_tmp, err_tmp := ioctx.Read(oid_tmp + "__" + strconv.FormatUint(oid_suffix+1,10), bytesOut_tmp, 0)
        for err_tmp == nil && len(bytesOut) < block_size {
            oid_suffix++
            bytesOut_tmp = bytesOut_tmp[:ret_tmp]
            var buffer bytes.Buffer
            buffer.Write(bytesOut)
            buffer.Write(bytesOut_tmp)
            bytesOut = buffer.Bytes()
            bytesOut_tmp = make([]byte, block_size)
            ret_tmp, err_tmp = ioctx.Read(oid_tmp + "__" + strconv.FormatUint(oid_suffix+1,10), bytesOut_tmp, 0)
        }
    }
    return bytesOut, err
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

    bytesOut, err := readObjectToBytes(ioctx, oid, block_size, offset)
    if err != nil {
        return false, []byte("error when read the object(oid:" + oid + ") to bytes:" + err.Error())
    }
    return true, bytesOut
}