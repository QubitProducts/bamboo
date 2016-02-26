import _ from 'lodash'
import React from 'react'
import cx from 'classnames'
import axios from 'axios'

import { NEW, EDIT } from './constants'

export default React.createClass({
  getInitialState () {
    const config = _.isEmpty(this.props.config) ? [["", ""]] : _.toPairs(this.props.config)

    return {
      id: this.props.id || "",
      config: config
    }
  },

  closeForm () {
    this.props.onRequestClose()
  },

  submitForm () {
    const payload = {
      id: this.state.id,
      config: _.fromPairs(this.state.config)
    }
    axios.post('/api/services', payload).then(() => {
      this.closeForm()
      this.props.onUpdate()
    })
  },

  handleIdChange (ev) {
    this.setState({id: ev.target.value})
  },

  removeConfigEntry (i) {
    return () => {
      let config = _.cloneDeep(this.state.config)
      config.splice(i, 1)
      if (config.length === 0) {
        config = [["", ""]]
      }
      this.setState({config})
    }
  },

  changeConfigKey (i) {
    return (ev) => {
      let config = _.cloneDeep(this.state.config)
      config[i][0] = ev.target.value

      this.setState({config})
    }
  },

  changeConfigValue (i) {
    return (ev) => {
      let config = _.cloneDeep(this.state.config)
      config[i][1] = ev.target.value

      this.setState({config})
    }
  },

  addConfigEntry () {
    let config = _.cloneDeep(this.state.config)
    config.push(["", ""])

    this.setState({config})
  },

  configTable () {
    return _.map(
      this.state.config,
      ([k, v], i) => {
        const kPlaceholder = i == 0 ? "Acl" : ""
        const vPlaceholder = i == 0 ? "hdr(host) -i app.example.com" : ""
        return (
          <tr key={i}>
            <td><input className="form-control"
                       type="text"
                       value={k}
                       placeholder={kPlaceholder}
                       onChange={this.changeConfigKey(i)}/></td>
            <td><input className="form-control"
                       type="text"
                       value={v}
                       placeholder={vPlaceholder}
                       onChange={this.changeConfigValue(i)}/></td>
            <td>
              <button className="btn btn-danger" onClick={this.removeConfigEntry(i)}>
                <i className="icon ion-android-trash"></i>
              </button>
            </td>
          </tr>
        )})
  },

  render () {
    return (
      <div className="modal-content">
        <div className="modal-header">
          <button type="button" className="close" onClick={this.closeForm}>x</button>
          <h4 className="modal-title">
            {this.props.mode === EDIT ? 'Update' : 'Create new'} Service Configuration
          </h4>
        </div>

        <div className="modal-body edit-form-body">
          <div className="inner-form">
            <div className="form-group">
              <label>Marathon ID</label>
              <input type="text"
                     className="form-control"
                     value={this.state.id}
                     onChange={this.handleIdChange}/>
            </div>

            <div className="form-group">
              <label>Config Values</label>
              <p className="help-block">
                Enter config values for the service. For the default HAProxy template, set the key to
                'Acl', and enter an ACL in the value field
                <br/>
                DNS approach: <code>hdr(host) -i app.example.com</code>
                <br/>
                Path prefix: <code>path_beg -i /app-group/app1</code>

                <button type="button" className="btn btn-primary pull-right" onClick={this.addConfigEntry}>
                  Add Entry
                </button>
              </p>

              <table className="table table-condensed">
                <thead>
                  <tr>
                    <th>Key</th>
                    <th>Value</th>
                    <th></th>
                  </tr>
                </thead>

                <tbody>
                  {this.configTable()}
                </tbody>
              </table>
            </div>
          </div>
        </div>

        <div className="modal-footer">
          <button type="button" className="btn btn-default" onClick={this.closeForm}>
            Close
          </button>
          <button type="button" className="btn btn-primary" onClick={this.submitForm}>
            {this.props.mode === EDIT ? 'Update' : 'Create'}
          </button>
        </div>
      </div>
    )
  }
})
