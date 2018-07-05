package ticker

import (
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/gomodule/redigo/redis"
)

var mux sync.Mutex
var shakes []int
var taps []int

var conn redis.Conn

func Add(shake int, tap int) {
	mux.Lock()
	log.Printf("Add %d %d\n", shake, tap)
	if shake > 0 {
		shakes = append(shakes, shake)
	}
	if tap > 0 {
		taps = append(taps, tap)
	}
	mux.Unlock()
}

// Lock must be taken
func ave(a []int) int {
	if len(a) == 0 {
		return 0
	}
	sum := 0
	for _, x := range a {
		sum += x
	}
	return sum / len(a)
}

type Point struct {
	Timestamp int64 `json:"timestamp"`
	Value     int   `json:"value"`
}

type Data struct {
	Shakes []Point `json:"shakes"`
	Taps   []Point `json:"taps"`
}

func to_points(ss [][]byte) []Point {
	n := len(ss)
	var ret []Point
	for i := 0; i < n; i += 2 {
		timestamp, err := strconv.ParseInt(string(ss[i]), 10, 64)
		if err != nil {
			log.Fatalf("Failed to parse int  %v", err.Error())
		}
		value, err := strconv.Atoi(string(ss[i+1]))
		if err != nil {
			log.Fatalf("Failed to atoi  %v", err.Error())
		}
		ret = append(ret, Point{timestamp, value})
	}
	return ret
}

func Retrieve() (*Data, error) {
	count := 24 * 60
	all_shakes, err := redis.ByteSlices(conn.Do("zrange", "dryer:shakes", -count, -1, "withscores"))
	if err != nil {
		log.Fatalln(err.Error())
		return nil, err
	}
	shake_points := to_points(all_shakes)
	all_taps, err := redis.ByteSlices(conn.Do("zrange", "dryer:taps", -count, -1, "withscores"))
	if err != nil {
		return nil, err
	}
	tap_points := to_points(all_taps)

	return &Data{
		shake_points,
		tap_points,
	}, nil
}

func Update() {
	timestamp := time.Now().Unix()

	mux.Lock()
	// score, value
	_, err := conn.Do("ZADD", "dryer:shakes", ave(shakes), timestamp)
	if err != nil {
		log.Printf("%v", err)
	}
	_, err = conn.Do("ZADD", "dryer:taps", ave(taps), timestamp)
	if err != nil {
		log.Printf("%v", err)
	}
	shakes = shakes[:0]
	taps = taps[:0]
	mux.Unlock()
}

func Start() {
	var err error
	conn, err = redis.Dial("tcp", ":6379")
	if err != nil {
		panic("Redis server is not ready.")
	}
	ticker := time.NewTicker(10 * time.Minute)
	go func() {
		for {
			select {
			case <-ticker.C:
				Update()
			}
		}
	}()
}
