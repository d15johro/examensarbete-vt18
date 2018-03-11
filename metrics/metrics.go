package metrics

import (
	"fmt"
	"log"
	"os"
)

func New() *Metrics {
	return &Metrics{}
}

type Metrics struct {
	ID                  uint32
	AccessTime          float64
	ResponseTime        float64
	SerializationTime   float64
	DeserializationTime float64
	StructuringTime     float64
	SerializedDataSize  int
	OriginalDataSize    uint64
	Filepath            string
}

func (m *Metrics) Log() {
	s := fmt.Sprintf("%d,%d,%d,%f,%f,%f,%f,%f\n",
		m.ID,
		m.OriginalDataSize,
		m.SerializedDataSize,
		m.AccessTime,
		m.ResponseTime,
		m.SerializationTime,
		m.StructuringTime,
		m.DeserializationTime)
	file, err := os.OpenFile(m.Filepath, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()
	if _, err = file.WriteString(s); err != nil {
		log.Fatalln(err)
	}
}

func (m *Metrics) Setup() error {
	// check if file exists. If it exists, remove it
	_, err := os.Stat(m.Filepath)
	if err != nil {
		if !os.IsNotExist(err) {
			if err := os.Remove(m.Filepath); err != nil {
				return err
			}
		}
	}
	// create file
	file, err := os.Create(m.Filepath)
	if err != nil {
		return err
	}
	defer file.Close()
	// write first line to file
	s := "ID,Orginal Datastorlek,Serialiserad Datastorlek,Accesstid,Svarstid,Serialiseringstid,Struktureringstid,Deserialiseringstid\n"
	if _, err = file.WriteString(s); err != nil {
		log.Fatalln(err)
	}
	return nil
}
