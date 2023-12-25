import * as React from "react"
import * as ReactDOM from "react-dom"
import { createEvent, createStore } from "effector"
import { useEvent, useStore } from "effector-react"

import { ConstantBackoff, Websocket, WebsocketBuilder, WebsocketEvents } from "websocket-ts"

import { applyPatch } from "fast-json-patch"
import MainField from "./field"
import Copatych from "./copatych"

// Вот тут всякая логика

const patchSnapEv = createEvent()
const $snap = createStore<any>({})
$snap.on(patchSnapEv, (state, patch:any) => {
    const doc = applyPatch(state, patch, false, false).newDocument
    return doc
})

const ws = new WebsocketBuilder("ws://localhost:8080/ws")
    .withBackoff(new ConstantBackoff(10*1000))
    .build()

const receiveMessage = (i: Websocket, ev: MessageEvent) => {
    const transaction = JSON.parse(ev.data)
    const patch = JSON.parse(transaction.Payload)
    patchSnapEv(patch)
}

ws.addEventListener(WebsocketEvents.message, receiveMessage)

const setPlayerNameEv = createEvent<string>()
const $playerName = createStore<string>("")
$playerName.on(setPlayerNameEv, (_, payload) => {
    return payload
})

function PlayerInput() {
    const setPlayerName = useEvent(setPlayerNameEv)
    const playerName = useStore($playerName)
    const handleChange = function(event: React.ChangeEvent<HTMLInputElement>) {
        setPlayerName(event.target.value)
    }

    const snap = useStore($snap)

    const handleClick = function() {
        fetch("http://127.0.0.1:8080/replace", {
        method: "POST",
        body:`[{"op":"add", "path": "/${playerName}", "value": {"x":50, "y":500}}]`})
    }

    const updatePlayer = function(player:any) {
        fetch("http://127.0.0.1:8080/replace", {
        method: "POST",
        body:`[{"op":"add", "path": "/${playerName}", "value": {"x":${player.x}, "y":${player.y}}}]`})
    }
    const handleUp = function() {
        let player = snap[playerName]
        player.y -= 5
        updatePlayer(player)
    }

    const handleDown = function() {
        let player = snap[playerName]
        player.y += 5
        updatePlayer(player)
    }

    const handleLeft = function() {
        let player = snap[playerName]
        player.x -= 5
        updatePlayer(player)
    }

    const handleRight = function() {
        let player = snap[playerName]
        player.x += 5
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
    for (const k in snap) {
        let player = snap[k]
        let param = 'translate(' + player.x + ' ' + player.y + ')'
        console.log(param)
        players.push(
            <g transform={param}>
                <Copatych name={k}></Copatych>
            </g>
        )
        // players.push(<circle cx={player.x} cy={player.y} r="20" stroke="white" fill="transparent" />)
    }
    return <>{players}</>
}

ReactDOM.render(
    <React.StrictMode>
      <PlayerInput></PlayerInput>
      <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 800 600" fill='none'>
        <MainField></MainField>
        {/* <rect x='0' y='0' width='100%' height='100%' fill='tomato' opacity='0.75' /> */}
        <Players></Players>
      </svg>
       {/* <MainField>
        <Players></Players>
      </MainField> */}
      {/* <MainField></MainField> */}
    </React.StrictMode>,
    document.getElementById("root")
  )