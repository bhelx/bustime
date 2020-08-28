package lib

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const CreateVehiclesSQL = `
create table vehicle_readings (
	id integer not null primary key,
	vid varchar(5),
	lat double,
	lon double,
	des text,
	tmstmp DATETIME,
	srvtmstmp DATETIME,
	hdg varchar(5),
	rt varchar(5),
	tatripid varchar(16),
	tablockid varchar(5),
	oid varchar(5),
	rid varchar(5),
	pdist int,
	pid int,
	spd int,
	blk int,
	tripid int8,
	dly boolean,
	ori boolean
);

CREATE INDEX vid_tmstmp ON vehicle_readings(vid, tmstmp);
`

type VehicleTimestamp struct {
	time.Time
}

// UnmarshalJSON
// We need a special unmarshal method for this string timetstamp. it's of the
// form "YYYYMMDD hh:mm"
func (t *VehicleTimestamp) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	parts := strings.Split(s, " ")
	year, err := strconv.Atoi(parts[0][0:4])
	if err != nil {
		return err
	}
	month, err := strconv.Atoi(parts[0][4:6])
	if err != nil {
		return err
	}
	day, err := strconv.Atoi(parts[0][6:8])
	if err != nil {
		return err
	}
	timeParts := strings.Split(parts[1], ":")
	hour, err := strconv.Atoi(timeParts[0])
	if err != nil {
		return err
	}
	mins, err := strconv.Atoi(timeParts[1])
	if err != nil {
		return err
	}

	//fmt.Println(year, time.Month(month), day, hour, mins, 0, 0, time.Local)
	time := time.Date(year, time.Month(month), day, hour, mins, 0, 0, time.Local)
	t.Time = time

	return nil
}

// Vehicle represents an individual reading of a vehicle and it's location
// at that point in time
// Example:
// {
//   "vid": "155",
//   "tmstmp": "20200827 11:51",
//   "lat": "29.962149326173048",
//   "lon": "-90.05214051918121",
//   "hdg": "357",
//   "pid": 275,
//   "rt": "5",
//   "des": "Saratoga at Canal",
//   "pdist": 10122,
//   "dly": false,
//   "spd": 20,
//   "tatripid": "3130339",
//   "tablockid": "15",
//   "zone": "",
//   "srvtmstmp": "20200827 11:51",
//   "oid": "445",
//   "or": true,
//   "rid": "501",
//   "blk": 2102,
//   "tripid": 982856020
// }
type Vehicle struct {
	Vid        string           `json:"vid"`
	Tmstmp     VehicleTimestamp `json:"tmstmp"`
	SrvTimstmp VehicleTimestamp `json:"srvtmstmp"`
	Lat        float64          `json:"lat,string"`
	Lon        float64          `json:"lon,string"`
	Hdg        string           `json:"hdg"`
	Rt         string           `json:"rt"`
	Tatripid   string           `json:"tatripid"`
	Tablockid  string           `json:"tablockid"`
	Zone       string           `json:"zone"`
	Oid        string           `json:"oid"`
	Rid        string           `json:"rid"`
	Des        string           `json:"des"`
	Pdist      int              `json:"pdist"`
	Pid        int              `json:"pid"`
	Spd        int              `json:"spd"`
	Blk        int              `json:"blk"`
	Tripid     int              `json:"tripid"`
	Dly        bool             `json:"dly"`
	Or         bool             `json:"or"`
}

func (v *Vehicle) ToSql() ([]string, []string) {
	columns := []string{
		"vid",
		"lat",
		"lon",
		"des",
		"tmstmp",
		"srvtmstmp",
		"hdg",
		"rt",
		"tatripid",
		"tablockid",
		"oid",
		"rid",
		"pdist",
		"pid",
		"spd",
		"blk",
		"tripid",
		"dly",
		"ori",
	}
	values := []string{
		fromString(v.Vid),
		fromFloat(v.Lat),
		fromFloat(v.Lon),
		fromString(v.Des),
		fromTimestamp(v.Tmstmp),
		fromTimestamp(v.SrvTimstmp),
		fromString(v.Hdg),
		fromString(v.Rt),
		fromString(v.Tatripid),
		fromString(v.Tablockid),
		fromString(v.Oid),
		fromString(v.Rid),
		fromInt(v.Pdist),
		fromInt(v.Pid),
		fromInt(v.Spd),
		fromInt(v.Blk),
		fromInt(v.Tripid),
		fromBool(v.Dly),
		fromBool(v.Or),
	}
	return columns, values
}

func fromFloat(v float64) string {
	return fmt.Sprintf("%.16f", v)
}

func fromTimestamp(v VehicleTimestamp) string {
	return fmt.Sprintf("%d", v.Time.Unix())
}

func fromString(v string) string {
	return fmt.Sprintf("\"%s\"", v)
}

func fromInt(v int) string {
	return strconv.Itoa(v)
}

func fromBool(v bool) string {
	if v {
		return "1"
	}
	return "0"
}
