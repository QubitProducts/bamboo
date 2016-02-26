import _ from 'lodash'
import React from 'react'
import cx from 'classnames'

import { NEW, EDIT, UNCONFIGURED, UNDEPLOYED, ACTIVE } from './constants'
import EditButton from './EditButton.jsx'
import DeleteButton from './DeleteButton.jsx'

const actionTypes = {
  UNCONFIGURED: 'service',
  UNDEPLOYED: 'marathon',
  ACTIVE: 'default'
}

export default React.createClass({
  status () {
    if (_.isUndefined(this.props.tasks)) {
      return UNDEPLOYED
    }
    if (_.isUndefined(this.props.config)) {
      return UNCONFIGURED
    }
    return ACTIVE
  },

  actions () {
    let items
    if (this.status() === UNDEPLOYED) {
      return (
        <span className="item-actions-group">
          <i className="message">Missing app in Marathon</i>
          <EditButton id={this.props.id}
                      config={this.props.config}
                      style={EDIT}
                      onUpdate={this.props.onUpdate} />
          <DeleteButton id={this.props.id}
                        onUpdate={this.props.onUpdate}/>
        </span>
      )
    }
    if (this.status() === UNCONFIGURED) {
      return (
        <span className="item-actions-group">
        <i className="message">Using default proxy rule</i>
        <EditButton id={this.props.id} style={NEW} onUpdate={this.props.onUpdate}/>
        </span>
    )

    }
    if (this.status() === ACTIVE) {
      return (
        <span className="item-actions-group">
          <EditButton id={this.props.id}
                      config={this.props.config}
                      style={EDIT}
                      onUpdate={this.props.onUpdate}/>
          <DeleteButton id={this.props.id}
                        onUpdate={this.props.onUpdate}/>
        </span>
      )
    }
  },

  render () {
    const acl = this.status() === UNCONFIGURED ? '' : this.props.config.Acl
    const taskCount = this.status() == UNDEPLOYED ? '-' : this.props.tasks.length
    const rowClasses = cx('row', 'service-item',
                          `service-action-type-${actionTypes[this.status()]}`)

    return (
      <div className={rowClasses}>
        <span className="col-xs-4">{this.props.id}</span>
        <span className="col-xs-4">{acl}</span>
        <span className="col-xs-1 col-instance-count">{taskCount}</span>
        <span className="col-xs-3 item-actions">{this.actions()}</span>
      </div>
    )
  }
})
