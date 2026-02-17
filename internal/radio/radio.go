package radio

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gopxl/beep"
	"github.com/gopxl/beep/mp3"
	"github.com/gopxl/beep/speaker"
	"github.com/maker2413/shellpod/internal/config"
	"github.com/maker2413/shellpod/internal/icyreader"
)

type Radio struct {
	SampleRate beep.SampleRate
	Streamer   beep.StreamSeekCloser
	Stations   []config.Station
	TitleChan  <-chan string
	Config     config.Config
	Resp       *http.Response
}

func NewRadio(c config.Config) (Radio, error) {
	client := &http.Client{Timeout: 0}
	req, err := http.NewRequest("GET", c.Stations[0].StreamURL, nil)
	if err != nil {
		return Radio{}, err
	}
	req.Header.Set("Icy-MetaData", "1")

	resp, err := client.Do(req)
	if err != nil {
		return Radio{}, fmt.Errorf("Failed to connect: %v", err)
	}

	// Check if the server actually supports ICY metadata
	icyIntStr := resp.Header.Get("icy-metaint")
	if icyIntStr == "" {
		return Radio{}, fmt.Errorf("Server did not return icy-metaint. This might not be a direct Icecast stream.")
	}

	// Get the interval from headers
	metaint, err := strconv.Atoi(resp.Header.Get("icy-metaint"))
	if err != nil {
		return Radio{}, err
	}

	reader := icyreader.NewIcyReader(resp.Body, metaint)
	titleChan := make(chan string, 10)
	reader.TitleChan = titleChan

	wrappedReader := icyreader.NewWrappedReader(reader, 32*1024) // 32KB buffer

	// We wrap in bufio to ensure the decoder gets enough data to identify the format
	streamer, format, err := mp3.Decode(wrappedReader)
	if err != nil {
		return Radio{}, fmt.Errorf("Failed to decode MP3: %v", err)
	}

	err = speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	if err != nil {
		return Radio{}, fmt.Errorf("Failed to initialize speaker: %v", err)
	}

	return Radio{
		format.SampleRate,
		streamer,
		c.Stations,
		titleChan,
		c,
		resp,
	}, nil
}
