// Package document 模板功能测试
package document

import (
	"testing"
)

// TestNewTemplateEngine 测试创建模板引擎
func TestNewTemplateEngine(t *testing.T) {
	engine := NewTemplateEngine()
	if engine == nil {
		t.Fatal("Expected template engine to be created")
	}

	if engine.cache == nil {
		t.Fatal("Expected cache to be initialized")
	}
}

// TestTemplateVariableReplacement 测试变量替换功能
func TestTemplateVariableReplacement(t *testing.T) {
	engine := NewTemplateEngine()

	// 创建包含变量的模板
	templateContent := "Hello {{name}}, welcome to {{company}}!"
	template, err := engine.LoadTemplate("test_template", templateContent)
	if err != nil {
		t.Fatalf("Failed to load template: %v", err)
	}

	// 验证模板变量解析
	if len(template.Variables) != 2 {
		t.Errorf("Expected 2 variables, got %d", len(template.Variables))
	}

	if _, exists := template.Variables["name"]; !exists {
		t.Error("Expected 'name' variable to be found")
	}

	if _, exists := template.Variables["company"]; !exists {
		t.Error("Expected 'company' variable to be found")
	}

	// 创建模板数据
	data := NewTemplateData()
	data.SetVariable("name", "张三")
	data.SetVariable("company", "WordZero公司")

	// 渲染模板
	doc, err := engine.RenderToDocument("test_template", data)
	if err != nil {
		t.Fatalf("Failed to render template: %v", err)
	}

	if doc == nil {
		t.Fatal("Expected document to be created")
	}

	// 检查文档内容
	if len(doc.Body.Elements) == 0 {
		t.Error("Expected document to have content")
	}
}

// TestTemplateConditionalStatements 测试条件语句功能
func TestTemplateConditionalStatements(t *testing.T) {
	engine := NewTemplateEngine()

	// 创建包含条件语句的模板
	templateContent := `{{#if showWelcome}}欢迎使用WordZero！{{/if}}
{{#if showDescription}}这是一个强大的Word文档操作库。{{/if}}`

	template, err := engine.LoadTemplate("conditional_template", templateContent)
	if err != nil {
		t.Fatalf("Failed to load template: %v", err)
	}

	// 验证条件块解析
	if len(template.Blocks) < 2 {
		t.Errorf("Expected at least 2 blocks, got %d", len(template.Blocks))
	}

	// 测试条件为真的情况
	data := NewTemplateData()
	data.SetCondition("showWelcome", true)
	data.SetCondition("showDescription", false)

	doc, err := engine.RenderToDocument("conditional_template", data)
	if err != nil {
		t.Fatalf("Failed to render template: %v", err)
	}

	if doc == nil {
		t.Fatal("Expected document to be created")
	}
}

// TestTemplateLoopStatements 测试循环语句功能
func TestTemplateLoopStatements(t *testing.T) {
	engine := NewTemplateEngine()

	// 创建包含循环语句的模板
	templateContent := `产品列表：
{{#each products}}
- 产品名称：{{name}}，价格：{{price}}元
{{/each}}`

	template, err := engine.LoadTemplate("loop_template", templateContent)
	if err != nil {
		t.Fatalf("Failed to load template: %v", err)
	}

	// 验证循环块解析
	foundEachBlock := false
	for _, block := range template.Blocks {
		if block.Type == "each" && block.Variable == "products" {
			foundEachBlock = true
			break
		}
	}

	if !foundEachBlock {
		t.Error("Expected to find 'each products' block")
	}

	// 创建列表数据
	data := NewTemplateData()
	products := []interface{}{
		map[string]interface{}{
			"name":  "iPhone",
			"price": 8999,
		},
		map[string]interface{}{
			"name":  "iPad",
			"price": 5999,
		},
	}
	data.SetList("products", products)

	doc, err := engine.RenderToDocument("loop_template", data)
	if err != nil {
		t.Fatalf("Failed to render template: %v", err)
	}

	if doc == nil {
		t.Fatal("Expected document to be created")
	}
}

