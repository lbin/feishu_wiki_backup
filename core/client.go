package core

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/chyroc/lark"
	"github.com/lbin/feishu_wiki_backup/utils"
)

type Client struct {
	larkClient *lark.Lark
}

func NewClient(appID, appSecret string) *Client {
	return &Client{
		larkClient: lark.New(
			lark.WithAppCredential(appID, appSecret),
			lark.WithTimeout(60*time.Second),
		),
	}
}

func (c *Client) GetWikiSpaceList(ctx context.Context) (*lark.GetWikiSpaceResp, error) {
	// var pageSize int64 = 20
	// var pageToken = ""
	var space_id = "6992758255126577155"

	// Test
	data, _, err := c.larkClient.Drive.GetWikiSpace(ctx, &lark.GetWikiSpaceReq{SpaceID: space_id})
	// fmt.Printf("resp %s\n", resp.message)

	// data, _, err := c.larkClient.Drive.GetWikiSpaceList(ctx, &lark.GetWikiSpaceListReq{PageSize: &pageSize, PageToken: &pageToken})

	return data, err
}

func (c *Client) GetDocContent(ctx context.Context, docToken string) (*lark.DocContent, error) {
	resp, _, err := c.larkClient.Drive.GetDriveDocContent(ctx, &lark.GetDriveDocContentReq{
		DocToken: docToken,
	})
	if err != nil {
		return nil, err
	}
	doc := &lark.DocContent{}
	err = json.Unmarshal([]byte(resp.Content), doc)
	if err != nil {
		return doc, err
	}

	if ctx.Value("Verbose").(bool) {
		pdoc := utils.PrettyPrint(doc)
		fmt.Println(pdoc)
		if err = os.WriteFile(fmt.Sprintf("%s_verbose.json", docToken), []byte(pdoc), 0o644); err != nil {
			return nil, err
		}
	}

	return doc, nil
}

func (c *Client) DownloadImage(ctx context.Context, imgToken string) (string, error) {
	resp, _, err := c.larkClient.Drive.DownloadDriveMedia(ctx, &lark.DownloadDriveMediaReq{
		FileToken: imgToken,
	})
	if err != nil {
		return imgToken, err
	}
	imgDir := ctx.Value("OutputConfig").(OutputConfig).ImageDir
	fileext := filepath.Ext(resp.Filename)
	filename := fmt.Sprintf("%s/%s%s", imgDir, imgToken, fileext)
	err = os.MkdirAll(filepath.Dir(filename), 0o755)
	if err != nil {
		return imgToken, err
	}
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0o666)
	if err != nil {
		return imgToken, err
	}
	defer file.Close()
	_, err = io.Copy(file, resp.File)
	if err != nil {
		return imgToken, err
	}
	return filename, nil
}

func (c *Client) GetDocxContent(ctx context.Context, docToken string) (*lark.DocxDocument, []*lark.DocxBlock, error) {
	resp, _, err := c.larkClient.Drive.GetDocxDocument(ctx, &lark.GetDocxDocumentReq{
		DocumentID: docToken,
	})
	if err != nil {
		return nil, nil, err
	}
	docx := &lark.DocxDocument{
		DocumentID: resp.Document.DocumentID,
		RevisionID: resp.Document.RevisionID,
		Title:      resp.Document.Title,
	}
	var blocks []*lark.DocxBlock
	var pageToken *string
	for {
		resp2, _, err := c.larkClient.Drive.GetDocxBlockListOfDocument(ctx, &lark.GetDocxBlockListOfDocumentReq{
			DocumentID: docx.DocumentID,
			PageToken:  pageToken,
		})
		if err != nil {
			return docx, nil, err
		}
		blocks = append(blocks, resp2.Items...)
		pageToken = &resp2.PageToken
		if !resp2.HasMore {
			break
		}
	}

	if ctx.Value("Verbose").(bool) {
		data := struct {
			Document *lark.DocxDocument `json:"document"`
			Blocks   []*lark.DocxBlock  `json:"blocks"`
		}{
			Document: docx,
			Blocks:   blocks,
		}
		pdata := utils.PrettyPrint(data)
		fmt.Println(pdata)
		if err = os.WriteFile(fmt.Sprintf("%s_verbose.json", docToken), []byte(pdata), 0o644); err != nil {
			return nil, nil, err
		}
	}

	return docx, blocks, nil
}

func (c *Client) GetWikiNodeInfo(ctx context.Context, token string) (*lark.GetWikiNodeRespNode, error) {
	resp, _, err := c.larkClient.Drive.GetWikiNode(ctx, &lark.GetWikiNodeReq{
		Token: token,
	})
	if err != nil {
		return nil, err
	}
	return resp.Node, nil
}
