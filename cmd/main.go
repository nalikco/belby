package main

import (
	"belby/pkg/vk"
	"fmt"
)

func main() {
	vkApi := vk.NewVk(
		"211151815",
		"vk1.a.3oBv1rvOtNfXiLw90nsLjvCCwq3TmBVz46PA-WdEqNZnHYZIzP6GZGY8BORFv3KziYAQOUfWBRXvauaFUlMjn-Ha6mzJ4Q6daH7wHXez1lefTuGH38GDsXirJaSttuM2I5VTTVc54rwZLwW8z_hbp1FaFdMgVBnqvIp-Fnf3_FQXCmEXzk1hxVqIytcTxkTt",
	)

	err := vkApi.Polling()
	if err != nil {
		fmt.Println(err)
	}
}
