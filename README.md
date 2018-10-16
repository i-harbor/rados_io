## Connect to the ceph cluster through the [go-ceph](https://github.com/ceph/go-ceph) and implement the I/O between redos and bytes(RAM)

### Compile the go scripts to dynamic library(.so) by the module Cgo in Goland.
`go build -buildmode=c-shared -o /python/xxx.so xxx.go `

#### Push Bytes to the rados object in python3

>Due to the encoding in C(.so) is ascii(1 character - 1 byte) and in python3 is unicode(1 character - 2 bytes), we use the utf-8(1 character - 1 byte in English, 1 character - 3 bytes in Chinese) encoding for compatibly

```
import ctypes

if __name__ =="__main__":
	ToObject = ctypes.CDLL('toobject.so').ToObject # CDLL
	ToObject.restype = ctypes.c_char_p # declare the expected type returned

	# parameters
	cluster_name = "ceph".encode('utf-8') # cluster name *string*
	user_name    = "client.objstore".encode('utf-8') # user name *string*
	conf_file    = "/etc/ceph/ceph.conf".encode('utf-8') # config file path *string*
	pool_name    = "objstore".encode('utf-8') # pool名称 *string*
	object_name  = "object_name".encode('utf-8') # object name *string*
	data         = "content".encode('utf-8') # data to be written *string*
	offset       = ctypes.c_ulonglong(0)) # write strat from where *ctypes.c_ulonglong*

	# execute
	result = ToObject(cluster_name, user_name, conf_file, pool_name, object_name, data, offset) # return *bytes*
	print(result.decode())
```