/*
* @Author: Ins
* @Date:   2018-10-10 09:54:12
* @Last Modified by:   Ins
* @Last Modified time: 2018-10-26 15:27:48
*/
package main

import "C"
import (
    "fmt"
    "unsafe"
    "github.com/ceph/go-ceph/rados"
)


func newConn(cluster_name, user_name, conf_file string) (*rados.Conn, error) {
    conn, err := rados.NewConnWithClusterAndUser(cluster_name, user_name)//"ceph","client.objstore"
    if err != nil {
        return nil, err
    }

    err = conn.ReadConfigFile(conf_file)//"/etc/ceph/ceph.conf"
    if err != nil {
        return nil, err
    }

    err = conn.Connect()
    if err != nil {
        return nil, err
    }

    return conn, nil
}

func ReadObjectToBytes(ioctx *rados.IOContext, oid string, block_size int, offset uint64) (int, []byte, error) {
    bytesOut := make([]byte, block_size)
    ret, err := ioctx.Read(oid, bytesOut, offset)
    if err != nil {
        return -1, bytesOut, err
    }
    bytesOut = bytesOut[:ret]
    return ret, bytesOut, err
}

func ObjectListFunc(oid string) {
    fmt.Println(oid)
}

//export ListObj
func ListObj(c_cluster_name *C.char, c_user_name *C.char, c_conf_file *C.char, c_pool_name *C.char) *C.char{
    cluster_name, user_name, conf_file, pool_name := C.GoString(c_cluster_name), C.GoString(c_user_name), C.GoString(c_conf_file), C.GoString(c_pool_name)
    conn, err := newConn(cluster_name, user_name, conf_file)
    if err != nil {
        return C.CString("error when invoke a new connection:" + err.Error())
    }
    defer conn.Shutdown()

    // open a pool handle
    ioctx, err := conn.OpenIOContext(pool_name)
    if err != nil {
        return C.CString("error when openIOContext" + err.Error())
    }
    defer ioctx.Destroy()

    // list the objects in pool just printed in terminal
    ioctx.ListObjects(ObjectListFunc)
    return C.CString("list the objects above in object:" + pool_name)
}

//export FromObj
func FromObj(c_cluster_name *C.char, c_user_name *C.char, c_conf_file *C.char, c_pool_name *C.char, block_size int, c_oid *C.char, offset uint64) (C._Bool, unsafe.Pointer, C.int){
    if block_size > 204800000 {
        result := "the block_size cannot be greater than 204800000"
        return false, unsafe.Pointer(C.CString(result)), C.int(len(result))
    }
    cluster_name, user_name, conf_file, pool_name, oid := C.GoString(c_cluster_name), C.GoString(c_user_name), C.GoString(c_conf_file), C.GoString(c_pool_name), C.GoString(c_oid)
    conn, err := newConn(cluster_name, user_name, conf_file)
    if err != nil {
        result := "error when invoke a new connection:" + err.Error()
        return false, unsafe.Pointer(C.CString(result)), C.int(len(result))
    }
    defer conn.Shutdown()

    // open a pool handle
    ioctx, err := conn.OpenIOContext(pool_name)
    if err != nil {
        result := "error when openIOContext" + err.Error()
        return false, unsafe.Pointer(C.CString(result)), C.int(len(result))
    }
    defer ioctx.Destroy()

    // read the data and write to the file

    ret, bytesOut, err := ReadObjectToBytes(ioctx, oid, block_size, offset)
    if ret == -1 {
        result := "error when read the object to bytes:" + err.Error()
        return false, unsafe.Pointer(C.CString(result)), C.int(len(result))
    }
    return true, C.CBytes(bytesOut), C.int(len(bytesOut))
}

//export ToObj
func ToObj(c_cluster_name *C.char, c_user_name *C.char, c_conf_file *C.char, c_pool_name *C.char, c_oid *C.char, c_bytesIn *C.char, bytesLen C.int, c_mode *C.char, offset uint64) (C._Bool, unsafe.Pointer, C.int){
    cluster_name, user_name, conf_file, pool_name, oid, mode := C.GoString(c_cluster_name), C.GoString(c_user_name), C.GoString(c_conf_file), C.GoString(c_pool_name), C.GoString(c_oid), C.GoString(c_mode)
    bytesIn := C.GoBytes(unsafe.Pointer(c_bytesIn),bytesLen)

    conn, err := newConn(cluster_name, user_name, conf_file)
    if err != nil {
        result := "error when invoke a new connection:" + err.Error()
        return false, unsafe.Pointer(C.CString(result)), C.int(len(result))
    }
    defer conn.Shutdown()

    // open a pool handle
    ioctx, err := conn.OpenIOContext(pool_name)
    if err != nil {
        result := "error when openIOContext:" + err.Error()
        return false, unsafe.Pointer(C.CString(result)), C.int(len(result))
    }
    defer ioctx.Destroy()

    // write data to object
    switch mode {
        case "w":
            err = ioctx.Write(oid, bytesIn, offset)
            if err != nil {
                result := "error when write to object:" + err.Error()
                return false, unsafe.Pointer(C.CString(result)), C.int(len(result))
            }
        case "wf":
            err = ioctx.WriteFull(oid, bytesIn)
            if err != nil {
                result := "error when write full to object:" + err.Error()
                return false, unsafe.Pointer(C.CString(result)), C.int(len(result))
            }
        case "wa":
            err = ioctx.Append(oid, bytesIn)
            if err != nil {
                result := "error when append to object:" + err.Error()
                return false, unsafe.Pointer(C.CString(result)), C.int(len(result))
            }
        default:
            result := "error when write to object: unknown wirte mode : " + mode + ", only ['w' : write; 'wf' :write full; 'wa':write append]"
            return false, unsafe.Pointer(C.CString(result)), C.int(len(result))
    }
    result := "successfully writed(mode : " + mode + ") to object:" + oid
    return true, unsafe.Pointer(C.CString(result)), C.int(len(result))
}



//export DelObj
func DelObj(c_cluster_name *C.char, c_user_name *C.char, c_conf_file *C.char, c_pool_name *C.char, c_oid *C.char) (C._Bool, unsafe.Pointer, C.int){
    cluster_name, user_name, conf_file, pool_name, oid := C.GoString(c_cluster_name), C.GoString(c_user_name), C.GoString(c_conf_file), C.GoString(c_pool_name), C.GoString(c_oid)
    conn, err := newConn(cluster_name, user_name, conf_file)
    if err != nil {
        result := "error when invoke a new connection:" + err.Error()
        return false, unsafe.Pointer(C.CString(result)), C.int(len(result))
    }
    defer conn.Shutdown()

    // open a pool handle
    ioctx, err := conn.OpenIOContext(pool_name)
    if err != nil {
        result := "error when openIOContext:" + err.Error()
        return false, unsafe.Pointer(C.CString(result)), C.int(len(result))
    }
    defer ioctx.Destroy()

    // delete a object 
    err = ioctx.Delete(oid)
    if err != nil {
        result := "error when delete the object:" + err.Error()
        return false, unsafe.Pointer(C.CString(result)), C.int(len(result))
    }
    result := "successfully delete the object:" + oid
    return false, unsafe.Pointer(C.CString(result)), C.int(len(result))
}
func main() {
    
}