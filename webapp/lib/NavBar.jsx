import React from 'react'

import EditButton from './EditButton.jsx'
import { NEW } from './constants'

export default React.createClass({
  render () {
    return (
      <nav className="navbar navbar-bamboo">
        <div className="container-fluid">
          <a className="navbar-brand">Bamboo</a>
          <EditButton style={NEW} label="New" onUpdate={this.props.onUpdate}/>
        </div>
      </nav>
    )
  }
})
