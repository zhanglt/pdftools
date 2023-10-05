package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/lxn/walk"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	pdfcpu "github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/types"
)

// func pathChange(root, input string) string {

func pathChange(root, file, insert string) string {
	//文件名称
	filename := filepath.Base(file)
	// 文件路径
	pathname := file[0 : len(file)-len(filepath.Base(file))]
	// 根目录后的子目录
	subpath := pathname[len(root):]
	//  组合新目录
	newpath := root + "\\" + insert + subpath
	// 判断目录
	if _, err := os.Stat(newpath); os.IsNotExist(err) {
		os.MkdirAll(newpath, os.ModePerm)
	}
	showMsg(newpath + filename)
	return newpath + filename
}

// 增加水印
func addWatermark(d waterDesc) {
	defer func() {
		showMsg("------水印追加完成!-----")
		showMsg("保存目录：" + d.pdfDir.Text() + "\\_out_")
	}()
	showMsg("------开始转换，请耐心等待-----")
	if d.waterTxt.Text() == "" {
		showMessage("水印文本为空", "请填写水印文本内容")
		d.waterTxt.SetFocus()
		return
	}
	files := getFilelist(d.pdfDir.Text(), "_out_")
	if len(files) == 0 {
		showMessage("目中没有文件", "目录中没有符合条件的文件")
		return
	} else {
		// 创建输出目录
		os.Mkdir(d.pdfDir.Text()+"\\_out_", os.ModePerm)
	}
	root := d.pdfDir.Text()
	for _, file := range files {
		pathname := file[0 : len(file)-len(filepath.Base(file))]
		filename := filepath.Base(file)
		if strings.Contains(filename, " ") {
			// 蹩脚的解决文件包含空格的问题，改名！！
			os.Rename(file, pathname+strings.Replace(filename, " ", "", -1))
			filename = strings.Replace(filename, " ", "", -1)
		}
		//cli := "wmark.exe  " + file + "  " + pathName + "_mark_" + fileName + "  " + strings.Replace(d.waterTxt.Text(), " ", "", -1) + "  " + d.opacity + " " + d.color + " " + d.rotation + " " + d.font
		cli := "wmark.exe" + " -i" + file + " -o" + pathChange(root, file, "_out_") + " -t" + strings.Replace(d.waterTxt.Text(), " ", "", -1) + " -p" + d.opacity + " -c" + d.color + " -r" + d.rotation + " -f" + d.font
		//showMsg(cli)
		c := exec.Command("cmd.exe", "/C", cli)
		c.SysProcAttr = &syscall.SysProcAttr{HideWindow: true} // 不显示命令窗口
		if err := c.Start(); err != nil {
			showMsg(fmt.Sprintf("%s增加水印失败, 错误信息: %s", filename, err))
		}

		c.Wait()
		/*
			p, err := os.FindProcess(c.Process.Pid)
			if err != nil {
				panic(err)
			}
			err = p.Kill()
			if err != nil {
				panic(err)
			}
		*/
	}
	app.pdfDir_export.SetText(root + "\\_out_")

}

// 水印预览
func markView(d waterDesc) {
	if d.waterTxt.Text() == "" {
		showMessage("水印文本为空", "请填写水印文本内容")
		d.waterTxt.SetFocus()
		return
	}
	var pathName, fileName, pid string
	// 目录中遍历查抄文件，排除_out_文件夹
	files := getFilelist(d.pdfDir.Text(), "_out_")

	if len(files) == 0 {
		showMessage("目中没有文件", "目录中没有符合条件的文件")
		return
	}
	// 从遍历结果随机取出一个文件
	for _, file := range files {
		pathName = file[0 : len(file)-len(filepath.Base(file))]
		fileName = filepath.Base(file)
		// 如果文件中包含空格，则去除空格改名
		if strings.Contains(fileName, " ") {
			os.Rename(file, pathName+strings.Replace(fileName, " ", "", -1))
			fileName = strings.Replace(fileName, " ", "", -1)
		}
		// 找到第一个文件就返回
		if fileName != "" {
			break
		}
	}
	// 给上面取出的文件增加水印
	//cli := "wmark.exe  " + file + "  " + pathName + "_mark_" + fileName + "  " + strings.Replace(d.waterTxt.Text(), " ", "", -1) + "  " + d.opacity + " " + d.color + " " + d.rotation + " " + d.font
	cli := "wmark.exe" + " -i" + pathName + fileName + " -o" + pathName + "_view_" + fileName + " -t" + strings.Replace(d.waterTxt.Text(), " ", "", -1) + " -p" + d.opacity + " -c" + d.color + " -r" + d.rotation + " -f" + d.font
	//showMsg(cli)
	c := exec.Command("cmd.exe", "/C", cli)
	c.SysProcAttr = &syscall.SysProcAttr{HideWindow: true} // 不显示命令窗口
	if err := c.Start(); err != nil {
		showMsg(fmt.Sprintf("%s增加水印失败, 错误信息: %s", fileName, err))
		return
	}
	c.Wait()
	go func() {
		// 预览结束后的善后工作，
		defer func() {
			// 删除水印文件
			os.Remove(pathName + "_view_" + fileName)
			// 从预览文件列表中删除
			removeElement(viewFile, pathName+"_view_"+fileName)
			// 从预览程序pid列表中删除
			removeElement(viewProcess, pid)
			// 设置预览按钮可用
			app.markView.SetEnabled(true)
		}()
		// 调用预览程序打开增加过水印的文件
		view := exec.Command("cmd.exe", "/C", "pdfview.exe  "+pathName+"_view_"+fileName)
		view.SysProcAttr = &syscall.SysProcAttr{HideWindow: true} // 不显示命令窗口
		if err := view.Start(); err != nil {
			showMsg(fmt.Sprintf("%s预览失败: %s", fileName, err))
		}
		// 获取预览程序的pid
		pid = strconv.Itoa(view.Process.Pid)
		// 将预览进程pid写入列表，用于退出时的清理
		viewProcess = append(viewProcess, pid)
		// 将预览文件名称写入预览文件列表，用于退出时的清理
		viewFile = append(viewFile, pathName+"_view_"+fileName)
		// showMsg(pathName + "_view_" + fileName)
		// 设置预览按钮为不可用
		app.markView.SetEnabled(false)
		view.Wait()
	}()
}

