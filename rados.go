/*
* @Author: Ins
* @Date:   2018-10-10 09:54:12
* @Last Modified by:   Ins
* @Last Modified time: 2018-10-30 22:02:39
*/
package main

import "C"
import (
    "unsafe"
    "strconv"
    "rados_io/rados_io_op"
)

const MAX_RADOS_BYTES uint64 = rados_io_op.MAX_RADOS_BYTES
const MAX_BLOCK_SIZE int = rados_io_op.MAX_BLOCK_SIZE

//export ListObj
func ListObj(cluster_name *C.char, user_name *C.char, conf_file *C.char, pool_name *C.char) (C._Bool, unsafe.Pointer, C.int) {
    stat, data := rados_io_op.RadosListObj(
        C.GoString(cluster_name),
        C.GoString(user_name),
        C.GoString(conf_file),
        C.GoString(pool_name))

    return C._Bool(stat), C.CBytes(data), C.int(len(data))
}

//export FromObj
func FromObj(cluster_name *C.char, user_name *C.char, conf_file *C.char, pool_name *C.char, block_size int, oid *C.char, offset uint64) (C._Bool, unsafe.Pointer, C.int) {
    if block_size > int(MAX_RADOS_BYTES) {
        result := []byte("the block_size cannot be greater than MAX_RADOS_BYTES:" + strconv.FormatUint(offset / MAX_RADOS_BYTES,10))
        return false, C.CBytes(result), C.int(len(result))
    }
    if block_size > MAX_BLOCK_SIZE {
        result := []byte("the block_size cannot be greater than MAX_BLOCK_SIZE:" + strconv.Itoa(MAX_BLOCK_SIZE))
        return false, C.CBytes(result), C.int(len(result))
    }

    stat, data := rados_io_op.RadosFromObj(
        C.GoString(cluster_name),
        C.GoString(user_name),
        C.GoString(conf_file),
        C.GoString(pool_name),
        block_size,
        C.GoString(oid),
        offset)

    return C._Bool(stat), C.CBytes(data), C.int(len(data))
}


//export ToObj
func ToObj(cluster_name *C.char, user_name *C.char, conf_file *C.char, pool_name *C.char, oid *C.char, bytesAddr unsafe.Pointer, bytesLen int, mode *C.char, offset uint64) (C._Bool, unsafe.Pointer, C.int) {
    if bytesLen > 2147483647 {
        result := []byte("the length of data cannot be greater than the positive range of C.int : 2147483647:")
        return false, C.CBytes(result), C.int(len(result))
    }
    bytesIn := C.GoBytes(bytesAddr,C.int(bytesLen))

    stat, data := rados_io_op.RadosToObj(
        C.GoString(cluster_name),
        C.GoString(user_name),
        C.GoString(conf_file),
        C.GoString(pool_name),
        C.GoString(oid),
        bytesIn,
        bytesLen,
        C.GoString(mode),
        offset)

    return C._Bool(stat), C.CBytes(data), C.int(len(data))
}

//export DelObj
func DelObj(cluster_name *C.char, user_name *C.char, conf_file *C.char, pool_name *C.char, oid *C.char) (C._Bool, unsafe.Pointer, C.int) {
    stat, data := rados_io_op.RadosDelObj(
        C.GoString(cluster_name), 
        C.GoString(user_name), 
        C.GoString(conf_file), 
        C.GoString(pool_name), 
        C.GoString(oid))

    return C._Bool(stat), C.CBytes(data), C.int(len(data))
}
func main() {

}