package main

import (
	"github.com/lxn/walk"
	pdfcpu "github.com/pdfcpu/pdfcpu/pkg/api"
)

type myApp struct {
	mw            *walk.MainWindow
	title         string
	pdfDir_export *walk.LineEdit // pdf文件目录
	imgDir        *walk.LineEdit // 图像文件目录
	msg           *walk.TextEdit //
	selectDir     *walk.PushButton
	markView      *walk.PushButton
	//tabPage       *walk.TabPage
}
type waterDesc struct {
	pdfDir   *walk.LineEdit // pdf文件目录
	waterTxt *walk.LineEdit // 水印文本
	opacity  string         // 透明度
	color    string         // 颜色
	bgcolor  *walk.ComboBox // 颜色
	rotation string         // 旋转
	font     string         // 水印文本字体
	points   string         // 水印文本字号
}
type cBox struct {
	opacity  *walk.ComboBox // 透明度
	color    *walk.ComboBox // 颜色
	font     *walk.ComboBox // 水印文本字体
	rotation *walk.ComboBox // 旋转
	points   *walk.ComboBox // 水印文本字号
}

// 关于对话框
func (mw *myApp) showAboutTriggered() {
	walk.MsgBox(mw.mw, "关于...", "泛生态行业线PDF处理工具", walk.MsgBoxIconInformation)
}

// 信息弹框
func (mw *myApp) showNoneMessage(message string) {
	walk.MsgBox(mw.mw, "提示", message, walk.MsgBoxIconInformation)
}

var app = new(myApp)
var desc = new(waterDesc)
var cbox = new(cBox)
var showAboutBoxAction *walk.Action

// var wg sync.WaitGroup
func init() {
	app.title = "文档水印处理-泛生态业务线"
	//desc.color = &walk.ComboBox{}

	//app.model = NewEnvModel()
	pdfcpu.LoadConfiguration()
	//importImagesFile()
	/*
		l, _ := pdfcpu.ListFonts()
		if len(l) < 18 {
			pdfcpu.InstallFonts([]string{"simsun.ttc"})
			pdfcpu.InstallFonts([]string{"simkai.ttf"})
		}*/

}
func main() {
	_ = getWindows()
	walk.App().SetProductName(app.title)
	walk.App().SetOrganizationName("泛生态业务线")
	app.mw.Show()
	app.mw.Run()
}
