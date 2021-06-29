> 本文主要介绍HTML的入门知识，以及常用的HTML标签与属性。

TODO: 看书重写

# 1. HTML简介

HTML是用来描述网页的一种语言，全称超文本标记语言(Hyper Text Markup Language)，它是一种标记语言，使用标记标签来描述网页，它的语法规则可以定义图片、表格、链接等。

## 1.1 标签

HTML标签是由尖括号包围、成对出现的关键词，如\<b\> 和 \</b\>，分别为开始标签和结束标签。HTML标签大小写不敏感，但推荐使用小写。

HTML文档包含HTML标签和文本，也被称为网页，HTML由web浏览器读取解析，并生成对应的网页形式。

一个基础的实例：

```html
<html>
<body>

<h1>我的第一个标题</h1>

<p>我的第一个段落。</p>

</body>
</html>
```

- \<html\> 与 \</html\> 之间的文本描述网页
- \<body\> 与 \</body\> 之间的文本是可见的页面内容
- \<h1\> 与 \</h1\> 之间的文本被显示为标题
- \<p\> 与 \</p\> 之间的文本被显示为段落

更多标签的详情可参考[HTML 标签参考手册](https://www.w3school.com.cn/tags/index.asp)。

## 1.2 元素

HTML元素是开始标签到对应的结束标签的所有代码，HTML元素可以嵌套包含其他HTML元素。

\<html\>元素定义整个HTML文档，\<body\>元素定义文档的主体。

空元素没有内容，在开始标签中关闭，如\<br\>，应当在开始标签添加斜杠：\<br /\>。

## 1.3 属性

HTML属性在开始标签中规定，以名称/值对name="value"的形式出现。属性和属性值大小写不敏感，但推荐小写。

以下是大多数HTML元素的属性：

1. class：规定元素的类名（classname）
2. id：规定元素的唯一 id
3. style：规定元素的行内样式（inline style）
4. title：规定元素的额外信息

以下是几个属性的实例：

1. HTML链接的地址：

   ```html
   <a href="http://www.w3school.com.cn">This is a link</a>
   ```

2. HTML标题的对齐方式：

   ```html
   <h1 align="center">
   ```

3. HTML表格的边框信息：

   ```html
   <table border="1">
   ```



# 2. 标签

## 2.1 基础标签

**注释**

```html
<!-- This is a comment -->
```

**标题**

通过 \<h1\> - \<h6\> 等标签进行定义的，\<h1\>最大，\<h6\>最小。

```html
<h1>This is a heading</h1>
<h2>This is a heading</h2>
<h3>This is a heading</h3>
```

**段落**

```html
<p>This is a paragraph.</p>
```

**换行**

显示页面时，浏览器会移除源代码多余的空格和空行，连续的空格或空行会被算作一个空格，因此需要用此标签来换行。

```html
<br />
```

**水平线**

创建水平线用于分割内容。

```html
<hr />
```

**链接**

href属性指定链接地址，标签中间为显示的文本。

target属性定义链接文档显示位置，`target="_blank"`表示在新窗口打开。

name属性或id属性定义锚（anchor），name="label"，在同一文档创建链接`<a href="#label">`，或者其他页面创建链接`<a href="url#label">`，点击后直接跳转到该锚。

```html
<a href="url">text</a>
```

**图像**

属性指定了图像的名称和尺寸。

alt属性定义替换的文本，无法载入图片时显示。

```
<img src="w3school.jpg" width="104" height="142" />
```

## 2.2 文本格式化

**粗体文本**

```html
<b>text</b>
```

**斜体字**

```html
<i>text</i>
```

**强调**

```html
<em>text</em>
```

**大号字**

```html
<big>text</big>
```

**小号字**

```html
<small>text</small>
```

**下标字**

```html
<sub>text</sub>
```

**上标字**

```html
<sup>text</sup>
```

**删除线**

```html
<del>text</del>
```

**下划线**

```html
<ins>text</ins>
```

## 2.3 引用

**短引用**

单行的引用，用引号包围。

```html
<q>text</q>
```

**长引用**

对整个元素进行缩进。

```html
<blockquote cite="http://www.worldwildlife.org/who/index.html">
五十年来，WWF 一直致力于保护自然界的未来。
WWF 工作于 100 个国家，并得到美国一百二十万会员及全球近五百万会员的支持。
</blockquote>
```

**缩略词**

对缩写标记，为浏览器、搜索引擎提供信息。

```html
<abbr title="World Health Organization">WHO</abbr>
```

**联系信息**

定义文档或文章的联系信息（作者/拥有者）。

```html
<address>
Written by Donald Duck.<br> 
Visit us at:<br>
Example.com<br>
Box 564, Disneyland<br>
USA
</address>
```

**著作标题**

```html
<cite>The Scream</cite>
```

## 2.4 表格

一个表格用\<table\>标签定义，每一行用\<tr\>标签定义，行的每个单元格用\<td\>标签定义。数据单元格可以包含文本、图片、列表、段落、表单、水平线、表格等等。第一行表头可用\<th\>标签定义，显示为粗体居中。

\<table\>标签的border属性定义边框。

一个表格的例子：

```html
<table border="1">
    <tr>
        <th>Heading</th>
        <th>Another Heading</th>
    </tr>
    <tr>
        <td>row 1, cell 1</td>
        <td>row 1, cell 2</td>
    </tr>
    <tr>
        <td>row 2, cell 1</td>
        <td>row 2, cell 2</td>
    </tr>
</table>
```

## 2.5 列表

不同的列表可以嵌套使用，即在\<li\>元素中定义新的列表。

**无序列表**

用圆点标记每个项目。

```html
<ul>
    <li>Coffee</li>
    <li>Milk</li>
</ul>
```

**有序列表**

用数字标记每个项目。

```html
<ol>
    <li>Coffee</li>
    <li>Milk</li>
</ol>
```

**自定义列表**

用给定文本标记每个项目。

```html
<dl>
    <dt>Coffee</dt>
    <dd>Black hot drink</dd>
    <dt>Milk</dt>
    <dd>White cold drink</dd>
</dl>
```



# 3. 属性

## 3.1 样式

可以通过HTML标签的style属性，将样式添加到HTML元素中，或者通过CSS进行定义。

**外部样式表**

当样式需要被应用到很多页面的时候，外部样式表将是理想的选择。定义公共的css文件来被多个HTML文件引用，在head部分引用css文件。

```html
<head>
<link rel="stylesheet" type="text/css" href="mystyle.css">
</head>
```

**内部样式表**

当单个文件需要特别样式时，就可以使用内部样式表。在head部分用\<style\>标签定义。

```html
<head>
<style type="text/css">
body {background-color: red}
p {margin-left: 20px}
</style>
</head>
```

**内联样式**

当特殊的样式需要应用到个别元素时，就可以使用内联样式。为对应的元素添加style属性。

一个样式的实例：通过style属性定义了各个层级元素的样式。

```html
<html>
<body style="background-color:yellow">
<h2 style="background-color:red">This is a heading</h2>
<p style="background-color:green">This is a paragraph.</p>
</body>
</html>
```

以下是几个样式的实例：

1. 通过`style="background-color:red"`淘汰bgcolor属性
2. 通过`style="font-family:arial;color:red;font-size:20px;"`淘汰\<font\>标签
3. 通过`style="text-align:center"`淘汰align属性

## 3.2 颜色

颜色定义为红绿蓝（RGB）的组合，每个颜色最小为0（十六进制：#00），最大为255（十六进制：#FF）。

#加上六个十六进制数表示一个RGB颜色，如 #FF0000表示rgb(255,0,0)也就是红色。

大多数浏览器支持的颜色名可以参考[HTML 颜色名](https://www.w3school.com.cn/html/html_colornames.asp)。



# 4. 参考

1. [W3School HTML 教程](https://www.w3school.com.cn/html/index.asp)

