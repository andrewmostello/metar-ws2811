package config

import (
	"github.com/andrewmostello/metar-ws2811/ws2811"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	cfgKeyLEDCount      = "led.count"
	cfgKeyLEDBrightness = "led.brightness"
	cfgKeyLEDGPIOPin    = "led.gpio_pin"
)

type LED struct {
	Count      int
	Brightness int
	GPIOPin    int
}

func GetLED() LED {
	return LED{
		Count:      viper.GetInt(cfgKeyLEDCount),
		Brightness: viper.GetInt(cfgKeyLEDBrightness),
		GPIOPin:    viper.GetInt(cfgKeyLEDGPIOPin),
	}
}

func AddLEDFlags(cmd *cobra.Command) {
	flag := "led-count"
	cmd.PersistentFlags().Int(flag, ws2811.DefaultLEDCount, "Total count of LEDs in the string.")
	viper.BindPFlag(cfgKeyLEDCount, cmd.PersistentFlags().Lookup(flag))

	flag = "led-brightness"
	cmd.PersistentFlags().Int(flag, ws2811.DefaultBrightness, "Brightness of the LEDs.")
	viper.BindPFlag(cfgKeyLEDBrightness, cmd.PersistentFlags().Lookup(flag))

	flag = "led-gpio-pin"
	cmd.PersistentFlags().Int(flag, ws2811.DefaultGPIOPin, "GPIO pin of the data input to the LEDs.")
	viper.BindPFlag(cfgKeyLEDGPIOPin, cmd.PersistentFlags().Lookup(flag))
}
