## Connect to the ceph cluster through the [go-ceph](https://github.com/ceph/go-ceph) and implement the I/O between redos and bytes(RAM)

### Compile the go scripts to dynamic library(.so) by the [Cgo](https://github.com/golang/go/wiki/cgo) in Goland.
`go build -buildmode=c-shared -o ./rados.so rados.go `

### In python3

>Due to the encoding in C(.so) is ascii(1 character - 1 byte) and in python3 is unicode(1 character - 2 bytes), we use the utf-8(1 character - 1 byte in English, 1 character - 3 bytes in Chinese) encoding for compatibly

- Push bytes to the rados object

```
import ctypes

# Return type for ToObj. 
class RetType(ctypes.Structure):
    _fields_ = [('x', ctypes.c_bool),('y', ctypes.c_void_p),('z', ctypes.c_int)]

rados = ctypes.CDLL('./rados.so')
ToObj = rados.ToObj # CDLL
ToObj.restype = RetType # declare the expected type returned
# parameters
cluster_name = "ceph".encode('utf-8') # cluster name. Type:bytes
user_name    = "client.objstore".encode('utf-8') # user name. Type:bytes
conf_file    = "/etc/ceph/ceph.conf".encode('utf-8') # config file path. Type:bytes
pool_name    = "objstore".encode('utf-8') # pool名称. Type:bytes
oid          = "oid".encode('utf-8') # object id. Type:bytes
data         = "content".encode('utf-8') # data to be written. Type:bytes
len          = ctypes.c_int(len(data)) # length of data to be written. Type:ctypes.c_int
mode         = "w".encode('utf-8') # write mode ['w':write,'wf':write full,'wa':write append]. Type:bytes
offset       = ctypes.c_ulonglong(0) # write strat from where(only effective in mode:'w'). Type:ctypes.c_ulonglong

# mode'w': writes len(data) bytes to the object with object id starting at byte offset offset. It returns an error, if any.
# mode'wf': writes len(data) bytes to the object with key oid. The object is filled with the provided data. If the object exists, it is atomically truncated and then written. It returns an error, if any.
# mode'wa': appends len(data) bytes to the object with key oid. The object is appended with the provided data. If the object exists, it is atomically appended to. It returns an error, if any.

# execute
result = ToObj(cluster_name, user_name, conf_file, pool_name, oid, data, len, mode, offset) # return. Type:RetType
stat = result.x # Whether the write operation executed successfully. Type:bool
data = ctypes.string_at(result.y,result.z) # The error(a string which can be convertd to a list, list[0]:code; list[1]:description) or success description. Type:bytes

```

- Get bytes from the rados object

```
import ctypes

# Return type for ToObj. 
class RetType(ctypes.Structure):
    _fields_ = [('x', ctypes.c_bool),('y', ctypes.c_void_p),('z', ctypes.c_int)]

rados = ctypes.CDLL('./rados.so')
FromObj = rados.FromObj # CDLL
FromObj.restype = RetType # declare the expected type returned

# parameters
cluster_name = "ceph".encode('utf-8') # cluster name. Type:bytes
user_name    = "client.objstore".encode('utf-8') # user name. Type:bytes
conf_file    = "/etc/ceph/ceph.conf".encode('utf-8') # config file path. Type:bytes
pool_name    = "objstore".encode('utf-8') # pool名称. Type:bytes
block_size   = ctypes.c_int(204800000) # Maximum number of bytes read each time. Type:ctypes.c_int
oid          = "oid".encode('utf-8') # object id. Type:bytes
offset       = ctypes.c_ulonglong(0) # where read strat from. Type:ctypes.c_ulonglong

# execute
result = FromObj(cluster_name, user_name, conf_file, pool_name, block_size, oid, offset) # return. Type:RetType
stat = result.x # Whether the read operation executed successfully. Type:bool
data = ctypes.string_at(result.y,result.z) # The error(a string which can be convertd to a list, list[0]:code; list[1]:description) or the bytes read out. Type:bytes
```

- Delete an object in pool

```
import ctypes

# Return type for ToObj. 
class RetType(ctypes.Structure):
    _fields_ = [('x', ctypes.c_bool),('y', ctypes.c_void_p),('z', ctypes.c_int)]

rados = ctypes.CDLL('./rados.so')
DelObj = rados.DelObj # CDLL
DelObj.restype = RetType # declare the expected type returned

# parameters
cluster_name = "ceph".encode('utf-8') # cluster name. Type:bytes
user_name    = "client.objstore".encode('utf-8') # user name. Type:bytes
conf_file    = "/etc/ceph/ceph.conf".encode('utf-8') # config file path. Type:bytes
pool_name    = "objstore".encode('utf-8') # pool名称. Type:bytes
oid          = "oid".encode('utf-8') # object id. Type:bytes

# execute
result = DelObj(cluster_name, user_name, conf_file, pool_name, oid) # return. Type:RetType
stat = result.x # Whether the delete operation executed successfully. Type:bool
data = ctypes.string_at(result.y,result.z) # The error(a string which can be convertd to a list, list[0]:code; list[1]:description) or success description. Type:bytes
```

- List the objects in pool

```
import ctypes

# Return type for ListObj. 
class RetType(ctypes.Structure):
    _fields_ = [('x', ctypes.c_bool),('y', ctypes.c_void_p),('z', ctypes.c_int)]

rados = ctypes.CDLL('./rados.so')
ListObj = rados.ListObj # CDLL
ListObj.restype = RetType # declare the expected type returned

# parameters
cluster_name = "ceph".encode('utf-8') # cluster name. Type:bytes
user_name    = "client.objstore".encode('utf-8') # user name. Type:bytes
conf_file    = "/etc/ceph/ceph.conf".encode('utf-8') # config file path. Type:bytes
pool_name    = "objstore".encode('utf-8') # pool名称. Type:bytes

# execute
result = ListObj(cluster_name, user_name, conf_file, pool_name) # return. Type:RetType
stat = result.x # Whether the list operation executed successfully. Type:bool
data = ctypes.string_at(result.y,result.z) # The error(a string which can be convertd to a list, list[0]:code; list[1]:description) or object list bytes. Type:bytes
objects_list = data.decode().split(',') # The list of the objects list. Type:list
# print(objects_list)
```

- Whether an object is existed in pool

```
import ctypes

# Return type for ToObj. 
class RetType(ctypes.Structure):
    _fields_ = [('x', ctypes.c_bool),('y', ctypes.c_void_p),('z', ctypes.c_int)]

rados = ctypes.CDLL('./rados.so')
ExistObj = rados.ExistObj # CDLL
ExistObj.restype = RetType # declare the expected type returned

# parameters
cluster_name = "ceph".encode('utf-8') # cluster name. Type:bytes
user_name    = "client.objstore".encode('utf-8') # user name. Type:bytes
conf_file    = "/etc/ceph/ceph.conf".encode('utf-8') # config file path. Type:bytes
pool_name    = "objstore".encode('utf-8') # pool名称. Type:bytes
oid          = "oid".encode('utf-8') # object id. Type:bytes

# execute
result = ExistObj(cluster_name, user_name, conf_file, pool_name, oid) # return. Type:RetType
stat = result.x # Whether the object is existed in pool. Type:bool
data = ctypes.string_at(result.y,result.z) # The error(a string which can be convertd to a list, list[0]:code; list[1]:description) or success description. Type:bytes
