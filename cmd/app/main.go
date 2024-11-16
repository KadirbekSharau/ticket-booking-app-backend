package main

import "ticket-booking-app-backend/internal"


const configsDir = "internal/infrastructure/configs"

func main() {
	internal.Run(configsDir)
}
