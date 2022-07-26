import json
from termcolor import colored
import os

groupNum=795478528
groupName="Curve"
qqNick='rua'
qqNum=1536441243
isBotAdmin=True
userDataPath="./userData.json"

getUser(795478528,1536441243)

def readUserJson(jsonPath):
    with open(jsonPath, 'r') as f:
        return json.load(f)

def saveUserJson(Data,jsonPath):
    with open(jsonPath, 'w') as f:
        json.dump(Data, f)
    print("json save successfully")

def getUser(groupNum,qqNum):
    if not os.path.exists(userDataPath):
        with open(userDataPath, 'w') as f:
            newUserData=[]
            json.dump(newUserData, f)
    with open(userDataPath, 'r') as f:
        oldJson= json.load(f)
        userList=getGroup(groupNum,oldJson)['users']
        for user in userList:
            if user['qqNum']==qqNum:
                print("read user("+qqNum+")["+qqNick+"] in group("+groupNum+")["+groupName+"] successfully")
                return user
            
def getGroup(groupNum,json):
    for groupData_i in json:
        if groupNum==groupData_i['groupNum'] :
            return groupData_i
    return createNewGroup(groupNum) 


    
def createNewUser(groupNum,qqNum):
    user={}
    user['qqNum']=qqNum
    user['qqNick']=qqNick
    user['botAdmin']=False
    getGroup(groupNum)['users'].append(user)        
    print("add a "+colored("new","blue")+" user("+qqNum+")["+qqNick+"] in a "+colored("new","blue")+" group("+groupNum+")["+groupName+"] successfully")
    return user


def createNewGroup(groupNum,json):
    groupData={}
    groupData['groupNum']=groupNum
    groupData['users']=[]
    json.append[groupData]
    print("create a "+colored("new","blue")+" group("+groupNum+")["+groupName+"] successfully")
    return groupData
 
 

def writeAllData(userData):
    with open(userDataPath,"w") as f:   
        json.dump(userData,f)
        print("加载入文件完成...")
