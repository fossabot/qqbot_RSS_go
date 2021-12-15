package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/mmcdole/gofeed"
	"log"
	"qqbot-RSS-go/db"
	"qqbot-RSS-go/modles/msg"
	"qqbot-RSS-go/services/bilibili"
	"qqbot-RSS-go/utils"
	"strconv"
)

func CommandAddRss(uri string, botUid int64, groupId int, userId int) string {
	if uri != "rss-add" && uri != "" {
		fp := gofeed.NewParser()
		rspCode := utils.CheckCode(uri)
		if rspCode == 200 {
			feed, rssErr := fp.ParseURL(uri)
			if rssErr != nil {
				log.Printf("非法RSS格式:%v", rssErr.Error())
				return "非法RSS格式"
			}
			result := db.InsertUrl(uri, feed.Title, botUid, groupId, userId)
			if result == true {
				return feed.Title + "订阅成功"
			} else {
				return feed.Title + "添加失败，订阅已经存在或注册异常"
			}
		} else {
			return uri + "无法访问或错误的URL"
		}
	} else {
		return "使用方法:rss-add RSS订阅URL\n建议使用https://rss.vark.fun获取RSS信息"
	}
}

func CommandAddLive(roomCode string, botUid int64, groupId int, userId int) string {
	if roomCode != "rss-live" && roomCode != "" {
		roomInfo := bilibili.LiveInfo(roomCode)
		var room msg.BiliLiveInfo
		fmt.Println(room.Code)
		err := json.Unmarshal(roomInfo, &room)
		if err != nil {
			log.Printf("序列化JSON发生异常:%v", err.Error())
			return "房间号码错误"
		}
		upData := bilibili.GetUpInfo(strconv.Itoa(room.Data.Uid))
		var upInfo msg.UpInfo
		err = json.Unmarshal(upData, &upInfo)
		if err != nil {
			log.Printf("序列化JSON发生异常:%v", err.Error())
			return "用户信息错误"
		}
		result := db.InsertRoom(room.Data.RoomId, upInfo.Data.Name, botUid, groupId, userId)
		if result == true {
			return upInfo.Data.Name + "直播间订阅成功"
		} else {
			return upInfo.Data.Name + "直播间订阅失败，已存在或注册异常"
		}
	} else {
		return "使用方法:rss-live bilibili直播间房间号"
	}
}

func CommandDelRss(botUid int64, groupId int, urlName string, createUserId int) string {
	if urlName != "rss-del" && urlName != "" {
		result := db.DelRss(botUid, groupId, urlName, createUserId)
		if result == true {
			return urlName + "取消订阅成功"
		} else {
			return urlName + "取消订阅失败，订阅不存在或权限不足"
		}
	} else {
		return "使用方法:rss-del 订阅名称"
	}
}

func CommandDelLive(botUid int64, groupId int, roomCode string, createUserId int) string {
	if roomCode != "rss-live-del" && roomCode != "" {
		result := db.DelLive(botUid, groupId, roomCode, createUserId)
		if result == true {
			return roomCode + "直播订阅取消成功"
		} else {
			return roomCode + "直播订阅取消失败，订阅不存在或权限不足"
		}
	} else {
		return "使用方法:rss-live-del bilibili直播房间号"
	}
}