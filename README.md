# als
a simple als server demo

[OTel Provider](https://istio.io/latest/docs/tasks/observability/logs/otel-provider/) is future.


# 安装 flatbuffer

```shell
# 下载源码
git clone https://github.com/google/flatbuffers.git

# 生产makefile文件
cd flatbuffers
cmake -G "Unix Makefiles"
# 安装
make
make install

# 添加到系统，方便以后使用
sudo cp flatc /usr/local/bin/flatc

# 安装成功后，查看版本号
flatc --version
```
