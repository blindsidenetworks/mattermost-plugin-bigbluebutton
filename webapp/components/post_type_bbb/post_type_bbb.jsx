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

const {
  Button,
  Modal,
  Thumbnail,
  Grid,
  Col,
  Row,
  Panel
} = window.ReactBootstrap;
import {bootstrapUtils} from 'react-bootstrap/lib/utils';
import {getChannel} from 'mattermost-redux/selectors/entities/channels';
import {formatDate} from '../../utils/date_utils';

import PropTypes from 'prop-types';
import {makeStyleFromTheme} from 'mattermost-redux/utils/theme_utils';

const PostUtils = window.PostUtils;

export default class PostTypebbb extends React.PureComponent {
  static propTypes = {
    post: PropTypes.object.isRequired,
    state: PropTypes.object.isRequired,
    currentUserId: PropTypes.string.isRequired,
    creatorId: PropTypes.string.isRequired,
    channelId: PropTypes.string.isRequired,
    username: PropTypes.string.isRequired,
    channel: PropTypes.object.isRequired,
    compactDisplay: PropTypes.bool,
    isRHS: PropTypes.bool,
    useMilitaryTime: PropTypes.bool,
    theme: PropTypes.object.isRequired,
    creatorName: PropTypes.string.isRequired,
    actions: PropTypes.shape({
      getJoinURL: PropTypes.func.isRequired,
      endMeeting: PropTypes.func.isRequired,
      getAttendees: PropTypes.func.isRequired,
      publishRecordings: PropTypes.func.isRequired,
      deleteRecordings: PropTypes.func.isRequired,
      isMeetingRunning: PropTypes.func.isRequired
    }).isRequired

  };

  static defaultProps = {
    mentionKeys: [],
    compactDisplay: false,
    isRHS: false
  };

  constructor(props) {
    super(props);

    this.state = {
      url: "#",
      users: {},
      userCount: 0,
      show: false,
      showWarning: true,
      showFullAttendees: false,
      showThumbnails: false
    };
  }

//Purpose of this code was for when someone manually deletes, we end meeting on
// react component unmount. However switching channels also ends unmounts the component
// and we dont want to end the meeting

  // componentWillUnmount() {
  //   this.endMeetingForUnmount()
  // }

  componentDidMount() {
	  this.loadJoinUrl();
  }

  handleClose = () => {
    this.setState({show: false});
  };

  handleShow = () => {
    this.setState({show: true});
  };

  toggleAttendeesInMeeting = () => {
    this.setState({
      showFullAttendees: !this.state.showFullAttendees
    });
  }

  openJoinUrl = async () => {

    var userAgent = navigator.userAgent.toLowerCase();
    var url = this.state.url;
    //for electron apps
    if (userAgent.indexOf(' electron/') > -1) {
      window.open(url);
    } else { //for webapps to circumvent popup blockers
      var newtab = await window.open('about:blank');
      newtab.location = url;
      newtab.focus();
    }

    await this.setState({
      users: [
        ...this.state.users,
        this.props.username
      ]
    });
  }

  async loadJoinUrl() {
    const getJoinUrlResp = await this.props.actions.getJoinURL(this.props.channelId, this.props.post.props.meeting_id, this.props.creatorId);
    const joinUrl = await getJoinUrlResp.data.joinurl.url;

    await this.setState({
      url: joinUrl,
    });
  }

  isMeetingRunning = async (id) => {
    var response = await this.props.actions.isMeetingRunning(id);
    return response.running;
  }

  endMeeting = async () => {
    await this.props.actions.endMeeting(this.props.channelId, this.props.post.props.meeting_id);
  }

  endMeetingForUnmount = async () => {
    var isRunning = await this.isMeetingRunning(this.props.post.props.meeting_id);
    if (isRunning) {
      await this.props.actions.endMeeting(this.props.channelId, this.props.post.props.meeting_id);
    }

  }

