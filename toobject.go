/*
* @Author: Ins
* @Date:   2018-10-10 09:54:12
* @Last Modified by:   Ins
* @Last Modified time: 2018-10-16 10:17:35
*/
package main
import "C"
import (
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


//export ToObject
func ToObject(c_cluster_name *C.char, c_user_name *C.char, c_conf_file *C.char, c_pool_name *C.char, c_oname *C.char, c_bytesIn *C.char, offset uint64) *C.char{
    cluster_name, user_name, conf_file, pool_name, oname, bytesIn := C.GoString(c_cluster_name), C.GoString(c_user_name), C.GoString(c_conf_file), C.GoString(c_pool_name), C.GoString(c_oname), C.GoString(c_bytesIn)
    conn, err := newConn(cluster_name, user_name, conf_file)
    if err != nil {
        return C.CString("error when invoke a new connection:" + err.Error())
    }
    defer conn.Shutdown()

    // open a pool handle
    ioctx, err := conn.OpenIOContext(pool_name)
    if err != nil {
        return C.CString("error when openIOContext:" + err.Error())
    }
    defer ioctx.Destroy()

    // write data to object
    err = ioctx.Write(oname, []byte(bytesIn), offset)
    if err != nil {
        return C.CString("error when write to object:" + err.Error())
    }

    return C.CString("successfully writed to object：" + oname)
}
func main() {
    
}