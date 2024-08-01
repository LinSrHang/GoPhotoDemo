<!DOCTYPE html>
<html>

<head>
    <meta charset="utf-8">
    <title>Document</title>
    <style>
        * {
            margin: 0;
            padding: 0;
        }

        li {
            list-style: none;
        }

        body {
            background-color: white;
            width: 100vw;
            height: 100vh;
            overflow: hidden;
        }

        h1 {
            text-align: center;
            height: 50px;
            line-height: 50px;
            color: #0059A4;
        }

        .box {
            width: 100%;
            height: 90vh;
            margin-top: 10px;
            margin-left: 10px;
        }

        .list1,
        .list2 {
            float: left;
            border: 5px solid #0059A4;
        }

        .list1 {
            width: 70%;
            height: 100%;
            overflow: hidden;
        }

        .list2 {
            width: 300px;
            height: 100%;
            margin-left: 0.75%;
            overflow-y: scroll;
            overflow-x: hidden;
        }

        .list1 img {
            width: 100%;
        }

        .list2 img {
            width: 300px;
        }

        button {
            background-color: white;
            width: 100px;
            height: 50px;
            font-size: larger;
            margin-left: 10px;
            margin-top: 10px;
            border-color: #0059A4;
        }

        button:hover {
            border-color: white;
            color: white;
            background-color: #0059A4;
        }
    </style>
</head>

<body>
    <h1>电子相册</h1>
    <div class="box">
        <ul class="list1">
            <!-- <li><img src="../img/danji.jpg" id="img1"></li> -->
        </ul>
        <ul class="list2">
            <li><img alt=""></li>
            <li><img alt=""></li>
            <li><img alt=""></li>
            <li><img alt=""></li>
        </ul>
        <button id="sub">-</button>
        <h1></h1>
        <button id="add">+</button>
    </div>

    <script src="https://unpkg.com/axios/dist/axios.min.js"></script>
    <script>
        let page = 1
        let maxPage = 0
        
        let load = function (page) {
            axios.get({{ .serverUrl }} + `/api/v1/static/fidList?page=${page}`)
            .then(fidListRes => {
                let fidList = fidListRes.data["data"]["fidList"]
                console.log(fidListRes)
                
                let queryStr = {{ .serverUrl }} + "/api/v1/static/thumbnail/get?"
                for (let i = 0; i < fidList.length; i++) {
                    if (i != 0) {
                        queryStr += "&"
                    }
                    queryStr += `fid=${fidList[i]}`
                }

                axios.get(queryStr)
                .then(thumbnailRes => {
                    console.log(thumbnailRes)

                    let list1 = document.getElementsByClassName("list1")[0]
                    let list2 = document.getElementsByClassName("list2")[0]
                    list1.innerHTML = "<li><span>请选择一张图片</span></li>"
                    list2.innerHTML = "";
                    for (let i = 0; i < thumbnailRes.data["data"]["thumbnail"].length; i++) {
                        list2.innerHTML += `<li><img src="${thumbnailRes.data["data"]["thumbnail"][i]}" alt="${fidList[i]}" id="img"></li>`
                    }
                    let imgs = document.getElementsByTagName("img")
                    for (var i = 0; i < imgs.length; i++) {
                        imgs[i].onclick = function () {
                            var name = this.alt
                            var list1 = document.getElementsByClassName("list1")[0]
                            list1.innerHTML = ""
                            axios.get({{ .serverUrl }} + `/api/v1/static/original/get?fid=${name}`).
                            then(originalRes => {
                                list1.innerHTML = `<li><img src="${originalRes.data["data"]["original"][0]}"></li>`
                            })
                        }
                    }
                })
                return fidListRes.data["data"]["total"]
            })
        }

        let add = document.getElementById("add");
        let sub = document.getElementById("sub");
        add.onclick = function () {
            page++
            if (page > maxPage) {
                page = maxPage
            }
            maxPage = load(page)
        }
        sub.onclick = function () {
            page--
            if (page < 1) {
                page = 1
            }
            maxPage = load(page)
        }

        maxPage = load(page);
    </script>
</body>

</html>