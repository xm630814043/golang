package pdf

import (
	"log"

	"github.com/signintech/gopdf"
)

func GeneratePdf() {
	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})
	pdf.AddPage()
	err := pdf.AddTTFFont("wts11", "NotoSansSC-Regular.ttf")
	if err != nil {
		log.Print(err.Error())
		return
	}

	err = pdf.SetFont("wts11", "", 14)
	if err != nil {
		log.Print(err.Error())
		return
	}
	err = pdf.Cell(nil, "您好")
	if err != nil {
		return
	}
	err = pdf.WritePdf("text.pdf")
	if err != nil {
		return
	}
}
