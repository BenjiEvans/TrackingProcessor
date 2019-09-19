package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/airmap/interfaces/src/go/tracking"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"math"
	"os"
	"strconv"
	"strings"
)

func main() {
	creds, err := credentials.NewClientTLSFromFile("cert.txt", "")
	fatalIfError(err, "Failed to read cert")
	cc, err := grpc.Dial("api.airmap.com:443", grpc.WithTransportCredentials(creds))
	fatalIfError(err, "Error on grpc dial")
	client := tracking.NewCollectorClient(cc)
	trackingDataProvider, err := client.ConnectProcessor(context.Background())
	fatalIfError(err, "Error connecting processor")
	bounds, err := generatePositionBounds()
	fatalIfError(err, "Error when parsing command line arguments")
	doneReading := make(chan bool)
	go ReadTrackingData(doneReading, trackingDataProvider, bounds)
	<-doneReading
}

func fatalIfError(err error, failMessage string) {
	if err != nil {
		fmt.Println(err)
		fmt.Println(failMessage)
		os.Exit(1)
	}
}

type PositionBounds struct {
	MaxLong float64
	MinLong float64
	MaxLat  float64
	MinLat  float64
}

func generatePositionBounds() (PositionBounds, error) {
	bounds := PositionBounds{
		MaxLong: math.MaxFloat64,
		MinLong: -math.MaxFloat64,
		MaxLat:  math.MaxFloat64,
		MinLat:  -math.MaxFloat64,
	}
	latBounds := flag.String("lat", "", "")
	longBounds := flag.String("long", "", "")
	flag.Parse()
	if *latBounds != "" {
		interval := strings.Split(*latBounds, ",")
		v1, err := strconv.ParseFloat(interval[0], 64)
		if err != nil {
			return bounds, err
		}
		v2, err := strconv.ParseFloat(interval[1], 64)
		if err != nil {
			return bounds, err
		}
		bounds.MaxLat = math.Max(v1, v2)
		bounds.MinLat = math.Min(v1, v2)
	}
	if *longBounds != "" {
		interval := strings.Split(*longBounds, ",")
		v1, err := strconv.ParseFloat(interval[0], 64)
		if err != nil {
			return bounds, err
		}
		v2, err := strconv.ParseFloat(interval[1], 64)
		if err != nil {
			return bounds, err
		}
		bounds.MaxLong = math.Max(v1, v2)
		bounds.MinLong = math.Min(v1, v2)
	}

	return bounds, nil
}
