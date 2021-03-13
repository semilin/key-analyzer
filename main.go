package main

import "fmt"

func main() {
	Data = TextData()
	LoadLayouts()
	for i:=len(Data)-1;i>0;i-=2 {
		if Data[i].Count == 1 {
			Data = Data[0:i-1]
		} else {
			break
		}
	}

	fmt.Println(Data[0:5])
	webServer()
}
