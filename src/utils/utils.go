package utils

import (
	"html/template"
	"os"
	"strings"

	"github.com/Ayobami0/personal_blog/src/constants"
	"github.com/Ayobami0/personal_blog/src/model"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

func GetBlogs() []model.Blog {
	blogs := make([]model.Blog, 0)

	files, err := os.ReadDir(constants.BLOG_ARTICLES_DIR)

	if err != nil {
		return blogs
	}

	for _, v := range files {
		if !v.IsDir() {
			m := strings.TrimSuffix(v.Name(), ".md")

			info, err := v.Info()
			if err != nil {
				continue
			}
			blogs = append(blogs, model.Blog{Title: m, DateModified: info.ModTime().UTC().Format(constants.BLOG_ARTICLE_DATE_FORMAT)})
		}
	}

	return blogs
}

func ParseBlogAsHTML(name string) (*model.Blog, error) {

	f := constants.BLOG_ARTICLES_DIR + strings.Trim(name, " ") + ".md"

	mdInfo, err := os.Stat(f)
	if err != nil {
		return nil, err
	}

	md, err := os.ReadFile(f)
	if err != nil {
		return nil, err
	}

	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse(md)

	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	var blog model.Blog

	blog.Content = template.HTML(string(markdown.Render(doc, renderer)))
	blog.DateModified = mdInfo.ModTime().UTC().Format(constants.BLOG_ARTICLE_DATE_FORMAT)
	blog.Title = name

	return &blog, nil
}
