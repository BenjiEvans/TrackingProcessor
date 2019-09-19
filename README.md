# Description 
Processor client that reads from a GRPC stream and prints out the following stats:
- Rolling average x and y velocities every 10 tracks 
- Rolling average latitude and longitude every 10 tracks
- Number of unique trackIds (printed every 10 tracks)
- Name of the most frequent trackId and how many times it has occured (printed every 10 tracks)

Note: each trackId that is processed is printed 

Console Output Format Example:

Track ID:  954<br/>
Track ID:  3868<br/>
Track ID:  131<br/>
Track ID:  2621<br/>
Track ID:  811<br/>
Track ID:  3762<br/>
Track ID:  CKS421-1129<br/>
Track ID:  GTI157-2527<br/>
Track ID:  flight|mL7MENBxOkBzmSX5dn0yYv0ZGBiR3KLW3g2AZX6IM303yAbeEEmn<br/>
Track ID:  UAL815-2578<br/>
Avg Lat: 43.371634, Avg Long: -30.149079

Avg X Velocity: -33.800000, Avg Y Velocity: -16.300000

Unique Flight Count: 10

Most Frequent:  SWA354-1560  Appeared:  1
__________________________________________________________

# Building Processor
Execute the following command in the root directory of the project

`go build`

# Runing the Processor
./Processor 

Note: Processor takes 2 optional command line args --lat --long that allow the user to bound/filter what tracks are processed. If not specified, there is no bounds

`--lat`: determines the latitude range for a track to be processed. If the lattitude of a track does not fall in the range it will not be processed 

Example: ./Processor --lat=100,-100

`--long`: determines the longitude range for a track to be processed. If the longitude of a track does not fall in the range it will not be processed 

Example: ./Processor --long=-100,100