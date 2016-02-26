import _ from 'lodash'
import axios from 'axios'
import React from 'react'

import Service from './Service.jsx'

function sortServices (services) {
  const active = _.filter(services, (s) =>
                          !_.isUndefined(s.config) && !_.isUndefined(s.tasks))
          .sort((a, b) => a.id < b.id)
  const undeployed = _.filter(services, (s) =>
                              !_.isUndefined(s.config) && _.isUndefined(s.tasks))
          .sort((a, b) => a.id < b.id)
  const unconfigured = _.filter(services, (s) =>
                              _.isUndefined(s.config) && !_.isUndefined(s.tasks))
          .sort((a, b) => a.id < b.id)

  return [].concat(active, undeployed, unconfigured)
}

export default React.createClass({
  getInitialState () {
    return {services: []}
  },

  loadStateFromAPI () {
    axios.get('/api/state')
      .then((res) => {
        const byId = {}
        _.forEach(res.data.Services, (s) => {
          byId[s.Id] = {id: s.Id, config: s.Config}
        })
        _.forEach(res.data.Apps, (a) => {
          if (_.isUndefined(byId[a.Id])) {
            byId[a.Id] = {id: a.Id}
          }
          byId[a.Id].tasks = a.Tasks
        })

        const services = sortServices(_.values(byId))

        this.setState({services})
      })
  },
  componentDidMount () {
    this.loadStateFromAPI()
    setInterval(this.loadStateFromAPI, this.props.pollInterval);
  },

  handleServiceUpdate () {
    this.loadStateFromAPI()
  },

  headerRow () {
    return (
      <div className="row service-list-title">
        <div className="col-xs-4">Marathon ID</div>
        <div className="col-xs-4">Configuration</div>
        <div className="col-xs-1">Instances</div>
        <div className="col-xs-3"></div>
      </div>
    )
  },

  render () {
    const services = _.map(
      this.state.services,
      (s) => (<Service key={s.id} id={s.id} config={s.config} tasks={s.tasks}
              onUpdate={this.handleServiceUpdate}/>))

    return (
      <div className="container-fluid service-list">
        {this.headerRow()}
        {services}
      </div>
    )
  }
})
