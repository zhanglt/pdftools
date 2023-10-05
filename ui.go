package main

import (
	"time"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"github.com/lxn/win"
)

// 初始始化窗体
func getWindows() error {
	icon, _ := walk.NewIconFromResourceId(3)
	err := MainWindow{
		//Background: SolidColorBrush{Color: walk.RGB(240, 240, 240)},
		MenuItems: []MenuItem{
			Menu{
				Text: "文件",
				Items: []MenuItem{
					Action{
						Text:        "退出",
						OnTriggered: func() { app.mw.Close() },
					},
				},
			},
			Menu{
				Text: "帮助",
				Items: []MenuItem{
					Action{
						AssignTo:    &showAboutBoxAction,
						Text:        "关于",
						OnTriggered: app.showAboutTriggered,
					},
				},
			},
		},
		Visible:  false,
		AssignTo: &app.mw,
		Title:    app.title,
		Size:     Size{Width: 300, Height: 160},
		Font:     Font{Family: "微软雅黑", PointSize: 9},
		Icon:     icon,
		Layout:   VBox{},
		Children: []Widget{
			Composite{
				Layout: VBox{MarginsZero: true},
				Children: []Widget{
					TabWidget{
						// ContentMarginsZero: true,
						//ContentMargins: Margins{Left: 400},
						//MinSize:        Size{Width: 500, Height: 600},
						Pages: []TabPage{
							// 水印操作
							TabPage{
								Title:      "水印操作",
								Font:       Font{Family: "微软雅黑", PointSize: 9},
								Background: SolidColorBrush{Color: walk.RGB(240, 240, 240)},
								Layout:     VBox{},
								Children: []Widget{
									GroupBox{
										Title:  "文件设置",
										Layout: Grid{Columns: 3},
										Children: []Widget{
											Label{
												Text: "文件目录:",
												//MinSize:   Size{220, 30},
												//TextColor: walk.RGB(255, 255, 0)
											},
											LineEdit{
												AssignTo: &desc.pdfDir,
												ReadOnly: true,

												MaxSize: Size{Width: 320, Height: 30}},

											Composite{
												Layout:  HBox{},
												MaxSize: Size{Width: 132, Height: 30},
												Children: []Widget{
													PushButton{
														AssignTo: &app.selectDir,
														MaxSize:  Size{Width: 60, Height: 30},
														Text:     "选择",
														OnClicked: func() {
															dlg := new(walk.FileDialog)

															dlg.Title = "选择需要增加水印的文件所在的目录"
															dlg.ShowBrowseFolder(app.mw)
															desc.pdfDir.SetText(dlg.FilePath)
															app.pdfDir_export.SetText(dlg.FilePath + "\\_out_")
														},
													},
												},
											},
											Label{
												Text: "水印文本:"},
											LineEdit{
												AssignTo: &desc.waterTxt,
												ReadOnly: false,
												//Text:     `业主重点合作伙伴沟通\n\n\n      业主重点合作伙伴沟通`,
											},
											//MaxSize:  Size{320, 30}},
											Composite{
												Layout:  HBox{},
												MaxSize: Size{Width: 132, Height: 30},
												Children: []Widget{
													PushButton{
														AssignTo: &app.markView,
														MaxSize:  Size{Width: 60, Height: 30},
														Text:     "预览",
														OnClicked: func() {
															//desc.waterTxt.SetText("")
															go markView(*desc)
														},
													},
												},
											},
											Label{
												Text: "导出目录:"},
											LineEdit{
												AssignTo: &app.pdfDir_export,
												ReadOnly: true,
												//Text:     `业主重点合作伙伴沟通\n\n\n      业主重点合作伙伴沟通`,
											},
											//MaxSize:  Size{320, 30}},
											Composite{
												Layout:  HBox{},
												MaxSize: Size{Width: 132, Height: 30},
												Children: []Widget{
													PushButton{
														AssignTo: &app.selectDir,
														MaxSize:  Size{Width: 60, Height: 30},
														Text:     "选择",
														OnClicked: func() {
															dlg := new(walk.FileDialog)

															dlg.Title = "选择需要增加水印的文件所在的目录"
															dlg.ShowBrowseFolder(app.mw)
															app.pdfDir_export.SetText(dlg.FilePath)
														},
													},
												},
											},
										},
									},
									GroupBox{
										Title:  "参数设置",
										Layout: Grid{Columns: 6},
										//MaxSize: Size{Width: 800, Height: 400},
										Children: []Widget{
											Label{Text: "颜 色："},
											ComboBox{
												AssignTo: &cbox.color,
												//	Model:    app.model,
												//Model:        []string{"1", "2", "3", "4", "5"},
												Model: []string{"黑色", "红色", "蓝色", "灰色", "棕色", "紫色", "绿色"},

												MinSize:      Size{Width: 120, Height: 30},
												CurrentIndex: 0,

												OnCurrentIndexChanged: func() {
													switch cbox.color.CurrentIndex() {
													case 0: //"openwrt":
														desc.color = "black"
													case 1: //"hiwifi":
														desc.color = "red"
													case 2: //"asus":
														desc.color = "blue"
													case 3: //"asus":
														desc.color = "gray"
													case 4: //"asus":
														desc.color = "brown"
													case 5: //"asus":
														desc.color = "purple"
													case 6: //"asus":
														desc.color = "green"
													default:
														desc.color = "black"
													}

												},
												//Font:           Font{PointSize: 1},
												//Value:         Bind("value", SelRequired{}),
												//BindingMember: "name",
												//DisplayMember: "value",
											},

											Label{Text: "透 明："},
											ComboBox{
												AssignTo:     &cbox.opacity,
												Model:        []string{"5%", "10%", "20%", "40%", "80%"},
												MinSize:      Size{Width: 120, Height: 30},
												CurrentIndex: 1,

												OnCurrentIndexChanged: func() {
													switch cbox.opacity.CurrentIndex() {
													case 0: //"openwrt":
														desc.opacity = "5%"
													case 1: //"openwrt":
														desc.opacity = "10%"
													case 2: //"hiwifi":
														desc.opacity = "20%"
													case 3: //"hiwifi":
														desc.opacity = "40%"
													case 4: //"hiwifi":
														desc.opacity = "90%"

													default:
														desc.opacity = "10%"
													}

												},
											},
											Label{Text: "背 景："},
											ComboBox{
												AssignTo:     &desc.bgcolor,
												Model:        []string{"0.1", "0.2", "0.4", "0.6"},
												MinSize:      Size{Width: 120, Height: 30},
												CurrentIndex: 0,
												Enabled:      false,
											},
											Label{Text: "旋 转："},

											ComboBox{
												AssignTo: &cbox.rotation,
												Model:    []string{"10", "45", "-20", "-45"},
												//MaxSize:      Size{Width: 320, Height: 30},
												CurrentIndex: 3,

												OnCurrentIndexChanged: func() {
													switch cbox.rotation.CurrentIndex() {
													case 0: //"openwrt":
														desc.rotation = "20"
													case 1: //"hiwifi":
														desc.rotation = "45"
													case 2: //"hiwifi":
														desc.rotation = "-20"
													case 3: //"hiwifi":
														desc.rotation = "-45"
													default:
														desc.rotation = "-45"
													}

												},
											},
											Label{Text: "字 体："},
											ComboBox{
												AssignTo: &cbox.font,
												//	Model:    app.model,
												//Model:        []string{"1", "2", "3", "4", "5"},
												Model: []string{"楷体", "仿宋", "黑体", "隶书"},
												OnCurrentIndexChanged: func() {
													switch cbox.font.CurrentIndex() {
													case 0:
														desc.font = "simkai"
													case 1:
														desc.font = "simfang"
													case 2:
														desc.font = "simhei"
													case 3:
														desc.font = "simli"
													default:
														desc.font = "simkai"
													}

												},
												MaxSize:      Size{Width: 320, Height: 30},
												CurrentIndex: 0,
											},

											Label{Text: "字号："},
											ComboBox{
												AssignTo: &cbox.points,
												Model:    []string{"10", "40", "80", "120"},
												//MaxSize:      Size{Width: 320, Height: 30},
												CurrentIndex: 0,
												Enabled:      false,
												OnCurrentIndexChanged: func() {
													switch cbox.points.CurrentIndex() {
													case 0: //"openwrt":
														desc.points = "10"
													case 1: //"hiwifi":
														desc.points = "40"
													case 2: //"hiwifi":
														desc.points = "80"
													case 3: //"hiwifi":
														desc.points = "120"

													default:
														desc.points = "10"
													}

												},
											},
										},
									},
									Composite{
										Layout: HBox{},
										Children: []Widget{
											PushButton{
												MinSize: Size{Width: 121, Height: 30},
												Enabled: true,
												Text:    "添加水印",
												OnClicked: func() {
													go addWatermark(*desc)
												},
											},
											PushButton{
												MinSize: Size{Width: 121, Height: 30},
												Text:    "导出PDF",
												Enabled: true,
												OnClicked: func() {
													go exportImagesFile()
												},
											},

											PushButton{
												MinSize: Size{Width: 121, Height: 30},
												Text:    "关闭退出",
												OnClicked: func() {
													reprocessing()
													// 无语的等待
													time.Sleep(time.Second * 1)
													walk.App().Exit(0)

												},
											},
										},
									},
								},
							},

							// 转换操作tab页
							TabPage{
								Title:      "转换操作",
								Font:       Font{Family: "微软雅黑", PointSize: 9},
								Layout:     Grid{Columns: 2},
								MaxSize:    Size{Width: 220, Height: 20},
								Background: SolidColorBrush{Color: walk.RGB(240, 240, 240)},
								Children: []Widget{
									GroupBox{
										Title:  "文件设置",
										Layout: Grid{Columns: 3},
										Children: []Widget{
											Label{
												Text: "文件目录:",
												//MinSize:   Size{220, 30},
												//TextColor: walk.RGB(255, 255, 0)
											},
											LineEdit{
												AssignTo: &app.imgDir,
												ReadOnly: true,
												MaxSize:  Size{Width: 320, Height: 30}},
											Composite{
												Layout:  HBox{},
												MaxSize: Size{Width: 132, Height: 30},
												Children: []Widget{
													PushButton{
														//AssignTo: &app.selectDir,
														MaxSize: Size{Width: 60, Height: 30},
														Text:    "选择",
														OnClicked: func() {
															dlg := new(walk.FileDialog)
															//dlg.FilePath = mw.prevFilePath
															dlg.Filter = "图像文件s (*.emf;*.bmp;*.exif;*.gif;*.jpeg;*.jpg;*.png;*.tiff)|*.emf;*.bmp;*.exif;*.gif;*.jpeg;*.jpg;*.png;*.tiff"
															dlg.Title = "选择图像文件"
															if ok, err := dlg.ShowOpen(app.mw); err != nil {
																showMsg("选择文件错误：" + err.Error())
															} else if !ok {
																showMsg("选择文件错误：" + err.Error())
															}

															//if err := importImagesFile([]string{dlg.FilePath}, dlg.FilePath+".pdf"); err != nil {
															//	showMsg("图像转换错误：" + err.Error())
															//} else {
															//	showMsg("图像转换完成：" + dlg.FilePath + ".pdf")
															//}
															app.imgDir.SetText(dlg.FilePath)
														},
													},
												},
											},
											Label{
												Text: "输出目录:"},
											TextEdit{
												//AssignTo: &desc.waterTxt,
												VScroll:  true,
												HScroll:  false,
												ReadOnly: false},
											//MaxSize:  Size{320, 30}},
											Composite{
												Layout:  HBox{},
												MaxSize: Size{Width: 132, Height: 30},
												Children: []Widget{
													PushButton{
														//AssignTo: &app.selectDir,
														MaxSize: Size{Width: 60, Height: 30},
														Text:    "选择",
														OnClicked: func() {
															//desc.waterTxt.SetText("")
														},
													},
												},
											},
										},
									},
									GroupBox{
										Title:  "参数设置",
										Layout: Grid{Columns: 6},
										//MaxSize: Size{Width: 800, Height: 400},
										Children: []Widget{
											Label{Text: "颜 色："},
											ComboBox{
												//AssignTo: &desc.color,
												//	Model:    app.model,
												//Model:        []string{"1", "2", "3", "4", "5"},
												Model: []string{"黑色", "红色", "蓝色", "灰色", "浅灰", "深灰", "绿色"},

												MinSize:      Size{Width: 120, Height: 30},
												CurrentIndex: 0,
												//Font:           Font{PointSize: 1},
												//Value:         Bind("value", SelRequired{}),
												//BindingMember: "name",
												//DisplayMember: "value",
											},
											Label{Text: "背 景："},
											ComboBox{
												//AssignTo:     &desc.bgcolor,
												Model:        []string{"0.2", "0.4", "0.6", "0.8"},
												MinSize:      Size{Width: 120, Height: 30},
												CurrentIndex: 0,
											},
											Label{Text: "透 明："},
											ComboBox{
												//AssignTo:     &desc.opacity,
												Model:        []string{"0.2", "0.4", "0.6", "0.8"},
												MinSize:      Size{Width: 120, Height: 30},
												CurrentIndex: 0,
											},

											Label{Text: "旋 转："},

											ComboBox{
												//AssignTo: &desc.rotation,
												Model: []string{"0.2", "0.4", "0.6", "0.8"},
												//MaxSize:      Size{Width: 320, Height: 30},
												CurrentIndex: 0,
											},
											Label{Text: "字 体："},
											ComboBox{
												//AssignTo: &desc.font,
												//	Model:    app.model,
												//Model:        []string{"1", "2", "3", "4", "5"},
												Model: []string{"宋体", "楷体"},

												MaxSize:      Size{Width: 320, Height: 30},
												CurrentIndex: 0,
											},

											Label{Text: "字号："},
											ComboBox{
												//	AssignTo: &desc.points,
												Model: []string{"0.2", "0.4", "0.6", "0.8"},
												//MaxSize:      Size{Width: 320, Height: 30},
												CurrentIndex: 0,
											},
										},
									},
									Composite{
										Layout: HBox{},
										Children: []Widget{
											PushButton{
												MinSize: Size{Width: 121, Height: 30},
												Enabled: true,
												Text:    "图像转pdf",
												OnClicked: func() {
													dlg := new(walk.FileDialog)
													//dlg.FilePath = mw.prevFilePath
													dlg.Filter = "图像文件s (*.emf;*.bmp;*.exif;*.gif;*.jpeg;*.jpg;*.png;*.tiff)|*.emf;*.bmp;*.exif;*.gif;*.jpeg;*.jpg;*.png;*.tiff"
													dlg.Title = "选择图像文件"
													if ok, err := dlg.ShowOpen(app.mw); err != nil {
														showMsg("选择文件错误：" + err.Error())
														return
													} else if !ok {
														return
													}

													if err := importImagesFile([]string{dlg.FilePath}, dlg.FilePath+".pdf"); err != nil {
														showMsg("图像转换错误：" + err.Error())
													} else {
														showMsg("图像转换完成：" + dlg.FilePath + ".pdf")
													}
													//app.pdfDir.SetText(dlg.FilePath)
												},
											},

											PushButton{
												MinSize: Size{Width: 121, Height: 30},
												Text:    "导出PDF",
												OnClicked: func() {
													go exportImagesFile()
												},
											},
											PushButton{
												MinSize: Size{Width: 121, Height: 30},
												Text:    "关闭退出",
												OnClicked: func() {
													walk.App().Exit(0)

												},
											},
										},
									},
								},
							},
							// 导出tab页
							TabPage{
								Title:      "导出操作",
								Font:       Font{Family: "微软雅黑", PointSize: 9},
								Layout:     VBox{},
								Background: SolidColorBrush{Color: walk.RGB(240, 240, 240)},
								Children: []Widget{
									GroupBox{
										Title:  "文件设置",
										Layout: Grid{Columns: 3},
										Children: []Widget{
											Label{
												Text: "文件目录:",
												//MinSize:   Size{220, 30},
												//TextColor: walk.RGB(255, 255, 0)
											},
											LineEdit{
												// 文件导出源目录
												//AssignTo: &app.pdfDir_export,
												ReadOnly: true,
												MaxSize:  Size{Width: 320, Height: 30}},
											Composite{
												Layout:  HBox{},
												MaxSize: Size{Width: 132, Height: 30},
												Children: []Widget{
													PushButton{
														AssignTo: &app.selectDir,
														MaxSize:  Size{Width: 60, Height: 30},
														Text:     "选择",
														OnClicked: func() {
															//walk.App().Exit(0)
															dlg := new(walk.FileDialog)
															//	dlg.FilePath = app.filePath
															//dlg.Filter = "Image Files (*.emf;*.bmp;*.exif;*.gif;*.jpeg;*.jpg;*.png;*.tiff)|*.emf;*.bmp;*.exif;*.gif;*.jpeg;*.jpg;*.png;*.tiff"
															dlg.Title = "选择PDF文件目录"
															dlg.ShowBrowseFolder(app.mw)

															app.pdfDir_export.SetText(dlg.FilePath)
														},
													},
												},
											},
											Label{
												Text: "水印文本:"},
											TextEdit{
												//AssignTo: &desc.waterTxt,
												VScroll:  true,
												HScroll:  false,
												ReadOnly: false},
											//MaxSize:  Size{320, 30}},
											Composite{
												Layout:  HBox{},
												MaxSize: Size{Width: 132, Height: 30},
												Children: []Widget{
													PushButton{
														//AssignTo: &app.selectDir,
														MaxSize: Size{Width: 60, Height: 30},
														Text:    "清空",
														OnClicked: func() {
															desc.waterTxt.SetText("")
														},
													},
												},
											},
										},
									},
									GroupBox{
										Title:  "参数设置",
										Layout: Grid{Columns: 6},
										//MaxSize: Size{Width: 800, Height: 400},
										Children: []Widget{
											Label{Text: "颜 色："},
											ComboBox{
												//AssignTo: &desc.color,
												//	Model:    app.model,
												//Model:        []string{"1", "2", "3", "4", "5"},
												Model: []string{"黑色", "红色", "蓝色", "灰色", "浅灰", "深灰", "绿色"},

												MinSize:      Size{Width: 120, Height: 30},
												CurrentIndex: 0,
												//Font:           Font{PointSize: 1},
												//Value:         Bind("value", SelRequired{}),
												//BindingMember: "name",
												//DisplayMember: "value",
											},
											Label{Text: "背 景："},
											ComboBox{
												//AssignTo:     &desc.bgcolor,
												Model:        []string{"0.2", "0.4", "0.6", "0.8"},
												MinSize:      Size{Width: 120, Height: 30},
												CurrentIndex: 0,
											},
											Label{Text: "透 明："},
											ComboBox{
												//AssignTo:     &desc.opacity,
												Model:        []string{"0.2", "0.4", "0.6", "0.8"},
												MinSize:      Size{Width: 120, Height: 30},
												CurrentIndex: 0,
											},

											Label{Text: "旋 转："},

											ComboBox{
												//AssignTo: &desc.rotation,
												Model: []string{"0.2", "0.4", "0.6", "0.8"},
												//MaxSize:      Size{Width: 320, Height: 30},
												CurrentIndex: 0,
											},
											Label{Text: "字 体："},
											ComboBox{
												//AssignTo: &desc.font,
												//	Model:    app.model,
												//Model:        []string{"1", "2", "3", "4", "5"},
												Model: []string{"宋体", "楷体"},

												MaxSize:      Size{Width: 320, Height: 30},
												CurrentIndex: 0,
											},

											Label{Text: "字号："},
											ComboBox{
												//AssignTo: &desc.points,
												Model: []string{"0.2", "0.4", "0.6", "0.8"},
												//MaxSize:      Size{Width: 320, Height: 30},
												CurrentIndex: 0,
											},
										},
									},
									Composite{
										Layout: HBox{},
										Children: []Widget{
											PushButton{
												MinSize: Size{Width: 121, Height: 30},
												Enabled: true,
												Text:    "图像转pdf",
												OnClicked: func() {
													dlg := new(walk.FileDialog)
													//dlg.FilePath = mw.prevFilePath
													dlg.Filter = "Image Files (*.emf;*.bmp;*.exif;*.gif;*.jpeg;*.jpg;*.png;*.tiff)|*.emf;*.bmp;*.exif;*.gif;*.jpeg;*.jpg;*.png;*.tiff"
													dlg.Title = "Select an Image"
													if ok, err := dlg.ShowOpen(app.mw); err != nil {
														showMsg("选择文件错误：" + err.Error())
													} else if !ok {
														showMsg("选择文件错误：" + err.Error())
													}

													if err := importImagesFile([]string{dlg.FilePath}, dlg.FilePath+".pdf"); err != nil {
														showMsg("图像转换错误：" + err.Error())
													} else {
														showMsg("图像转换完成：" + dlg.FilePath + ".pdf")
													}
													//app.pdfDir.SetText(dlg.FilePath)
												},
											},
											PushButton{
												MinSize: Size{Width: 121, Height: 30},
												Text:    "导出PDF",
												OnClicked: func() {
													go exportImagesFile()
												},
											},
											PushButton{
												MinSize: Size{Width: 121, Height: 30},
												Text:    "关闭退出",
												OnClicked: func() {
													walk.App().Exit(0)

												},
											},
										},
									},
								},
							},
						},
					},
					GroupBox{
						Title:  "操作输出：",
						Layout: Grid{Columns: 1},

						Children: []Widget{
							PushButton{
								MinSize: Size{Width: 80, Height: 30},
								Text:    "清除日志",
								OnClicked: func() {
									app.msg.SetText("")

								},
							},
							TextEdit{AssignTo: &app.msg, HScroll: true, VScroll: true, ReadOnly: true, Font: Font{PointSize: 14}, MinSize: Size{Width: 132, Height: 200}},
						},
					},
				},
			},
		},
		OnSizeChanged: func() {
			_ = app.mw.SetSize(walk.Size(Size{Width: 500, Height: 360}))
		},
	}.Create()
	winLong := win.GetWindowLong(app.mw.Handle(), win.GWL_STYLE)
	// 不能调整窗口大小，禁用最大化按钮,取消一切操作
	win.SetWindowLong(app.mw.Handle(), win.GWL_STYLE, winLong & ^win.WS_SIZEBOX & ^win.WS_MAXIMIZEBOX & ^win.WS_SIZEBOX & ^win.WS_SYSMENU)
	// 设置窗体生成在屏幕的正中间，并处理高分屏的情况
	// 窗体横坐标 = ( 屏幕宽度 - 窗体宽度 ) / 2
	// 窗体纵坐标 = ( 屏幕高度 - 窗体高度 ) / 2
	_ = app.mw.SetX((int(win.GetSystemMetrics(0)) - app.mw.Width()) / 2 / app.mw.DPI() * 96)
	_ = app.mw.SetY((int(win.GetSystemMetrics(1)) - app.mw.Height()) / 2 / app.mw.DPI() * 96)
	return err
}
