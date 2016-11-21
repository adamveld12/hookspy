import React, { Component } from 'react'
import { store } from './store.js'
import { createSession, openSession } from './actions.js'

import { newSession } from './api.js'

import RequestList from './RequestList.js'
import HookInfo from './HookInfo.js'

import './App.css'

const dispatcher = store.dispatcher()

export default class App extends Component {
  constructor(){
    super()
    console.log(store.store())
    this.state = store.store()
  }

  componentWillMount(){
    store.onDispatchComplete((state) => this.setState(state))

    // doesn't cut it at the moment, need to grab from the URL instead
    const { hookId } = this.state

    dispatcher(createSession(hookId))
  }

  render() {
    const { __SMAN__, entries, hookId } = this.state

    return (
      <div className="App">
        <HookInfo hookId={hookId} />
        <RequestList requests={ entries } />
        <StateDebugger actions={(__SMAN__ || {}).actions} />
      </div>
    )
  }
}


class StateDebugger extends Component {
  render(){
    const { actions } = this.props

    return (
      <ul style={{ listStyleType: "none" }}>
        { (actions || []).map((x, i) => (<li key={i} > {JSON.stringify(x.action)} </li>)) }
      </ul>
    )
  }
}
