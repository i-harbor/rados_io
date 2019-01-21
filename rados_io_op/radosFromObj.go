/*
* @Author: Ins
* @Date:   2018-10-30 16:21:00
* @Last Modified by:   Ins
* @Last Modified time: 2019-01-21 16:10:00
*/
package rados_io_op

import (
    "strconv"
    "bytes"
    "github.com/ceph/go-ceph/rados"
)

func readObjToBytes(ioctx *rados.IOContext, oid string, block_size int, offset uint64) ([]byte, error) {
    var err error = nil
    oid_suffix := offset / MAX_RADOS_BYTES
    oid_tmp := oid
    if oid_suffix > 0 {
        oid += "_" + strconv.FormatUint(oid_suffix,10)
    }
    offset %= MAX_RADOS_BYTES
    
    bytesOut := make([]byte, block_size)
    ret, err := ioctx.Read(oid, bytesOut, offset)
    if err != nil {
        return bytesOut, err
    }
    bytesOut = bytesOut[:ret]
    block_size_tmp := int(offset + uint64(block_size) - MAX_RADOS_BYTES)
    // read the bytes cyclically if the data size user want greater than MAX_RADOS_BYTES
    if(block_size_tmp > 0){
        bytesOut_tmp := make([]byte, block_size_tmp)
        ret_tmp, err_tmp := ioctx.Read(oid_tmp + "_" + strconv.FormatUint(oid_suffix+1,10), bytesOut_tmp, 0)
        for err_tmp == nil && len(bytesOut) < block_size {
            bytesOut_tmp = bytesOut_tmp[:ret_tmp]
            var buffer bytes.Buffer
            buffer.Write(bytesOut)
            buffer.Write(bytesOut_tmp)
            bytesOut = buffer.Bytes()

            oid_suffix++
            block_size_tmp -= int(MAX_RADOS_BYTES)
            if block_size_tmp <= 0 {
                break
            }

            bytesOut_tmp = make([]byte, block_size_tmp)
            ret_tmp, err_tmp = ioctx.Read(oid_tmp + "_" + strconv.FormatUint(oid_suffix+1,10), bytesOut_tmp, 0)
        }
    }
    return bytesOut, err
}

func RadosFromObj(cluster_name string, user_name string, conf_file string, pool_name string, block_size int, oid string, offset uint64) (bool, []byte) {
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

    bytesOut, err := readObjToBytes(ioctx, oid, block_size, offset)
    if err != nil {
        return false, []byte(err.Error() + ",error when read the object(oid:" + oid + ") to bytes:" + ERR_INFO[err.Error()])
    }
    return true, bytesOut
}