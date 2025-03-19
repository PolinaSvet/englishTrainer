package main

import cmd "dictionary/cmd"

func init() {
	cmd.InitData()
}

func main() {
	cmd.Handler()
}
