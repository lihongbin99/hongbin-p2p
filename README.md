## 关于 P2P.java 中 getSocket 方法里调用的 NATTest.test(localPort);
在未知自己网络地址交换类型且第一次使用的时候需要执行此代码判断是否支持穿透  
已知网络交换地址为 非Symmetric 或曾经使用此代码验证过一次后即可注销此代码

## 异常: Address already in use: connect

修改两个client的本地端口  
如果修改端口还是出现这个异常则增加一下 P2P.SOCKET_BIND_WAIT_TIME 的值