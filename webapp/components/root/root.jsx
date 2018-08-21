/*
Copyright 2018 Blindside Networks

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// const React = window.React;
import React from 'react';

import PropTypes from 'prop-types';
import {makeStyleFromTheme, changeOpacity} from 'mattermost-redux/utils/theme_utils';
import {Link} from 'react-router-dom'
import {viewChannel, getChannelStats} from 'mattermost-redux/actions/channels';
import {isDirectChannel} from 'mattermost-redux/utils/channel_utils';
import {Tooltip, OverlayTrigger, Modal} from 'react-bootstrap';
import {Client4} from 'mattermost-redux/client';
import {getUser} from 'mattermost-redux/selectors/entities/users';
import {getChannel} from 'mattermost-redux/selectors/entities/channels';

export default class Root extends React.PureComponent {
  static propTypes = {

    cur_user: PropTypes.object.isRequired,
    state: PropTypes.object.isRequired,
    teamname: PropTypes.string.isRequired,
    lastpostperchannel: PropTypes.object.isRequired,
    unreadChannelIds: PropTypes.array.isRequired,
    theme: PropTypes.object.isRequired,
    actions: PropTypes.shape({getJoinURL: PropTypes.func.isRequired}).isRequired

  }

  constructor(props) {
    super(props);

    this.state = {
      ignoredPosts: [],
      show: false,
      channelId: "",
      channelName: "",
      meetingId: "",
      profilePicUrl: "",
      channelURL: ""
    };
  }

  handleClose = () => {
    this.setState({show: false});
  };

  openmodal = async (postid, channelid, meetingId, src) => {
    var channel = getChannel(this.props.state, channelid);
    var channelurl;
    if (channel.type === "D") {
      channelurl = "/messages/@" + channel.display_name
    } else if (channel.type === "G") {
      channelurl = "/messages/" + channel.name
    }
    await this.setState({
      ignoredPosts: [
        ...this.state.ignoredPosts,
        postid
      ],
      show: true,
      channelId: channelid,
      meetingId: meetingId,
      profilePicUrl: src,
      channelName: channel.display_name,
      channelURL: channelurl
    });

    console.log("open modal is called")
  };

  getJoinURL = async () => {
    var newtab = await window.open('https://blindsidenetworks.com/', '_blank');
    var myurl = await this.props.actions.getJoinURL(this.state.channelId, this.state.meetingId, "");
    var myvar = await myurl.data.joinurl.url;
    newtab.location.href = myvar;
  }
  getSiteUrl = () => {
    if (window.location.origin) {
      return window.location.origin;
    }
    return window.location.protocol + '//' + window.location.hostname + (
      window.location.port
      ? ':' + window.location.port
      : '');
  }

  render() {
    var gotoButton;
    var renderchannelid;
    var meetingid;
    var inviteuserid;
    var src = "";

    for (var i = 0; i < this.props.unreadChannelIds.length; i++) {

      console.log("root stuff is running")

      var channelid = this.props.unreadChannelIds[i];
      if (channelid in this.props.lastpostperchannel) {
        if (this.props.lastpostperchannel[channelid].type === "custom_bbb" && !this.state.ignoredPosts.includes(this.props.lastpostperchannel[channelid].id) && (Date.now() - this.props.lastpostperchannel[channelid].create_at < 2000) && this.props.lastpostperchannel[channelid].user_id != this.props.cur_user.id) {
          const postid = this.props.lastpostperchannel[channelid].id;
          const user = getUser(this.props.state, this.props.lastpostperchannel[channelid].user_id);
          src = Client4.getProfilePictureUrl(user.id, user.last_picture_update);
          renderchannelid = channelid;
          var message = this.props.lastpostperchannel[channelid].message;
          var index = message.indexOf('#ID');
          meetingid = message.substr(index + 3)
          console.log("screened for meeting")
          this.openmodal(postid, channelid, meetingid, src);

        }
      }
    }

    const style = getStyle(this.props.theme);
    const myteam = this.props.teamname
    const tooltip = (<Tooltip id="tooltip">
      Go to this channel
    </Tooltip>);
    if (!this.state.show){
      return (null);
    }
    return (
        <div
            style={style.backdrop}
            onClick={()=>{this.state.show = false}}
        >
            <div style={style.modal}>
                <div>
                        <div >
                          <img src={this.getSiteUrl() + this.state.profilePicUrl} class="img-responsive img-circle center-block "/>
                        </div>
                        
                        <div>
                          { 'You have triggered the root component of the demo plugin.' }
                          <br/>
                          <br/>
                          { 'Click anywhere to close.' }
                        </div>
                </div>

            </div>
        </div>
    );

    // return (<Modal show={this.state.show} onHide={this.handleClose}>
    //
    //   <Modal.Header closeButton={true} style={style.header}></Modal.Header>
    //
    //   <Modal.Body style={style.body}>
    //     <div >
    //       <div >
    //         <img src={this.getSiteUrl() + this.state.profilePicUrl} class="img-responsive img-circle center-block "/>
    //       </div>
    //       <div style={style.bodyText}>
    //         <span >
    //           BigBlueButton meeting request from
    //           <strong>
    //             <OverlayTrigger placement="top" overlay={tooltip}>
    //               <Link to={"/" + this.props.teamname + this.state.channelURL}>
    //                 {this.state.channelName}
    //               </Link>
    //             </OverlayTrigger>
    //
    //           </strong>
    //         </span>
    //       </div>
    //     </div>
    //   </Modal.Body>
    //   <Modal.Footer>
    //     <button type='button' className='btn btn-default' onClick={this.handleClose}>
    //       Close
    //
    //     </button>
    //
    //     <button type='button' className='btn btn-primary pull-left' onClick={this.getJoinURL}>
    //       Join Meeting
    //     </button>
    //
    //   </Modal.Footer>
    // </Modal>);
  }
}

/* Define CSS styles here */
const getStyle = makeStyleFromTheme((theme) => {
  return {

    header: {
      background: '#FFFFFF',
      color: '#0059A5',
      borderStyle: "none",
      height: "10px",
      minHeight: "28px"
    },
    body: {
      padding: '0px 0px 10px 0px'
    },
    bodyText: {
      textAlign: 'center',
      margin: '20px 0 0 0',
      fontSize: '17px',
      lineHeight: '19px'
    },
    meetingId: {
      marginTop: '55px'
    },
    backdrop: {
      position: 'absolute',
      display: 'flex',
      top: 0,
      left: 0,
      right: 0,
      bottom: 0,
      backgroundColor: 'rgba(0, 0, 0, 0.50)',
      zIndex: 2000,
      alignItems: 'center',
      justifyContent: 'center',
    },
    modal: {
      height: '250px',
      width: '400px',
      padding: '1em',
      color: theme.centerChannelColor,
      backgroundColor: theme.centerChannelBg,
    },
  };
});
