package calendar

import (
	"L2_18/configs"
	"L2_18/configs/loader/dotEnvLoader"
	"L2_18/internal/repository"
	"L2_18/internal/usecase"
)

func Run() {
	cfg := configs.MustLoad(dotEnvLoader.DotEnvLoader{})
	repo := repository.NewEventRepository()
	uc := usecase.NewEventUC(repo)
}
