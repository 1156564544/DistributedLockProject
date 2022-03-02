package main

type Leader struct {
	Followers []string
	Clients   []string
	Lock      map[string]string
}