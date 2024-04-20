package week01

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestIntDel(t *testing.T) {
	ints := []int{1, 2, 3, 4, 5}
	t.Logf("ints:%v, len:%d, cap:%d", ints, len(ints), cap(ints))
	SliceDelIdx(&ints, 2)
	t.Logf("ints:%v, len:%d, cap:%d", ints, len(ints), cap(ints))
}

func TestStringDel(t *testing.T) {
	strings := []string{"a", "b", "c", "d", "e"}
	t.Logf("strings:%v, len:%d, cap:%d", strings, len(strings), cap(strings))
	SliceDelIdx(&strings, 2)
	t.Logf("strings:%v, len:%d, cap:%d", strings, len(strings), cap(strings))
}

func TestSliceReduce(t *testing.T) {
	i := make([]int, 5, 257)
	i[0] = 1
	i[1] = 2
	i[2] = 3
	t.Logf("i:%v, len:%d, cap:%d", i, len(i), cap(i))
	SliceDelIdx(&i, 1)
	t.Logf("i:%v, len:%d, cap:%d", i, len(i), cap(i))

	i = append(i, 4)
	SliceDelIdx(&i, 1)
	t.Logf("i:%v, len:%d, cap:%d", i, len(i), cap(i))

	SliceDelIdx(&i, 0)
	t.Logf("i:%v, len:%d, cap:%d", i, len(i), cap(i))
}

func TestSliceReduce2(t *testing.T) {
	ints := make([]int, 1, 1)
	ints[0] = 1
	t.Logf("ints:%v, len:%d, cap:%d", ints, len(ints), cap(ints))
	SliceDelIdx(&ints, 0)
	t.Logf("ints:%v, len:%d, cap:%d", ints, len(ints), cap(ints))
	ints = append(ints, 2)
	t.Logf("ints:%v, len:%d, cap:%d", ints, len(ints), cap(ints))
}

func BenchmarkSliceDelIdx(b *testing.B) {
	//创造一个10w长度的slice
	ints := make([]int, 100000, 100000)
	for i := 0; i < 100000; i++ {
		ints[i] = i
	}
	for i := 0; i < b.N; i++ {
		go func(i int) {
			SliceDelIdx(&ints, uint(i%100000))
		}(i)
	}
}

