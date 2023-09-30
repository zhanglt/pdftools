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
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/types"
)

// pdfcpu 增加水印--已经准备废弃--用c++的小工具代替
func textWatermark(d waterDesc) {
	//desc := "font:SimSun,points:48, sc:1, fillc:.8 .8 .4, op:.3"
	desc := "font:" + d.font + ", points:" + d.points + ", scale:1.0 rel" + " , rot:" + d.rotation + ", fillc:" + d.color + ", margin:10, border:10 round, opacity:" + d.opacity
	//desc := "font:" + d.font + ", points:" + d.points + ", rtl:off, sc:0.5 rel, pos:c, off:0 0, align:c, fillc:#808080, strokec:#808080, rot:0, d:1, op:1, mo:0, ma:0, bo:0"
	//desc := "font:" + d.font + ", points:" + d.points + ", scale:1.0 rel, rot:0, fillc:#000000, bgcol:#ab6f30, margin:10, border:10 round, opacity:.7"
	//desc := "font:" + d.font + ", points:" + d.points
	showMsg(desc)
	files := getFilelist(d.pdfDir.Text())
	wm, _ := pdfcpu.TextWatermark(d.waterTxt.Text(), desc, true, false, types.POINTS)
	//wm, _ := pdfcpu.ImageWatermark("image.png", "scalefactor:.5 a, rot:-90", true, false, types.POINTS)
	for _, file := range files {
		pathName := file[0 : len(file)-len(filepath.Base(file))]
		fileName := filepath.Base(file)
		go stmp(file, pathName, fileName, wm)
		/*
			err := pdfcpu.AddWatermarksFile(file, pathName+"/_mark_"+fileName, nil, wm, nil)
			if err != nil {
				showMsg("水印错误：" + err.Error())
			}
		*/
	}

	// Stamp all odd pages of in.pdf in red "Confidential" in 48 point Courier
	// using a rotation angle of 45 degrees and an absolute scalefactor of 1.0.
	//onTop = true
	//wm, _ = pdfcpu.TextWatermark(app.waterTxt.Text(), "font:"+app.fontname.Text()+", points:"+app.fontsize.Text()+", col: 1 0 0, rot:45, sc:1 abs ", onTop, update, types.POINTS)
	//pdfcpu.AddWatermarksFile("1.pdf", "3.pdf", nil, wm, nil)

}

func stmp(file, pathName, fileName string, wm *model.Watermark) {
	err := pdfcpu.AddWatermarksFile(file, pathName+"/_mark_"+fileName, nil, wm, nil)
	if err != nil {
		showMsg("水印错误：" + err.Error())
	}
}

