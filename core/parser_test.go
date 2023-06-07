package core_test

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/88250/lute"
	"github.com/lbin/feishu_wiki_backup/core"
	"github.com/lbin/feishu_wiki_backup/utils"

	"github.com/chyroc/lark"
	"github.com/stretchr/testify/assert"
)

func TestParseDocContent(t *testing.T) {
	root := utils.RootDir()
	engine := lute.New(func(l *lute.Lute) {
		l.RenderOptions.AutoSpace = true
	})

	testdata := []string{
		"testdocs.1",
	}
	for _, td := range testdata {
		t.Run(td, func(t *testing.T) {
			jsonFile, err := os.Open(path.Join(root, "testdata", td+".json"))
			utils.CheckErr(err)
			defer jsonFile.Close()

			var docs lark.DocContent
			byteValue, _ := ioutil.ReadAll(jsonFile)
			json.Unmarshal(byteValue, &docs)

			parser := core.NewParser(context.Background())
			mdParsed := parser.ParseDocContent(&docs)
			mdParsed = engine.FormatStr("md", mdParsed)

			mdFile, err := ioutil.ReadFile(path.Join(root, "testdata", td+".md"))
			utils.CheckErr(err)
			mdExpected := string(mdFile)

			assert.Equal(t, mdExpected, mdParsed)
		})
	}
}

func TestParseDocxContent(t *testing.T) {
	root := utils.RootDir()
	engine := lute.New(func(l *lute.Lute) {
		l.RenderOptions.AutoSpace = true
	})

	testdata := []string{
		"testdocx.1",
		"testdocx.2",
	}
	for _, td := range testdata {
		t.Run(td, func(t *testing.T) {
			jsonFile, err := os.Open(path.Join(root, "testdata", td+".json"))
			utils.CheckErr(err)
			defer jsonFile.Close()

			data := struct {
				Document *lark.DocxDocument `json:"document"`
				Blocks   []*lark.DocxBlock  `json:"blocks"`
			}{}
			byteValue, _ := ioutil.ReadAll(jsonFile)
			json.Unmarshal(byteValue, &data)

			parser := core.NewParser(context.Background())
			mdParsed := parser.ParseDocxContent(data.Document, data.Blocks)
			mdParsed = engine.FormatStr("md", mdParsed)

			mdFile, err := ioutil.ReadFile(path.Join(root, "testdata", td+".md"))
			utils.CheckErr(err)
			mdExpected := string(mdFile)

			assert.Equal(t, mdExpected, mdParsed)
		})
	}
}
