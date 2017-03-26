# metar-fetcher

Simple program to illustrate concurrency and reading JSON data from a URL.

The program takes a list of airport codes as command line arguments and concurrenlty fetches the current weather
conditions (METARs) for each airport from https://avwx.rest/ and outputs the decoded information to to STDOUT.

Example:

~~~~
$ go run main.go ksea kmem phnl
........
Fetched: KSEA in 0.41s
.
Fetched: PHNL in 0.43s
...
Fetched: KMEM in 0.60s
All stations fetched in 0.60s

Station:	KSEA --  Seattle, WA -- Seattle-Tacoma International Airport
Time:     	262309Z
Temperature:	48.2°F
Dew Point:	44.6°F
Wind:     	100° @ 05KT
Conditions:	light rain
Cloud Layers:	SCT @ 800FT	BKN @ 3200FT	OVC @ 3900FT
Visibility:	10sm
Pressure:	29.75inHg
Flight Rules:	VFR
Raw Report:	KSEA 262309Z 10005KT 10SM -RA SCT008 BKN032 OVC039 09/07 A2975 RMK AO2 P0001 T00890072


Station:	KMEM --  Memphis, TN -- International Airport
Time:     	262254Z
Temperature:	73.4°F
Dew Point:	51.8°F
Wind:     	250° @ 08KT
Cloud Layers:	FEW @ 5000FT	SCT @ 25000FT
Visibility:	10sm
Pressure:	29.93inHg
Flight Rules:	VFR
Raw Report:	KMEM 262254Z 25008KT 10SM FEW050 SCT250 23/11 A2993 RMK AO2 SLP132 T02280111


Station:	PHNL --  Honolulu, HI -- International Airport
Time:     	262253Z
Temperature:	84.2°F
Dew Point:	64.4°F
Wind:     	060° @ 13KT Gusts to 24KT
Cloud Layers:	FEW @ 2500FT	SCT @ 3500FT	SCT @ 5000FT
Visibility:	10sm
Pressure:	30.15inHg
Flight Rules:	VFR
Raw Report:	PHNL 262253Z 06013G24KT 10SM FEW025 SCT035 SCT050 29/18 A3015 RMK AO2 SLP209 VCSH NE T02890178 $
~~~~
