package greet

import "gopkg.in/gookit/color.v1"

func HomeScreen() {
	color.Cyan.Println("\t\t\tHello Voter, Welcome to voting system.\n---------------------------------------------------------------------------------")
	color.Red.Print("NOTE: Don't ever use '")
	color.BgMagenta.Print("space")
	color.Red.Print("' as input in entire program. It may lead to unintended execution. However, if necessary use '_' instead.\n")
	color.LightMagenta.Println("Choose:")
	color.Yellow.Print("\n1.\tNew User?\n2.\tLogin\n3.\t")
	color.BgRed.Println("Admin Login.")
}
