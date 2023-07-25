package gwg

import "testing"

func TestStruct(t *testing.T) {
	s := Struct{
		Name:         "Hello",
		Combinations: []string{"good", "baoo"},
		Properties: []Property{
			{
				Label: "Good",
				Type:  "*gin.Context",
				Tags: []Tag{
					{
						Label:   "json",
						Content: "id",
					},
					{
						Label:   "gorm",
						Content: "id,index",
					},
				},
			},
		},
	}
	t.Log(s)
	i := Import{
		Packages: []string{"github.com/gin/gorm", "github.com/gin/gorm2"},
	}
	t.Log(i)
	f := Func{
		Name: "HelloWorld",
		Parameters: Parameters{Pairs: []Pair{
			{
				Left:  "hello",
				Right: "int",
			},
			{
				Left:  "hello2",
				Right: "int",
			}}},
		Outputs: Outputs{
			Pairs: []Pair{
				{
					Right: "int",
				},
			},
		},
	}
	f.AddLine(Line{"return 1"})
	t.Log(f)
	p := Package{
		Name: "model",
	}
	p.AddImport(i)
	inter := Interface{
		Name: "AccountApiInterface",
	}
	inter.AddMethod(Method{
		Name:       "CreateAccount",
		Parameters: Parameters{[]Pair{{Left: "good", Right: "int"}}},
	})
	p.AddCode(f)
	p.AddCode(inter)
	p.Wirte("./test")
}
