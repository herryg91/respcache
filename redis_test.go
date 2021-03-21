package respcache

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/gomodule/redigo/redis"
	"github.com/rafaeljusto/redigomock/v3"
)

type TestStruct struct {
	Id   int
	Name string
}

type TestStructCustJson struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func Test_resp_cache_redis_Run(t *testing.T) {
	mockRedisConn := redigomock.NewConn()
	mockRedisPool := &redis.Pool{
		Dial: func() (redis.Conn, error) {
			return mockRedisConn, nil
		}}

	var out_happy_scenario_uncached TestStruct
	mockRedisConn.Command("GET", "happy-scenario-uncached").ExpectError(redis.ErrNil)
	mockRedisConn.Command("SET", "happy-scenario-uncached", `{"Id":1,"Name":"test 1"}`, "EX", 10).ExpectError(nil)

	var out_happy_scenario_cached TestStruct
	mockRedisConn.Command("GET", "happy-scenario-cached").Expect(`{"Id":1,"Name":"test 1"}`).ExpectError(nil)

	var out_happy_scenario_cached_string string
	mockRedisConn.Command("GET", "happy-scenario-cached-string").Expect(`"test string"`).ExpectError(nil)

	var out_happy_scenario_cached_int int
	mockRedisConn.Command("GET", "happy-scenario-cached-int").Expect(`1`).ExpectError(nil)

	var out_happy_scenario_cached_float float64
	mockRedisConn.Command("GET", "happy-scenario-cached-float").Expect(`1.23`).ExpectError(nil)

	var out_happy_scenario_cached_cust_json TestStructCustJson
	mockRedisConn.Command("GET", "happy-scenario-cached-cust-json").Expect(`{"id":1,"name":"test 1"}`).ExpectError(nil)

	var out_happy_scenario_cached_arr_struct []TestStruct
	mockRedisConn.Command("GET", "happy-scenario-cached-arr-struct").Expect(`[{"Id":1,"Name":"test 1"},{"Id":2,"Name":"test 2"}]`).ExpectError(nil)

	var out_happy_scenario_cached_arr_ptr_struct []*TestStruct
	mockRedisConn.Command("GET", "happy-scenario-cached-arr-ptr-struct").Expect(`[{"Id":1,"Name":"test 1"},{"Id":2,"Name":"test 2"}]`).ExpectError(nil)

	var out_fail_fallback_error TestStruct
	mockRedisConn.Command("GET", "fail-fallback-error").ExpectError(redis.ErrNil)

	var out_fail_scenario_unmarshal TestStruct
	mockRedisConn.Command("GET", "fail-scenario-unmarshal").Expect(`a`).ExpectError(nil)
	mockRedisConn.Command("SET", "fail-scenario-unmarshal", `{"Id":1,"Name":"test 1"}`).ExpectError(nil)

	var out_fail_not_pointer TestStruct
	mockRedisConn.Command("GET", "fail-not-pointer").ExpectError(redis.ErrNil)

	var out_fallback_nil TestStruct
	mockRedisConn.Command("GET", "fallback-nil").ExpectError(nil)

	type fields struct {
		rdsPool *redis.Pool
	}
	type args struct {
		key        string
		ttl        int
		out        interface{}
		fallbackFn CacheFallback
	}
	tests := []struct {
		name         string
		fields       fields
		args         args
		wantIscached bool
		wantErr      bool
	}{
		// TODO: Add test cases.
		{
			name:   "happy scenario uncached",
			fields: fields{rdsPool: mockRedisPool},
			args: args{
				"happy-scenario-uncached",
				10,
				&out_happy_scenario_uncached,
				func() (interface{}, error) {
					return TestStruct{Id: 1, Name: "test 1"}, nil
				}},
			wantIscached: false,
			wantErr:      false,
		},
		{
			name:   "happy scenario cached",
			fields: fields{rdsPool: mockRedisPool},
			args: args{
				"happy-scenario-cached",
				10,
				&out_happy_scenario_cached,
				func() (interface{}, error) {
					return TestStruct{Id: 1, Name: "test 1"}, nil
				}},
			wantIscached: true,
			wantErr:      false,
		},
		{
			name:   "happy scenario cached string",
			fields: fields{rdsPool: mockRedisPool},
			args: args{
				"happy-scenario-cached-string",
				10,
				&out_happy_scenario_cached_string,
				func() (interface{}, error) {
					return "test string", nil
				}},
			wantIscached: true,
			wantErr:      false,
		},
		{
			name:   "happy scenario cached int",
			fields: fields{rdsPool: mockRedisPool},
			args: args{
				"happy-scenario-cached-int",
				10,
				&out_happy_scenario_cached_int,
				func() (interface{}, error) {
					return 1, nil
				}},
			wantIscached: true,
			wantErr:      false,
		},
		{
			name:   "happy scenario cached float",
			fields: fields{rdsPool: mockRedisPool},
			args: args{
				"happy-scenario-cached-float",
				10,
				&out_happy_scenario_cached_float,
				func() (interface{}, error) {
					return 1.23, nil
				}},
			wantIscached: true,
			wantErr:      false,
		},
		{
			name:   "happy scenario cached cust json",
			fields: fields{rdsPool: mockRedisPool},
			args: args{
				"happy-scenario-cached-cust-json",
				10,
				&out_happy_scenario_cached_cust_json,
				func() (interface{}, error) {
					return TestStructCustJson{Id: 1, Name: "test 1"}, nil
				}},
			wantIscached: true,
			wantErr:      false,
		},
		{
			name:   "happy scenario cached arr struct",
			fields: fields{rdsPool: mockRedisPool},
			args: args{
				"happy-scenario-cached-arr-struct",
				10,
				&out_happy_scenario_cached_arr_struct,
				func() (interface{}, error) {
					return []TestStruct{
						TestStruct{Id: 1, Name: "test 1"},
						TestStruct{Id: 2, Name: "test 2"},
					}, nil
				}},
			wantIscached: true,
			wantErr:      false,
		},
		{
			name:   "happy scenario cached arr ptr struct",
			fields: fields{rdsPool: mockRedisPool},
			args: args{
				"happy-scenario-cached-arr-ptr-struct",
				10,
				&out_happy_scenario_cached_arr_ptr_struct,
				func() (interface{}, error) {
					return []*TestStruct{
						&TestStruct{Id: 1, Name: "test 1"},
						&TestStruct{Id: 2, Name: "test 2"},
					}, nil
				}},
			wantIscached: true,
			wantErr:      false,
		},
		{
			name:   "fail fallback error",
			fields: fields{rdsPool: mockRedisPool},
			args: args{
				"fail-fallback-error",
				10,
				&out_fail_fallback_error,
				func() (interface{}, error) {
					return TestStruct{Id: 1, Name: "test 1"}, fmt.Errorf("fallback error")
				}},
			wantIscached: false,
			wantErr:      true,
		},
		{
			name:   "fail scenario unmarshal",
			fields: fields{rdsPool: mockRedisPool},
			args: args{
				"fail-scenario-unmarshal",
				0,
				&out_fail_scenario_unmarshal,
				func() (interface{}, error) {
					return TestStruct{Id: 1, Name: "test 1"}, nil
				}},
			wantIscached: false,
			wantErr:      false,
		},
		{
			name:   "fail not pointer",
			fields: fields{rdsPool: mockRedisPool},
			args: args{
				"fail-not-pointer",
				10,
				out_fail_not_pointer,
				func() (interface{}, error) {
					return TestStruct{Id: 1, Name: "test 1"}, nil
				}},
			wantIscached: false,
			wantErr:      true,
		},
		{
			name:   "fallback nil",
			fields: fields{rdsPool: mockRedisPool},
			args: args{
				"fallback-nil",
				10,
				&out_fallback_nil,
				func() (interface{}, error) {
					return nil, nil
				}},
			wantIscached: false,
			wantErr:      false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rc := &resp_cache_redis{
				rdsPool: tt.fields.rdsPool,
			}
			gotIscached, err := rc.Run(tt.args.key, tt.args.ttl, tt.args.out, tt.args.fallbackFn)
			if (err != nil) != tt.wantErr {
				t.Errorf("resp_cache_redis.Run() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotIscached != tt.wantIscached {
				t.Errorf("resp_cache_redis.Run() = %v, want %v", gotIscached, tt.wantIscached)
			}
		})
	}
}

