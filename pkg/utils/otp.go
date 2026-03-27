package utils

import (
	"backend/config"
	"backend/internal/dto"
	"bytes"
	"crypto/rand"
	"fmt"
	"html/template"
	"net/smtp"
	"time"
)

const OTPTTL = 10 * time.Minute

type CompanyEmailData struct {
	CompanyName      string
	ActivePackage    string
	TrialStart       string
	TrialEnd         string
	AvailablePackages []struct {
		Name        string
		DurationDays int
	}
	LogoURL string
	Year    int
}

func SendEmail(to, subject, body string) error {
	cfg := config.LoadConfig().SMTP

	if cfg.Host == "" || cfg.Port == "" || cfg.User == "" || cfg.Pass == "" {
		return fmt.Errorf("SMTP configuration is missing")
	}

	addr := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)
	auth := smtp.PlainAuth("", cfg.User, cfg.Pass, cfg.Host)

	msg := fmt.Sprintf(
		"From: %s\r\n"+
			"To: %s\r\n"+
			"Subject: %s\r\n"+
			"MIME-Version: 1.0\r\n"+
			"Content-Type: text/html; charset=\"UTF-8\"\r\n\r\n"+
			"%s",
		cfg.User, to, subject, body,
	)

	if err := smtp.SendMail(addr, auth, cfg.User, []string{to}, []byte(msg)); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

func SendOTPEmail(email, otp string, ttl time.Duration) error {
	cfg := config.LoadConfig()

	templatePath := fmt.Sprintf("%s/otp.html", cfg.App.TemplatePath)

	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return fmt.Errorf("parse template otp: %w", err)
	}

	logoURL := "http://localhost:5173/src/assets/HumadifySecondary.png"
	data := struct {
		OTP     string
		TTL     int
		Year    int
		LogoURL string
	}{
		OTP:     otp,
		TTL:     int(ttl.Minutes()),
		Year:    time.Now().Year(),
		LogoURL: logoURL,
	}

	var body bytes.Buffer
	if err := tmpl.Execute(&body, data); err != nil {
		return fmt.Errorf("execute otp template: %w", err)
	}

	subject := "Kode OTP Verifikasi Anda"

	if err := SendEmail(email, subject, body.String()); err != nil {
		return fmt.Errorf("send otp email: %w", err)
	}

	return nil
}

func GenerateOTP(length int) (string, error) {
	const digits = "0123456789"
	otp := make([]byte, length)
	_, err := rand.Read(otp)
	if err != nil {
		return "", err
	}
	for i := 0; i < 6; i++ {
		otp[i] = digits[int(otp[i])%len(digits)]
	}
	return string(otp), nil
}


func SendCompanyNotificationEmail(company dto.CompanyInfo, recipientEmail string) error {
	tmpl, err := template.ParseFiles("internal/template/companyNotification.html")
	if err != nil {
		return fmt.Errorf("parse company template: %w", err)
	}

	start, err := time.Parse("2006-01-02", company.StartDate)
	if err != nil {
		start = time.Now()
	}
	end, err := time.Parse("2006-01-02", company.EndDate)
	if err != nil {
		end = time.Now().AddDate(0, 1, 0)
	}

	data := CompanyEmailData{
		CompanyName:   company.CompanyName,
		ActivePackage: company.PackageName,
		TrialStart:    start.Format("2006-01-02"),
		TrialEnd:      end.Format("2006-01-02"),
		Year:          time.Now().Year(),
		AvailablePackages: []struct {
			Name        string
			DurationDays int
		}{
			{"Basic", 30},
			{"Pro", 90},
			{"Enterprise", 180},
		},
	}

	var body bytes.Buffer
	if err := tmpl.Execute(&body, data); err != nil {
		return fmt.Errorf("execute company template: %w", err)
	}

	subject := "Selamat! Profil Perusahaan Anda Telah Dibuat"

	if err := SendEmail(recipientEmail, subject, body.String()); err != nil {
		return fmt.Errorf("send company email: %w", err)
	}

	return nil
}

