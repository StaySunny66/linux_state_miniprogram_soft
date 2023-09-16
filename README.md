# linux_state_miniprogram_soft
Linux 性能检测接口
# 建议的使用方法
1. 在你的Linux 设备上安装GO 环境 前往安全组开启服务器1323端口
2. 下载 src 文件夹到你的服务器
3. 执行 go mod tidy 下载依赖包
4. 执行 go run main.go / go build main.go  运行或者生成你所需要的二进制文件
# 或者可以这样(不一定可以用哦~)
你可以选择 使用realse 的二进制直接运行
如下

wget https://github.com/StaySunny66/linux_state_miniprogram_soft/blob/main/releases/main_Linux_X64 

sudo chmod A+X main_Linux_X64  

./main_Linux_X64

# 后续 。。。
你可以使用微信小程序 矢光小屋 扫描 终端的二维码进行绑定，即可轻松在移动设备上查看服务器的实时性能  
你也可以修改源代码中的端口号  
绑定码为web服务器的接口地址的Base64编码  
绑定码 = base64（http://你的ip:1323）  
你也可以使用反向代理到你的域名 进行base64编码后手动添加到微信小程序  
# 参考的项目

https://github.com/vaxilu/x-ui/
