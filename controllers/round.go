package controllers

type RoundController interface {
}

type roundController struct{}

func NewRoundController() RoundController {
	return &roundController{}
}
