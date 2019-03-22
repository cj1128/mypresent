# Mypresent

在 `golang/tools/present` 基础上修改而来，增加了以下功能：

1. 语法高亮
2. 导出为一个独立的文件夹

- 移除 playground 功能
- 移除 article 模板，只渲染 slide 模板

## Format

### slide format

title
[subtitle]
[time](format: "15:04 2 Jan 2006" or "2 Jan 2006")
<blank>
[misc info]
[sections]

`*` title
`: ` speaker note
`#` comments

## Static Resource

- action.tmpl
- index.tmpl
- slide.tmpl

- dir.css
- dir.js
- favicon.ico
- jquery-ui.js
- notes.css
- slides.js
- styles.css
