## 通过go-ceph模块连接ceph集群，并实现rados对象和文件之间的读写。

### 通过go语言的Cgo模块，将go脚本编译成动态库.so文件，用于其他语言的调用
`go build -buildmode=c-shared -o xxx.so xxx.go `

#### 存储字节流到rados对象
```
import ctypes

if __name__ =="__main__":
	ToObject = ctypes.CDLL('toobject.so').ToObject # 建立CDLL对象
	ToObject.restype = ctypes.c_char_p # 设置返回数据类型

	# 参数
	cluster_name = "ceph".encode('utf-8') # 集群名称 string
	user_name    = "client.objstore".encode('utf-8') # 用户名称 string
	conf_file    = "/etc/ceph/ceph.conf".encode('utf-8') # 配置文件地址 string
	pool_name    = "objstore".encode('utf-8') # pool名称 string
	object_name  = "object_name".encode('utf-8') # 写入对象名 string
	content      = "content".encode('utf-8') # 写入数据 string
	offset       = ctypes.c_ulonglong(0)) # 偏移量，从第几字节开始写 ctypes.c_ulonglong
	
	#执行
	result = ToObject() # 返回写入结果 bytes()
	print(result.decode())
```