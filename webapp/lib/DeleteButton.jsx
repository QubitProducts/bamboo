import axios from 'axios'
import React from 'react'
import Modal from 'react-modal'

export default React.createClass({
  getInitialState () {
    return {isWarningOpen: false}
  },

  performDeletion () {
    axios.delete(`/api/services/${this.props.id}`)
  },

  openWarning () {
    this.setState({isWarningOpen: true})
  },
  closeWarning () {
    this.setState({isWarningOpen: false})
  },
  acceptWarning () {
    this.closeWarning()
    this.performDeletion()
    this.props.onUpdate()
  },

  warningModal () {
    // The default modal stylings are rather... opinionated
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
         isOpen={this.state.isWarningOpen}
         onRequestClose={this.closeWarning}>

        <div className="modal-content">
          <div className="modal-header">
            <button type="button" className="close" onClick={this.closeWarning}>x</button>
            <h4 className="modal-title">
              Are you sure?
            </h4>
          </div>
          <div className="modal-body">
            Delete Marathon ID {this.props.id}
          </div>
          <div className="modal-footer">
            <button type="button" className="btn btn-default" onClick={this.closeWarning}>
              Close
            </button>
            <button type="button" className="btn btn-primary" onClick={this.acceptWarning}>
              Delete it
            </button>
          </div>
        </div>

      </Modal>
    )
  },

  render () {
    return (
        <button className="btn btn-danger" onClick={this.openWarning}>
          <i className="icon ion-android-trash"></i>

          {this.warningModal()}
        </button>
    )
  }
})
