# vaga.go

Another personal video site written by golang.

> version: 0.1  
> ci: none  
> release: none  

## Journal

-- update

最近偶尔会用用, 但是 window 机器在卧室里略吵, 所以打算移到(储物室的)服务器上去. 稍微把配置提取了一下. 

也方便有相同需求人用(如果有的话).

-- init

春节回家想起家里的机器上还有些视频想看, 所以干脆写个在线点播的网站.

暂时没打算完善所以配置写死的, 细节也没费心只有简单整理, 毕竟想起来写这个的时候离高铁发车只有三四个钟头了.

## Install & Config

1. 安装 `go get github.com/bagaking/vaga.go` , 并配置好 PATH
2. 在准备运行 vaga 的路径配置 `vaga.json` 文件, 配置详见 `vaga.json.sample` 文件  
    > 运行路径不一定需要是视频存储路径  
    不过如果视频在运行路径下, vaga.json 中可以配置相对路径  
3. 将 tpl 拷贝到运行路径下: `cp -R $GOPATH/src/github.com/bagaking/vaga.go/tpl ./`
4. 运行 `vaga.go` 即可, 由于是直接推流, 所以支持的格式只有 `.mp4 | .ogg | .webm`

## Deploy

考虑家用视频一般在本地网络的Windows机器下, NAT几乎是必须的, 以我自己的配置为例:  

```
香港VPS -FRP-> 
- 内网NAT中继服务器 -NGINX-> 
 - WINDOWS机器 -> 
  - vaga.go(*:9001)
```

当然如果有外网ip和网关权限的情况下, 事情会简单很多, 一个 DDNS + 端口转发 就搞定了  
如果用的是阿里云的DNS服务, 可以用 [ddns-aliyun](https://www.npmjs.com/package/ddns-aliyun) 快速搭一个

## Feature

1. 支持指定 `VideoBlob` 根目录, 所指定的目录将被显示在首页, 
目前的配置方式是增删改 `/def.go` 中的 `blob<Name,Desc,RootPath>`, 并填在`AllAvailableVideoBlobs`中.
适当配置后, 自动生成 `index`, `tree view`, `watch view`, 可以直接在线点播.

2. 目前仅支持 MP4 格式

3. 视频播放器使用 `video.js` 而非简单视频流, 因此可以支持视频的变速播放. 
播放器配置可以在 `/tpl/watch.html` 中修改, 具体可参照 `video.js` 的官方文档.



