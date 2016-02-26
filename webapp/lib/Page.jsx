import React from 'react'

import NavBar from './NavBar.jsx'
import ServiceList from './ServiceList.jsx'

export default React.createClass({
  render () {
    const serviceList = (<ServiceList pollInterval={5000} />)
    return (
      <div>
        <NavBar onUpdate={serviceList.handleServiceUpdate}/>
        {serviceList}
      </div>
    )
  }
})
