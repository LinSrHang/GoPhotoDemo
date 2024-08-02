# GoPhotoDemo

基于go的简易照片直播，仅实现了前端上传、前端查看、手动删除图片文件并刷新，删除功能在新建文件了...

访问 http://localhost:8080/page/index 打开照片展示页面

访问 http://localhost:8080/page/upload 打开图片上传页面

访问 http://localhost:8080/api/v1/static/refresh -Method POST {user,password} 刷新

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
│   ├── refresh.go
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
