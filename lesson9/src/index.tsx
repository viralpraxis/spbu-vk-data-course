import * as React from "react"
import * as ReactDOM from "react-dom"
import { createEvent, createStore } from "effector"
import { useEvent, useStore } from "effector-react"

import { ConstantBackoff, Websocket, WebsocketBuilder, WebsocketEvent } from "websocket-ts"

import { applyPatch } from "fast-json-patch"

ReactDOM.render(
  <React.StrictMode>
  </React.StrictMode>,
  document.getElementById("root")
)

const patchSnapEv = createEvent()
const $snap = createStore<any>({})
$snap.on(patchSnapEv, (state, patch:any) => {
  const doc = applyPatch(state, patch, false, false).newDocument
  return doc
})

const ws = new WebsocketBuilder("ws://localhost:8080/ws")
  .withBackoff(new ConstantBackoff(10 * 1000))
  .build()

const receiveMessage = (i: Websocket, ev: MessageEvent) => {
  const transaction = JSON.parse(ev.data)
  const patch = JSON.parse(transaction.Payload)
  patchSnapEv(patch)
}

ws.addEventListener(WebsocketEvent.message, receiveMessage)

const setPlayerNameEv = createEvent<string>()
const $playerName = createStore<string>("")
$playerName.on(setPlayerNameEv, (_, payload) => {
  return payload
})

const LOGOS = new Map([
  ["ruby", "https://upload.wikimedia.org/wikipedia/commons/f/f1/Ruby_logo.png"],
  ["golang", "https://go.dev/blog/go-brand/Go-Logo/PNG/Go-Logo_Aqua.png"],
  ["haskell", "https://upload.wikimedia.org/wikipedia/commons/thumb/1/1c/Haskell-Logo.svg/2560px-Haskell-Logo.svg.png"],
  ["lisp", "https://upload.wikimedia.org/wikipedia/commons/f/f4/Lisplogo.png"]
])

function PlayerInput() {
  const setPlayerName = useEvent(setPlayerNameEv)
  const playerName = useStore($playerName)
  const handleChange = function(event: React.ChangeEvent<HTMLInputElement>) {
    setPlayerName(event.target.value)
  }

  const snap = useStore($snap)

  const randomPosition = function() {
    return [
      Math.floor(Math.random() * window.innerWidth * 0.5),
      Math.floor(Math.random() * window.innerHeight * 0.5)
    ]
  }

  const updateRemoteState = function(body:string) {
    fetch("http://localhost:8080/replace", {
      method: "POST",
      body:body
    })
  }

  const handleClick = function() {
    const [x, y] = randomPosition()
    const size = 20
    const logoIndex = Math.floor(Math.random() * Array.from(LOGOS.keys()).length)
    const logo = Array.from(LOGOS.keys())[logoIndex]

    updateRemoteState(
      `[{"op":"add", "path": "/${playerName}", "value": {"x":${x}, "y":${y}, "size":${size}, "logo":"${logo}"}}]`
    )
  }

  const updatePlayer = function(player:any) {
    updateRemoteState(
      `[{"op":"add", "path": "/${playerName}", "value": {"x":${player.x}, "y":${player.y}, "size":${player.size}, "logo":"${player.logo}"}}]`
    )
  }

  const shouldUpdatePlayerSize = function(player: any): Boolean {
    return player.size < 200 && (player.x % 42 == 0 || Math.random() > .9)
  }

  const handleUp = function() {
    let player = snap[playerName]
    player.y--
    if (shouldUpdatePlayerSize(player)) {player.size++}
    updatePlayer(player)
  }

  const handleDown = function() {
    let player = snap[playerName]
    player.y++
    if (shouldUpdatePlayerSize(player)) {player.size++}
    updatePlayer(player)
  }

  const handleLeft = function() {
    let player = snap[playerName]
    player.x--
    if (shouldUpdatePlayerSize(player)) {player.size++}
    updatePlayer(player)
  }

  const handleRight = function() {
    let player = snap[playerName]
    player.x++
    if (shouldUpdatePlayerSize(player)) {player.size++}
    updatePlayer(player)
  }

  return (<>
    <input type="text"
      onChange={handleChange}
      value={playerName}/>
    <input type="button"
      value={"Установить"}
      onClick={handleClick}/>
    <input type="button"
      value={"Вверх"}
      onClick={handleUp}/>
    <input type="button"
      value={"Вниз"}
      onClick={handleDown}/>
    <input type="button"
      value={"Влево"}
      onClick={handleLeft}/>
    <input type="button"
      value={"Вправо"}
      onClick={handleRight}/>
  </>)
}

function Players() {
  const snap = useStore($snap)
  let players = []
  for (const key in snap) {
    let player = snap[key]

    players.push(<circle cx={player.x} cy={player.y} r={player.size} stroke="white" fill="transparent" />)
    players.push(<text x={player.x} y={player.y} font-size="8">{key}</text>)
    players.push(<image href={LOGOS.get(String(player.logo))} x={player.x} y={player.y} height={player.size} width={player.size}/>)
  }
  return <>{players}</>
}

function GameResult() {
  const snap = useStore($snap)
  let players = []
  for (const key in snap) {
    let player = snap[key]

    if (player.size > 30) {
      return <><div>{key} won!</div></>
    }
  }

  return <></>
}

ReactDOM.render(
  <React.StrictMode>
    <GameResult></GameResult>
    <PlayerInput></PlayerInput>
    <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 800 600">
      <rect x='0' y='0' width='100%' height='100%' fill='tomato' opacity='0.75' />
      <Players></Players>
    </svg>
  </React.StrictMode>,
  document.getElementById("root")
)
