package main

import (
	"os"
	"strconv"

	"github.com/lxn/walk"
	pdfcpu "github.com/pdfcpu/pdfcpu/pkg/api"
)

type myApp struct {
	mw             *walk.MainWindow
	title          string
	pdfDir_export  *walk.LineEdit // pdf文件目录
	imgDir_single  *walk.LineEdit // 图像文件目录
	imgDir_batch   *walk.LineEdit // 图像文件目录
	msg            *walk.TextEdit //
	pdfpath        *walk.PushButton
	pdfpath_export *walk.PushButton
	markView       *walk.PushButton
	exitbtn        *walk.PushButton
	tabwatermark   *walk.TabPage
	tabimage       *walk.TabPage
}
type waterDesc struct {
	pdfDir    *walk.LineEdit // pdf文件目录
	waterTxt  *walk.LineEdit // 水印文本
	opacity   string         // 透明度
	color     string         // 颜色
	bgcolor   *walk.ComboBox // 颜色
	rotation  string         // 旋转
	font      string         // 水印文本字体
	points    string         // 水印文本字号
	addbtn    *walk.PushButton
	exportbtn *walk.PushButton
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
var viewProcess []string
var viewFile []string
var showAboutBoxAction *walk.Action

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
func reprocessing() {
	// 清理预览进程 ，无语问苍天
	for _, pid := range viewProcess {
		p, _ := strconv.Atoi(pid)
		KillAll(p)
	}
	// 清理预览临时文件，无语问苍天
	for _, file := range viewFile {
		os.RemoveAll(file)
	}
}
func main() {
	InitLogger()
	//SugarLogger.Info("程序启动了......")
	_ = getWindows()
	walk.App().SetProductName(app.title)
	walk.App().SetOrganizationName("泛生态业务线")
	//brush, _ := walk.NewSolidColorBrush(walk.RGB(0, 255, 0))
	//app.mw.SetBackground(brush)
	app.mw.Show()
	app.mw.Run()
}
