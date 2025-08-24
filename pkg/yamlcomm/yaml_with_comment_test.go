package yamlcomm

import (
	"fmt"
	"github.io/uberate/hcli/pkg/config"
	"testing"
)

func TestMarshalWithComments(t *testing.T) {
	c := config.CliConfig{
		Templates: []config.TemplateConfig{
			config.TemplateConfig{
				Name:       "test name",
				Categories: []string{"test category"},
				Template:   "test template",
			},
		},
	}

	v, err := MarshalWithComments(c)
	fmt.Println("=== CliConfig Test ===")
	fmt.Println(string(v), err)

	// 测试简单的字符串切片
	strSlice := []string{"item1", "item2", "item3"}
	v, err = MarshalWithComments(strSlice)
	fmt.Println("\n=== String Slice Test ===")
	fmt.Println(string(v), err)

	// 测试带注释的切片结构体
	type SliceStruct struct {
		Items []string `yaml:"items" comment:"This is a list of items.\nEach item is a string."`
	}

	sliceStruct := SliceStruct{
		Items: []string{"item1", "item2"},
	}
	v, err = MarshalWithComments(sliceStruct)
	fmt.Println("\n=== Slice Struct Test ===")
	fmt.Println(string(v), err)

	// 测试长注释自动换行
	type LongCommentStruct struct {
		LongList []string `yaml:"longList" comment:"This is a very long comment that should be automatically wrapped to multiple lines when it exceeds the maximum line length limit of 100 characters. This is a very long comment that should be automatically wrapped to multiple lines when it exceeds the maximum line length limit of 100 characters."`
	}

	longCommentStruct := LongCommentStruct{
		LongList: []string{"item1", "item2"},
	}
	v, err = MarshalWithComments(longCommentStruct)
	fmt.Println("\n=== Long Comment Struct Test ===")
	fmt.Println(string(v), err)
}
