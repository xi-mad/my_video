# windows使用教程


### 下载

在[release](https://github.com/xi-mad/my_video/releases)中下载新版本 my_video_full.rar  
下载[mtn](https://moviethumbnail.sourceforge.net/)，这个是生成缩略图的工具

### 解压
解压my_video_full.rar到任意目录，比如：D:\soft\my_video\my_video_full  
mtn也解压到任意目录，比如：D:\soft\mtn-200808a-win32

### 配置
在config文件夹下的config.yaml修改配置信息，配置mth路径，字体路径
```
thumbnail:
  mtn: "D:/soft/mtn-200808a-win32/mtn.exe" 
  font: "C:/Windows/Fonts/STSONG.TTF"
```

### 运行
运行 my_video.exe  

----

# macos使用教程

### 下载

在[release](https://github.com/xi-mad/my_video/releases)中下载新版本 my_video_macos_full.zip  
安装[mtn](https://gitlab.com/movie_thumbnailer/mtn/-/wikis/home#macos)，这个是生成缩略图的工具  
需要使用[Homebrew](https://brew.sh/)安装

安装Homebrew和mtn可能需要使用代理，并且耗时比较长。这两个软件安装较为复杂，不建议没有经验的用户使用macos版本  
可以使用如下命令设置bash的代理，并注意修改端口号，前提是有代理软件在运行
```
 export no_proxy="localhost,127.0.0.1,localaddress,.localdomain.com"
 export https_proxy=http://127.0.0.1:10811 http_proxy=http://127.0.0.1:10811 all_proxy=socks5://127.0.0.1:10810
```

### 解压
解压my_video_macos_full.zip到任意目录，比如：/Users/xxx/my_video_full

### 配置
修改config.yaml配置文件，配置mth路径，字体路径 一行的最后有# macos 的是mac的配置

### 运行
打开命令行，先cd到my_video目录，然后执行`./my_video`，不要关闭命令行

---

# 通用

### 使用
打开浏览器输入：http://127.0.0.1:8080/ (端口号默认是8080，配置文件中可以修改)  
登陆用户名密码随便写，没有校验

### 视频扫描
点击导航栏对象管理，点击扫描路径，输入视频路径，比如：D:\video，点击确定  
扫描完成后关闭日志，然后刷新  
![运行截图](/img/usage/1.png)
![运行截图](/img/usage/2.png)
![运行截图](/img/usage/3.png)

### 文件夹介绍
config          配置文件夹，里面有config.yaml  
data            数据文件夹，数据库文件  
log             日志文件夹，记录运行日志  
static          前端页面文件夹  
temp            临时文件夹，缩略图会临时存放，然后删除
my_video.exe    主程序

### debug
如果程序出错，请贴日志文件中的错误信息  
如果扫描视频出错，请贴扫描日志中的信息