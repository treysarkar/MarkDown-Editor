// Markdown editor I am building to test my knowledge,  like what I have learnt so far

package main

import(
	"fmt"
	"log"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"io/ioutil"
	"fyne.io/fyne/v2/storage"
	"strings"
)

type our_app struct{
	Edit *widget.Entry
	Preview *widget.RichText
	CurrentFile fyne.URI
	SaveMenuItem *fyne.MenuItem
}

var cfg our_app

func main(){
		// create the app
		a:=app.New()

		 //create the window
		 w:=a.NewWindow("Markdown Editor")

		 //generate the UI

		 edit,preview:=cfg.MakeUI()

		 // Set the contents

		 w.SetContent(container.NewHSplit(edit,preview))

		 //create the menu

		 openMenuItem:=fyne.NewMenuItem("Open",cfg.openFunc(w))
		 saveMenuItem:=fyne.NewMenuItem("Save",cfg.saveFunction(w))
		 saveAsMenuItem:=fyne.NewMenuItem("Saveas",cfg.saveasfunc(w))

		 fileMenu:=fyne.NewMenu("File",openMenuItem,saveMenuItem,saveAsMenuItem)

		 menu:=fyne.NewMainMenu(fileMenu)

		 w.SetMainMenu(menu)

		 w.Resize(fyne.Size{Width:500,Height:500})


		 //display the window

		 w.ShowAndRun()
}

func(cfg *our_app) MakeUI() (*widget.Entry, *widget.RichText){
		edit:=widget.NewMultiLineEntry()
		preview:=widget.NewRichTextFromMarkdown("")
		edit.OnChanged=preview.ParseMarkdown
		cfg.Edit=edit
		cfg.Preview=preview
		return edit,preview
}

var filter=storage.NewExtensionFileFilter([]string{".md",".MD"})

func(cfg *our_app) saveasfunc(window fyne.Window) func(){ //this part needs to be understood a bi
  return func(){
	dialog:=dialog.NewFileSave(func(file fyne.URIWriteCloser, err error){
		if err!=nil{
			log.Println("Error Saving the file:",err)
			return
		}

		if file==nil{
			fmt.Println("Cancelled")
			return
		}

		if !strings.HasSuffix(strings.ToLower(file.URI().String()),".md"){
			dialog.ShowInformation("Error","Please use .md extension",window)
			return
		}

		defer file.Close()

		_,err= file.Write([]byte(cfg.Edit.Text))
		if err!=nil{
			log.Println("Error while writing to the file:",err)
			return
		}

		fmt.Println("Saved Successfully")
		window.SetTitle("Markdown "+" - "+file.URI().Name())
		cfg.CurrentFile=file.URI()
	//fmt.Println("Here Here :",cfg)

		//cfg.SaveMenuItem.Disabled=false
	},window)

	dialog.SetFileName("example.md")
	dialog.SetFilter(filter)

	dialog.Show()
}

}

func(cfg *our_app) saveFunction(win fyne.Window) func(){
	return func(){
		if cfg.CurrentFile!=nil{
			write,err :=storage.Writer(cfg.CurrentFile)
			if err!=nil{
				dialog.ShowError(err,win)
				return
			}

			write.Write([]byte(cfg.Edit.Text))
			defer write.Close()
		} else if cfg.CurrentFile==nil{
			fmt.Println(cfg)
			dialog.ShowInformation("File Havent been created yet","Please create the file using the saveas option first",win)
			return
		}
	}
}




func (cfg *our_app) openFunc(win fyne.Window) func() {
	return func() {
		openDialog := dialog.NewFileOpen(func(read fyne.URIReadCloser, err error) {
			if err != nil {
				dialog.ShowError(err, win)
				return
			}

			if read == nil {
				return
			}

			defer read.Close()

			data, err := ioutil.ReadAll(read)
			if err != nil {
				dialog.ShowError(err, win)
				return
			}

			cfg.Edit.SetText(string(data))

			cfg.CurrentFile = read.URI()
			//fmt.Println("Here:",cfg)
			win.SetTitle("Markdown " + " - " + read.URI().Name())
			//cfg.SaveMenuItem.Disabled = false

		}, win)

		openDialog.SetFilter(filter)

		openDialog.Show()
	}
}