package handlers

import (
	"io"
	"main/internal/models"
	"main/internal/utils"
	"net/http"
	"os"

	"github.com/jordan-wright/email"
	"net/smtp"
)

func SendEmailHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		utils.SendErrorResponse(w, "Unable to process form", http.StatusBadRequest)
		return
	}

	from := r.FormValue("email")
	message := r.FormValue("message")

	e := email.NewEmail()
	e.From = from
	e.To = []string{"kh.abukhassym@gmail.com"}
	e.Subject = "Support Request"
	e.Text = []byte(message)

	file, header, err := r.FormFile("attachment")
	if err == nil {
		defer file.Close()
		attachmentPath := "./" + header.Filename
		out, err := os.Create(attachmentPath)
		if err != nil {
			utils.SendErrorResponse(w, "Error saving attachment", http.StatusInternalServerError)
			return
		}
		defer out.Close()

		if _, err = io.Copy(out, file); err != nil {
			utils.SendErrorResponse(w, "Error processing file", http.StatusInternalServerError)
			return
		}

		e.AttachFile(attachmentPath)
		defer os.Remove(attachmentPath)
	}

	auth := smtp.PlainAuth("", "kh.abukhassym@gmail.com", "bsml lwzy akas ezfm", "smtp.gmail.com")
	if err := e.Send("smtp.gmail.com:587", auth); err != nil {
		utils.SendErrorResponse(w, "Error sending email", http.StatusInternalServerError)
		return
	}

	utils.SendJSONResponse(w, http.StatusOK, models.ResponseData{
		Status:  "success",
		Message: "Email sent successfully",
	})
}
