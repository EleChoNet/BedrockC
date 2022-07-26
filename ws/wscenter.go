package ws

import (
	"BedrockC/bedrock"
	"github.com/tidwall/gjson"
)

var bedrockServer bedrock.BedrockS
func InitCenter(b bedrock.BedrockS){
	bedrockServer = b
	bedrockServer.Run()
}
//订阅事件响应
func onEventGot(msg string){

}
//CommandRequest 返回
func onCommandResponseGot(msg string){

}
func isWouldToSendMsg() bool{

}
func addSubscription(eventName string){

}
func removeSubscription(eventName string){

}
func commandRequest(command string){
	
}

func onMessageIn(msg string){
	messagePurpose:=gjson.Get(msg,"header.messagePurpose").String()
	switch messagePurpose{
	case "event":onEventGot(msg)
	case "commandResponse":onCommandResponseGot(msg)
	}
}