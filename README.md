# ngweb


# golang lib

### github.com/rcrowley/goagain
    在Hacker News看到用Go达到Zero-downtime restarts，意思大概为零下线时间式重启，很早就知道nginx可以轻松做到平滑重启，一直都想用go来实现这样的功能。看了一下它的代码实现，所以有了这篇博文。

项目名字叫goagain，地址在：https://github.com/rcrowley/goagain。该项目是一个类库，也就是package，在go开发的程序中添加这个package就可以轻松地重启程序。

goagain会监控2个系统信号，一个为SIGTERM，接收到这个信号，程序就停止运行。另一个信号为SIGUSR2，接收到这个信号的行为是，当前进程，也就是父进程会新建一个子进程，然后把父进程的pid保存到一个名为GOAGAIN_PPID的环境变量；子进程启动的时候会检索GOAGAIN_PPID这个变量，来判断程序是否要重启，通过这个变量来关闭父进程，来达到平滑重启的效果


### redis tools
redis--使用redis-rdb-tools分析redis的内存使用情况

早就听说redis性能卓越，不过难以使用，但看了Hacker News的报道后，还是被它的难度吓倒了。

原文标题：From 1.5 GB to 50 MB: The Story of My Redis Database，链接在这里：http://davidcel.is/blog/2013/03/20/the-story-of-my-redis-database/

文章的内容简单的来说（英文水平不够，有错莫怪），刚开始，由于设计不好，redis的内存使用达到了1.5GB，服务器出现崩溃的问题。开始优化，把类似recommendable:users:1234:liked_beers这样的长组合健简化为u:1234:lb短组合的健，这样的改变只节省了10MB的内存。后来使用了redis-rdb-tools，分析了一下redis的内存使用情况，才得出了服务器内存使用过高的真正瓶颈所在。简化业务后，redis的使用情况变为50MB，惊人的变化。作者在文章后面说，在以后的一段时间里，应该不会出现redis内存使用过高的问题了。具体的内容，大家还是可以去看看，了解一下。

redis-rdb-tools可以把redis的数据库转变为json文件，生成redis的内存使用报告。通过它，可以对redis的使用情况有个大致的了解。redis-rdb-tools的代码在这里：

https://github.com/sripathikrishnan/redis-rdb-tools


### linux system info to influx

github.com/novaquark/sysinfo_influxdb



### influx manager webUI

https://github.com/influxdb/grafana.git
https://github.com/nareix/curl

github.com/e-dard/netbug


### research

https://github.com/etsy/Hound
https://github.com/tobegit3hub/seagull
https://github.com/pksunkara/alpaca

http://www.noodlesoft.com/blog/

https://github.com/mikespook/gorbac

### tools
https://github.com/Induction/Induction

https://github.com/dropbox/godropbox/

https://github.com/wangtuanjie/ip17mon  // ip地址查询

https://github.com/blevesearch/bleve    // golang fulltext search engine


https://github.com/oschwald/geoip2-golang

### data struct and sortting

https://github.com/ITCase/sqlalchemy_mptt


### postgres tools

https://github.com/jackc/tern   迁移工具

