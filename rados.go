/*
* @Author: Ins
* @Date:   2018-10-10 09:54:12
* @Last Modified by:   Ins
* @Last Modified time: 2019-01-21 16:10:18
*/
package main

import "C"
import (
    "unsafe"
    "rados_io/rados_io_op"
)

const MAX_RADOS_BYTES uint64 = rados_io_op.MAX_RADOS_BYTES

//export ExistObj
func ExistObj(cluster_name *C.char, user_name *C.char, conf_file *C.char, pool_name *C.char, oid *C.char) (C._Bool, unsafe.Pointer, C.int) {
    stat, data := rados_io_op.RadosExistObj(
        C.GoString(cluster_name), 
        C.GoString(user_name), 
        C.GoString(conf_file), 
        C.GoString(pool_name), 
        C.GoString(oid))

    return C._Bool(stat), C.CBytes(data), C.int(len(data))
}

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
func FromObj(cluster_name *C.char, user_name *C.char, conf_file *C.char, pool_name *C.char, block_size C.int, oid *C.char, offset C.ulonglong) (C._Bool, unsafe.Pointer, C.int) {
    stat, data := rados_io_op.RadosFromObj(
        C.GoString(cluster_name),
        C.GoString(user_name),
        C.GoString(conf_file),
        C.GoString(pool_name),
        int(block_size),
        C.GoString(oid),
        uint64(offset))

    return C._Bool(stat), C.CBytes(data), C.int(len(data))
}


//export ToObj
func ToObj(cluster_name *C.char, user_name *C.char, conf_file *C.char, pool_name *C.char, oid *C.char, bytesAddr unsafe.Pointer, bytesLen C.int, offset C.ulonglong) (C._Bool, unsafe.Pointer, C.int) {
    stat, data := rados_io_op.RadosToObj(
        C.GoString(cluster_name),
        C.GoString(user_name),
        C.GoString(conf_file),
        C.GoString(pool_name),
        C.GoString(oid),
        C.GoBytes(bytesAddr,bytesLen),
        uint64(offset))

    return C._Bool(stat), C.CBytes(data), C.int(len(data))
}

//export DelObj
func DelObj(cluster_name *C.char, user_name *C.char, conf_file *C.char, pool_name *C.char, oid *C.char, osize C.ulonglong) (C._Bool, unsafe.Pointer, C.int) {
    stat, data := rados_io_op.RadosDelObj(
        C.GoString(cluster_name), 
        C.GoString(user_name), 
        C.GoString(conf_file), 
        C.GoString(pool_name), 
        C.GoString(oid),
        uint64(osize))

    return C._Bool(stat), C.CBytes(data), C.int(len(data))
}
func main() {

}