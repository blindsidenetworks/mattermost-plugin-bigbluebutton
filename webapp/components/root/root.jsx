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

import React from 'react';
import PopoverListMembersItem from './popover_list_members_item.jsx';
import PropTypes from 'prop-types';
import {makeStyleFromTheme, changeOpacity} from 'mattermost-redux/utils/theme_utils';
import {Link} from 'react-router-dom'
import {viewChannel, getChannelStats} from 'mattermost-redux/actions/channels';
import {isDirectChannel} from 'mattermost-redux/utils/channel_utils';
const {Tooltip,Popover, OverlayTrigger, Modal,Overlay} = window.ReactBootstrap
import {Client4} from 'mattermost-redux/client';
import {getUser} from 'mattermost-redux/selectors/entities/users';
import {getChannel} from 'mattermost-redux/selectors/entities/channels';
import {searchPosts} from 'mattermost-redux/actions/search'
import * as UserUtils from 'mattermost-redux/utils/user_utils';

export default class Root extends React.PureComponent {
  static propTypes = {

    cur_user: PropTypes.object.isRequired,
    state: PropTypes.object.isRequired,
    teamname: PropTypes.string.isRequired,
    lastpostperchannel: PropTypes.object.isRequired,
    unreadChannelIds: PropTypes.array.isRequired,
    theme: PropTypes.object.isRequired,
    channelName: PropTypes.string.isRequired,
    channel: PropTypes.object.isRequired,
    actions: PropTypes.shape({getJoinURL: PropTypes.func.isRequired,
      channelId: PropTypes.string.isRequired,
      directChannels: PropTypes.array.isRequired,
      teamId: PropTypes.string.isRequired,
      visible: PropTypes.bool.isRequired,
      actions: PropTypes.shape({startMeeting: PropTypes.func.isRequired, showRecordings: PropTypes.func.isRequired, closePopover: PropTypes.func.isRequired}).isRequired,

    startMeeting: PropTypes.func.isRequired, showRecordings: PropTypes.func.isRequired, closePopover: PropTypes.func.isRequired}).isRequired

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

  searchRecordings = () => {
    this.props.actions.showRecordings();
  }

  startMeeting = async () => {
    await this.props.actions.startMeeting(this.props.channelId, "", this.props.channel.display_name);
    this.close_the_popover()
  }
  close_the_popover = () =>{
    this.props.actions.closePopover();
    this.setState({showPopover: false});
  }

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

  };

  getJoinURL = async () => {
    var userAgent = navigator.userAgent.toLowerCase();
    var myurl;
    var myvar;
    //for electron apps
    if (userAgent.indexOf(' electron/') > -1) {
      var myurl = await this.props.actions.getJoinURL(this.state.channelId, this.state.meetingId, "");
      myvar = await myurl.data.joinurl.url;
      window.open(myvar);
    }else{ //for webapps to circumvent popup blockers
      var newtab = await window.open('https://blindsidenetworks.com/', '_blank');
      var myurl = await this.props.actions.getJoinURL(this.state.channelId, this.state.meetingId, "");
      myvar = await myurl.data.joinurl.url;
      newtab.location.href = myvar;
    }
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
          this.openmodal(postid, channelid, meetingid, src);

        }
      }
    }
    let popoverButton = (<div className='more-modal__button'>

      <a className='btn  btn-link' onClick={this.searchRecordings}>

        {'View Recordings'}
      </a>

    </div>);

    var pos_width = (window.innerWidth - 400 + "px");
    var style = getStyle(pos_width,this.props.theme);

    style.popover["marginLeft"] = pos_width
    style.popoverDM["marginLeft"] = pos_width

    const myteam = this.props.teamname
    const tooltip = (<Tooltip id="tooltip">
      Go to this channel
    </Tooltip>);

    var channel = getChannel(this.props.state, this.props.channelId);
    // console.log(channel)
    var channelName = "";
    var ownChannel = false;
    if (channel == undefined){
       ownChannel = true;
    }
    else{
      channelName = channel.display_name;
    }


    return (
      <div>
      { !ownChannel && //shows popup when not on own channel
      <Overlay rootClose={true} show={this.props.visible}  onHide={this.close_the_popover} placement='bottom'>
        <Popover id='bbbPopover' style={this.props.channel.type === "D"
            ? style.popoverDM
            : style.popover}>
          <div style={this.props.channel.type === "D"
              ? style.popoverBodyDM
              : style.popoverBody}>
            {
              this.props.channel.type === "D"
                ? <PopoverListMembersItem onItemClick={this.startMeeting} cam={1} text={<span> {
                      'Call '
                    }
                    <strong>{channelName}</strong>
                  </span>} theme={this.props.theme}/>
                : <PopoverListMembersItem onItemClick={this.startMeeting} cam={1} text={<span> {
                      'Create a BigBlueButton Meeting'
                    }
                    </span>} theme={this.props.theme}/>
            }
          </div>
          {popoverButton}
        </Popover>
      </Overlay>
    }

      <Modal show={this.state.show} onHide={this.handleClose}>
      <Modal.Header closeButton={true} style={style.header}></Modal.Header>

      <Modal.Body style={style.body}>
        <div >
          <div >
            <img src={this.getSiteUrl() + this.state.profilePicUrl} class="img-responsive img-circle center-block "/>
          </div>
          <div style={style.bodyText}>
            <span >
              {"BigBlueButton meeting request from "}
              <strong>
                <OverlayTrigger placement="top" overlay={tooltip}>
                  <Link to={"/" + this.props.teamname + this.state.channelURL}>
                    {this.state.channelName}
                  </Link>
                </OverlayTrigger>
              </strong>
            </span>
          </div>
        </div>
      </Modal.Body>
      <Modal.Footer>
        <button type='button' className='btn btn-default' onClick={this.handleClose}>
          Close

        </button>

        <button type='button' className='btn btn-primary pull-left' onClick={this.getJoinURL}>
          Join Meeting
        </button>

      </Modal.Footer>
    </Modal>
  </div>);
  }
}

/* Define CSS styles here */
var getStyle = makeStyleFromTheme((theme) => {
  var x_pos = (window.innerWidth - 400 + "px"); //shouldn't be set here as it doesn't rerender
  return {
    popover: {
      marginLeft: x_pos,
      marginTop: "50px",
      maxWidth: '300px',
      height: '105px',
      width: '300px',
      background: theme.centerChannelBg
    },
    popoverDM: {
      marginLeft: x_pos,
      marginTop: "50px",
      maxWidth: '220px',
      height: '105px',
      width: '220px',
      background: theme.centerChannelBg
    },
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
    iconStyle: {
      position: 'relative',
      top: '-1px'
    },

    popoverBody: {
      maxHeight: '305px',
      overflow: 'auto',
      position: 'relative',
      width: '298px',
      left: '-14px',
      top: '-9px',
      borderBottom: '1px solid #D8D8D9'
    },

    popoverBodyDM: {
      maxHeight: '305px',
      overflow: 'auto',
      position: 'relative',
      width: '218px',
      left: '-14px',
      top: '-9px',
      borderBottom: '1px solid #D8D8D9'
    }
  };
});
