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
    if (this.status() === UNDEPLOYED) {
      return (
        <span className="item-actions-group">
          <i className="message">Missing app in Marathon</i>
          <EditButton mode={EDIT}
                      {...this.props} />
          <DeleteButton {...this.props} />
        </span>
      )
    }
    if (this.status() === UNCONFIGURED) {
      return (
        <span className="item-actions-group">
        <i className="message">Using default proxy rule</i>
        <EditButton mode={NEW} {...this.props}/>
        </span>
    )

    }
    if (this.status() === ACTIVE) {
      return (
        <span className="item-actions-group">
          <EditButton mode={EDIT}
                      {...this.props}/>
          <DeleteButton {...this.props}/>
        </span>
      )
    }
  },

  render () {
    let config
    if (this.status() === UNCONFIGURED) {
      config = ""
    } else {
      config = _.map(_.toPairs(this.props.config),
                     ([k, v]) => `${k}=${v}`
                    ).join(', ')
    }
    const taskCount = this.status() == UNDEPLOYED ? '-' : this.props.tasks.length
    const rowClasses = cx('row', 'service-item',
                          `service-action-type-${actionTypes[this.status()]}`)

    return (
      <div className={rowClasses}>
        <span className="col-xs-4">{this.props.id}</span>
        <span className="col-xs-4">{config}</span>
        <span className="col-xs-1 col-instance-count">{taskCount}</span>
        <span className="col-xs-3 item-actions">{this.actions()}</span>
      </div>
    )
  }
})
