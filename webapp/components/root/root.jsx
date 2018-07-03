const React = window.react;
import PropTypes from 'prop-types';
import {makeStyleFromTheme,changeOpacity} from 'mattermost-redux/utils/theme_utils';
import { Link } from 'react-router-dom'
import {viewChannel, getChannelStats} from 'mattermost-redux/actions/channels';
import {isDirectChannel} from 'mattermost-redux/utils/channel_utils';
import {Well,Glyphicon,Button, ButtonGroup,ButtonToolbar, Tooltip, OverlayTrigger,Modal,Thumbnail,
Grid, Col, Row, Image}  from 'react-bootstrap';
import {browserHistory} from '../../utils/browser_history.jsx';
import {Client4} from 'mattermost-redux/client';
import {getUser} from 'mattermost-redux/selectors/entities/users';
import {Svgs} from '../../constants';
import {getChannel} from 'mattermost-redux/selectors/entities/channels';


export default class Root extends React.PureComponent {
    static propTypes = {


        cur_user:  PropTypes.object.isRequired,
        state: PropTypes.object.isRequired,
        teamname: PropTypes.string.isRequired,
        lastpostperchannel: PropTypes.object.isRequired,
        unreadChannelIds: PropTypes.array.isRequired,
        theme: PropTypes.object.isRequired,
         actions: PropTypes.shape({
             getJoinURL: PropTypes.func.isRequired,

         }).isRequired




    }

    constructor(props) {
        super(props);


        this.state = {
          ignoredPosts : [],
          show : false,
          channelId: "",
          channelName: "",
          meetingId: "",
          profilePicUrl: "",
          channelURL: "",
        };
    }


    handleClose = () => {
      this.setState({ show: false });
    };

    openmodal = async (postid,channelid,meetingId,src) => {
      var channel = getChannel(this.props.state, channelid);
      var channelurl;
      if (channel.type === "D"){
        channelurl = "/messages/@" + channel.display_name
      }else if (channel.type === "G"){
        channelurl = "/messages/" +channel.name
      }
      await this.setState({
        ignoredPosts:[...this.state.ignoredPosts, postid],
        show: true,
        channelId:channelid,
        meetingId: meetingId,
        profilePicUrl:src,
        channelName : channel.display_name,
        channelURL : channelurl,
      });
    };

    getJoinURL = async () => {
      var newtab =  await window.open('', '_blank');
      var myurl = await this.props.actions.getJoinURL(this.state.channelId, this.state.meetingId,"");
      var myvar = await myurl.data.joinurl.url;
      newtab.location.href = myvar;
    }
    getSiteUrl = () => {
      if (window.location.origin) {
      return window.location.origin;
      }
      return window.location.protocol + '//' + window.location.hostname + (window.location.port ? ':' + window.location.port : '');
    }


    render() {
       console.log(this.props.lastpostperchannel);
      console.log(this.props.unreadChannelIds);
        var gotoButton;
        var renderchannelid;
        var meetingid;
        var inviteuserid;
        var src = "";

      for (var i = 0; i < this.props.unreadChannelIds.length; i++){
        var channelid = this.props.unreadChannelIds[i];
        if (channelid in this.props.lastpostperchannel ){
          if (this.props.lastpostperchannel[channelid].type === "custom_bbb" && !this.state.ignoredPosts.includes(this.props.lastpostperchannel[channelid].id)
        && (Date.now()-this.props.lastpostperchannel[channelid].create_at < 2000) &&this.props.lastpostperchannel[channelid].user_id != this.props.cur_user.id ){
              const postid = this.props.lastpostperchannel[channelid].id;
              const user = getUser(this.props.state, this.props.lastpostperchannel[channelid].user_id);
              console.log("ussrrr " + user);
              src =Client4.getProfilePictureUrl(user.id, user.last_picture_update);
              console.log("srrrccc " + src);
              renderchannelid = channelid;
              var message = this.props.lastpostperchannel[channelid].message;
              var index = message.indexOf('#ID');
               meetingid = message.substr(index+3)
              console.log("MEETING ID " + meetingid);
            this.openmodal(postid,channelid,meetingid,src);

          }
        }
      }


      const style = getStyle(this.props.theme);
      const myteam = this.props.teamname
      const tooltip = (
        <Tooltip id="tooltip">
          Go to this channel
          </Tooltip>
        );

      return (
        <Modal
          show={this.state.show}
          onHide={this.handleClose}
          >
          <Modal.Header  closeButton={true} style={style.header}>
        </Modal.Header>

        <Modal.Body style = {style.body}>
          <div >
            <div >
              <img
                src={this.getSiteUrl() + this.state.profilePicUrl}
                class="img-responsive img-circle center-block " />
            </div>
            <div style={style.bodyText}>
              <span >
                BigBlueButton meeting request from <strong>
                <OverlayTrigger placement="top" overlay={tooltip}>
                  <Link to= {"/"+this.props.teamname + this.state.channelURL}>
                        {this.state.channelName}
                    </Link>
                </OverlayTrigger>


              </strong>
            </span>
            </div>
          </div>
        </Modal.Body>
        <Modal.Footer>
          <button
            type='button'
            className='btn btn-default'
            onClick={this.handleClose}
            >
            Close

          </button>

          <button
            type='button'
            className='btn btn-primary pull-left'
            onClick={this.getJoinURL}
            >
            Join Meeting
          </button>

        </Modal.Footer>
      </Modal>
      );
    }
}

/* Define CSS styles here */
const getStyle = makeStyleFromTheme((theme) => {
  return {

      header:{
        background:'#FFFFFF',
        color: '#0059A5',
        borderStyle: "none",
        height: "10px",
        minHeight: "28px",
      },
      body:{
        padding: '0px 0px 10px 0px',
      },
      bodyText:{
        textAlign: 'center',
        margin: '20px 0 0 0',
        fontSize: '17px',
        lineHeight: '19px'
      },
      meetingId: {
          marginTop: '55px'
      }
  };
});
