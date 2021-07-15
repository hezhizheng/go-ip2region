# go-ip2region

> 基于 [ip2region](https://github.com/lionsoul2014/ip2region)  使用 go 扩展的一个简单的 IP 地址归属地查询服务

## 功能
- 提供 http 接口查询 IP 地址归属地
- 提供命令行 下载/更新 最新的 `ip2region.db` ip 库 (数据来源原仓库)

## 使用
可直接下载 [releases](https://github.com/hezhizheng/go-ip2region/releases) 文件启动即可，可选参数说明：
```
./go-ip2region_windows_amd64.exe -h
Usage of D:\go-ip2region\go-ip2region_windows_amd64.exe:
  -d string
        仅用于下载最新的ip地址库，保存在当前目录 (default "0")
  -p string
        本地监听的端口 (default "9090")
```

## 启动http服务
```
// 没有IP地址库会自动下载，保存在当前目录
./go-ip2region_windows_amd64.exe

// 没有指定IP会获取当前客户端IP
curl http://127.0.0.1:9090
curl http://127.0.0.1:9090?ip=59.42.37.186

// 返回数据格式
{
    "code": 200,
    "msg": "",
    "data": {
        "ip": "59.42.37.186",
        "country": "中国",
        "province": "广东",
        "city": "广州",
        "county": "0",
        "isp": "电信"
    }
}
```
![](https://files.catbox.moe/q36ces.png)
![](https://files.catbox.moe/n4j1h5.png)

## 下载/更新 IP 地址库
```
// 仅用于下载/更新 IP 地址库
./go-ip2region_windows_amd64.exe -d 1
```

## 自行编译
```
// 跨平台编译
gox -osarch="windows/amd64" -ldflags "-s -w" -gcflags="all=-trimpath=${PWD}" -asmflags="all=-trimpath=${PWD}"

gox -osarch="darwin/amd64" -ldflags "-s -w" -gcflags="all=-trimpath=${PWD}" -asmflags="all=-trimpath=${PWD}"

gox -osarch="linux/amd64" -ldflags "-s -w" -gcflags="all=-trimpath=${PWD}" -asmflags="all=-trimpath=${PWD}"
```

## 部署Nginx
```nginx
server
{
    listen 80;
    listen 443 ssl http2;
    server_name ip.hzz.cool;
    
    ssl_certificate /path/fullchain.cer;   
    ssl_certificate_key /path/hzz.cool.key;   
    ssl_session_timeout  5m;  
    ssl_protocols TLSv1.1 TLSv1.2 TLSv1.3;  
    ssl_ciphers  ECDHE-RSA-AES128-GCM-SHA256:HIGH:!aNULL:!MD5:!RC4:!DHE;  
    ssl_prefer_server_ciphers  on;
    
    
    location / {
        try_files /_not_exists_ @backend;
    }
    
    location @backend {
        proxy_set_header X-Forwarded-For $remote_addr;
        proxy_set_header Host            $http_host;
        proxy_pass http://127.0.0.1:9090;
    }
    
    access_log  /www/wwwlogs/ip2region.log;
    error_log  /www/wwwlogs/ip2region.error.log;
}

```

## License
[MIT](./LICENSE.txt)