func Test_resp_cache_redis_set(t *testing.T) {
	mockRedisConn := redigomock.NewConn()
	mockRedisPool := &redis.Pool{
		Dial: func() (redis.Conn, error) {
			return mockRedisConn, nil
		}}

	rc := &resp_cache_redis{
		rdsPool: mockRedisPool,
	}
	// scenario: happy scenario with ttl
	mockRedisConn.Command("SET", "happy-scenario-with-ttl", `{"Id":1,"Name":"test 1"}`, "EX", 10).ExpectError(nil)
	rc.set("happy-scenario-with-ttl", 10, TestStruct{Id: 1, Name: "test 1"})

	// scenario: happy scenario without ttl
	mockRedisConn.Command("SET", "happy-scenario-without-ttl", `{"Id":1,"Name":"test 1"}`).ExpectError(nil)
	rc.set("happy-scenario-without-ttl-1", 0, TestStruct{Id: 1, Name: "test 1"})
	rc.set("happy-scenario-without-ttl-2", -10000, TestStruct{Id: 1, Name: "test 1"})

	// scenario: fail scenario set
	mockRedisConn.Command("SET", "fail-scenario-set", `{"Id":1,"Name":"test 1"}`).ExpectError(fmt.Errorf("redis error"))
	rc.set("fail-scenario-set", 0, TestStruct{Id: 1, Name: "test 1"})

	// scenario: fail scenario marshal
	rc.set("happy-scenario-without-ttl-1", 0, make(chan (string)))
}