  getAttendees = async () => {
    var isRunning = await this.isMeetingRunning(this.props.post.props.meeting_id);
    if (isRunning) {
      var resp = await this.props.actions.getAttendees(this.props.channelId, this.props.post.props.meeting_id);
      await this.setState({users: resp.attendees, userCount: resp.num});
      return resp.num;
    }return 0;

  }

  publishRecordings = async () => {
    await this.props.actions.publishRecordings(this.props.channelId, this.props.post.props.record_id, "true", this.props.post.props.meeting_id);
  }
  unpublishRecordings = async () => {
    await this.props.actions.publishRecordings(this.props.channelId, this.props.post.props.record_id, "false", this.props.post.props.meeting_id);
  }
  toggleThumbnails = () => {
    this.setState({
      showThumbnails: !this.state.showThumbnails
    })
  };

  deleteRecordings = async () => {
    await this.props.actions.deleteRecordings(this.props.channelId, this.props.post.props.record_id, this.props.post.props.meeting_id);
    this.setState({show: false});
  }

  render() {

    //overrides default Mattermost style with out own
    bootstrapUtils.addStyle(Button, 'custom');
    var arrayAttendants = [];
    const style = getStyle(this.props.theme);
    const post = this.props.post;
    const props = post.props || {};

    var attendees = "";
    var attendeesFull = "";
    var otherWords = "";

    if (props.attendees == undefined || props.attendees === "") {
      attendees = "there are no attendees in this session";
      //if we're on a direct message channel
      if (this.props.channel.type === "D") {
        var channel = getChannel(this.props.state, this.props.channelId);
        var channelName = channel.display_name;
        if (this.props.currentUserId === this.props.creatorId) {
          attendees = "Invited " + channelName + " to this meeting";
        }
      }
    } else {
      arrayAttendants = props.attendees.split(",");
      if (arrayAttendants != null && props.user_count > 0) {
        for (var i = 0; i < arrayAttendants.length; i++) {
          if (i <= 3) {

            attendees += arrayAttendants[i];
            if (i != arrayAttendants.length - 1) {
              attendees += ", ";
            }
          }
          attendeesFull += arrayAttendants[i];
          if (i != arrayAttendants.length - 1) {
            attendeesFull += ", ";
          }
        }
        if (arrayAttendants.length > 4) {
          otherWords = "and " + (
          arrayAttendants.length - 4) + " others";
        }

      } else {
        attendees = "there are no attendees in this session";
      }
    }

    let preText;
    let content;
    let subtitle;
    let activeUsers;
    let recordingstuff;
    var userlist = [];
    if (arrayAttendants != null || arrayAttendants != []) {
      for (var i = 0; i < arrayAttendants.length; i++) {
        userlist.push(<li>{arrayAttendants[i]}</li>);
      }
    }

    this.setState({userCount: props.user_count})

    if (props.meeting_status === 'STARTED') {

      preText = PostUtils.formatText("Meeting created by @" + this.props.creatorName, {
        mentionHighlight: false,
        atMentions: true
      });
      let attendeestext;
      if (this.state.showFullAttendees) {
        attendeestext = (<span onDoubleClick={this.toggleAttendeesInMeeting}>
          {attendeesFull}</span>);
      } else {
        attendeestext = (<span>
          <span>{attendees}</span>
          <span onClick={this.toggleAttendeesInMeeting}>{otherWords}</span>
        </span>);
      }
      content = (<div onMouseEnter={this.getAttendees}>
        <div>
          <span style={style.summary}>
            Attendees:
          </span>
          {
            (arrayAttendants != null && props.user_count > 0)
              ? <span style={style.summaryItem}>&ensp; {attendeestext}</span>
              : <span style={style.summaryItemGreyItalics}>
                  &ensp; {attendees}</span>
          }

        </div>
        <span >
          <a className='btn btn-lg btn-primary' style={style.button} onClick={this.openJoinUrl} href={this.state.url}>

            {'Join Meeting'}
          </a>
          {
            this.props.currentUserId == this.props.creatorId && <a className='btn btn-lg btn-link' style={style.buttonEnd} onClick={this.endMeeting} href='#'>
                <i style={style.buttonIcon}/> {'End meeting'}
              </a>
          }

        </span>
      </div>);
      if (props.meeting_desc != "") {
        subtitle = (<span>
          {'Description : '}

          {props.meeting_desc}

        </span>);
      }

    } else if (props.meeting_status === 'ENDED') {
      preText = PostUtils.formatText("@" + this.props.creatorName + " has ended the meeting", {
        mentionHighlight: false,
        atMentions: true
      });
      if (props.ended_by === "" || props.ended_by === undefined) {
        preText = `Meeting ended`;
      }
      if (props.meeting_desc != "") {
        subtitle = 'Description : ' + props.meeting_desc;
      }

      const startDate = new Date(post.create_at);
      const start = formatDate(startDate);
      const length = Math.ceil((new Date(post.update_at) - startDate) / 1000 / 60);
      var attendeestext;
      if (props.attendents == undefined || props.attendents === "") {
        attendees = "there were no attendees in this session";
      } else {
        var arrayAttendants = props.attendents.split(",");
        attendees = "";
        attendeesFull = "";
        otherWords = "";

        if (arrayAttendants != null && arrayAttendants.length > 0) {
          for (var i = 0; i < arrayAttendants.length; i++) {
            if (i <= 3) {
              attendees += arrayAttendants[i];
              if (i != arrayAttendants.length - 1) {
                attendees += ", ";
              }
            }
            attendeesFull += arrayAttendants[i];
            if (i != arrayAttendants.length - 1) {
              attendeesFull += ", ";
            }
          }
          if (arrayAttendants.length > 4) {
            otherWords = "and " + (
            arrayAttendants.length - 4) + " others";
          }
        } else {
          attendees = "there were no attendees in this session";
        }

        if (this.state.showFullAttendees) {
          attendeestext = (<span onDoubleClick={this.toggleAttendeesInMeeting}>
            {attendeesFull}</span>);
        } else {
          attendeestext = (<span>
            <span>{attendees}</span>
            <span onClick={this.toggleAttendeesInMeeting}>{otherWords}</span>
          </span>);
        }
      }

      content = (<div>
        <span>
          <span style={style.summary}>{'Date: '}</span>
          <span style={style.summaryItem}>{'Started at ' + start}</span>
        </span>
        &emsp;&emsp;
        <span>
          <span style={style.summary}>{'Meeting Length: '}</span>
          <span style={style.summaryItem}>{props.duration}</span>
        </span>
        &emsp;&emsp;
        <span style={style.summary}>
          Attendees:
        </span>
        {
          (props.attendents != undefined && arrayAttendants != null && arrayAttendants.length > 0)
            ? <span style={style.summaryItem}>&ensp; {attendeestext}</span>
            : <span style={style.summaryItemGreyItalics}>
                &ensp; {attendees}</span>
        }

      </div>);

      if (props.recording_status === 'COMPLETE' && (props.is_deleted == undefined || props.is_deleted != "true")) {

        var images = [];
        if (props.images != undefined && props.images != "" && typeof props.images === 'string') {
          var imagesArray = props.images.split(",");

          for (var i = 0; i < imagesArray.length; i++) {
            images.push(<Col sm={3} xs={3} md={2} lg={2}>
              <Thumbnail href={props.recording_url} responsive="responsive" src={imagesArray[i]}/>
            </Col>);
          }

        }
        recordingstuff = (<div>
          <div style={style.summaryRecording}>Recording
          </div>
          <div style={style.recordingBody}>
            <div>
              {
                props.is_published === "true"
                  ? <a href={props.recording_url} target="_blank">
                      {'Click to view recording'}
                    </a>
                  : <span style={style.summaryItemGreyItalics}>
                      Recording currently not viewable
                    </span>
              }

            </div>
            <div style={style.extraPadding}>
              {
                (this.props.currentUserId == this.props.creatorId)
                  ? <div>
                      {
                        props.is_published === "true"
                          ? <a onClick={this.unpublishRecordings}>
                              <span>
                                Make recording invisible
                              </span>
                            </a>
                          : <a onClick={this.publishRecordings}>
                              <span>
                                Show recording
                              </span>
                            </a>
                      }
                      <span style={style.bluebar}>
                        {'   |   '}
                      </span>
                      <a onClick={this.handleShow}>
                        <span>
                          Delete recording
                        </span>
                      </a>
                      {
                        props.is_published === "true" && <span>
                            <span style={style.bluebar}>
                              {'   |   '}
                            </span>

                            <a onClick={this.toggleThumbnails}>
                              <span>
                                Thumbnails
                              </span>
                            </a>
                          </span>
                      }
                    </div>

                  : <div>
                      {
                        props.is_published === "true" && <span>
                            <a onClick={this.toggleThumbnails}>
                              <span>
                                Thumbnails
                              </span>
                            </a>
                          </span>
                      }
                    </div>
              }
            </div>
          </div>
          {
            props.is_published === "true" && <div>
                <Panel expanded={this.state.showThumbnails}>
                  <Panel.Collapse>
                    <Panel.Body>
                      <span>
                        <Grid>
                          <Row>
                            {images}
                          </Row>
                        </Grid>
                      </span>
                    </Panel.Body>
                  </Panel.Collapse>
                </Panel>
              </div>
          }

        </div>);

      }
    }

    return (<div style={style.attachment}>
      {PostUtils.messageHtmlToComponent(preText)}
      <div style={style.content}>
        <div style={style.container}>
          <div style={style.body}>
            {content}
            {recordingstuff}
          </div>
        </div>
      </div>

      <Modal show={this.state.show} onHide={this.handleClose} bsSize="small">
        <Modal.Header closeButton="closeButton">
          <Modal.Title>Are You Sure?</Modal.Title>
        </Modal.Header>
        <Modal.Body>
          <p>
            Once deleted, the recording will be gone forever.
          </p>
        </Modal.Body>
        <Modal.Footer>
          <span>
            <Button onClick={this.handleClose}>Close</Button>
          </span>
          <span>
            <Button bsStyle="danger" onClick={this.deleteRecordings}>Delete Recording</Button>
          </span>
        </Modal.Footer>
      </Modal>
    </div>);
  }
}

