# JetIot 捷特物联网接入服务器
> ###为什么叫 Jet 
> Jet(捷特)指 游戏《Fallout》系列中的一种药物，在游戏设定中是一种相当强大的兴奋剂。
> 它可以刺激中央神经系统。当使用时，会感到一股力量，但是只会持续几分钟而已。
> 此处寓意是希望以最快的速度搭建物联平台并接入设备和终端。

## 一、运行环境：
> 推荐使用docker搭建环境
- mysql 5.8
- redis-server
- EMQ X Broker
## 二、基本原理
架构主要分为两个部分：
- **HttpServer**：用来与用户终端进行交互使用gin框架
- **MqttClient**：用来和EMQ进行交互用来获取设备信息使用使用paho.mqtt库构建了一个基于事件触发和函数回调的消息系统

---
~~别的等我下次上班摸鱼再写吧~~

硬件接入方案暂定 [看这里](https://github.com/xsj321/JetIot/blob/master/apis/MQTT%E7%89%A9%E9%80%9A%E4%BF%A1%E6%A0%BC%E5%BC%8F/%E7%89%A9%E5%9F%BA%E6%9C%AC%E9%80%9A%E4%BF%A1%E6%A0%BC%E5%BC%8F%E8%AF%B4%E6%98%8E.md)
