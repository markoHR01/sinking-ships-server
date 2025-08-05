package main

import "strings"

type Message map[string]string

func Deserialize(data []byte) Message {
	message := make(map[string]string)
	input := string(data)

	msgLines := strings.Split(input, "\n")
	for _, msgLine := range msgLines {
		keyValue := strings.SplitN(msgLine, "=", 2)
		if len(keyValue) == 2 {
			key := keyValue[0]
			value := keyValue[1]
			message[key] = value
		}
	}

	return message
}

func Serialize(message Message) string {
	var serialized strings.Builder

	for key, value := range message {
		serialized.WriteString(key)
		serialized.WriteString("=")
		serialized.WriteString(value)
		serialized.WriteString("\n")
	}
	serialized.WriteString("STOP\n")

	return serialized.String()
}
