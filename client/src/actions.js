import { newSession } from './api.js'

export const createSession = function(hookId){
  return (dispatch) => {
    newSession(hookId).then((socket) => {
      dispatch({ type: 'BEGIN_CREATE_SESSION' })

      socket.onmessage = (e) => {
        const evtData  = JSON.parse(e.data)
        if (evtData.type) {
          dispatch({ type: evtData.type, payload: evtData.payload })
        }
      }
    })

  }
}

function createSessionReducer(state, action){
  const newState = {...state}

  switch (action.type) {
    case "CREATE_SESSION":
      const { hookId, entries } = action.payload
      if (hookId !== newState.hookId) {
        window.location.hash = `#/${hookId}`
      }

      newState.hookId = hookId
      newState.entries = !entries || entries == null ? [] : entries
      break
  }

  return newState
}

function addRequestReducer(state, action){
  const newState = {...state}

  switch (action.type) {
    case "NEW_REQUEST":
      newState.entries.push(action.payload)
      break
  }

  return newState
}

export const reducers = [createSessionReducer, addRequestReducer]

function unwrapAndDispatch(type, dispatch, fetchP, map){
  dispatch({ type: `BEGIN_${type}` });
  fetchP.then(
      (data) => {
        if (data.success)
          dispatch({ type, payload: map ? map(data.payload) : data.payload })
        else
          dispatch({ type: `FAILED_${type}`, error: data.error })
      }
    ).catch((error) => dispatch({ type: `FAILED_NETWORK_${type}`, error }))
}
