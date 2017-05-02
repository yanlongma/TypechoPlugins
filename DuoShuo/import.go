package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

// 多说文章数据struct
type Threads struct {
	// 和json字段格式不一样 就要指定
	ThreadId  int    `json:"thread_id"`  // 多说文章ID
	ThreadKey string `json:"thread_key"` // typecho文章ID
}

// 多说评论数据struct
type Posts struct {
	PostId      int    `json:"post_id"`    // 多说评论ID
	ThreadId    int    `json:"thread_id"`  // 多说文章ID
	CreatedAt   string `json:"created_at"` // 创建日期
	AuthorName  string `json:"author_name"`
	AuthorEmail string `json:"author_email"`
	AuthorUrl   string `json:"author_url"`
	Ip          string
	Message     string
	Parents     []int
}

// 多说最终返回数据struct
type Response struct {
	Generator string    //`json:"generator"`
	Version   string    //`json:"version"`
	Threads   []Threads //`json:"threads"`
	Posts     []Posts
}

func main() {

	resp, err := readFile("export.json")
	if err != nil {
		fmt.Println("readFile:", err.Error())
		return
	}

	//var threadIdRelationCid map[int]string
	threadIdRelationCid := make(map[int]string)
	for _, item := range resp.Threads {
		threadIdRelationCid[item.ThreadId] = item.ThreadKey
	}

	coid := 1001
	postIdRelationCoid := make(map[int]int)
	for _, item := range resp.Posts {
		postIdRelationCoid[item.PostId] = coid
		coid += 1
	}

	sql := ""
	for _, item := range resp.Posts {
		tcoid := postIdRelationCoid[item.PostId]
		cid := threadIdRelationCid[item.ThreadId]
		t, _ := time.Parse("2006-01-02T15:04:05Z07:00", item.CreatedAt)
		created := t.Unix()
		author := item.AuthorName
		mail := item.AuthorEmail
		url := item.AuthorUrl
		ip := item.Ip
		text := item.Message
		parent := 0
		if len(item.Parents) > 0 {
			parent = postIdRelationCoid[item.Parents[0]]
		}

		sql += "INSERT INTO `typecho_comments` (`coid`, `cid`, `created`, `author`, `authorId`, `ownerId`, `mail`, `url`, `ip`, `agent`, `text`, `type`, `status`, `parent`) VALUES"
		sql += fmt.Sprintf("(%d, %s, %d, '%s', 0, 1, '%s', '%s', '%s', '', '%s', 'comment', 'approved', %d);\n", tcoid, cid, created, author, mail, url, ip, text, parent)
	}

	//fmt.Println(sql)
	writeFile("insert.sql", sql)
	fmt.Println("end")
}

// 读取文件的内容并返回
func readFile(filename string) (resp Response, err error) {

	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return
	}

	if err := json.Unmarshal(bytes, &resp); err != nil {
		fmt.Println(err)
	}

	return resp, err
}

// 向文件写入内容
func writeFile(filename string, text string) {

	fout, err := os.Create(filename)
	if err != nil {
		fmt.Println(filename, err)
		return
	}
	defer fout.Close()

	fout.WriteString(text)
}

// func read() {
// 	filename := "export.json"
// 	fout, err := os.Open(filename)
// 	if err != nil {
// 		fmt.Println(filename, err)
// 		return
// 	}
// 	defer fout.Close()

// 	var jsonStr string
// 	n, _ := fout.Read(jsonStr)

// 	fmt.Println(n, jsonStr)
// }
