import React, { Component } from 'react'
import './RequestInfo.css'

export default class RequestInfo extends Component {
  render(){
    const {
      created,
      proto,
      remoteAddr,
      method,
      headers,
      body
    } = this.props.request

    const time = new Date(created);
    const raw = false;

    return (
      <div className="RequestInfo">
        <h2 className="header">
          <p className="time">Recieved on: {time.toUTCString()}</p>
          <span>{method}&nbsp;{remoteAddr}&nbsp;{proto}</span>
        </h2>
        <div className="info">
          <div>
            <h3 className="label">Headers</h3>
            <pre>
             {
               Object.keys(headers).map((m) =>
                 <span key={`header-${m}`} id={`header-${m}`}>
                   <strong>{m}</strong>: <span>{headers[m] + "\n"}</span>
                 </span>
               )
             }
            </pre>
          </div>
          <div className="bodyInfo" style={ body === "" ? { "display": "none"} : {}}>
              <input type="checkbox" id={"body-" + created} defaultChecked={ true }/>
              <label className="label" htmlFor={"body-" + created}>
                <h3>Body</h3>
              </label>
              <pre>{ prettyPrint(body, headers['Content-Type'], raw) }</pre>
          </div>
        </div>
      </div>
    )
  }
}

function prettyPrint(body, contentType, raw){
  if (raw)
    return window.atob(body)

  if (contentType){
    if (contentType.match(/json/i) !== null)
      return JSON.stringify(JSON.parse(window.atob(body)), null, 2);
    else if (contentType.match(/form-urlencoded/i) !== null)
      return decodeURIComponent(window.atob(body))
  }
  
  return "Body was empty.";
}
