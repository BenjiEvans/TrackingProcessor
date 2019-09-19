package main

import (
	"fmt"
	"github.com/airmap/interfaces/src/go/measurements"
	"github.com/airmap/interfaces/src/go/tracking"
)

type TrackId_StatTracker struct {
	UniqueCount  int
	MaxCount     int
	MostFrequent string
	TrackIdCount map[string]int
}

func ProcessTrackingData(data chan *tracking.Track, bounds PositionBounds) {
	trackCount := 0
	var lat float64 = 0
	var long float64 = 0
	var velocityX float64 = 0
	var velocityY float64 = 0
	trackIdStats := TrackId_StatTracker{
		TrackIdCount: map[string]int{},
	}
	for {
		select {
		case track := <-data:
			if !isInRange(track.Position, bounds) {
				continue
			}
			processPositionData(track.Position, &lat, &long)
			processVelocityData(track, &velocityX, &velocityY)
			processTrackIds(track, &trackIdStats)
			trackCount++
			if (trackCount % 10) == 0 {
				fmt.Printf("Avg Lat: %f, Avg Long: %f\n\n", (lat / float64(trackCount)), (long / float64(trackCount)))
				fmt.Printf("Avg X Velocity: %f, Avg Y Velocity: %f\n\n", (velocityX / float64(trackCount)), (velocityY / float64(trackCount)))
				fmt.Printf("Unique Flight Count: %d\n\n", trackIdStats.UniqueCount)
				fmt.Println("Most Frequent: ", trackIdStats.MostFrequent, " Appeared: ", trackIdStats.MaxCount)
				fmt.Println("_____________________________________________________________")
				lat = 0
				long = 0
				velocityX = 0
				velocityY = 0
				trackCount = 0
			}
		default:
		}
	}
}

func processPositionData(position *measurements.Position, lat *float64, long *float64) {
	absolutePosition := position.GetAbsolute().Coordinate
	if absolutePosition != nil {
		/*
			For simplicity sake I'm going to assume that
			absolutePosition.Latitude and absolutePosition.Longitude are not null when absolutePosition is
			not null. Have not ran into any issues with that assumption
			so far but should still check whether absolutePosition.Longitude and
			absolutePosition.Latitude are valid pointers before accessing them
		*/
		*lat += absolutePosition.Latitude.Value
		*long += absolutePosition.Longitude.Value
	}
}

func processVelocityData(track *tracking.Track, velocityX *float64, velocityY *float64) {
	if track.Velocity == nil {
		return
	}

	velocity := track.Velocity.GetCartesian()
	if velocity != nil {
		/*
			For simplicity sake I'm going to assume that
			velocity.X and velocity.Y are not null when velocity is
			not null. Have not ran into any issues with this assumption
			so far but should still checkout that I'm accessing valid
			pointers
		*/
		*velocityX += velocity.X.Value
		*velocityY += velocity.Y.Value
	}
}

func processTrackIds(track *tracking.Track, stats *TrackId_StatTracker) {
	for _, id := range track.Identities {
		trackId := id.GetTrackId()
		if trackId == nil {
			continue
		}
		fmt.Println("Track ID: ", trackId.AsString)
		val := stats.TrackIdCount[trackId.AsString]
		if val == 0 {
			stats.UniqueCount++
		}
		updateAmount := val + 1
		if updateAmount > stats.MaxCount {
			stats.MaxCount = updateAmount
			stats.MostFrequent = trackId.AsString
		}
		stats.TrackIdCount[trackId.AsString] = updateAmount
	}
}

func isInRange(position *measurements.Position, bounds PositionBounds) bool {
	if position == nil {
		return false
	}
	absolutePosition := position.GetAbsolute().Coordinate
	if absolutePosition != nil {
		/*
			For simplicity sake I'm going to assume that
			absolutePosition.Latitude and absolutePosition.Longitude are not null when absolutePosition is
			not null. Have not ran into any issues with that assumption
			so far but should still check whether absolutePosition.Longitude and
			absolutePosition.Latitude are valid pointers before accessing them
		*/

		isInLatRange := absolutePosition.Latitude.Value >= bounds.MinLat && absolutePosition.Latitude.Value <= bounds.MaxLat
		isInLongRange := absolutePosition.Longitude.Value >= bounds.MinLong && absolutePosition.Longitude.Value <= bounds.MaxLong
		return isInLatRange && isInLongRange
	}

	return false

}
