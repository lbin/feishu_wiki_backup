# feishu_wiki_backup

## 获取 API Token

配置文件需要填写 APP ID 和 APP SECRET 信息，请参考 [飞书官方文档](https://open.feishu.cn/document/ukTMukTMukTM/ukDNz4SO0MjL5QzM/get-) 获取。推荐设置为

- 进入飞书[开发者后台](https://open.feishu.cn/app)
- 创建企业自建应用，信息随意填写
- 选择测试企业和人员，创建测试企业，绑定应用，切换至测试版本
- （重要）打开权限管理，云文档，开通所有只读权限
  - 「查看、评论和导出文档」权限 `docs:doc:readonly`
  - 「查看 DocX 文档」权限 `docx:document:readonly`
  - 「查看、评论和下载云空间中所有文件」权限 `drive:drive:readonly`
  - 「查看和下载云空间中的文件」权限 `drive:file:readonly`

- 打开凭证与基础信息，获取 App ID 和 App Secret

### 生成配置文件**

通过 `./bin/lark_backup config --appId <your_id> --appSecret <your_secret>` 命令即可生成该工具的配置文件。

通过 `./bin/lark_backup config` 命令可以查看配置文件路径以及是否成功配置。

更多的配置选项请手动打开配置文件更改。

## 感谢

- [Wsine/feishu2md](https://github.com/Wsine/feishu2md)
- [chyroc/lark](https://github.com/chyroc/lark)
- [chyroc/lark_docs_md](https://github.com/chyroc/lark_docs_md)