// 图像转pdf
func importImagesFile(files []string, pdfFile string) error {
	imp, _ := pdfcpu.Import("form:A3, pos:c, s:1.0", types.POINTS)
	err := pdfcpu.ImportImagesFile(files, pdfFile, imp, nil)
	if err != nil {
		showMsg(err.Error())
		return err
	}
	return nil
	//pdfcpu.ImportImagesFile([]string{"a1.png", "a2.png", "a3.png"}, "out2.pdf", imp, nil)
}

// PDF合并
func MergeAppendFile(inFiles []string, fileName string) {
	defer showMsg("合并结束-------------")
	pdfcpu.MergeAppendFile(inFiles, fileName, nil)
}

// pdf--图像--pdf
func exportImagesFile() {
	var pageCount int
	// 遍历文件夹获取pdf文件，排除"_img"文件夹
	pdfFiles := getFilelist(app.pdfDir_export.Text(), "_img")
	if len(pdfFiles) == 0 {
		showMessage("目中没有文件", "目录中没有符合条件的文件")
		return
	}
	root := app.pdfDir_export.Text()
	for _, file := range pdfFiles {
		// 获取pdf文件的页数
		pageCount = getPageCount(file)
		// 获取pdf所在目录
		pathName := file[0 : len(file)-len(filepath.Base(file))]
		// 获取pdf文件名称
		fileName := filepath.Base(file)
		// 临时目录"_img"不存在则创建
		if _, err := os.Stat(pathName + "\\_img"); os.IsNotExist(err) {
			os.MkdirAll(pathName+"_img", os.ModePerm)
		}
		// 导出文件为图片，导出文件名为数字，为了后续合并方便。
		c := exec.Command("cmd", "/C", "mutool draw -r 120 -o "+pathName+"/_img/"+"%d.png  -F png  ", file)
		c.SysProcAttr = &syscall.SysProcAttr{HideWindow: true} // 不显示命令窗口
		if err := c.Start(); err != nil {
			showMsg(fmt.Sprintf("转换数据失败, 错误信息: %s", err))
		}
		c.Wait()
		// 图片转为pdf
		for i := 1; i < pageCount+1; i++ {
			importImagesFile([]string{pathName + "/_img/" + strconv.Itoa(i) + ".png"}, pathName+strconv.Itoa(i)+fileName)
		}
		// 合并pdf
		var pdf []string
		for i := 1; i < pageCount+1; i++ {
			pdf = append(pdf, pathName+strconv.Itoa(i)+fileName)
		}
		MergeAppendFile(pdf, pathChange(root, file, "_img_")+fileName)
		// 删除临时文件夹
		os.RemoveAll(pathName + "/_img")
		// 删除临时pdf文件
		for i := 1; i < pageCount+1; i++ {
			os.Remove(pathName + strconv.Itoa(i) + fileName)
		}
	}
	showMsg("----------任务结束-----------")
	time.Sleep(time.Second * 2)
}
func showMsg(msg string) {
	app.msg.AppendText(time.Now().Format("2006-01-02 15:04:05 "))
	app.msg.AppendText(msg)
	app.msg.AppendText("\r\n")
}

// 获取PDF文件页数
func getPageCount(inFile string) int {
	conf := pdfcpu.LoadConfiguration()
	f, err := os.Open(inFile)
	if err != nil {
		showMsg("打开文件失败：" + inFile + err.Error())
	}
	defer f.Close()

	info, err := api.PDFInfo(f, inFile, nil, conf)
	if err != nil {
		showMsg("获取文件信息失败：" + inFile + err.Error())
	} //else {
	//	showMsg("文件页数：" + strconv.Itoa(info.PageCount))
	//}
	return info.PageCount
}

// 获取PDF文件列表
func getFilelist(filePath, skip string) []string {
	var files []string
	var root string
	if filePath != "" {
		root = filePath
	} else {
		showMessage("文件目录为空", "请选择需要处理的文件所在目录")
		return nil
	}
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		//showMsg(info.Name())
		if info.IsDir() && info.Name() == skip {
			return filepath.SkipDir
		}
		if strings.Contains(strings.ToLower(path), ".pdf") {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		showMsg("获取PDF文件列表失败:" + err.Error())
	}
	return files
}

func showMessage(title, msg string) {
	walk.MsgBox(app.mw,
		title,
		msg,
		walk.MsgBoxOK|walk.MsgBoxIconInformation)
}

// 删除slice中的指定项
func removeElement(slice []string, elem string) []string {
	for i, v := range slice {
		if v == elem {
			return append(slice[:i], slice[i+1:]...)
		}
	}
	return slice
}