func Test002(t *testing.T) {
	// 获取 access_token
	resp, err := http.Post("https://ai.fakeopen.com/auth/session", "application/x-www-form-urlencoded", bytes.NewBufferString("session_token=eyJhbGciOiJkaXIiLCJlbmMiOiJBMjU2R0NNIn0..L1FLcLfF3s3ihVTe.kQJZ2tJ0P6VdjJvrUdW67V2noVYMX6VhdmRNW1SnB0SnL2psRjIiWrGGUtqGY7pZOS2TdOga8WqA4cOSLDd4B5TgoOOyxcxCE599OVhhfpysNdpWkSOJJO9HGgR40ISAT4AKsuFYb0bbCkQtG46B69VhgHpe4UIJLbgNimoJK6KgICEFM8SuPo2ypBAvDv_S_NugcLe3dhSFSVX7-sk-uXLQ-4F0QD3Fwa1jcieUtcGDVqyeuFJjnB5pOQBmm6DaHDhp1ZLSXU9rQXd3USizbK2dV0rGPqXc-aATMgTLo38T7Q49A44ebcp3uo2Gvp8gN7AQFyDPKrMNQ4fcRQDXPt8yHO5efjN87ftjTWboyRadntY8zpLJEHV8dd8g3uC_0viED18rCo7LN0GHSyO2DoQLLrqPyJqGMy00f57ZNAZYLHEQNHCwrL8t1WsKi5n-9NsW8h2xC7aRVMjYz6bawwXzV3_s9K1fW-HBN2Zhu1NBBaFXT1XiyQVlU2rZIYH4hUOTibeFsyxaRKB9Ao5eeVe9pgPLgiYfe2AhAlhXyx3nfC_OjNlNbls8EbVWkPswQyzKGr4nhdJMGjRjs-xQw6wXG0NFqhEoXKS-j59s9NFPt-_iPZwCKW3u00LD6W8X9jGq1t49kFkKmgcIDjS_JK5lNyUqGt07P3QBNX0eMLiwUJAeyoGSHA4CRdgbu0zmAd_89PMJxKfB21NYg1Zs-ocSYhSD675xdmO7OGIYuQ8o1qXS-zE0I1IjEge8vrK7JGUK8Gh-Z07x0IN4p72u0QhknwFHLaJO_gGAv-tGFdqAGo5ZVOTTlFBJHI8I9jJ9-33UcSmwGZ5XX_H94cgxZRSUHtKvQqMbx8GuQC91f76h3ZCRIbMexWfaJYcyk6EK8w-DRlaWOo4c6zZx2du9m-tken659uW4KrNglnZvpqW1BHzta-liAXXLA6CiBWSpdLdNvX4IOWdCbBn23Rt3Y-9hB74dPs0UcdEzw-Y009SYTIVQzE0DzC_5wkC8HVp2Yjm1dyqz3oDRO_2kUUy8Go1pPRK3tK7i3QqHKIiu54Kl8sW4QP49zC6HzsKwCziJ8IMn7nYVBP73VyetUMu1Vj1ihlefI5z8bW77eDoOu4HiYvHIBvc6fw-C1bJYIRzL9TrWVfrxL3qlDsvUwKowo4aYtk7x8VoqdZY0zkcr5_apyEWqmlgqYtbrC69WDRT653XzNhf5CMnIMnpR5tp4TY2LmAhDsjskGpI-nd3pCzmH0f02eaqKL26_eY7UAMO0CWEUnlYZw6pq92JgkHdgvahXc6-9acoSi0ItycaBy8r6M2rVohw1eUTyzCo_SxHtc6PrsqIgANp6D-u7aClwEXEdUnQSJRqe2exCHKEDSZeGvEDpHMsoz6XKOC9nNKJJ89yCO1cXEm50kTqG9ZRM0rTxCVIdyQwUBaHQQ9bGLfeR2zRvRTiuvMKOenGL9jFqok1y6ed4xj4Zs8Lb01DS6Hm31S4LsOub5wAibZjhgfwcOFYt4fM8ne_p5ItVt7_rAgi-urHnwhrfqR_BQT_XMtbqB_6WbtVpf-t-lBIT24DxyNvJbz9bW-w2ueCYQag2VcBr9-yVV-6dBBWavKObI3GhiiRWoJ4Y4X5eYiEN3Fs0cJK-jtIgT3lUg1QoyPSSCL8XtcbAd8IXp-2DYoQmw4S_M4UPTOjExqyFyWtuAc3oBvJr-LTaSyJVk3tPUu5NIKAyBDoX3U1xgoFVIiIwiA-3KFtdcauFxEXjblQFN5rtKRE5YpnGH7zGQUafthUd-rquvGfyLhyNNbnAQRLxTnsh5hLefmYyOhhlbvExqCzmn-zTIMXK-stW93gPIEvuYO-DtJbIIw_EmCTmPaSQ02psMKu_cXBcD4JA4_fBWG8ACsdD4cgDSbYR9KqpDMqu3fYX_tOJYsAOlbPC4NOkNXJ7Y676_cOpE5ukGXZq4SmJGc_h7RMXXThFrkT9B6USVbzsINgdX4Ys20khtxh6b8axi6ZEG5vtSM8AfcyKQsO5CVQEQq4WQ8yF3qV2P4PAVQF4f4p5RDZx21d9eTrAKL8KbuMLbbDP04ix6YAgTqZFR6KLSqQIj69LMohB3VYT9r28o-21dTQC2Vc4KKdJmVMqVGDh80Nx9xLwXnd1PoEG43ZA1lKwDmh97Ob28tJIFP4BQlIE8UixQv3PjpBakoyz_h-glZKmEHIUhv-gEZSN68waWXgVsQhG9gJVw8W3oY3OvAHzGYmbFd9NHIPgcjMWd8VbilhFbHeZnm4BsJOr_sfRtui-YeGwIeIrbWj748JfawHaqau7xBA1oKI4fvaUfkLNysmKoe4UAGioHiNTqbuKSiG2Q0Zm2PSVu9DFdYdmAkgXBi7VPzhhYMd1vsxaPXk_Z6Dqx4w04Y3CMMoO_6-4xK_rP19q2ptM9oqH7zaoAwB8fLb8v0nm4_dnkv-_ZA1mYM3wVcpYQv4-E3V_TJZc2PD6j8ZKfO-ex3IPIEge214tRQeW-VwnN1RvugkJeNfq6Usg6ikGEYh0ktbxeZHIudiZcuUlUglmqBJ3HlgUdhjZwfs_rYzWrEErH5Kd5h6LlSxQU9a0sMry70usgzV-Jqq64g6m8xJcnb97EV6WzKcmh_elqkwW_c83fkdBKRAY6MZ2_c_KlD95V01BPceq35Ayy2nL7tSnMw2RK801ciqCsWJYT6rNJ8dwprmJgU5d3DwzfTFCxmG0I1EwKbC9EduXR67T8VJrvi4Ped7SSp5XW4Kv0k-ShcmiJyjqebY5caOET_f6o35-RnPsm556KeskGkcqnaaiL050oxFkLYfw3L9nORDaHPboT4sUGO5RHbvlEuDyU3iH82AarH35ew5yDlW17UFp3O_pF-9eNLjFg7D0WJIMggMNSgYKzp1bfA.QHB9jqnZBqgQIVAyYb72vg"))
	if err != nil {
		fmt.Println("发送请求失败:", err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("读取响应失败:", err)
		return
	}

	var result map[string]interface{}
	json.Unmarshal(body, &result)

	accessToken, ok := result["access_token"].(string)
	if !ok || accessToken == "" {
		fmt.Println("获取 access_token 失败")
		return
	}

	names := []string{"bonjour"}
	for _, name := range names {
		// 设置请求体
		data := fmt.Sprintf("unique_name=%s&access_token=%s&expires_in=0&site_limit=https://ai.ycvkzzz.life&show_conversations=false&show_userinfo=true", name, accessToken)
		tokenResp, err := http.Post("https://ai.fakeopen.com/token/register", "application/x-www-form-urlencoded", bytes.NewBufferString(data))
		if err != nil {
			fmt.Println("发送注册请求失败:", err)
			continue
		}
		defer tokenResp.Body.Close()

		tokenBody, err := ioutil.ReadAll(tokenResp.Body)
		if err != nil {
			fmt.Println("读取注册响应失败:", err)
			continue
		}

		var tokenResult map[string]interface{}
		json.Unmarshal(tokenBody, &tokenResult)

		tokenKey, ok := tokenResult["token_key"].(string)
		if !ok {
			fmt.Println(name, "注册失败")
			continue
		}

		fmt.Println(name, tokenKey)
	}
}