const getStyle = makeStyleFromTheme((theme) => {
  return {

    attachment: {
      marginLeft: '-5px',
      position: 'relative'
    },
    content: {
      marginTop: '8px',
      borderRadius: '4px',
      borderStyle: 'solid',
      borderWidth: '0px',
      borderColor: '#BDBDBF',

    },
    container: {
      borderLeftStyle: 'solid',
      borderLeftWidth: '2px',
      paddingLeft: '10px',
      paddingBottom: '5px',
      paddingTop: '5px',
      borderLeftColor: theme.buttonBg
    },
    body: {
      overflowX: 'auto',
      overflowY: 'hidden',
      paddingRight: '5px',

      width: '100%'
    },

    button: {
      fontFamily: 'Open Sans',
      fontSize: '13px',
      lineHeight: '13px',
      marginTop: '10px',
      marginRight: '2px',
      borderRadius: '4px',
      color: theme.buttonColor
    },

    buttonEnd: {
      fontFamily: 'Open Sans',
      fontSize: '13px',
      lineHeight: '13px',
      marginTop: '10px',
      marginRight: '2px',
      borderRadius: '4px',
      color: theme.buttonBg
    },
    extraPadding: {
      marginTop: '10px'
    },

    summary: {
      fontFamily: 'Open Sans',
      fontSize: '14px',
      fontWeight: '600'
    },
    summaryRecording: {
      fontFamily: 'Open Sans',
      fontSize: '14px',
      fontWeight: '600',
      lineHeight: '26px'
    },
    summaryItem: {
      fontFamily: 'Open Sans',
      fontSize: '14px',
    },
    recordingBody: {
      lineHeight: '26px'
    },

    bluebar: {
      color: '#008BD2',
      fontWeight: '400'
    },
    summaryItemGreyItalics: {
      fontFamily: 'Open Sans',
      fontSize: '14px',
      fontStyle: 'italic',
      color: '#8D8D94'
    }
  };
});
