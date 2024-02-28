package main

import (
	"os"
	"testing"
)

func TestMain(t *testing.T) {
	// Se asume que tienes un directorio de prueba llamado "test_directory" con archivos válidos
	directory := "C:/ZONAS/TRUORA/0 PRUEBA/enron_mail_20110402/maildir/allen-p/"

	// Ejecutar la función main con el directorio de prueba
	os.Args = []string{"main", directory}
	main()
}