// TestTemplateInheritance 测试模板继承功能
func TestTemplateInheritance(t *testing.T) {
	engine := NewTemplateEngine()

	// 创建基础模板
	baseTemplateContent := `文档标题：{{title}}
基础内容：这是基础模板的内容。`

	_, err := engine.LoadTemplate("base_template", baseTemplateContent)
	if err != nil {
		t.Fatalf("Failed to load base template: %v", err)
	}

	// 创建继承模板
	childTemplateContent := `{{extends "base_template"}}
扩展内容：这是子模板的内容。`

	childTemplate, err := engine.LoadTemplate("child_template", childTemplateContent)
	if err != nil {
		t.Fatalf("Failed to load child template: %v", err)
	}

	// 验证继承关系
	if childTemplate.Parent == nil {
		t.Error("Expected child template to have parent")
	}

	if childTemplate.Parent.Name != "base_template" {
		t.Errorf("Expected parent template name to be 'base_template', got '%s'", childTemplate.Parent.Name)
	}
}

// TestTemplateValidation 测试模板验证功能
func TestTemplateValidation(t *testing.T) {
	engine := NewTemplateEngine()

	// 测试有效模板
	validTemplate := `Hello {{name}}!
{{#if showMessage}}This is a message.{{/if}}
{{#each items}}Item: {{this}}{{/each}}`

	template, err := engine.LoadTemplate("valid_template", validTemplate)
	if err != nil {
		t.Fatalf("Failed to load valid template: %v", err)
	}

	err = engine.ValidateTemplate(template)
	if err != nil {
		t.Errorf("Expected valid template to pass validation, got error: %v", err)
	}

	// 测试无效模板 - 括号不匹配
	invalidTemplate1 := `Hello {{name}!`
	template1, err := engine.LoadTemplate("invalid_template1", invalidTemplate1)
	if err != nil {
		t.Fatalf("Failed to load invalid template: %v", err)
	}

	err = engine.ValidateTemplate(template1)
	if err == nil {
		t.Error("Expected invalid template (mismatched brackets) to fail validation")
	}

	// 测试无效模板 - if语句不匹配
	invalidTemplate2 := `{{#if condition}}Hello`
	template2, err := engine.LoadTemplate("invalid_template2", invalidTemplate2)
	if err != nil {
		t.Fatalf("Failed to load invalid template: %v", err)
	}

	err = engine.ValidateTemplate(template2)
	if err == nil {
		t.Error("Expected invalid template (mismatched if statements) to fail validation")
	}
}

// TestTemplateData 测试模板数据功能
func TestTemplateData(t *testing.T) {
	data := NewTemplateData()

	// 测试设置和获取变量
	data.SetVariable("name", "测试")
	value, exists := data.GetVariable("name")
	if !exists {
		t.Error("Expected variable 'name' to exist")
	}
	if value != "测试" {
		t.Errorf("Expected variable value to be '测试', got '%v'", value)
	}

	// 测试设置和获取列表
	items := []interface{}{"item1", "item2", "item3"}
	data.SetList("items", items)
	list, exists := data.GetList("items")
	if !exists {
		t.Error("Expected list 'items' to exist")
	}
	if len(list) != 3 {
		t.Errorf("Expected list length to be 3, got %d", len(list))
	}

	// 测试设置和获取条件
	data.SetCondition("enabled", true)
	condition, exists := data.GetCondition("enabled")
	if !exists {
		t.Error("Expected condition 'enabled' to exist")
	}
	if !condition {
		t.Error("Expected condition value to be true")
	}

	// 测试批量设置变量
	variables := map[string]interface{}{
		"title":   "测试标题",
		"content": "测试内容",
	}
	data.SetVariables(variables)

	title, exists := data.GetVariable("title")
	if !exists || title != "测试标题" {
		t.Error("Expected batch set variables to work")
	}
}

