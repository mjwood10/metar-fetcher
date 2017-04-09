# metar-fetcher

Simple program to illustrate concurrency and reading JSON data from a URL.

The program takes a list of airport codes as command line arguments and concurrenlty fetches the current weather
conditions (METARs) for each airport from https://avwx.rest/ and prints the decoded information to STDOUT.

Example:

~~~~
$ go run main.go kdfw ksea kmem phnl
..................
Fetched: KDFW in 0.91s

Fetched: KSEA in 0.91s

Fetched: PHNL in 0.91s
.......................................
Fetched: KMEM in 3.02s
All stations fetched in 3.02s

Station:        KDFW --  Dallas, TX -- Dallas-Fort Worth International Airport
Time:           270053Z
Temperature:    80.6°F / 27.0°C
Dew Point:      59.0°F / 15.0°C
Wind:           160° (SSE) @22KT Gusts to 30KT
Conditions:     thunderstorm
Cloud Layers:   FEW @5000FT (cumulonimbus) -- SCT @7500FT -- BKN @17000FT -- BKN @28000FT
Visibility:     10sm
Pressure:       29.67inHg
Flight Rules:   VFR
Raw Report:     KDFW 270053Z COR 16022G30KT 10SM TS FEW050CB SCT075 BKN170 BKN280 27/15 A2967 RMK AO2 PK WND 15030/0053 LTG DSNT W-NE TS
B29 SLP040 FRQ LTGICCG NW-N TS NW-N MOV NE T02670150


Station:        KSEA --  Seattle, WA -- Seattle-Tacoma International Airport
Time:           270053Z
Temperature:    48.2°F / 9.0°C
Dew Point:      44.6°F / 7.0°C
Wind:           120° (ESE) @6KT
Cloud Layers:   BKN @3000FT -- OVC @3600FT
Visibility:     10sm
Pressure:       29.74inHg
Flight Rules:   MVFR
Raw Report:     KSEA 270053Z 12006KT 10SM BKN030 OVC036 09/07 A2974 RMK AO2 SLP080 T00890067


Station:        KMEM --  Memphis, TN -- International Airport
Time:           270054Z
Temperature:    69.8°F / 21.0°C
Dew Point:      53.6°F / 12.0°C
Wind:           160° (SSE) @5KT
Cloud Layers:   FEW @20000FT -- SCT @25000FT
Visibility:     10sm
Pressure:       29.94inHg
Flight Rules:   VFR
Raw Report:     KMEM 270054Z 16005KT 10SM FEW200 SCT250 21/12 A2994 RMK AO2 SLP137 T02060117


Station:        PHNL --  Honolulu, HI -- International Airport
Time:           270053Z
Temperature:    84.2°F / 29.0°C
Dew Point:      60.8°F / 16.0°C
Wind:           050° (NE) @18KT Gusts to 25KT
Cloud Layers:   FEW @2600FT -- SCT @5000FT -- SCT @6000FT -- BKN @20000FT
Visibility:     10sm
Pressure:       30.11inHg
Flight Rules:   VFR
Raw Report:     PHNL 270053Z 05018G25KT 10SM FEW026 SCT050 SCT060 BKN200 29/16 A3011 RMK AO2 SLP194 T02890161 $
~~~~
