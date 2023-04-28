# 本地影视资源管理软件
### 简介
本地影视资源管理

---  
### 使用教程
[简单使用教程](./usage.md)


---
### 下一步计划
- [ ] 记录导入记录，提供更精细化的管理方式
- [ ] 配置移动到项目中 [issue #1](https://github.com/xi-mad/my_video/issues/1)
- [ ] 评分+排序
- [ ] 增加合集管理
- [ ] 代码优化 



--- 
### 支持以下等功能：
- [x] 登陆，账号密码随便输，请不要部署在公网，没有做任何鉴权  
- [x] 对象管理，支持按照路径、演员、标签、分类检索，没有做全文检索
- [x] 缩略图预览、上一张、下一张
- [x] 路径扫描、路径扫描日志，支持对扫描的对象设置分组、演员、标签  
  扫描的时候不要刷新页面，可能会导致锁库，需要删除./data文件夹下的，除了my_video.db以外的文件，注意是除了my_video.db以外的文件  
- [x] 分组管理，可以这样用，比如：动漫、电影、电视剧、综艺  
- [x] 演员管理，可以这样用，比如：张三、李四  
- [x] 标签管理，自由发挥
- [x] 支持在线播放
- [x] 记录播放次数
- [x] 将程序移动到托盘栏 [issue #1](https://github.com/xi-mad/my_video/issues/1)
- [x] 根据路径去重
- [x] 支持nfo文件 [issue #2](https://github.com/xi-mad/my_video/issues/2)
- [x] 记录日志，存入到log文件夹
- [x] 启动通知，尽力而为

### mtn工具
https://moviethumbnail.sourceforge.net/

### 本地运行
```
clone 本项目
cd my_video     // 进入项目目录
cd backend      // 进入后端目录
go mod tidy     // 下载依赖
go run main.go  // 运行后端
go build        // 打包后端, 生成my_video.exe

另开一个cmd
cd frontend     // 进入前端目录
npm install     // 下载依赖
npm run dev     // 运行前端
npm run build   // 打包前端, 打包后的文件在'../backend/static'目录下
```

### 软件分发
```
创建新的文件夹，比如：my_video
复制：config、data、static、my_video.exe 到my_video文件夹下
运行 my_video.exe
浏览器打开：http://127.0.0.1:8080/  // 注意配置文件中的端口号
```


### 一些预览图
![运行截图](/img/1.png)
![运行截图](/img/2.png)
![运行截图](/img/3.png)
![运行截图](/img/4.png)
![运行截图](/img/5.png)
![运行截图](/img/6.png)
![运行截图](/img/7.png)
![运行截图](/img/8.png)
