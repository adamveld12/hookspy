import Weedux, { middleware } from 'weedux'
const { thunk, logger } = middleware
import { reducers as actionReducers } from './actions.js'

const hookMatches = window.location.hash.match(/\#\/([\w\d-]+)\/?/)

var hookId = "";
if (hookMatches !== null && hookMatches.length > 0){
  hookId = hookMatches[1]
}

const initialState = {
  hookId,
  entries: []
}

export const store = new Weedux(initialState, actionReducers, [thunk, logger])
