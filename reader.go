package elecreader

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/tarm/serial"
)

const (
	logPkg string = "reader"
)

type Connection struct {
	Port *serial.Port
}

func Extract(answer string) (Result, error) {
	res := Result{}
	lines := strings.Split(answer, "\n")
	if len(lines) < 2 {
		lines = strings.Split(answer, "\r")
	}
	for _, line := range lines {
		log.Tracef("line: %s", line)
		if strings.Contains(line, "1.8.0") {
			res.TotalActiveEnergyImport, _ = ExtractNumber(line)
		} else if strings.Contains(line, "1.8.1") {
			res.ActiveEnergyImportRate1, _ = ExtractNumber(line)
		} else if strings.Contains(line, "1.8.2") {
			res.ActiveEnergyImportRate2, _ = ExtractNumber(line)
		} else if strings.Contains(line, "2.8.0") && !strings.Contains(line, "82.8.0") {
			res.TotalActiveEnergyExport, _ = ExtractNumber(line)
		} else if strings.Contains(line, "2.8.1") && !strings.Contains(line, "82.8.1") {
			res.ActiveEnergyExportRate1, _ = ExtractNumber(line)
		} else if strings.Contains(line, "2.8.2") && !strings.Contains(line, "82.8.2") {
			res.ActiveEnergyExportRate2, _ = ExtractNumber(line)
		} else if strings.Contains(line, "15.8.0") {
			res.ActiveEnergyAbsolute, _ = ExtractNumber(line)
		} else if strings.Contains(line, "31.7") {
			res.CurrentLine1, _ = ExtractNumber(line)
		} else if strings.Contains(line, "51.7") {
			res.CurrentLine2, _ = ExtractNumber(line)
		} else if strings.Contains(line, "71.7") {
			res.CurrentLine3, _ = ExtractNumber(line)
		} else if strings.Contains(line, "32.7") {
			res.VoltageLine1, _ = ExtractNumber(line)
		} else if strings.Contains(line, "52.7") {
			res.VoltageLine2, _ = ExtractNumber(line)
		} else if strings.Contains(line, "72.7") {
			res.VoltageLine3, _ = ExtractNumber(line)
		} else {
			// todo
		}
	}
	log.Infof("Result: %+v", res)
	return res, nil
}

func ExtractNumber(line string) (float64, error) {
	log.Infof("Line: %s", line)
	startIndex := strings.Index(line, "(") + 1
	endIndex := strings.Index(line, "*")
	numberStr := line[startIndex:endIndex]
	log.Tracef("%s", numberStr)
	nr, err := strconv.ParseFloat(numberStr, 64)
	if err != nil {
		log.Errorf("Failed to parse line %s: %v", line, err)
	}
	return nr, nil
}

func (conn *Connection) Read() (string, error) {
	// Opening string: /?!<CR><LF> (in hex: 2F 3F 21 0D 0A)
	msg := "/?!\x0D\x0A"
	resp, err := conn.Send([]byte(msg), []byte("\x0A"))
	if err != nil {
		log.Errorf("Failed to send: %v", err)
	}
	log.Tracef("%X", resp)
	log.Infof("%s", resp)

	// Acknowledgement; Z: Transmission rate <ACK>000<CR><LF>
	msg = "\x06000\x0D\x0A"
	endOfTransmition := []byte("\x03") // <ETX>
	resp, err = conn.Send([]byte(msg), endOfTransmition)
	if err != nil {
		log.Errorf("Failed to Acknowledge (000) %v", err)
		time.Sleep(time.Second * 1)
		return conn.Read()
	}

	log.Tracef("Response was %s\n", resp)
	return string(resp), nil
}

func OpenPort(config *serial.Config) (*serial.Port, error) {
	openedPort, err := serial.OpenPort(config)
	if err != nil {
		log.WithField("package", logPkg).Errorf("Failed to open: %s", err)
		return openedPort, err
	}
	log.WithField("package", logPkg).Tracef("Open RS232 '%s' with Baud: %v, Size: %v, Parity: %v, StopBits: %v \n",
		config.Name, config.Baud, config.Size, config.Parity, config.StopBits)
	return openedPort, nil
}

func (conn *Connection) ClosePort() error {
	if conn.Port == nil {
		log.WithField("package", logPkg).Errorf("Port is nil")
		return fmt.Errorf("port is nil")
	}
	err := conn.Port.Close()
	if err != nil {
		log.WithField("package", logPkg).Errorf("Failed to close port: %s ", err)
		return err
	}
	log.WithField("package", logPkg).Traceln("RS232 port closed")
	return nil
}

func (conn *Connection) Send(msg []byte, endBytes []byte) ([]byte, error) {
	if conn.Port == nil {
		log.WithField("package", logPkg).Errorln("Port not opened")
		return []byte{}, fmt.Errorf("port not opened")
	}

	log.WithField("package", logPkg).Tracef("Command Text %s Hex: %X", msg, msg)
	_, err := conn.Port.Write(msg)
	if err != nil {
		log.WithField("package", logPkg).Error(err)
		return []byte{}, err
	}

	buf := make([]byte, 1024)
	var resp []byte
	length := 0
	n := 1
	endTime := time.Now().Add(time.Minute * 2)
	for n > 0 {
		n, err = conn.Port.Read(buf)
		if err != nil {
			log.WithField("package", logPkg).Error(err)
			return resp, err
		}
		length += n
		resp = append(resp, buf[:n]...)
		if bytes.Contains(buf[:n], endBytes) {
			log.WithField("package", logPkg).Traceln("Found End Byte")
			break
		}
		if time.Now().After(endTime) {
			log.WithField("package", logPkg).Errorf("Timeout while reading. Maybe there is noise? %s ", resp)
			return resp[:length], fmt.Errorf("time out while reading. Maybe there is noise? %s ", resp)
		}
	}

	log.WithField("package", logPkg).Tracef("Length of response %v", length)
	log.WithField("package", logPkg).Tracef("Response in Text: '%s', Hex: %X ", resp[:length], resp[:length])
	return resp[:length], nil
}
