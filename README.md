# Mypresent

在 `golang/tools/present` 基础上做了如下修改：

- 增加语法高亮
- 增加导出为独立文件夹功能
- 移除 playground 功能
- 移除 article 模板，只保留 slide 模板

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

- index.css

tmpl
  - action.tmpl
  - index.tmpl
  - slide.tmpl

hljs
  - hljs.js
  - hljs.css
