import React, { Component } from 'react'
import RequestInfo from './RequestInfo.js'

import './RequestList.css'

export default class RequestList extends Component {
  render(){
    const { requests } = this.props;
    return (
        <ul className="RequestList">
          { (requests || []).length <= 0 ? (<li>
              <p>Send HTTP requests to the link at the top in the white box.</p>
              <p>Requests will show up here, newest at the top.</p>
             </li>) : "" }
          {
              (requests || []).map((r, idx) => (
                <li key={idx}>
                  <RequestInfo request={r} />
                </li>
              )).reverse()
          }
        </ul>
    )

  }
}
