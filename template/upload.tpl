<!DOCTYPE html>
<html lang="en">

<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Document</title>

    <style>
        /*外层div*/
        .input-file-box {
            border: 1px solid gray;
            width: 150px;
            height: 150px;
            position: relative;
            text-align: center;
            border-radius: 8px;
        }

        /*文字描述*/
        .input-file-box>span {
            display: block;
            width: 100px;
            height: 30px;
            position: absolute;
            top: 0px;
            bottom: 0;
            left: 0;
            right: 0;
            margin: auto;
            color: gray;
        }

        /*input框*/
        .input-file-box #uploadImage, .input-file-box  #uploadZip {
            opacity: 0;
            width: 100%;
            height: 100%;
            cursor: pointer;
        }
    </style>
</head>

<body>
    <div class="input-file-box">
        <span>点击上传图片</span>
        <input type="file" name="" id="uploadImage" accept=".jpeg, .jpg, .png" multiple>
    </div>
    <div class="input-file-box">
        <span>点击上传图片压缩包</span>
        <input type="file" name="" id="uploadZip" accept=".zip" multiple>
    </div>

    <script src="https://unpkg.com/axios/dist/axios.min.js"></script>
    <script>
        window.onload = function () {
            let input = document.getElementById("uploadImage")
            // 当用户上传时触发事件
            input.onchange = function () {
                let fileList = this.files
                for (let i = 0; i < fileList.length; i++) {
                    let form = new FormData()
                    form.append("image", fileList[i])
                    // axios({
                    //     method: "POST",
                    //     url: {{ .serverUrl }} + '/api/v1/static/post',
                    //     data: form,
                    // })
                    axios.post({{ .serverUrl }} + '/api/v1/static/post', form)
                }
                alert('Upload successful.')
            }
            let zipInput = document.getElementById("uploadZip")
            zipInput.onchange = function () {
                let zipFile = this.files
                for (let i = 0; i < zipFile.length; i++) {
                    let form = new FormData()
                    form.append("imageZip", zipFile[i])
                    // axios({
                    //     method: 'POST',
                    //     url: {{ .serverUrl }} + '/api/v1/static/postZip',
                    //     data: form,
                    // })
                    axios.post({{ .serverUrl }} + '/api/v1/static/postZip', form)
                }
            }
        }

    </script>
</body>

</html>