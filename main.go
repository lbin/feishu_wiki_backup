package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	lark "github.com/larksuite/oapi-sdk-go/v3"
	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	larkwiki "github.com/larksuite/oapi-sdk-go/v3/service/wiki/v2"

	"github.com/lbin/feishu_wiki_backup/core"
	"github.com/lbin/feishu_wiki_backup/utils"
	"github.com/urfave/cli/v2"
)

var version = "v0.1.0"

func handleConfigCommand(appId, appSecret string) error {
	configPath, err := core.GetConfigFilePath()
	if err != nil {
		return err
	}
	fmt.Println("Configuration file on: " + configPath)
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		config := core.NewConfig(appId, appSecret)
		if err = config.WriteConfig2File(configPath); err != nil {
			return err
		}
		fmt.Println(utils.PrettyPrint(config))
	} else {
		config, err := core.ReadConfigFromFile(configPath)
		if err != nil {
			return err
		}
		if appId != "" {
			config.Feishu.AppId = appId
		}
		if appSecret != "" {
			config.Feishu.AppSecret = appSecret
		}
		if appId != "" || appSecret != "" {
			if err = config.WriteConfig2File(configPath); err != nil {
				return err
			}
		}
		fmt.Println(utils.PrettyPrint(config))
	}
	return nil
}

func handleBackupCommand() error {
	configPath, err := core.GetConfigFilePath()
	if err != nil {
		return err
	}
	config, err := core.ReadConfigFromFile(configPath)
	if err != nil {
		return err
	}

	client := lark.NewClient(config.Feishu.AppId, config.Feishu.AppSecret)
	// req := larkwiki.NewGetSpaceReqBuilder().
	// 	SpaceId("6992758255126577155").
	// 	Build()
	// // 发起请求
	// resp, err := client.Wiki.Space.Get(context.Background(), req)

	// req := larkwiki.NewCreateSpaceMemberReqBuilder().
	// 	SpaceId("6992758255126577155").
	// 	NeedNotification(true).
	// 	Member(larkwiki.NewMemberBuilder().
	// 		MemberType("openid").
	// 		MemberId("ou_92bbf97867d1fa85a591513eb9a87ff1").
	// 		MemberRole("admin").
	// 		Build()).
	// 	Build()
	// // 发起请求
	// resp, err := client.Wiki.SpaceMember.Create(context.Background(), req)

	req := larkwiki.NewGetNodeSpaceReqBuilder().
		Token("wikcn3eKUgmqiYc3GYaGentOSBh").
		Build()
	// 发起请求
	resp, err := client.Wiki.Space.GetNode(context.Background(), req)

	// 处理错误
	if err != nil {
		fmt.Println(err)
		return err
	}

	// 服务端错误处理
	if !resp.Success() {
		fmt.Println(resp.Code, resp.Msg, resp.RequestId())
		return err
	}

	// 业务处理
	fmt.Println(larkcore.Prettify(resp))

	return nil
}

func main() {
	app := &cli.App{
		Name:    "feishu_wiki_backup",
		Version: strings.TrimSpace(string(version)),
		Usage:   "download feishu/larksuite document to markdown file",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "verbose",
				Aliases: []string{"vv"},
				Usage:   "verbose the intermediate output",
			},
		},
		Action: func(ctx *cli.Context) error {
			if ctx.NArg() == 0 {
				return handleBackupCommand()
			} else {
				cli.ShowAppHelp(ctx)
			}
			return nil
		},
		Commands: []*cli.Command{
			{
				Name:  "config",
				Usage: "read config file or set field(s) if provided",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "appId",
						Value: "",
						Usage: "set app id for the OPEN api",
					},
					&cli.StringFlag{
						Name:  "appSecret",
						Value: "",
						Usage: "set app secret for the OPEN api",
					},
				},
				Action: func(ctx *cli.Context) error {
					return handleConfigCommand(
						ctx.String("appId"), ctx.String("appSecret"),
					)
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
