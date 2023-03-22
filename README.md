# 网络监测小工具 (Net-Test)

网络监测小工具是一个用Golang编写的网络性能监测工具，它可以定时测试网络延迟、丢包率和下载速度，并将结果记录到Excel文件中。

![%E7%BD%91%E7%BB%9C%E6%A3%80%E6%B5%8B%E8%AE%B0%E5%BD%95.png](https://cdn.wyr.me/post-files/2023-03-22/1679493322261/网络检测记录.png)

## 功能

- 每5分钟自动测试网络状态
- 记录主路由、网关、百度网站以及服务器的平均延迟和丢包率
- 记录两个下载链接的下载速度
- 将结果保存到Excel文件中

## 开发

确保已安装Go，然后在项目根目录运行以下命令：

```sh
go get -u github.com/xuri/excelize/v2
go get -u github.com/prometheus-community/pro-bing
```

## 开发环境运行

在项目根目录下运行以下命令：

```sh
go run main.go
```

程序运行后，结果将记录在当前目录下的`网络检测记录.xlsx`文件中。

## 编译

在项目根目录下运行以下命令：

```sh
go build .
```

将会在`bin`目录下生成`net-test`和`net-test.exe`两个文件。你需要根据对应操作系统和架构编译程序。

## 运行

### Unix

```sh
./net-test
```

### Windows

双击`net-test.exe`文件。

## 开源许可

本项目采用MIT许可证。有关详细信息，请参阅[LICENSE](LICENSE)文件。
