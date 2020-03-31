# wechatrobot-sender

Nightingale的理念，是将告警事件扔到redis里就不管了，接下来由各种sender来读取redis里的事件并发送，毕竟发送报警的方式太多了，适配起来比较费劲，希望社区同仁能够共建。

这里提供一个微信机器人的sender，参考了[https://github.com/n9e/wechat-sender](https://github.com/wechat-sender)。

夜莺后台创建虚拟user，企业微信机器人webhook的key值作为虚拟user的im，需要在夜莺项目monapi.yaml设置里的notify添加im告警，如下：
```
notify:
  p1: ["im"]
  p2: ["im"]
  p3: ["im"]
```

## compile

```bash
cd $GOPATH/src
mkdir -p github.com/n9e
cd github.com/n9e
git clone https://github.com/n9e/wechatrobot-sender.git
cd wechatrobot-sender
./control build
```

如上编译完就可以拿到二进制了。

## configuration

直接修改etc/wechatrobot-sender.yml即可

## pack

编译完成之后可以打个包扔到线上去跑，将二进制和配置文件打包即可：

```bash
tar zcvf wechatrobot-sender.tar.gz wechatrobot-sender etc/wechatrobot-sender.yml etc/wechatrobot.tpl
```

## test

配置etc/wechatrobot-sender.yml，相关配置修改好，我们先来测试一下是否好使， `./wechatrobot-sender -t <toUser>`，程序会自动读取etc目录下的配置文件，发一个测试消息给`toUser`

## run

如果测试发送没问题，扔到线上跑吧，使用systemd或者supervisor之类的托管起来，systemd的配置实例：


```
$ cat wechatrobot-sender.service
[Unit]
Description=Nightingale wechatrobot sender
After=network-online.target
Wants=network-online.target

[Service]
User=root
Group=root

Type=simple
ExecStart=/home/n9e/wechatrobot-sender
WorkingDirectory=/home/n9e

Restart=always
RestartSec=1
StartLimitInterval=0

[Install]
WantedBy=multi-user.target
```