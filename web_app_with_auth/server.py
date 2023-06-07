#!/usr/bin/env python
# -*- coding: UTF-8 -*-
import os
import logging

from auth import Auth
from dotenv import load_dotenv, find_dotenv
from flask import Flask, render_template, session, jsonify, request

# 日志格式设置
LOG_FORMAT = "%(asctime)s - %(message)s"
DATE_FORMAT = "%m/%d/%Y %H:%M:%S"
logging.basicConfig(level=logging.DEBUG, format=LOG_FORMAT, datefmt=DATE_FORMAT)

# const
# 在session中存储用户信息 user info 所需要的对应 session key
USER_INFO_KEY = "UserInfo"
# secret_key 是使用 flask session 所必须有的
SECRET_KEY = "ThisIsSecretKey"

# 从 .env 文件加载环境变量参数
load_dotenv(find_dotenv())

# 初始化 flask 网页应用
app = Flask(__name__, static_url_path="/public", static_folder="./public")
app.secret_key = SECRET_KEY
app.debug = True

# 获取环境变量值
APP_ID = os.getenv("APP_ID")
APP_SECRET = os.getenv("APP_SECRET")
FEISHU_HOST = os.getenv("FEISHU_HOST")

# 用获取的环境变量初始化免登流程类Auth
auth = Auth(FEISHU_HOST, APP_ID, APP_SECRET)

# 业务逻辑类
class Biz(object):
    @staticmethod
    def home_handler():
        # 主页加载流程
        return Biz._show_user_info()

    @staticmethod
    def login_handler():
        # 需要走免登流程
        return render_template("index.html", user_info={"name": "unknown"}, login_info="needLogin")

    @staticmethod
    def login_failed_handler(err_info):
        # 出错后的页面加载流程
        return Biz._show_err_info(err_info)

    # Session in Flask has a concept very similar to that of a cookie, 
    # i.e. data containing identifier to recognize the computer on the network, 
    # except the fact that session data is stored in a server.
    @staticmethod
    def _show_user_info():
        # 直接展示session中存储的用户信息
        return render_template("index.html", user_info=session[USER_INFO_KEY], login_info="alreadyLogin")

    @staticmethod
    def _show_err_info(err_info):
        # 将错误信息展示在页面上
        return render_template("err_info.html", err_info=err_info)

# 出错时走错误页面加载流程Biz.login_failed_handler(err_info)
@app.errorhandler(Exception)
def auth_error_handler(ex):
    return Biz.login_failed_handler(ex)


# 默认的主页路径
@app.route("/", methods=["GET"])
def get_home():
    # 打开本网页应用会执行的第一个函数

    # 如果session当中没有存储user info，则走免登业务流程Biz.login_handler()
    if USER_INFO_KEY not in session:
        logging.info("need to get user information")
        return Biz.login_handler()
    else:
        # 如果session中已经有user info，则直接走主页加载流程Biz.home_handler()
        logging.info("already have user information")
        return Biz.home_handler()

@app.route("/callback", methods=["GET"])
def callback():
    # 获取 user info

    # 拿到前端传来的临时授权码 Code
    code = request.args.get("code")
    # 先获取 user_access_token
    auth.authorize_user_access_token(code)
    # 再获取 user info
    user_info = auth.get_user_info()
    # 将 user info 存入 session
    session[USER_INFO_KEY] = user_info
    return jsonify(user_info)

@app.route("/get_appid", methods=["GET"])
def get_appid():
    # 获取 appid
    # 为了安全，app_id不应对外泄露，尤其不应在前端明文书写，因此此处从服务端传递过去
    return jsonify(
        {
            "appid": APP_ID
        }
    )


if __name__ == "__main__":
    # 以debug模式运行本网页应用
    # debug模式能检测服务端模块的代码变化，如果有修改会自动重启服务
    app.run(host="0.0.0.0", port=3000, debug=True)
    
