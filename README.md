# SmartPingMonitorSystem
SmartPing监控报警通知

## 功能（配置）

- 支持配置监控频率
- 支持配置报警阈值
- 支持配置邮件发送间隔

## 起因
> [SmartPing](https://github.com/smartping/smartping)  是一款 开源、高效、便捷的网络质量监控神器

可惜我在使用的过程中发现其**报警通知**功能非常脆弱，只能在web上面通知，不会发邮件通知或者IM。
这就造成了有些报警通知不能及时收到，未能及时处理故障。于是我利用业余时间做了个邮件报警通知，
就是你现在看到的这个**SmartPingMonitorSystem**！

## 使用
在[RELEASE](https://github.com/cnlyzy/SmartPingMonitorSystem/releases)页面下载系统对应的包 或者 自行clone源码编译（编译前请注意安装完整依赖）。

编辑并保存**conf.ini**存置文件 打开对应的程序即可 Linux 用户可用 **nohup** 命令 让程序在后台执行。