func addWatermark(d waterDesc) {
	if d.waterTxt.Text() == "" {
		showMessage("水印文本为空", "请填写水印文本内容")
		d.waterTxt.SetFocus()
		return
	}
	files := getFilelist(d.pdfDir.Text())
	for _, file := range files {
		pathName := file[0 : len(file)-len(filepath.Base(file))]
		fileName := filepath.Base(file)
		if strings.Contains(fileName, " ") {
			showMsg(fileName)
			os.Rename(file, pathName+strings.Replace(fileName, " ", "", -1))
			fileName = strings.Replace(fileName, " ", "", -1)
		}

		//cli := "wmark.exe  " + file + "  " + pathName + "_mark_" + fileName + "  " + strings.Replace(d.waterTxt.Text(), " ", "", -1) + "  " + d.opacity + " " + d.color + " " + d.rotation + " " + d.font
		cli := "wmark.exe" + " " + pathName + fileName + " " + pathName + "_mark_" + fileName + " " + strings.Replace(d.waterTxt.Text(), " ", "", -1) + " " + d.opacity + " " + d.color + " " + d.rotation + " " + d.font

		showMsg(cli)
		c := exec.Command("cmd.exe", "/C", cli)
		c.SysProcAttr = &syscall.SysProcAttr{HideWindow: true} // 不显示命令窗口

		if err := c.Start(); err != nil {
			showMsg(fmt.Sprintf("%s增加水印失败, 错误信息: %s", fileName, err))
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

}

func markView(d waterDesc) {
	if d.waterTxt.Text() == "" {
		showMessage("水印文本为空", "请填写水印文本内容")
		d.waterTxt.SetFocus()
		return
	}
	var pathName, fileName string
	files := getFilelist(d.pdfDir.Text())
	for _, file := range files {
		pathName = file[0 : len(file)-len(filepath.Base(file))]
		fileName = filepath.Base(file)
		if strings.Contains(fileName, " ") {
			os.Rename(file, pathName+strings.Replace(fileName, " ", "", -1))
			fileName = strings.Replace(fileName, " ", "", -1)
		}
		if fileName != "" {
			break
		}
	}

	//cli := "wmark.exe  " + file + "  " + pathName + "_mark_" + fileName + "  " + strings.Replace(d.waterTxt.Text(), " ", "", -1) + "  " + d.opacity + " " + d.color + " " + d.rotation + " " + d.font
	cli := "wmark.exe" + " " + pathName + fileName + " " + pathName + "_mark_" + fileName + " " + strings.Replace(d.waterTxt.Text(), " ", "", -1) + " " + d.opacity + " " + d.color + " " + d.rotation + " " + d.font
	showMsg(cli)
	c := exec.Command("cmd.exe", "/C", cli)
	c.SysProcAttr = &syscall.SysProcAttr{HideWindow: true} // 不显示命令窗口
	if err := c.Start(); err != nil {
		showMsg(fmt.Sprintf("%s增加水印失败, 错误信息: %s", fileName, err))
	}
	c.Wait()
	go func() {
		defer func() {
			os.Remove(pathName + "_mark_" + fileName)
			app.markView.SetEnabled(true)
		}()

		view := exec.Command("cmd.exe", "/C", "mupdf.exe  "+pathName+"_mark_"+fileName)
		view.SysProcAttr = &syscall.SysProcAttr{HideWindow: true} // 不显示命令窗口
		if err := view.Start(); err != nil {
			showMsg(fmt.Sprintf("%s预览失败: %s", fileName, err))
		}
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
	pdfFiles := getFilelist(app.pdfDir_export.Text())
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
		c := exec.Command("cmd", "/C", "mutool draw -r 120 -o "+pathName+"/_img/"+"%d.png  -F png  ", file)
		c.SysProcAttr = &syscall.SysProcAttr{HideWindow: true} // 不显示命令窗口
		if err := c.Start(); err != nil {
			showMsg(fmt.Sprintf("转换数据失败, 错误信息: %s", err))
		}
		c.Wait()
		for i := 1; i < pageCount+1; i++ {
			importImagesFile([]string{pathName + "/_img/" + strconv.Itoa(i) + ".png"}, pathName+strconv.Itoa(i)+fileName)
		}
		var pdf []string
		for i := 1; i < pageCount+1; i++ {
			pdf = append(pdf, pathName+strconv.Itoa(i)+fileName)
		}
		MergeAppendFile(pdf, pathName+"_img_"+fileName)
		os.RemoveAll(pathName + "/_img")
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
	//getFilelist()
	conf := pdfcpu.LoadConfiguration()
	//inFile := filepath.Join("e:/test/", file)
	f, err := os.Open(inFile)
	if err != nil {
		showMsg(err.Error())
	}
	defer f.Close()

	info, err := api.PDFInfo(f, inFile, nil, conf)
	if err != nil {
		showMsg(err.Error())
	} else {
		showMsg("文件页数：" + strconv.Itoa(info.PageCount))
	}
	return info.PageCount
}

// 获取PDF文件列表
func getFilelist(filePath string) []string {
	var files []string
	var root string
	if filePath != "" {
		root = filePath
	} else {
		showMessage("文件目录为空", "请选择需要处理的文件所在目录")
		return nil
	}
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if strings.Contains(strings.ToLower(path), ".pdf") {
			files = append(files, path)
			//showMsg(filepath.Base(path))
			//showMsg(path[0 : len(path)-len(filepath.Base(path))])
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
