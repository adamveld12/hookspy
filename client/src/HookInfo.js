import React, { Component } from 'react'
import './HookInfo.css'


export default class HookInfo extends Component {
  render(){
    const { hookId } = this.props

    return (
      <div className="HookInfo">
        <textarea ref="input"
                  readOnly={true}
                  onClick={() =>{this.refs.input.select()}}
                  value={`${window.location.origin}/hook/${hookId}`}>
        </textarea>
      </div>
    )
  }
}
