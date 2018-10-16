## Connect to the ceph cluster through the [go-ceph](https://github.com/ceph/go-ceph) and implement the I/O between redos and bytes(RAM)

### Compile the go scripts to dynamic library(.so) by the [Cgo](https://github.com/golang/go/wiki/cgo) in Goland.
`go build -buildmode=c-shared -o /python/xxx.so xxx.go `

### In python3

>Due to the encoding in C(.so) is ascii(1 character - 1 byte) and in python3 is unicode(1 character - 2 bytes), we use the utf-8(1 character - 1 byte in English, 1 character - 3 bytes in Chinese) encoding for compatibly

#### Push bytes to the rados object

```
import ctypes

if __name__ =="__main__":
	ToObject = ctypes.CDLL('toobject.so').ToObject # CDLL
	ToObject.restype = ctypes.c_char_p # declare the expected type returned

	# parameters
	cluster_name = "ceph".encode('utf-8') # cluster name. type:string
	user_name    = "client.objstore".encode('utf-8') # user name. type:string
	conf_file    = "/etc/ceph/ceph.conf".encode('utf-8') # config file path. type:string
	pool_name    = "objstore".encode('utf-8') # pool名称. type:string
	object_name  = "object_name".encode('utf-8') # object name. type:string
	data         = "content".encode('utf-8') # data to be written. type:string
	offset       = ctypes.c_ulonglong(0)) # write strat from where. type:ctypes.c_ulonglong

	# execute
	result = ToObject(cluster_name, user_name, conf_file, pool_name, object_name, data, offset) # return. type:bytes
	# print(result.decode())
```

#### Get bytes from the rados object

```
import ctypes

if __name__ =="__main__":
	FromObj = ctypes.CDLL('fromobject.so').FromObj # CDLL
	FromObj.restype = ctypes.c_char_p # declare the expected type returned

	# parameters
	cluster_name = "ceph".encode('utf-8') # cluster name. type:string
	user_name    = "client.objstore".encode('utf-8') # user name. type:string
	conf_file    = "/etc/ceph/ceph.conf".encode('utf-8') # config file path. type:string
	pool_name    = "objstore".encode('utf-8') # pool名称. type:string
	block_size   = 204800000 # Maximum number of bytes read each time. type:string
	object_name  = "object_name".encode('utf-8') # object name. type:string
	offset       = ctypes.c_ulonglong(0)) # where read strat from. type:ctypes.c_ulonglong

	# execute
	bytesOut = FromObj(cluster_name, user_name, conf_file, pool_name, block_size, object_name, offset) # return. type:bytes
	# print(bytesOut.decode())
```