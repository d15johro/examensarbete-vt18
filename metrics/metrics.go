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
	DataSize            int
	Filepath            string
}

func (m *Metrics) Log() {
	s := fmt.Sprintf("%d,%f,%f,%f,%f,%f,%d\n",
		m.ID,
		m.AccessTime,
		m.ResponseTime,
		m.SerializationTime,
		m.DeserializationTime,
		m.StructuringTime,
		m.DataSize)
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
	_, err := os.Stat(m.Filepath)
	if err != nil {
		if !os.IsNotExist(err) {
			if err := os.Remove(m.Filepath); err != nil {
				return err
			}
		}
	}
	_, err = os.Create(m.Filepath)
	return err
}
