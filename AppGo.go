package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"io"
	"io/ioutil"
	"path/filepath"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// Função para criptografar dados
func encrypt(data []byte) ([]byte, error) {
	key := []byte("chaveexemplo1234") // Chave fixa de 16 bytes
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	ciphertext := make([]byte, aes.BlockSize+len(data))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], data)

	return ciphertext, nil
}

// Função para descriptografar dados
func decrypt(data []byte) ([]byte, error) {
	key := []byte("chaveexemplo1234") // Chave fixa de 16 bytes
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	if len(data) < aes.BlockSize {
		return nil, errors.New("ciphertext too short")
	}
	iv := data[:aes.BlockSize]
	data = data[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(data, data)

	return data, nil
}

func main() {
	myApp := app.NewWithID("com.exemplo.meuapp")
	myApp.Settings().SetTheme(theme.DarkTheme())
	myWindow := myApp.NewWindow("Aplicativo em Golang")

	// Estilos e layout aprimorados
	title := widget.NewLabelWithStyle("Aplicativo de Criptografia em Go", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	title.TextStyle.Bold = true

	// Botão para criptografar
	encryptButton := widget.NewButtonWithIcon("Criptografar", theme.ConfirmIcon(), func() {
		dialog.ShowFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil || reader == nil {
				return
			}
			data, err := ioutil.ReadAll(reader)
			if err != nil {
				dialog.ShowError(err, myWindow)
				return
			}

			encryptedData, err := encrypt(data)
			if err != nil {
				dialog.ShowError(err, myWindow)
				return
			}

			// Caminho do arquivo criptografado
			originalFilePath := reader.URI().Path()
			extension := filepath.Ext(originalFilePath)
			newFilePath := strings.TrimSuffix(originalFilePath, extension) + "_criptografado" + extension

			err = ioutil.WriteFile(newFilePath, encryptedData, 0644)
			if err != nil {
				dialog.ShowError(err, myWindow)
				return
			}
			dialog.ShowInformation("Sucesso", "Arquivo criptografado com sucesso!", myWindow)
		}, myWindow)
	})

	// Botão para descriptografar
	decryptButton := widget.NewButtonWithIcon("Descriptografar", theme.ViewRefreshIcon(), func() {
		dialog.ShowFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil || reader == nil {
				return
			}
			data, err := ioutil.ReadAll(reader)
			if err != nil {
				dialog.ShowError(err, myWindow)
				return
			}

			decryptedData, err := decrypt(data)
			if err != nil {
				dialog.ShowError(err, myWindow)
				return
			}

			// Caminho do arquivo descriptografado
			originalFilePath := reader.URI().Path()
			extension := filepath.Ext(originalFilePath)
			newFilePath := strings.TrimSuffix(originalFilePath, extension) + "_decriptografado" + extension

			err = ioutil.WriteFile(newFilePath, decryptedData, 0644)
			if err != nil {
				dialog.ShowError(err, myWindow)
				return
			}
			dialog.ShowInformation("Sucesso", "Arquivo descriptografado com sucesso!", myWindow)
		}, myWindow)
	})

	// Layout estilizado
	buttons := container.NewHBox(encryptButton, decryptButton)
	content := container.NewVBox(
		title,
		widget.NewSeparator(),
		widget.NewLabel("Deseja:"),
		buttons,
	)

	myWindow.SetContent(content)
	myWindow.Resize(fyne.NewSize(400, 200))
	myWindow.ShowAndRun()
}
