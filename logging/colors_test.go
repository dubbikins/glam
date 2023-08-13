package logging

import "testing"

func TestColor(t *testing.T) {
	t.Log("TestColor")
	colorString := color("test", "\033[1;31m")
	if colorString != "\033[1;31mtest\033[0m" {
		t.Error("color() failed")
	}
}

func TestRed(t *testing.T) {
	t.Log("TestRed")
	colorString := Red("test")
	if colorString != "\033[1;31mtest\033[0m" {
		t.Error("Red() failed")
	}
}

func TestGreen(t *testing.T) {
	t.Log("TestGreen")
	colorString := Green("test")
	if colorString != "\033[1;32mtest\033[0m" {
		t.Error("Green() failed")
	}
}

func TestYellow(t *testing.T) {
	t.Log("TestYellow")
	colorString := Yellow("test")
	if colorString != "\033[1;33mtest\033[0m" {
		t.Error("Yellow() failed")
	}
}

func TestBlue(t *testing.T) {
	t.Log("TestBlue")
	colorString := Blue("test")
	if colorString != "\033[1;34mtest\033[0m" {
		t.Error("Blue() failed")
	}
}

func TestPurple(t *testing.T) {
	t.Log("TestPurple")
	colorString := Purple("test")
	if colorString != "\033[1;35mtest\033[0m" {
		t.Error("Purple() failed")
	}
}

func TestCyan(t *testing.T) {
	t.Log("TestCyan")
	colorString := Cyan("test")
	if colorString != "\033[1;36mtest\033[0m" {
		t.Error("Cyan() failed")
	}
}

func TestGray(t *testing.T) {
	t.Log("TestGray")
	colorString := Gray("test")
	if colorString != "\033[1;37mtest\033[0m" {
		t.Error("Gray() failed")
	}
}

func TestMagenta(t *testing.T) {
	t.Log("TestMagenta")
	colorString := Magenta("test")
	if colorString != "\033[1;95mtest\033[0m" {
		t.Error("Magenta() failed")
	}
}

func TestWhite(t *testing.T) {
	t.Log("TestWhite")
	colorString := White("test")
	if colorString != "\033[1;97mtest\033[0m" {
		t.Error("White() failed")
	}
}
