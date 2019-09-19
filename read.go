package main

import (
	"fmt"
	"github.com/airmap/interfaces/src/go/tracking"
)

func ReadTrackingData(done chan<- bool, dataSource tracking.Collector_ConnectProcessorClient, bounds PositionBounds) {
	data := make(chan *tracking.Track, 30)
	go ProcessTrackingData(data, bounds)
	for {
		tracks, err := extractTracks(dataSource)
		if err != nil {
			fmt.Println(err)
			continue
		}
		for _, track := range tracks {
			data <- track
		}
	}
}

func extractTracks(dataSource tracking.Collector_ConnectProcessorClient) ([]*tracking.Track, error) {
	update, err := dataSource.Recv()
	if err != nil {
		return nil, err
	}
	return update.GetBatch().Tracks, nil
}
