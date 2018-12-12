/*
* @Author: Ins
* @Date:   2018-12-12 10:21:00
* @Last Modified by:   Ins
* @Last Modified time: 2018-12-12 14:46:32
*/
package rados_io_op

func RadosExistObj(cluster_name string, user_name string, conf_file string, pool_name string, oid string) (bool, []byte) {
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

    //open a iter
    iter, err := ioctx.Iter()
    if err != nil {
        return false, []byte(err.Error() + ",error when openIter:" + ERR_INFO[err.Error()])
    }
    defer iter.Close()

    for iter.Next() {
        if iter.Value() == oid {
            return true, []byte("The object is really existed")
        }
    }
    return false, []byte("2,error:" + ERR_INFO["2"])
}