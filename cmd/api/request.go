package main

type ApproveRequest struct {
	Email    string `json:"email"`
	Approved bool   `json:"approved"`
}
