package tests

import (
	"testing"
	".."
	"fmt"
)

func TestEchos(tt *testing.T) {
	t := s.T(tt)

	s.Register("/echo1", Echo1)
	s.Register("/echo2", Echo2)
	s.Register("/echo3", Echo3)
	s.SetTestHeader("Cid", "test-client")

	s.StartTestService()
	defer s.StopTestService()
	s.EnableLogs(false)

	data := s.TestService("/echo1?aaa=11&bbb=_o_", s.Map{
		"ccc": "ccc",
		"DDD": 101.123,
		"eEe": true,
		"fff": nil,
		"ggg": "223",
	})

	datas, ok := data.(map[string]interface{})
	t.Test(ok, "[Echo1] Data1", data)
	d1, ok := datas["in"].(map[string]interface{})
	t.Test(ok, "[Echo1] Data2", data)
	d2, ok := datas["headers"].(map[string]interface{})
	t.Test(ok, "[Echo1] Data3", data)
	t.Test(d1["aaa"].(float64) == 11 && d1["bbb"] == "_o_" && d1["ddd"] == 101.123 && d1["eee"] == true && d1["fff"] == nil, "[Echo1] In", data)
	t.Test(d2["cid"] == "test-client", "[Echo1] Headers", data)

	data = s.TestService("/echo2?aaa=11&bbb=_o_", s.Map{
		"ccc": "ccc",
		"DDD": 101.123,
		"eEe": true,
		"fff": nil,
		"ggg": 223,
	})

	d, ok := data.(map[string]interface{})
	t.Test(ok, "[Echo2] Data1", data)
	t.Test(d["aaa"].(float64) == 11 && d["bbb"] == "_o_" && d["ddd"] == 101.123 && d["eee"] == true && d["fff"] == nil, "[Echo2] Data2", data)

	data = s.TestService("/echo3?a=1", s.Map{"name": "Star"})
	a, ok := data.([]interface{})
	t.Test(ok, "[Echo3] Data1", data)
	t.Test(a[0] == "Star", "[Echo3] Data2", data)
	t.Test(a[1] == "/echo3", "[Echo3] Data3", data)
}

func TestFilters(tt *testing.T) {
	t := s.T(tt)
	s.Register("/echo2", Echo2)

	s.StartTestService()
	defer s.StopTestService()
	s.EnableLogs(false)

	data := s.TestService("/echo2?aaa=11&bbb=_o_", s.Map{"ccc": "ccc"})
	d, _ := data.(map[string]interface{})
	t.Test(d["filterTag"] == "", "[Test InFilter 1] Response", data)

	s.SetInFilter(func(in map[string]interface{}) interface{} {
		in["filterTag"] = "Abc"
		in["filterTag2"] = 1000
		return nil
	})
	data = s.TestService("/echo2?aaa=11&bbb=_o_", s.Map{"ccc": "ccc"})
	d, _ = data.(map[string]interface{})
	t.Test(d["filterTag"] == "Abc" && d["filterTag2"].(float64) == 1000, "[Test InFilter 2] Response", data)

	s.SetOutFilter(func(in map[string]interface{}, result interface{}) (interface{}, bool) {
		data := result.(echo2Args)
		data.FilterTag2 = data.FilterTag2 + 100
		return data, false
	})

	data = s.TestService("/echo2?aaa=11&bbb=_o_", s.Map{"ccc": "ccc"})
	d, _ = data.(map[string]interface{})
	t.Test(d["filterTag"] == "Abc" && d["filterTag2"].(float64) == 1100, "[Test OutFilters 1] Response", data)

	s.SetOutFilter(func(in map[string]interface{}, result interface{}) (interface{}, bool) {
		data := result.(echo2Args)
		fmt.Println(" ***************", data.FilterTag2 + 100)
		return s.Map{
			"filterTag":  in["filterTag"],
			"filterTag2": data.FilterTag2 + 100,
		}, true
	})
	s.EnableLogs(true)
	data = s.TestService("/echo2?aaa=11&bbb=_o_", s.Map{"ccc": "ccc"})
	d, _ = data.(map[string]interface{})
	t.Test(d["filterTag"] == "Abc" && d["filterTag2"].(float64) == 1200, "[Test OutFilters 2] Response", data)


	s.SetInFilter(func(in map[string]interface{}) (interface{}) {
		return echo2Args{
			FilterTag:  in["filterTag"].(string),
			FilterTag2: in["filterTag2"].(int) + 100,
		}
	})
	data = s.TestService("/echo2?aaa=11&bbb=_o_", s.Map{"ccc": "ccc"})
	d, _ = data.(map[string]interface{})
	t.Test(d["filterTag"] == "Abc" && d["filterTag2"].(float64) == 1300, "[Test InFilter 3] Response", data)
}

func BenchmarkEchosForStruct(tb *testing.B) {
	tb.StopTimer()
	s.Register("/echo1", Echo1)
	s.EnableLogs(false)

	s.StartTestService()
	defer s.StopTestService()

	s.TestService("/echo1?aaa=11&bbb=_o_", s.Map{})

	tb.StartTimer()

	tb.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			s.TestService("/echo1?aaa=11&bbb=_o_", s.Map{
				"ccc": "ccc",
				"DDD": 101.123,
				"eEe": true,
				"fff": nil,
				"ggg": 223,
			})
		}
	})
}

func BenchmarkEchosForMap(tb *testing.B) {
	tb.StopTimer()
	s.Register("/echo2", Echo2)
	s.EnableLogs(false)
	s.SetTestHeader("Cid", "test-client")

	s.StartTestService()
	defer s.StopTestService()

	s.TestService("/echo2?aaa=11&bbb=_o_", s.Map{})

	tb.StartTimer()

	for i := 0; i < tb.N; i++ {

		s.TestService("/echo2?aaa=11&bbb=_o_", s.Map{
			"ccc": "ccc",
			"DDD": 101.123,
			"eEe": true,
			"fff": nil,
			"ggg": 223,
		})

	}
}
