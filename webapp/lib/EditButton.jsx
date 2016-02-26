import _ from 'lodash'
import React from 'react'
import Modal from 'react-modal'

import { NEW } from './constants'
import EditForm from './EditForm.jsx'

export default React.createClass({
  getInitialState () {
    return {
      isFormOpen: false
    }
  },

  openForm () {
    this.setState({isFormOpen: true})
  },
  closeForm () {
    this.setState(this.getInitialState())
  },

  formModal () {
    const modalStyle = {
      overlay: {
        backgroundColor: 'rgba(0, 0, 0, 0.5)'
      },
      content : {
        position: 'absolute',
        top: '0px',
        left: '0px',
        right: '0px',
        bottom: 'auto',
        border: '0px solid #ccc',
        borderRadius: '0px',
        outline: 'none',
        padding: '0px'
      }
    }

    return (
      <Modal
         style={modalStyle}
         className="Modal__Bootstrap modal-dialog"
         isOpen={this.state.isFormOpen}
         onRequestClose={this.closeForm}>

        <EditForm {...this.props}
                  onRequestClose={this.closeForm}/>
      </Modal>
    )

  },

  render () {
    let btnClasses
    let iconClasses
    if (this.props.mode === NEW) {
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
