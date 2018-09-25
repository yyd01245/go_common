package common

import(
	"fmt"
	"os"
	"errors"
	"io/ioutil"
	"strings"
	"strconv"
	"syscall"
	log "github.com/Sirupsen/logrus"

)

func CheckFileIsExist(filename string) bool{
	var exist =true
	if _,err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	log.Debug("file:",filename," exist=",exist)

	return exist
}

func WriteListLineToFile(filename string, paramList[]string) error{

	if len(paramList) == 0 {
		log.Warnf("write firle paramList len =0")
		return errors.New("write firle paramList len =0")
	}
	file, err := os.OpenFile(filename,os.O_CREATE|os.O_TRUNC|os.O_WRONLY,0644)
	if err != nil {
		log.Warnf("open file failed !")
		return err
	}
	defer file.Close()

	for k,v := range paramList {
		log.Debug("key=",k," value=",v)
		text := fmt.Sprintf("%s\n",v)
		log.Debug("ready to write file=",filename," text=",text)
		_,err := file.WriteString(text)
		if err != nil {
			return err
		}
	}
	return nil
}

func AppendListLineToFile(filename string, paramList[]string) error{

	if len(paramList) == 0 {
		log.Warnf("write firle paramList len =0")
		return errors.New("write firle paramList len =0")
	}
	file, err := os.OpenFile(filename,os.O_CREATE|os.O_RDWR|os.O_APPEND,0644)
	if err != nil {
		log.Warnf("open file failed !")
		return err
	}
	defer file.Close()

	for k,v := range paramList {
		log.Debug("key=",k," value=",v)
		text := fmt.Sprintf("%s\n",v)
		log.Debug("ready to write file=",filename," text=",text)
		_,err := file.WriteString(text)
		if err != nil {
			log.Errorf("write string:%v Error:%v",text,err)
			return err
		}
	}
	return nil
}

func ReadFileAll(filename string) (string,error) {

	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return "",err
	}
	return string(b),nil
}

func CreateDir(path string) error {

	return os.MkdirAll(path,0744)

}

func WritePidToFile(filename string) error {
	file, err := os.OpenFile(filename,os.O_CREATE|os.O_TRUNC|os.O_WRONLY,0644)
	if err != nil {
		log.Warnf("open file failed !")
		return err
	}
	defer file.Close()
	txt := fmt.Sprintf("pid=%d",os.Getpid())
	log.Infof(txt)
	file.WriteString(txt)
	return nil
}

// CheckPidFromFile return nil success, error no pid
func CheckPidFromFile(filename string) error {
	file, err := os.OpenFile(filename,os.O_RDONLY,0644)
	if err != nil {
		log.Warnf("open file failed !")
		return err
	}
	defer file.Close()
	body, err := ioutil.ReadAll(file)
	if err != nil {
		log.Warnf("ReadAll", err)
			return err
	}
	data := strings.Split(string(body),"=")
	log.Infof("get pid file: %v",data)
	pid,err := strconv.Atoi(data[1])
	log.Infof("get pid = %d",pid)

	return FindProcess(pid)
}

func FindProcess(pid int) error{
	err := syscall.Kill(pid, 0);
	if err == nil {
		log.Infof("find process success ")
    return nil
	}else {
		log.Infof("Failed to find process: %v\n", err)
		return err	
	}
	

	// _, err := os.FindProcess(pid)
	// if err != nil {
	// 		log.Infof("Failed to find process: %v\n", err)
	// 		return err
	// } else {
	// 		log.Infof("find process success ")
	// }
	// return nil
}