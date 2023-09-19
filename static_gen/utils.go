package static_gen

import (
	"io"
	"os"
	"path"
	"path/filepath"
	"regexp"
)

func (sG StaticGen) getAbsPath(relPath string) string {
	articlePath := path.Join(
		sG.config.OutputPaths.Root,
		relPath,
	)
	return articlePath
}

func (sG StaticGen) getAbsUrl(relPath string) string {
	articleUrl := sG.config.RootUrl + "/" + relPath
	return articleUrl
}

func urlSafeName(s string) string {
	specialChars := regexp.MustCompile(`([^A-z0-9])`)
	s = specialChars.ReplaceAllString(s, `_`)
	return s
}

func copyFile(srcPath string, destPath string) error {
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	fMode := os.FileMode(0777)
	dir := filepath.Dir(destPath)
	if err := os.MkdirAll(dir, fMode); err != nil {
		return err
	}
	defer srcFile.Close()
	dstFile, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer dstFile.Close()
	_, err = io.Copy(dstFile, srcFile)
	return err
}

func writeFile(filePath string, data []byte) error {
	fMode := os.FileMode(0777)
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, fMode); err != nil {
		return err
	}
	dstFile, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer dstFile.Close()
	_, err = dstFile.Write(data)
	return err
}
