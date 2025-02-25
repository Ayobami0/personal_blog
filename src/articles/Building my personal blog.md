So, Earlier this week, i thought to my self '🤔 i should have a personal blog'. So i set out to build said blog.

### The problem statement
Most of my work is done in the terminal so i'd prefer to write the articles as a document. I choose markdown because it allow customization and i already use it for note taking [Obsidian](https://obsidian.md/). The notes are then rendered as html which can then be viewed on a browser. My initial workflow was this: Create article in markdown, push created article to github, create a github action that uses pandoc or some other html to markdown converter, store the blog files as pages in directory which will then be rendered in the client.
After a deep thought and a lot of searching, i deceided to just use a library(package) that allows me convert markdown to html then inject the converted markdown into a template instead of creating a html page for each blog article.

With that out of the way, my next task was identifying the technologies to use. My picks were python+flask, go and javascript+svelte (not react 🤮). I went with my favorite, golang. I love go because of it's simplicity and rich standard library (the net/http library was more than enough for the project). After picking my language of choice, i considered using tailwind css for styling as i am primarily a Backend developer so styling wasn't my strong suit. Decided to settle with plain css to learn and also so i would have better control over the styles (not that tailwind limits my control). Also, i would also need a library for converting markdown to html text for go. That was a relatively quick find, [gomarkdown was a google search away](https://github.com/gomarkdown/markdown).

Starting the project was relativly easy, though i needed to check out go's excelent documentation for the template/html package. My first road block was templating. In flask for instance, you use jinja for templating and extending and template in heritance is relatively eaiser to implement.

For Example, this is a base template using jinja syntax
```html
<!DOCTYPE html>
<html lang="en">
<head>
    {% block head %}
    <link rel="stylesheet" href="style.css" />
    <title>{% block title %}{% endblock %} - My Webpage</title>
    {% endblock %}
</head>
<body>
    <div id="content">{% block content %}{% endblock %}</div>
    <div id="footer">
        {% block footer %}
        &copy; Copyright 2008 by <a href="http://domain.invalid/">you</a>.
        {% endblock %}
    </div>
</body>
</html>
```
and this is the child template, it's job is to fill the empty blocks with content
```html
{% extends "base.html" %}
{% block title %}Index{% endblock %}
{% block head %}
    {{ super() }}
    <style type="text/css">
        .important { color: #336699; }
    </style>
{% endblock %}
{% block content %}
    <h1>Index</h1>
    <p class="important">
      Welcome to my awesome homepage.
    </p>
{% endblock %}
```
All you need to do in flask is just return the html file in your response, jinja would take care of the rest. In go, it is not that straight foreward.

### Go's template/text and template/html
In order to understand the problem, we must first understand the template packages. Tempalating in go is not all that different and it has a similar syntax for inheritance.  Lets go through this block of html together.
```html
{{define "base"}}
<!doctype html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    {{template "head" .}}
  </head>
  <body>
    <header></header>
    <section id="body">{{template "body" .}}</section>
    <footer>© Ayobami Oludemi. 2025</footer>
  </body>
</html>
{{end}}
```
First, we define a name for a template 'base', this is be what we would be calling when rendering the output. The next line that is not html syntax is the `{{templete "head" .}}` line. What this does is to call an already registered template name 'head' and pass any argument to it (hence the '.'). A similar thing occurs for the 'body' template. What the 'call' does is to replace that part with what is included in the template registration. This would be much clearer once we take a look at where 'head' and 'body' are defined. The `{{end}}` signifies the end of the template defination.
```html
{{define "head"}}
<link rel="stylesheet" type="text/css" href="/static/css/blog.css" />
{{end}}
{{define "body"}}
<div class="blog c-col">
  <div class="blog-title c-col">
    <h1>{{.Title}}</h1>
    <p class="outline-text">{{.DateModified}}</p>
  </div>
  <div class="blog-content">{{.Content}}</div>
</div>
{{end}}

```
We can see both the 'head' and 'body' templates are defined here. This would replace the template calls in the earlier html file. Don't mind the `{{.Title}}`, `{{.DateModified}}` and `{{.Content}}`, the represent struct fields from the variables passed to the template when it is executed.

You then might think "this is not too different from flask's jinja". Yes, they are similar, but where the difference is is how they are implemented in code. With python+flask, all you typically need to do is just return the correct html file as the response, in go there is an extra step.  
First, you need to create the template from the name then to get the output, you call a function to execute the template with the required data or nil. Stay with me, this example should clearify that.
```go
tmpl, err := template.New("name").Parse(...) // html will be supplied as a string
// Error checking elided
err = tmpl.Execute(out, data) // out is the stdout it can be a file, the terminal output or a network response
```
Hope the problem is becoming clearer. By default, the new function creates a template using a name and just one html string where both the defination and the call of the template is in. But that isn't what we want. Lucklily, the template/html package gives us a function just for that, the **[ParseFiles](https://cs.opensource.google/go/go/+/go1.24.0:src/html/template/template.go;l=382)** function. It takes a list of file names and parses the template definations from the named files. The template created would have the name of the first file. This pratically solves our problems and we now have access to the base template by just supplying the required file.
``` go
tmpl, err := template.ParseFiles("index.html", "base.html")
// ...
err = tmpl.ExecuteTemplate(out, "base", data)
```
Pheww, That was exicting. We solve our main problem and can execute the 'base' template without errors. But we quickly run into a minor hickup. Can you guess what?
```go
tmpl1, err := template.ParseFiles("index.html", "base.html")
tmpl2, err := template.ParseFiles("login.html", "base.html")
tmpl3, err := template.ParseFiles("home.html", "base.html")
tmpl4, err := template.ParseFiles("settings.html", "base.html")
```
We see that we would need a way to keep track and manage the multiple templates we create, also, writing `templates.ParseFiles` multiple times could by tiring. This problem is easy solveble by using a data structure to store the templates and access them when needed. Go has slices, maps and struct and all are possible fit for the structure but personally, i prefer using a struct, you'd see why in a moment.

```go
type Templ struct {
	templates map[string]*template.Template
}

func New() *Templ {
  return &Templ{
    templates: make(map[string]*template.Template),
  }
}

func (t *Templ) Register(name string, files ...string) {
  var f []string

  for _, v := range files {
    f = append(f, string(fmt.Sprintf("%s/%s", constants.TEMPLATES_DIR, v)))
  }

  t.templates[name] = template.Must(template.ParseFiles(f...))
}

func (t Templ) Get(name string) *template.Template {
  tmp, ok := t.templates[name]

  if !ok {
    return nil
  }

  return tmp
}
```
First we define a struct `Templ` to manage templates. Internally, the struct uses a private field templates which is a map with string keys, and template values. We then have a `Register` mehod that stores the templates into the underlying map. Retrival is done by calling the `Get` method with the name of the template. You can see why i prefer structs, it is cleaner and properly abstracts away the internal stucture.

With that resolved, using the templates was a breeze.
```go
// creation
t := templ.New()

// defination
t.Register("index", "index.html", "base.html")
t.Register("settings", "settings.html", "base.html")
t.Register("login", "login.html", "base.html")
t.Register("home", "home.html", "base.html")

// usage
template := t.Get("index")
// ....
```
Finally, the major hurdle is cleared, the next details were just the frontend ui with i did my best at. Also the names of the blog articles were created based on the name of the file.

### Conclusion & Lesson learnt
Building this blog was an exciting challenge for me. I gave me a deeper understanding of go's templating system and the standard library. I especially love that go requires a more hands-on approach to managing templates, it gave me more control over how i choose to structure the application. 

For future improvements, i'm considering:
- Implementing a caching system for the markdown files and their templates
- Use a better styling for the ui
- Add a way to include some additional metadata about the articles (like summary)
- Include local images to the blog

The site github is [My Blog](https://github.com/Ayobami0/personal_blog)

### Resources
- [Go's beautiful standard library](https://pkg.go.dev/std)
- [This question on Stackoverflow](https://stackoverflow.com/questions/11467731/is-it-possible-to-have-nested-templates-in-go-using-the-standard-library)
- [A great explaination on the template packake](https://developer.hashicorp.com/nomad/tutorials/templates/go-template-syntax)
