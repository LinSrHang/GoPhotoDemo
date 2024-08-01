# GoPhotoDemo
基于go的简易照片直播，仅实现了前端上传和前端查看，删除功能在新建文件了...

访问 http://localhost:8080/page/index 打开照片展示页面
访问 http://localhost:8080/page/upload 打开图片上传页面

```js
v1
│
├── conf/
│   └── config.yaml
│
├── config/
│   └── config.go
│
├── handler/
│   ├── image.go
│   ├── index.go
│   └── upload.go
├── logger/
│   └── logger.go
│
├── template/
│   ├── index.tpl
│   └── upload.tpl
│
├── utils/
│   ├── getImageUrlBase64.go
│   ├── imageCompress.go
│   ├── removeImage.go
│   └── utils.go
│
└── main.go
```
