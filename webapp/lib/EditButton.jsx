import _ from 'lodash'
import React from 'react'
import Modal from 'react-modal'
import cx from 'classnames'
import axios from 'axios'

import { NEW, EDIT } from './constants'

export default React.createClass({
  getInitialState () {
    const config = _.isEmpty(this.props.config) ? [["", ""]] : _.toPairs(this.props.config)

    return {
      isFormOpen: false,
      id: this.props.id || "",
      config: config
    }
  },

  openForm () {
    this.setState({isFormOpen: true})
  },
  closeForm () {
    this.setState(this.getInitialState())
  },

  submitForm () {
    console.log(this.state)
    const payload = {
      id: this.state.id,
      config: _.fromPairs(this.state.config)
    }
    axios.post('/api/services', payload).then((res) => {
      this.closeForm()
      this.props.onUpdate(res.data)
    })
  },

  handleIdChange (ev) {
    this.setState({id: ev.target.value})
  },

  removeConfigEntry (i) {
    return (ev) => {
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
    return _.map(this.state.config,
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

  formModal () {
    const modalStyle = {
      overlay: {
        backgroundColor: 'rgba(0, 0, 0, 0.5)'
      },
      content : {
        position                   : 'absolute',
        top                        : '0px',
        left                       : '0px',
        right                      : '0px',
        bottom                     : 'auto',
        border                     : '0px solid #ccc',
        borderRadius               : '0px',
        outline                    : 'none',
        padding                    : '0px'
      }
    }

    return (
      <Modal
         style={modalStyle}
         className="Modal__Bootstrap modal-dialog"
         isOpen={this.state.isFormOpen}
         onRequestClose={this.closeForm}>

        <div className="modal-content">
          <div className="modal-header">
            <button type="button" className="close" onClick={this.closeForm}>x</button>
            <h4 className="modal-title">
              {this.props.style === EDIT ? 'Update' : 'Create new'} Service Configuration
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
              {this.props.style === EDIT ? 'Update' : 'Create'}
            </button>
          </div>
        </div>

      </Modal>
    )

  },

  render () {
    let btnClasses
    let iconClasses
    if (this.props.style === NEW) {
      btnClasses = 'btn btn-primary btn-create-service'
      iconClasses = 'icon ion-plus-round'
    } else {
      btnClasses = 'btn btn-default'
      iconClasses = 'icon ion-compose'
    }

    let label = ''
    if (!_.isUndefined(this.props.label)) {
      label = ` ${this.props.label}`
    }

    return (
      <button className={btnClasses} onClick={this.openForm}>
        <i className={iconClasses}></i>
        {label}

        {this.formModal()}
      </button>
    )
  }
})
