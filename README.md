## Bustime

This program scrapes vehicle locations from Clever Devices' bustime product https://www.cleverdevices.com/products/bustime/
which the New Orleans RTA uses to track their bus and streetcar fleet. It scrapes all the locations in the fleet and
puts them into a sqlite database.

### Running

First build:

```
go build cmd/bustime.go
```

Copy the example config file and set each property:

```
cp .config.example.yaml .config.yaml
```

I got the url and the key by performing a MITM on the mobile app. It's not encrypted. The full url is of this form:

```
http://<ip-address>/bustime/api/v3/getvehicles?key=<apikey>
```

Set `url` to the base url, and `key` to the value of the `key` query parameter.
Set `interval` to the time in seconds you want to wait before polling new fleet locations. 5 seconds
seems like a good middle ground.

Pass the config file as an arg to the program:

```
$ ./bustime .config.yaml                                                                                                            18:43:44
Found 115 vehicles
Found 115 vehicles
```

It will run forever unless you `ctrl-c` out or it crashes.

You'll see a file now in the current folder `vehicle_readings.db`. This is the sqlite database. You can run
queries and perform stats on it. Example:

```
$ sqlite3 vehicle_readings.db "select tmstmp, vid, lat, lon from vehicle_readings where vid = '276'"
1598566560|276|29.9677383333333|-90.088285
1598566620|276|29.9676466666667|-90.08864
1598566680|276|29.9676466666667|-90.08864
```


