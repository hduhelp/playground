package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/Mrs4s/MiraiGo/client"
	"github.com/Mrs4s/MiraiGo/message"
	"github.com/spf13/viper"
)

func init() {
	viper.SetConfigName("bot")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	bot := client.NewClient(viper.GetInt64("bot.uin"), viper.GetString("bot.passwd"))
	client.GenRandomDevice()
	ioutil.WriteFile("device.json", client.SystemDeviceInfo.ToJson(), os.FileMode(0755))
	b, _ := os.ReadFile("device.json")
	client.SystemDeviceInfo.ReadJson(b)
	bot.Login()

	groupId, friendId :=
		viper.GetInt64("group.demo_QQ_Group"),
		viper.GetInt64("friend.demo_QQ_Friend")

	// 发送私聊消息
	bot.SendPrivateMessage(friendId, message.NewSendingMessage().Append(message.NewText("你好！ (*´▽｀)ノノ")))

	// 发送群聊消息
	bot.SendGroupMessage(groupId, message.NewSendingMessage().Append(message.NewText("大家好，我是一个机器人！Hi~ o(*￣▽￣*)ブ")))

	// 在有权限的情况下 `@全体成员`
	bot.SendGroupMessage(groupId, message.NewSendingMessage().Append(message.AtAll()))

	// 在群成员加入时 at 并欢迎
	bot.GroupMemberJoinEvent.Subscribe(func(c *client.QQClient, e *client.MemberJoinGroupEvent) {
		if e.Group.Code == groupId {
			c.SendGroupMessage(e.Group.Code,
				message.NewSendingMessage().Append(message.NewAt(e.Member.Uin, "欢迎━(*｀∀´*)ノ亻!")))
		}
	})

	// 在被 At 的情况下返回当天天气
	bot.GroupMessageEvent.Subscribe(func(c *client.QQClient, e *message.GroupMessage) {
		beenAt, city := false, ""
		if e.GroupCode == groupId {
			for _, msgElem := range e.Elements {
				if msgElem.Type() == message.At {
					beenAt = true
				}
				if msgElem.Type() != message.At && beenAt {
					city = getCity(msgElem)
				}
			}
		}

		if city != "" {
			weatherInfo := getWeather(city)
			c.SendGroupMessage(e.GroupCode, message.NewSendingMessage().Append(message.NewText(weatherInfo)))
		}
	})

	// ...可选功能实现
	_block()
}

func _block() {
	select {}
}

func getCity(msgElem message.IMessageElement) string {
	s := fmt.Sprintf("%s", msgElem)
	return strings.TrimSpace(s[2 : len(s)-1])
}

func getWeather(city string) string {
	return weatherIn(city)
}
