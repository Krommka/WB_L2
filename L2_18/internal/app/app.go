package calendar

import (
	"L2_18/configs"
	"L2_18/configs/loader/dotEnvLoader"
)

func Run() {
	cfg := configs.MustLoad(dotEnvLoader.DotEnvLoader{})
}
