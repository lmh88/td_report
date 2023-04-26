# GoFrame Project

https://goframe.org

## 目录说明
本项目有四个入口，main.go 是针对web入口， cmdmain.go 是报表拉取脚本入口， toolmain.go 是小工具集合，
other是其他的调研，测试验证等

### 1 入口文件!

##### 1.1 拉取报表的脚本
report_job

###### addbatch: 手动添加调度的批次,手动快速添加批次，区别于product是可以自定义开始时间和结束时间，报表类型，名称，profileid等
###### product: linux crontab 自动添加调度的批次，参数基本固定
###### dsp_consumer, sd_consumer, sp_consumer, sb_consumer 消费者消费脚本，4种报表类型每种一个
###### dsp_last_month 这个分别是快速拉取dsp一个月的脚本，
###### dsp_divide_file 将dsp快速拉取的脚本中redis消费队列消费处理（文件分割处理，方便下游入库）
###### profile_token 分别获取数据库的ppc和dsp的profile 存放到redis里面，linux上定时执行
###### root 是cobra的入口根命令
###### new_customer 是处理新客户的脚本，开发优化中。将数据库里面的数据拉取出来执行拉取亚马逊的接口拉取数据，同时异步检测拉取结果如果结果不存在则放于消费队列中再次执行拉取亚马逊接口
###### varibles 是公共的变量

##### 1.2 辅助工具
tools
###### varibles 是公共的变量
###### root 是cobra的入口根命令
###### singel 生成命令行shell脚本，快速执行命令拉取数据（基本废弃）
###### singelclient 是调用singelprofile 命令
###### statis 是统计某一个具体的日期拉取报表文件的数量
###### statis_tool 检测14天拉取报表的详细情况，存于数据库中，并且对于拉取失败没有对应的报表文件的，和相隔天数差距大于配置文件的某一参数的重新入列等待消费

##### 1.3 调研的或者测试验证的工具集，可忽略

##### 1.4 s3 上传问题
由于亚马逊s3上传对应的账号是读取环境变量，或者读取特定目录下的账号文件配置，没法读取本地配置文件的账号信息
除非去修改亚马逊提供的go的api。 测试环境中根据配置文件默认是关闭了上传到s3的，避免本地测试上传了文件影响生产环境。
测试环境中s3的桶对应的路劲和生产环境稍有区别。