// TestTemplateDataFromStruct 测试从结构体创建模板数据
func TestTemplateDataFromStruct(t *testing.T) {
	type TestStruct struct {
		Name    string
		Age     int
		Enabled bool
	}

	testData := TestStruct{
		Name:    "张三",
		Age:     30,
		Enabled: true,
	}

	templateData := NewTemplateData()
	err := templateData.FromStruct(testData)
	if err != nil {
		t.Fatalf("Failed to create template data from struct: %v", err)
	}

	// 验证变量是否正确设置
	name, exists := templateData.GetVariable("name")
	if !exists || name != "张三" {
		t.Error("Expected 'name' variable to be set correctly")
	}

	age, exists := templateData.GetVariable("age")
	if !exists || age != 30 {
		t.Error("Expected 'age' variable to be set correctly")
	}

	enabled, exists := templateData.GetVariable("enabled")
	if !exists || enabled != true {
		t.Error("Expected 'enabled' variable to be set correctly")
	}
}

// TestTemplateMerge 测试模板数据合并
func TestTemplateMerge(t *testing.T) {
	data1 := NewTemplateData()
	data1.SetVariable("name", "张三")
	data1.SetCondition("enabled", true)

	data2 := NewTemplateData()
	data2.SetVariable("age", 30)
	data2.SetList("items", []interface{}{"item1", "item2"})

	// 合并数据
	data1.Merge(data2)

	// 验证合并结果
	name, exists := data1.GetVariable("name")
	if !exists || name != "张三" {
		t.Error("Expected original variable to remain")
	}

	age, exists := data1.GetVariable("age")
	if !exists || age != 30 {
		t.Error("Expected merged variable to be present")
	}

	enabled, exists := data1.GetCondition("enabled")
	if !exists || !enabled {
		t.Error("Expected original condition to remain")
	}

	items, exists := data1.GetList("items")
	if !exists || len(items) != 2 {
		t.Error("Expected merged list to be present")
	}
}

// TestTemplateCache 测试模板缓存功能
func TestTemplateCache(t *testing.T) {
	engine := NewTemplateEngine()

	// 加载模板
	templateContent := "Hello {{name}}!"
	template1, err := engine.LoadTemplate("cached_template", templateContent)
	if err != nil {
		t.Fatalf("Failed to load template: %v", err)
	}

	// 从缓存获取模板
	template2, err := engine.GetTemplate("cached_template")
	if err != nil {
		t.Fatalf("Failed to get template from cache: %v", err)
	}

	// 验证是同一个模板实例
	if template1 != template2 {
		t.Error("Expected to get same template instance from cache")
	}

	// 清空缓存
	engine.ClearCache()

	// 尝试获取已清空的模板
	_, err = engine.GetTemplate("cached_template")
	if err == nil {
		t.Error("Expected error when getting template after cache clear")
	}
}

// TestComplexTemplateRendering 测试复杂模板渲染
func TestComplexTemplateRendering(t *testing.T) {
	engine := NewTemplateEngine()

	// 创建复杂模板
	complexTemplate := `报告标题：{{title}}
作者：{{author}}

{{#if showSummary}}
概要：{{summary}}
{{/if}}

详细内容：
{{#each sections}}
章节 {{@index}}: {{title}}
内容：{{content}}

{{/each}}

{{#if showFooter}}
报告完毕。
{{/if}}`

	_, err := engine.LoadTemplate("complex_template", complexTemplate)
	if err != nil {
		t.Fatalf("Failed to load complex template: %v", err)
	}

	// 创建复杂数据
	data := NewTemplateData()
	data.SetVariable("title", "WordZero功能测试报告")
	data.SetVariable("author", "开发团队")
	data.SetVariable("summary", "本报告测试了WordZero的模板功能。")

	data.SetCondition("showSummary", true)
	data.SetCondition("showFooter", true)

	sections := []interface{}{
		map[string]interface{}{
			"title":   "基础功能",
			"content": "测试了基础的文档操作功能。",
		},
		map[string]interface{}{
			"title":   "模板功能",
			"content": "测试了模板引擎的各种功能。",
		},
	}
	data.SetList("sections", sections)

	// 渲染复杂模板
	doc, err := engine.RenderToDocument("complex_template", data)
	if err != nil {
		t.Fatalf("Failed to render complex template: %v", err)
	}

	if doc == nil {
		t.Fatal("Expected document to be created")
	}

	// 验证文档有内容
	if len(doc.Body.Elements) == 0 {
		t.Error("Expected document to have content")
	}
}

