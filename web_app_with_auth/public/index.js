let lang = window.navigator.language;
console.log(lang);


function apiAuth() {
    if (!window.h5sdk) {
        console.log('invalid h5sdk')
        alert('please open in feishu')
        return
    }

    // 通过服务端的Route: get_appid获取app_id
    // 服务端Route: get_appid的具体内容请参阅服务端模块server.py的get_appid()函数
    // 为了安全，app_id不应对外泄露，尤其不应在前端明文书写，因此此处从服务端获取
    fetch(`/get_appid`).then(response1 => response1.json().then(res1 => {
        console.log("get appid succeed: ", res1.appid);
        // 通过error接口处理API验证失败后的回调
        window.h5sdk.error(err => {
            throw('h5sdk error:', JSON.stringify(err));
        });
        // 通过ready接口确认环境准备就绪后才能调用API
        window.h5sdk.ready(() => {
            console.log("window.h5sdk.ready");
            console.log("url:", window.location.href);
            // 调用JSAPI tt.requestAuthCode 获取 authorization code
            tt.requestAuthCode({
                appId: res1.appid,
                // 获取成功后的回调
                success(res) {
                    console.log("getAuthCode succeed");
                    //authorization code 存储在 res.code
                    // 此处通过fetch把code传递给接入方服务端Route: callback，并获得user_info
                    // 服务端Route: callback的具体内容请参阅服务端模块server.py的callback()函数
                    fetch(`/callback?code=${res.code}`).then(response2 => response2.json().then(res2 => {
                        console.log("getUserInfo succeed");
                        // 示例Demo中单独定义的函数showUser，用于将用户信息展示在前端页面上
                        showUser(res2);}
                        )
                    ).catch(function (e) {console.error(e)})
                },
                // 获取失败后的回调
                fail(err) {
                    console.log(`getAuthCode failed, err:`, JSON.stringify(err));
                }
            })
        }
        )
    })).catch(function (e) { // 从服务端获取app_id失败后的处理
        console.error(e)
        })
}

function showUser(res) {
    // 展示用户信息
    // 头像
    $('#img_div').html(`<img src="${res.avatar_url}" width="100%" height=""100%/>`);
    // 名称
    $('#hello_text_name').text(lang === "zh_CN" || lang === "zh-CN" ? `${res.name}` : `${res.en_name}`);
    // 欢迎语
    $('#hello_text_welcome').text(lang === "zh_CN" || lang === "zh-CN" ? "欢迎使用飞书" : "welcome to Feishu");
}
