package core_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/lbin/feishu_wiki_backup/core"
	"github.com/lbin/feishu_wiki_backup/utils"
)

func getIdAndSecretFromEnv() (string, string) {
	utils.LoadEnv()
	appID := os.Getenv("FEISHU_APP_ID")
	appSecret := os.Getenv("FEISHU_APP_SECRET")
	return appID, appSecret
}

func TestNewClient(t *testing.T) {
	appID, appSecret := getIdAndSecretFromEnv()
	c := core.NewClient(appID, appSecret)
	if c == nil {
		t.Errorf("Error creating DocClient")
	}
}

func TestGetDocContent(t *testing.T) {
	appID, appSecret := getIdAndSecretFromEnv()
	c := core.NewClient(appID, appSecret)
	content, err := c.GetDocContent(
		context.WithValue(context.Background(), "Verbose", false),
		"doccnzSnjwt7Bd01MMcnNWpwV4d",
	)
	if err != nil {
		t.Error(err)
	}
	title := content.Title.Elements[0].TextRun.Text
	fmt.Println(title)
	if title == "" {
		t.Errorf("Error: parsed title is empty")
	}
}

func TestDownloadImage(t *testing.T) {
	appID, appSecret := getIdAndSecretFromEnv()
	c := core.NewClient(appID, appSecret)
	imgToken := "boxcnA1QKPanfMhLxzF1eMhoArM"
	filename, err := c.DownloadImage(
		context.WithValue(
			context.Background(),
			"OutputConfig", core.OutputConfig{
				ImageDir:        "static",
				TitleAsFilename: true,
			},
		),
		imgToken,
	)
	if err != nil {
		t.Error(err)
	}
	if filename != "static/"+imgToken+".png" {
		fmt.Println(filename)
		t.Errorf("Error: not expected file extension")
	}
}

func TestGetDocxContent(t *testing.T) {
	appID, appSecret := getIdAndSecretFromEnv()
	c := core.NewClient(appID, appSecret)
	docx, blocks, err := c.GetDocxContent(
		context.WithValue(context.Background(), "Verbose", false),
		"doxcnXhd93zqoLnmVPGIPTy7AFe",
	)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(docx.Title)
	if docx.Title == "" {
		t.Errorf("Error: parsed title is empty")
	}
	fmt.Printf("number of blocks: %d\n", len(blocks))
	if len(blocks) == 0 {
		t.Errorf("Error: parsed blocks are empty")
	}
}

func TestGetWikiNodeInfo(t *testing.T) {
	appID, appSecret := getIdAndSecretFromEnv()
	c := core.NewClient(appID, appSecret)
	const token = "wikcnLgRX9AMtvaB5x1cl57Yuah"
	node, err := c.GetWikiNodeInfo(context.Background(), token)
	if err != nil {
		t.Error(err)
	}
	if node.ObjType != "docx" {
		t.Errorf("Error: node type incorrect")
	}
}