// TestImagePlaceholder 测试图片占位符功能
func TestImagePlaceholder(t *testing.T) {
	engine := NewTemplateEngine()

	// 测试基础图片占位符解析
	t.Run("解析图片占位符", func(t *testing.T) {
		templateContent := `文档标题

这里有一个图片：
{{#image testImage}}

更多内容...`

		template, err := engine.LoadTemplate("image_test", templateContent)
		if err != nil {
			t.Fatalf("加载模板失败: %v", err)
		}

		// 检查是否正确解析了图片占位符
		hasImageBlock := false
		for _, block := range template.Blocks {
			if block.Type == "image" && block.Name == "testImage" {
				hasImageBlock = true
				break
			}
		}

		if !hasImageBlock {
			t.Error("模板应该包含图片块")
		}
	})

	// 测试图片占位符渲染（字符串模板）
	t.Run("渲染图片占位符到字符串", func(t *testing.T) {
		templateContent := `产品介绍：{{productName}}

产品图片：
{{#image productImage}}

描述：{{description}}`

		_, err := engine.LoadTemplate("product", templateContent)
		if err != nil {
			t.Fatalf("加载模板失败: %v", err)
		}

		data := NewTemplateData()
		data.SetVariable("productName", "测试产品")
		data.SetVariable("description", "这是一个测试产品")

		// 创建图片配置
		imageConfig := &ImageConfig{
			Position:  ImagePositionInline,
			Alignment: AlignCenter,
			Size: &ImageSize{
				Width:           100,
				KeepAspectRatio: true,
			},
			AltText: "测试图片",
			Title:   "测试产品图片",
		}

		// 设置图片数据（使用示例二进制数据）
		imageData := createTestImageData()
		data.SetImageFromData("productImage", imageData, imageConfig)

		// 渲染模板
		doc, err := engine.RenderToDocument("product", data)
		if err != nil {
			t.Fatalf("渲染模板失败: %v", err)
		}

		if doc == nil {
			t.Error("渲染结果不应为空")
		}
	})

	// 测试从文档模板渲染图片占位符
	t.Run("从文档模板渲染图片占位符", func(t *testing.T) {
		// 创建基础文档
		baseDoc := New()
		baseDoc.AddParagraph("报告标题：{{title}}")
		baseDoc.AddParagraph("{{#image reportChart}}")
		baseDoc.AddParagraph("总结：{{summary}}")

		// 从文档创建模板
		template, err := engine.LoadTemplateFromDocument("report_template", baseDoc)
		if err != nil {
			t.Fatalf("从文档创建模板失败: %v", err)
		}

		if len(template.Variables) == 0 {
			t.Error("模板应该包含变量")
		}

		// 准备数据
		data := NewTemplateData()
		data.SetVariable("title", "月度报告")
		data.SetVariable("summary", "数据显示增长趋势良好")

		chartConfig := &ImageConfig{
			Position:  ImagePositionInline,
			Alignment: AlignCenter,
			Size: &ImageSize{
				Width: 120,
			},
		}

		imageData := createTestImageData()
		data.SetImageFromData("reportChart", imageData, chartConfig)

		// 使用RenderTemplateToDocument方法（推荐用于文档模板）
		doc, err := engine.RenderTemplateToDocument("report_template", data)
		if err != nil {
			t.Fatalf("渲染文档模板失败: %v", err)
		}

		if doc == nil {
			t.Fatal("渲染结果不应为空")
		}

		// 检查文档中是否有元素
		if len(doc.Body.Elements) == 0 {
			t.Error("文档应该包含元素")
		}
	})

	// 测试图片数据管理方法
	t.Run("测试图片数据管理", func(t *testing.T) {
		data := NewTemplateData()

		// 测试SetImage方法
		config := &ImageConfig{
			Position: ImagePositionInline,
			Size:     &ImageSize{Width: 50},
		}
		data.SetImage("test1", "path/to/image.jpg", config)

		// 测试SetImageFromData方法
		imageData := createTestImageData()
		data.SetImageFromData("test2", imageData, config)

		// 测试SetImageWithDetails方法
		data.SetImageWithDetails("test3", "path/to/image2.jpg", imageData, config, "alt text", "title")

		// 测试GetImage方法
		img1, exists1 := data.GetImage("test1")
		if !exists1 || img1.FilePath != "path/to/image.jpg" {
			t.Error("图片1数据不正确")
		}

		img2, exists2 := data.GetImage("test2")
		if !exists2 || len(img2.Data) == 0 {
			t.Error("图片2数据不正确")
		}

		img3, exists3 := data.GetImage("test3")
		if !exists3 || img3.AltText != "alt text" || img3.Title != "title" {
			t.Error("图片3数据不正确")
		}

		// 测试不存在的图片
		_, exists4 := data.GetImage("nonexistent")
		if exists4 {
			t.Error("不存在的图片不应该返回true")
		}
	})

	// 测试图片占位符与其他模板语法的兼容性
	t.Run("图片占位符与其他语法兼容性", func(t *testing.T) {
		templateContent := `{{#if showImage}}
图片标题：{{imageTitle}}
{{#image dynamicImage}}
{{/if}}

{{#each items}}
项目：{{name}}
{{#image itemImage}}
描述：{{description}}
{{/each}}`

		_, err := engine.LoadTemplate("complex", templateContent)
		if err != nil {
			t.Fatalf("加载复杂模板失败: %v", err)
		}

		data := NewTemplateData()
		data.SetCondition("showImage", true)
		data.SetVariable("imageTitle", "主要图片")

		items := []interface{}{
			map[string]interface{}{
				"name":        "项目1",
				"description": "项目1描述",
			},
		}
		data.SetList("items", items)

		config := &ImageConfig{Position: ImagePositionInline}
		imageData := createTestImageData()
		data.SetImageFromData("dynamicImage", imageData, config)
		data.SetImageFromData("itemImage", imageData, config)

		// 渲染不应该出错
		doc, err := engine.RenderToDocument("complex", data)
		if err != nil {
			t.Fatalf("渲染复杂模板失败: %v", err)
		}

		if doc == nil {
			t.Error("渲染结果不应为空")
		}
	})
}

// createTestImageData 创建测试用的图片数据
func createTestImageData() []byte {
	// 创建一个最小的PNG图片数据用于测试
	return []byte{
		0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A,
		0x00, 0x00, 0x00, 0x0D, 0x49, 0x48, 0x44, 0x52,
		0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01,
		0x08, 0x06, 0x00, 0x00, 0x00, 0x1F, 0x15, 0xC4,
		0x89, 0x00, 0x00, 0x00, 0x0A, 0x49, 0x44, 0x41,
		0x54, 0x78, 0x9C, 0x63, 0x00, 0x01, 0x00, 0x00,
		0x05, 0x00, 0x01, 0x0D, 0x0A, 0x2D, 0xB4, 0x00,
		0x00, 0x00, 0x00, 0x49, 0x45, 0x4E, 0x44, 0xAE,
		0x42, 0x60, 0x82,
	}
}
