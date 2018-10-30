/*
* @Author: Ins
* @Date:   2018-10-30 16:31:57
* @Last Modified by:   Ins
* @Last Modified time: 2018-10-30 17:15:43
*/
package rados_io_op

import (
    "github.com/ceph/go-ceph/rados"
)
func NewConn(cluster_name, user_name, conf_file string) (*rados.Conn, error) {
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