// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

// Hub maintains the set of active clients and broadcasts messages to the
// clients.

var questionCycle = []string{"Q1", "Q2", "Q3"}
var answerCycle = []string{"A1", "A2", "A3"}
var cycleIndex = 0

type Hub struct {
    // Registered clients.
    clients map[*Client]bool

    // Inbound messages from the clients.
    broadcast chan []byte

    // Register requests from the clients.
    register chan *Client

    // Unregister requests from clients.
    unregister chan *Client

    // Question to be answered
    question string

    // Answer to question
    answer string
}

func newHub() *Hub {
    return &Hub{
        broadcast:  make(chan []byte),
        register:   make(chan *Client),
        unregister: make(chan *Client),
        clients:    make(map[*Client]bool),
        question: "Send 'play' for a new question",
        answer: "initialanswer",
    }
}

func (h *Hub) newQuestion() {
    h.question = questionCycle[cycleIndex]
    h.answer = answerCycle[cycleIndex]
    cycleIndex = (cycleIndex + 1) % len(questionCycle)
}

func (h *Hub) run() {
    for {
        select {
        case client := <-h.register:
            h.clients[client] = true
        case client := <-h.unregister:
            if _, ok := h.clients[client]; ok {
                delete(h.clients, client)
                close(client.send)
            }
        case message := <-h.broadcast:
            for client := range h.clients {
                select {
                case client.send <- message:
                default:
                    close(client.send)
                    delete(h.clients, client)
                }
            }
        }
    }
}