func Test_resp_cache_redis_get(t *testing.T) {
	mockRedisConn := redigomock.NewConn()
	mockRedisPool := &redis.Pool{
		Dial: func() (redis.Conn, error) {
			return mockRedisConn, nil
		}}

	mockRedisConn.Command("GET", "happy-scenario-not-cached-yet").ExpectError(redis.ErrNil)
	mockRedisConn.Command("GET", "happy-scenario-cached").Expect([]byte(`{"Id":1,"Name":"test 1"}`)).ExpectError(nil)
	mockRedisConn.Command("GET", "fail-scenario-err-redis").ExpectError(redis.ErrPoolExhausted)

	type fields struct {
		rdsPool *redis.Pool
	}
	type args struct {
		key string
	}
	tests := []struct {
		name         string
		fields       fields
		args         args
		wantIscached bool
		wantResp     string
	}{
		// TODO: Add test cases.
		{
			name:         "happy scenario not cached yet",
			fields:       fields{rdsPool: mockRedisPool},
			args:         args{key: "happy-scenario-not-cached-yet"},
			wantIscached: false,
			wantResp:     "",
		},
		{
			name:         "happy scenario cached",
			fields:       fields{rdsPool: mockRedisPool},
			args:         args{key: "happy-scenario-cached"},
			wantIscached: true,
			wantResp:     `{"Id":1,"Name":"test 1"}`,
		},
		{
			name:         "fail scenario err redis",
			fields:       fields{rdsPool: mockRedisPool},
			args:         args{key: "fail-scenario-err-redis"},
			wantIscached: false,
			wantResp:     "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rc := &resp_cache_redis{
				rdsPool: tt.fields.rdsPool,
			}
			gotIscached, gotResp := rc.get(tt.args.key)
			if gotIscached != tt.wantIscached {
				t.Errorf("resp_cache_redis.get() gotIscached = %v, want %v", gotIscached, tt.wantIscached)
			}
			if gotResp != tt.wantResp {
				t.Errorf("resp_cache_redis.get() gotResp = %v, want %v", gotResp, tt.wantResp)
			}
		})
	}
}

func TestNewRedisCache(t *testing.T) {
	type args struct {
		rdsPool *redis.Pool
	}
	tests := []struct {
		name string
		args args
		want RespCache
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewRedisCache(tt.args.rdsPool); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewRedisCache() = %v, want %v", got, tt.want)
			}
		})
	}
}
